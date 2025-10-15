-- ==================================================================================
-- PRODUCTION HOTFIX: Apply migration 010 (canonical package references)
-- ==================================================================================
--
-- This script manually applies migration 010 to fix the package format issue.
--
-- CONTEXT:
-- - Production has migration 9 = "009_separate_official_metadata" (old numbering)
-- - Code expects migration 9 = "009_migrate_canonical_package_refs"
-- - This migration was never applied due to version number mismatch
--
-- SAFETY: This migration is IDEMPOTENT and safe to run multiple times
--
-- TO APPLY:
-- Run this entire script in production database as user 'app'
-- ==================================================================================

BEGIN;

-- Helper function to convert OCI package to canonical reference format
CREATE OR REPLACE FUNCTION convert_oci_package_to_canonical(pkg jsonb)
RETURNS jsonb
LANGUAGE plpgsql
AS $$
DECLARE
    result jsonb;
    registry_url text;
    registry_host text;
    identifier text;
    version_str text;
    canonical_ref text;
BEGIN
    -- Start with the original package
    result := pkg;

    -- Only process OCI packages
    IF pkg->>'registryType' != 'oci' THEN
        RETURN result;
    END IF;

    -- Get current values
    registry_url := pkg->>'registryBaseUrl';
    identifier := pkg->>'identifier';
    version_str := pkg->>'version';

    -- Skip if already in canonical format (no registryBaseUrl and no separate version)
    -- Canonical format has registry in identifier like "ghcr.io/owner/repo:tag"
    IF registry_url IS NULL AND (identifier LIKE '%:%' OR identifier LIKE '%@sha256:%') THEN
        RETURN result;
    END IF;

    -- Extract registry host from registryBaseUrl
    IF registry_url IS NOT NULL THEN
        registry_host := regexp_replace(registry_url, '^https?://', '');
    ELSE
        -- Default to docker.io if no registry specified
        registry_host := 'docker.io';
    END IF;

    -- Build canonical reference
    -- Format: registry/namespace/image:tag
    IF version_str IS NOT NULL AND version_str != '' THEN
        canonical_ref := registry_host || '/' || identifier || ':' || version_str;
    ELSE
        -- Default to 'latest' tag if no version specified
        canonical_ref := registry_host || '/' || identifier || ':latest';
    END IF;

    -- Update identifier to canonical reference
    result := jsonb_set(result, '{identifier}', to_jsonb(canonical_ref));

    -- Remove registryBaseUrl field (no longer needed for OCI)
    result := result - 'registryBaseUrl';

    -- Remove version field (now part of identifier for OCI)
    result := result - 'version';

    RETURN result;
END;
$$;

-- Helper function to convert MCPB package to canonical reference format
CREATE OR REPLACE FUNCTION convert_mcpb_package_to_canonical(pkg jsonb)
RETURNS jsonb
LANGUAGE plpgsql
AS $$
DECLARE
    result jsonb;
BEGIN
    -- Start with the original package
    result := pkg;

    -- Only process MCPB packages
    IF pkg->>'registryType' != 'mcpb' THEN
        RETURN result;
    END IF;

    -- Remove version field if it exists (version is embedded in the download URL)
    result := result - 'version';

    -- Remove registryBaseUrl field if it exists (not needed for MCPB)
    result := result - 'registryBaseUrl';

    RETURN result;
END;
$$;

-- Helper function to remove package-type-specific forbidden fields
CREATE OR REPLACE FUNCTION remove_forbidden_fields(pkg jsonb)
RETURNS jsonb
LANGUAGE plpgsql
AS $$
DECLARE
    result jsonb;
    registry_type text;
BEGIN
    result := pkg;
    registry_type := pkg->>'registryType';

    -- Remove fileSha256 from non-MCPB packages (it's MCPB-only)
    IF registry_type != 'mcpb' AND result ? 'fileSha256' THEN
        result := result - 'fileSha256';
    END IF;

    RETURN result;
END;
$$;

-- Helper function to ensure transport field exists
CREATE OR REPLACE FUNCTION ensure_transport_field(pkg jsonb)
RETURNS jsonb
LANGUAGE plpgsql
AS $$
DECLARE
    result jsonb;
    registry_type text;
    runtime_hint text;
BEGIN
    result := pkg;

    -- If transport already exists, return as-is
    IF result ? 'transport' THEN
        RETURN result;
    END IF;

    -- Get registry type and runtime hint
    registry_type := pkg->>'registryType';
    runtime_hint := pkg->>'runtimeHint';

    -- Add default transport based on package type
    -- Most packages use stdio transport
    result := jsonb_set(result, '{transport}', '{"type": "stdio"}'::jsonb);

    RETURN result;
END;
$$;

-- Helper function to convert all packages in a packages array
CREATE OR REPLACE FUNCTION convert_packages_array(packages jsonb)
RETURNS jsonb
LANGUAGE plpgsql
AS $$
DECLARE
    result jsonb := '[]'::jsonb;
    pkg jsonb;
BEGIN
    -- Handle null or empty arrays
    IF packages IS NULL OR jsonb_array_length(packages) = 0 THEN
        RETURN packages;
    END IF;

    -- Process each package
    FOR pkg IN SELECT * FROM jsonb_array_elements(packages)
    LOOP
        -- First convert OCI packages to canonical format
        pkg := convert_oci_package_to_canonical(pkg);

        -- Then convert MCPB packages to canonical format
        pkg := convert_mcpb_package_to_canonical(pkg);

        -- Remove package-type-specific forbidden fields
        pkg := remove_forbidden_fields(pkg);

        -- Finally ensure transport field exists
        pkg := ensure_transport_field(pkg);

        -- Add to result array
        result := result || jsonb_build_array(pkg);
    END LOOP;

    RETURN result;
END;
$$;

-- Show affected servers before migration
DO $$
DECLARE
    affected_count int;
BEGIN
    SELECT COUNT(*) INTO affected_count
    FROM servers
    WHERE value ? 'packages'
      AND value->'packages' IS NOT NULL
      AND jsonb_array_length(value->'packages') > 0
      AND EXISTS (
        SELECT 1 FROM jsonb_array_elements(value->'packages') AS pkg
        WHERE pkg ? 'registryBaseUrl' OR pkg ? 'version'
      );

    RAISE NOTICE 'About to migrate % servers with old package format', affected_count;
END $$;

-- Migrate all server packages to canonical format
UPDATE servers
SET value = jsonb_set(
    value,
    '{packages}',
    convert_packages_array(value->'packages')
)
WHERE value ? 'packages'
  AND value->'packages' IS NOT NULL
  AND jsonb_array_length(value->'packages') > 0;

-- Show result
DO $$
DECLARE
    remaining_count int;
BEGIN
    SELECT COUNT(*) INTO remaining_count
    FROM servers
    WHERE value ? 'packages'
      AND value->'packages' IS NOT NULL
      AND jsonb_array_length(value->'packages') > 0
      AND EXISTS (
        SELECT 1 FROM jsonb_array_elements(value->'packages') AS pkg
        WHERE pkg ? 'registryBaseUrl' OR pkg ? 'version'
      );

    RAISE NOTICE 'Migration complete. Remaining servers with old format: %', remaining_count;
END $$;

-- Clean up helper functions
DROP FUNCTION IF EXISTS convert_oci_package_to_canonical(jsonb);
DROP FUNCTION IF EXISTS convert_mcpb_package_to_canonical(jsonb);
DROP FUNCTION IF EXISTS remove_forbidden_fields(jsonb);
DROP FUNCTION IF EXISTS ensure_transport_field(jsonb);
DROP FUNCTION IF EXISTS convert_packages_array(jsonb);

-- Record migration as applied (using version 10 to avoid conflict with existing version 9)
INSERT INTO schema_migrations (version, name, applied_at)
VALUES (10, '010_migrate_canonical_package_refs', NOW())
ON CONFLICT (version) DO NOTHING;

COMMIT;

-- Show migration status
SELECT version, name, applied_at
FROM schema_migrations
ORDER BY version;
