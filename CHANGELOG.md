# MCP Registry Migration Changelog

This release implements two major breaking changes to improve the MCP Registry's consistency and usability. **All users must update their implementations to use the new field names and API endpoints.**

## 🔄 **BREAKING CHANGES**

### 1. Server ID Consistency ([#396](https://github.com/modelcontextprotocol/registry/issues/396))

**Problem:** Each server version had a unique ID, preventing version history tracking and server renaming.

**Solution:** Introduced consistent server identification across versions.

#### API Changes - **ACTION REQUIRED**

| **Old Endpoint** | **New Endpoint** | **Migration Required** |
|------------------|------------------|------------------------|
| `GET /v0/servers/{id}` | `GET /v0/servers/{server_id}` | ✅ Update all API calls to use `server_id` parameter |
| N/A | `GET /v0/servers/{server_id}/versions` | ✅ New endpoint to list all versions |
| N/A | `GET /v0/servers/{server_id}?version=1.0.0` | ✅ New query parameter for specific versions |

#### Registry Metadata Changes - **ACTION REQUIRED**

**Old Structure:**
```json
{
  "_meta": {
    "io.modelcontextprotocol.registry/official": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "published_at": "2024-01-01T00:00:00Z",
      "is_latest": true
    }
  }
}
```

**New Structure:**
```json
{
  "_meta": {
    "io.modelcontextprotocol.registry/official": {
      "serverId": "550e8400-e29b-41d4-a716-446655440000",
      "versionId": "773f9b2e-1a47-4c8d-b5e6-2f8d9c4a7b3e",
      "published_at": "2024-01-01T00:00:00Z",
      "is_latest": true
    }
  }
}
```

**Migration Actions:**
- ✅ Update code reading `_meta.io.modelcontextprotocol.registry/official.id` → use `serverId`
- ✅ Add support for `versionId` if you need version-specific operations
- ✅ Update API clients to use new endpoint URLs with `server_id`

### 2. JSON Schema Standardization to camelCase ([#428](https://github.com/modelcontextprotocol/registry/issues/428))

**Problem:** Inconsistent field naming between snake_case and camelCase across server.json.

**Solution:** Standardized all fields to camelCase per MCP specification.

#### Complete Field Migration Table - **ACTION REQUIRED**

| **Old Field (snake_case)** | **New Field (camelCase)** | **Context** |
|----------------------------|---------------------------|-------------|
| `registry_type` | `registryType` | Package configuration |
| `registry_base_url` | `registryBaseUrl` | Package configuration |
| `file_sha256` | `fileSha256` | Package configuration |
| `runtime_hint` | `runtimeHint` | Package configuration |
| `runtime_arguments` | `runtimeArguments` | Package configuration |
| `package_arguments` | `packageArguments` | Package configuration |
| `environment_variables` | `environmentVariables` | Package configuration |
| `is_required` | `isRequired` | Input/Argument configuration |
| `is_secret` | `isSecret` | Input configuration |
| `value_hint` | `valueHint` | Argument configuration |
| `is_repeated` | `isRepeated` | Argument configuration |
| `website_url` | `websiteUrl` | Server metadata |

#### Migration Examples

**Package Configuration:**
```json
// OLD - Will be rejected
{
  "package": {
    "registry_type": "npm",
    "registry_base_url": "https://registry.npmjs.org",
    "file_sha256": "abc123...",
    "runtime_hint": "node",
    "runtime_arguments": [...],
    "package_arguments": [...],
    "environment_variables": [...]
  }
}

// NEW - Required format
{
  "package": {
    "registryType": "npm",
    "registryBaseUrl": "https://registry.npmjs.org",
    "fileSha256": "abc123...",
    "runtimeHint": "node",
    "runtimeArguments": [...],
    "packageArguments": [...],
    "environmentVariables": [...]
  }
}
```

**Arguments Configuration:**
```json
// OLD - Will be rejected
{
  "runtime_arguments": [
    {
      "name": "port",
      "is_required": true,
      "is_repeated": false,
      "value_hint": "8080"
    }
  ]
}

// NEW - Required format
{
  "runtimeArguments": [
    {
      "name": "port",
      "isRequired": true,
      "isRepeated": false,
      "valueHint": "8080"
    }
  ]
}
```

**Environment Variables:**
```json
// OLD - Will be rejected
{
  "environment_variables": [
    {
      "name": "API_KEY",
      "is_required": true,
      "is_secret": true
    }
  ]
}

// NEW - Required format
{
  "environmentVariables": [
    {
      "name": "API_KEY",
      "isRequired": true,
      "isSecret": true
    }
  ]
}
```

**Server Metadata:**
```json
// OLD - Will be rejected
{
  "name": "my-server",
  "website_url": "https://example.com"
}

// NEW - Required format
{
  "name": "my-server",
  "websiteUrl": "https://example.com"
}
```

## 📋 **MIGRATION CHECKLIST**

### For Server Publishers:
- [ ] Update your `server.json` files to use camelCase field names
- [ ] Test server publishing with new CLI version
- [ ] Update any automation scripts that reference old field names
- [ ] Update documentation referencing old field names

### For API Consumers:
- [ ] Update API endpoint URLs from `/v0/servers/{id}` to `/v0/servers/{server_id}`
- [ ] Update code reading registry metadata from `id` to `serverId`/`versionId`
- [ ] Add support for new `/v0/servers/{server_id}/versions` endpoint if needed
- [ ] Update JSON parsing to expect camelCase field names
- [ ] Test with new API responses

### For Subregistry Operators:
- [ ] Update ETL processes for new field names
- [ ] Update server identification logic to use `serverId`
- [ ] Implement version tracking using `versionId`
- [ ] Update any custom validation to expect camelCase

## 🛠 **MIGRATION TOOLS**

The registry includes automatic migration for existing data, but you need to update your client code:

1. **Publisher CLI**: Update to latest version for camelCase support
2. **API Clients**: Update endpoint URLs and field name parsing
3. **Validation**: Update schemas to expect camelCase fields

## 🔧 **BACKWARD COMPATIBILITY**

⚠️ **No backward compatibility** - old field names and endpoints will return errors.

**Timeline:**
- Publishing with snake_case fields: **Rejected immediately**
- Old API endpoints: **Return 404 errors**
- Old registry metadata fields: **No longer populated**

## ✨ **NEW CAPABILITIES**

With these changes, you can now:
- Track all versions of a server with consistent identification
- Query specific server versions via `?version=` parameter
- List all versions of a server via new `/versions` endpoint
- Maintain cleaner, spec-compliant JSON with camelCase naming

## 📚 **UPDATED DOCUMENTATION**

- [OpenAPI Specification](docs/reference/api/openapi.yaml) - Updated endpoints and schemas
- [Server JSON Schema](docs/reference/server-json/server.schema.json) - New camelCase fields
- [Publishing Guide](docs/guides/publishing/publish-server.md) - Updated examples
- [API Usage Guide](docs/guides/consuming/use-rest-api.md) - New endpoint patterns