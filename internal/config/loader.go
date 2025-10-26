package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig loads configuration from files and environment variables
func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Configure file paths
	v.SetConfigName("registry")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")                   // Current directory
	v.AddConfigPath("/etc/mcp-registry/")  // System-wide config
	v.AddConfigPath("$HOME/.mcp-registry") // User-specific config

	// Read config file (optional - don't error if not found)
	if err := v.ReadInConfig(); err != nil {
		var configNotFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &configNotFoundErr) {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; proceed with env vars and defaults
	}

	// Environment variables
	v.SetEnvPrefix("MCP_REGISTRY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind environment variables explicitly for nested structures
	// This is necessary because AutomaticEnv() doesn't work with Unmarshal()
	bindEnvVars(v)

	// Unmarshal into config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Apply backwards compatibility transformations
	applyBackwardsCompatibility(&cfg)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default values for all configuration options
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.address", ":8080")
	v.SetDefault("server.readTimeout", "30s")
	v.SetDefault("server.writeTimeout", "30s")
	v.SetDefault("server.shutdownTimeout", "10s")

	// Database defaults
	v.SetDefault("database.type", "postgres")
	v.SetDefault("database.url", "postgres://localhost:5432/mcp-registry?sslmode=disable")
	v.SetDefault("database.maxConnections", 30)
	v.SetDefault("database.minConnections", 5)
	v.SetDefault("database.maxIdleTime", "30m")
	v.SetDefault("database.maxLifetime", "2h")
	v.SetDefault("database.connectionTimeout", "10s")
	v.SetDefault("database.autoMigrate", true)
	v.SetDefault("database.migrationsPath", "")
	v.SetDefault("database.sslMode", "prefer")
	v.SetDefault("database.encrypt", true)
	v.SetDefault("database.region", "us-east-1")
	v.SetDefault("database.tableName", "mcp-registry-servers")
	v.SetDefault("database.endpoint", "")
	v.SetDefault("database.seedFrom", "")

	// Auth defaults
	v.SetDefault("auth.github.clientID", "")
	v.SetDefault("auth.github.clientSecret", "")
	v.SetDefault("auth.jwt.privateKey", "")
	v.SetDefault("auth.oidc.enabled", false)
	v.SetDefault("auth.oidc.issuer", "")
	v.SetDefault("auth.oidc.clientID", "")
	v.SetDefault("auth.oidc.extraClaims", "")
	v.SetDefault("auth.oidc.editPerms", "")
	v.SetDefault("auth.oidc.publishPerms", "")

	// Feature flags defaults
	v.SetDefault("features.enableAnonymousAuth", false)
	v.SetDefault("features.enableRegistryValidation", true)

	// Telemetry defaults
	v.SetDefault("telemetry.version", "dev")
}

// bindEnvVars explicitly binds environment variables to config fields
// This is required for Viper to properly read env vars when using Unmarshal()
// We map camelCase config keys to SNAKE_CASE environment variables
func bindEnvVars(v *viper.Viper) {
	// Server
	_ = v.BindEnv("server.address", "MCP_REGISTRY_SERVER_ADDRESS")
	_ = v.BindEnv("server.readTimeout", "MCP_REGISTRY_SERVER_READ_TIMEOUT")
	_ = v.BindEnv("server.writeTimeout", "MCP_REGISTRY_SERVER_WRITE_TIMEOUT")
	_ = v.BindEnv("server.shutdownTimeout", "MCP_REGISTRY_SERVER_SHUTDOWN_TIMEOUT")

	// Database
	_ = v.BindEnv("database.type", "MCP_REGISTRY_DATABASE_TYPE")
	_ = v.BindEnv("database.url", "MCP_REGISTRY_DATABASE_URL")
	_ = v.BindEnv("database.maxConnections", "MCP_REGISTRY_DATABASE_MAX_CONNECTIONS")
	_ = v.BindEnv("database.minConnections", "MCP_REGISTRY_DATABASE_MIN_CONNECTIONS")
	_ = v.BindEnv("database.maxIdleTime", "MCP_REGISTRY_DATABASE_MAX_IDLE_TIME")
	_ = v.BindEnv("database.maxLifetime", "MCP_REGISTRY_DATABASE_MAX_LIFETIME")
	_ = v.BindEnv("database.connectionTimeout", "MCP_REGISTRY_DATABASE_CONNECTION_TIMEOUT")
	_ = v.BindEnv("database.autoMigrate", "MCP_REGISTRY_DATABASE_AUTO_MIGRATE")
	_ = v.BindEnv("database.migrationsPath", "MCP_REGISTRY_DATABASE_MIGRATIONS_PATH")
	_ = v.BindEnv("database.sslMode", "MCP_REGISTRY_DATABASE_SSL_MODE")
	_ = v.BindEnv("database.encrypt", "MCP_REGISTRY_DATABASE_ENCRYPT")
	_ = v.BindEnv("database.region", "MCP_REGISTRY_DATABASE_REGION")
	_ = v.BindEnv("database.tableName", "MCP_REGISTRY_DATABASE_TABLE_NAME")
	_ = v.BindEnv("database.endpoint", "MCP_REGISTRY_DATABASE_ENDPOINT")
	_ = v.BindEnv("database.seedFrom", "MCP_REGISTRY_DATABASE_SEED_FROM")

	// Auth - GitHub
	_ = v.BindEnv("auth.github.clientID", "MCP_REGISTRY_AUTH_GITHUB_CLIENT_ID")
	_ = v.BindEnv("auth.github.clientSecret", "MCP_REGISTRY_AUTH_GITHUB_CLIENT_SECRET")

	// Auth - JWT
	_ = v.BindEnv("auth.jwt.privateKey", "MCP_REGISTRY_AUTH_JWT_PRIVATE_KEY")

	// Auth - OIDC
	_ = v.BindEnv("auth.oidc.enabled", "MCP_REGISTRY_AUTH_OIDC_ENABLED")
	_ = v.BindEnv("auth.oidc.issuer", "MCP_REGISTRY_AUTH_OIDC_ISSUER")
	_ = v.BindEnv("auth.oidc.clientID", "MCP_REGISTRY_AUTH_OIDC_CLIENT_ID")
	_ = v.BindEnv("auth.oidc.extraClaims", "MCP_REGISTRY_AUTH_OIDC_EXTRA_CLAIMS")
	_ = v.BindEnv("auth.oidc.editPerms", "MCP_REGISTRY_AUTH_OIDC_EDIT_PERMS")
	_ = v.BindEnv("auth.oidc.publishPerms", "MCP_REGISTRY_AUTH_OIDC_PUBLISH_PERMS")

	// Features
	_ = v.BindEnv("features.enableAnonymousAuth", "MCP_REGISTRY_FEATURES_ENABLE_ANONYMOUS_AUTH")
	_ = v.BindEnv("features.enableRegistryValidation", "MCP_REGISTRY_FEATURES_ENABLE_REGISTRY_VALIDATION")

	// Telemetry
	_ = v.BindEnv("telemetry.version", "MCP_REGISTRY_TELEMETRY_VERSION")
	_ = v.BindEnv("telemetry.gitCommit", "MCP_REGISTRY_TELEMETRY_GIT_COMMIT")
	_ = v.BindEnv("telemetry.buildTime", "MCP_REGISTRY_TELEMETRY_BUILD_TIME")
}

// applyBackwardsCompatibility applies transformations for backwards compatibility
func applyBackwardsCompatibility(_ *Config) {
	// No transformations needed for initial migration
	// This function is a placeholder for future compatibility needs
}
