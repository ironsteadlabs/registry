package config

import "time"

// Config is the root configuration structure
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Auth      AuthConfig
	Features  FeatureFlags
	Telemetry TelemetryConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	// Address is the address the HTTP server listens on
	// Format: host:port or :port
	// Example: ":8080" or "localhost:8080"
	Address string `mapstructure:"address"`

	// ReadTimeout is the maximum duration for reading the entire request
	ReadTimeout time.Duration `mapstructure:"readTimeout"`

	// WriteTimeout is the maximum duration before timing out writes
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`

	// ShutdownTimeout is the maximum duration to wait for graceful shutdown
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	// Type is the database backend to use
	// Supported: postgres, sqlserver, mysql, sqlite, dynamodb
	Type string `mapstructure:"type"`

	// URL is the database connection string
	// Format varies by database type:
	// - postgres: postgres://user:pass@host:5432/dbname?sslmode=require
	// - sqlserver: sqlserver://user:pass@host:1433?database=dbname
	// - mysql: mysql://user:pass@host:3306/dbname
	// - sqlite: file:path/to/db.sqlite
	URL string `mapstructure:"url"`

	// MaxConnections is the maximum number of open connections to the database
	MaxConnections int `mapstructure:"maxConnections"`

	// MinConnections is the minimum number of idle connections
	MinConnections int `mapstructure:"minConnections"`

	// MaxIdleTime is the maximum time a connection can be idle
	MaxIdleTime time.Duration `mapstructure:"maxIdleTime"`

	// MaxLifetime is the maximum time a connection can be reused
	MaxLifetime time.Duration `mapstructure:"maxLifetime"`

	// ConnectionTimeout is the timeout for establishing connections
	ConnectionTimeout time.Duration `mapstructure:"connectionTimeout"`

	// AutoMigrate runs database migrations on startup
	AutoMigrate bool `mapstructure:"autoMigrate"`

	// MigrationsPath is the path to migration files
	MigrationsPath string `mapstructure:"migrationsPath"`

	// SSLMode controls SSL/TLS for PostgreSQL and MySQL
	// Values: disable, require, prefer, verify-ca, verify-full
	SSLMode string `mapstructure:"sslMode"`

	// Encrypt enables encryption for SQL Server connections
	Encrypt bool `mapstructure:"encrypt"`

	// Region is the AWS region for DynamoDB
	Region string `mapstructure:"region"`

	// TableName is the table name for DynamoDB
	TableName string `mapstructure:"tableName"`

	// Endpoint is the DynamoDB endpoint (for local development)
	Endpoint string `mapstructure:"endpoint"`

	// SeedFrom is the path to seed data (legacy, kept for backwards compatibility)
	SeedFrom string `mapstructure:"seedFrom"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	GitHub GithubAuthConfig `mapstructure:"github"`
	OIDC   OIDCConfig       `mapstructure:"oidc"`
	JWT    JWTConfig        `mapstructure:"jwt"`
}

// GithubAuthConfig holds GitHub OAuth configuration
type GithubAuthConfig struct {
	ClientID     string `mapstructure:"clientID"`
	ClientSecret string `mapstructure:"clientSecret"`
}

// OIDCConfig holds OIDC configuration
type OIDCConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	Issuer       string `mapstructure:"issuer"`
	ClientID     string `mapstructure:"clientID"`
	ExtraClaims  string `mapstructure:"extraClaims"`
	EditPerms    string `mapstructure:"editPerms"`
	PublishPerms string `mapstructure:"publishPerms"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	PrivateKey string `mapstructure:"privateKey"`
}

// FeatureFlags holds feature flag configuration
type FeatureFlags struct {
	EnableAnonymousAuth      bool `mapstructure:"enableAnonymousAuth"`
	EnableRegistryValidation bool `mapstructure:"enableRegistryValidation"`
}

// TelemetryConfig holds telemetry configuration
type TelemetryConfig struct {
	Version string `mapstructure:"version"`
}

// DatabaseType represents supported database types
type DatabaseType string

const (
	DatabaseTypePostgreSQL DatabaseType = "postgres"
	DatabaseTypeSQLServer  DatabaseType = "sqlserver"
	DatabaseTypeMySQL      DatabaseType = "mysql"
	DatabaseTypeSQLite     DatabaseType = "sqlite"
	DatabaseTypeDynamoDB   DatabaseType = "dynamodb"
)
