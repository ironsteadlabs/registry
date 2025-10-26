package config

// NewConfig creates a new configuration
// This is the main entry point for loading configuration
// It maintains backwards compatibility with the previous env-only approach
func NewConfig() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		// Maintain backwards compatibility: panic on error like before
		panic(err)
	}
	return cfg
}

// Backwards compatibility: Support for legacy field access
// These methods allow existing code to access fields using the old flat structure

// ServerAddress returns the server address (backwards compatible)
func (c *Config) ServerAddress() string {
	return c.Server.Address
}

// DatabaseURL returns the database URL (backwards compatible)
func (c *Config) DatabaseURL() string {
	return c.Database.URL
}

// SeedFrom returns the seed data path (backwards compatible)
func (c *Config) SeedFrom() string {
	return c.Database.SeedFrom
}

// Version returns the application version (backwards compatible)
func (c *Config) Version() string {
	return c.Telemetry.Version
}

// GithubClientID returns the GitHub client ID (backwards compatible)
func (c *Config) GithubClientID() string {
	return c.Auth.GitHub.ClientID
}

// GithubClientSecret returns the GitHub client secret (backwards compatible)
func (c *Config) GithubClientSecret() string {
	return c.Auth.GitHub.ClientSecret
}

// JWTPrivateKey returns the JWT private key (backwards compatible)
func (c *Config) JWTPrivateKey() string {
	return c.Auth.JWT.PrivateKey
}

// EnableAnonymousAuth returns whether anonymous auth is enabled (backwards compatible)
func (c *Config) EnableAnonymousAuth() bool {
	return c.Features.EnableAnonymousAuth
}

// EnableRegistryValidation returns whether registry validation is enabled (backwards compatible)
func (c *Config) EnableRegistryValidation() bool {
	return c.Features.EnableRegistryValidation
}

// OIDCEnabled returns whether OIDC is enabled (backwards compatible)
func (c *Config) OIDCEnabled() bool {
	return c.Auth.OIDC.Enabled
}

// OIDCIssuer returns the OIDC issuer (backwards compatible)
func (c *Config) OIDCIssuer() string {
	return c.Auth.OIDC.Issuer
}

// OIDCClientID returns the OIDC client ID (backwards compatible)
func (c *Config) OIDCClientID() string {
	return c.Auth.OIDC.ClientID
}

// OIDCExtraClaims returns the OIDC extra claims (backwards compatible)
func (c *Config) OIDCExtraClaims() string {
	return c.Auth.OIDC.ExtraClaims
}

// OIDCEditPerms returns the OIDC edit permissions (backwards compatible)
func (c *Config) OIDCEditPerms() string {
	return c.Auth.OIDC.EditPerms
}

// OIDCPublishPerms returns the OIDC publish permissions (backwards compatible)
func (c *Config) OIDCPublishPerms() string {
	return c.Auth.OIDC.PublishPerms
}
