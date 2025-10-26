package config_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/registry/internal/config"
)

func TestNewConfig_WithDefaults(t *testing.T) {
	// Clear any existing env vars
	clearEnvVars()

	cfg := config.NewConfig()

	// Test server defaults
	if cfg.Server.Address != ":8080" {
		t.Errorf("expected server address :8080, got %s", cfg.Server.Address)
	}

	// Test database defaults
	if cfg.Database.Type != "postgres" {
		t.Errorf("expected database type postgres, got %s", cfg.Database.Type)
	}

	if cfg.Database.MaxConnections != 30 {
		t.Errorf("expected max connections 30, got %d", cfg.Database.MaxConnections)
	}
}

func TestNewConfig_WithEnvVars(t *testing.T) {
	// Clear any existing env vars
	clearEnvVars()

	// Set test env vars
	os.Setenv("MCP_REGISTRY_SERVER_ADDRESS", ":9090")
	os.Setenv("MCP_REGISTRY_DATABASE_TYPE", "mysql")
	os.Setenv("MCP_REGISTRY_DATABASE_MAX_CONNECTIONS", "50")
	defer clearEnvVars()

	cfg := config.NewConfig()

	if cfg.Server.Address != ":9090" {
		t.Errorf("expected server address :9090, got %s", cfg.Server.Address)
	}

	if cfg.Database.Type != "mysql" {
		t.Errorf("expected database type mysql, got %s", cfg.Database.Type)
	}

	if cfg.Database.MaxConnections != 50 {
		t.Errorf("expected max connections 50, got %d", cfg.Database.MaxConnections)
	}
}

func TestDatabaseConfig_Validate_Success(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.DatabaseConfig
	}{
		{
			name: "valid postgres config",
			cfg: config.DatabaseConfig{
				Type:           "postgres",
				URL:            "postgres://localhost:5432/test",
				MaxConnections: 30,
				MinConnections: 5,
				SSLMode:        "require",
			},
		},
		{
			name: "valid dynamodb config",
			cfg: config.DatabaseConfig{
				Type:           "dynamodb",
				Region:         "us-east-1",
				TableName:      "test-table",
				MaxConnections: 30,
				MinConnections: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); err != nil {
				t.Errorf("expected validation to pass, got error: %v", err)
			}
		})
	}
}

func TestDatabaseConfig_Validate_Errors(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config.DatabaseConfig
		expectedErr string
	}{
		{
			name: "unsupported database type",
			cfg: config.DatabaseConfig{
				Type:           "mongodb",
				MaxConnections: 30,
				MinConnections: 5,
			},
			expectedErr: "unsupported database type",
		},
		{
			name: "missing URL for postgres",
			cfg: config.DatabaseConfig{
				Type:           "postgres",
				URL:            "",
				MaxConnections: 30,
				MinConnections: 5,
			},
			expectedErr: "database URL is required",
		},
		{
			name: "max connections less than min",
			cfg: config.DatabaseConfig{
				Type:           "postgres",
				URL:            "postgres://localhost:5432/test",
				MaxConnections: 5,
				MinConnections: 10,
			},
			expectedErr: "max connections",
		},
		{
			name: "missing region for dynamodb",
			cfg: config.DatabaseConfig{
				Type:           "dynamodb",
				Region:         "",
				TableName:      "test",
				MaxConnections: 30,
				MinConnections: 5,
			},
			expectedErr: "region is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if err == nil {
				t.Errorf("expected validation error, got nil")
				return
			}
			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("expected error containing %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestServerConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.ServerConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: config.ServerConfig{
				Address:         ":8080",
				ReadTimeout:     30 * time.Second,
				WriteTimeout:    30 * time.Second,
				ShutdownTimeout: 10 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "missing address",
			cfg: config.ServerConfig{
				Address: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBackwardsCompatibility(t *testing.T) {
	// Test that backwards compatibility methods work
	clearEnvVars()

	os.Setenv("MCP_REGISTRY_DATABASE_URL", "postgres://test:5432/db")
	os.Setenv("MCP_REGISTRY_AUTH_GITHUB_CLIENT_ID", "test-client-id")
	defer clearEnvVars()

	cfg := config.NewConfig()

	// Test backwards compatible methods
	if cfg.DatabaseURL() != "postgres://test:5432/db" {
		t.Errorf("DatabaseURL() failed, got %s", cfg.DatabaseURL())
	}

	if cfg.GithubClientID() != "test-client-id" {
		t.Errorf("GithubClientID() failed, got %s", cfg.GithubClientID())
	}
}

// Helper functions

func clearEnvVars() {
	envVars := []string{
		"MCP_REGISTRY_SERVER_ADDRESS",
		"MCP_REGISTRY_DATABASE_TYPE",
		"MCP_REGISTRY_DATABASE_URL",
		"MCP_REGISTRY_DATABASE_MAX_CONNECTIONS",
		"MCP_REGISTRY_AUTH_GITHUB_CLIENT_ID",
		"MCP_REGISTRY_AUTH_GITHUB_CLIENT_SECRET",
	}

	for _, v := range envVars {
		os.Unsetenv(v)
	}
}
