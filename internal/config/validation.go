package config

import (
	"errors"
	"fmt"
)

// Validate validates the entire configuration
func (c *Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config: %w", err)
	}

	if err := c.Auth.Validate(); err != nil {
		return fmt.Errorf("auth config: %w", err)
	}

	return nil
}

// Validate validates server configuration
func (sc *ServerConfig) Validate() error {
	if sc.Address == "" {
		return errors.New("server address is required")
	}

	if sc.ReadTimeout < 0 {
		return errors.New("server read timeout must be non-negative")
	}

	if sc.WriteTimeout < 0 {
		return errors.New("server write timeout must be non-negative")
	}

	if sc.ShutdownTimeout < 0 {
		return errors.New("server shutdown timeout must be non-negative")
	}

	return nil
}

// Validate validates database configuration
func (dc *DatabaseConfig) Validate() error {
	// Validate database type
	validTypes := map[DatabaseType]bool{
		DatabaseTypePostgreSQL: true,
		DatabaseTypeSQLServer:  true,
		DatabaseTypeMySQL:      true,
		DatabaseTypeSQLite:     true,
		DatabaseTypeDynamoDB:   true,
	}

	dbType := DatabaseType(dc.Type)
	if !validTypes[dbType] {
		return fmt.Errorf("unsupported database type: %s (supported: postgres, sqlserver, mysql, sqlite, dynamodb)", dc.Type)
	}

	// Validate URL is present for SQL databases
	if dbType != DatabaseTypeDynamoDB && dc.URL == "" {
		return errors.New("database URL is required")
	}

	// Validate connection pool settings
	if dc.MaxConnections < dc.MinConnections {
		return fmt.Errorf("max connections (%d) must be >= min connections (%d)",
			dc.MaxConnections, dc.MinConnections)
	}

	if dc.MinConnections < 0 {
		return errors.New("min connections must be non-negative")
	}

	if dc.MaxConnections < 1 {
		return errors.New("max connections must be at least 1")
	}

	// Validate DynamoDB-specific settings
	if dbType == DatabaseTypeDynamoDB {
		if dc.Region == "" {
			return errors.New("database region is required for DynamoDB")
		}
		if dc.TableName == "" {
			return errors.New("database table name is required for DynamoDB")
		}
	}

	// Validate SSL mode for databases that support it
	if dbType == DatabaseTypePostgreSQL || dbType == DatabaseTypeMySQL {
		validSSLModes := map[string]bool{
			"":            true, // Allow empty for optional setting
			"disable":     true,
			"require":     true,
			"prefer":      true,
			"verify-ca":   true,
			"verify-full": true,
		}
		if !validSSLModes[dc.SSLMode] {
			return fmt.Errorf("invalid SSL mode: %s (valid: disable, require, prefer, verify-ca, verify-full)", dc.SSLMode)
		}
	}

	return nil
}

// Validate validates auth configuration
func (ac *AuthConfig) Validate() error {
	// OIDC validation
	if ac.OIDC.Enabled {
		if ac.OIDC.Issuer == "" {
			return errors.New("OIDC issuer is required when OIDC is enabled")
		}
		if ac.OIDC.ClientID == "" {
			return errors.New("OIDC client ID is required when OIDC is enabled")
		}
	}

	return nil
}
