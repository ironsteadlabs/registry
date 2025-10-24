package commands

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/registry/cmd/publisher/auth"
)

const (
	DefaultRegistryURL = "https://registry.modelcontextprotocol.io"
	TokenFileName      = ".mcp_publisher_token" //nolint:gosec // Not a credential, just a filename
	MethodGitHub       = "github"
	MethodGitHubOIDC   = "github-oidc"
	MethodDNS          = "dns"
	MethodHTTP         = "http"
	MethodNone         = "none"
)

type CryptoAlgorithm auth.CryptoAlgorithm

func (c *CryptoAlgorithm) String() string {
	return string(*c)
}

func (c *CryptoAlgorithm) Set(v string) error {
	switch v {
	case string(auth.AlgorithmEd25519), string(auth.AlgorithmECDSAP384):
		*c = CryptoAlgorithm(v)
		return nil
	}
	return fmt.Errorf("invalid algorithm: %q (allowed: ed25519, ecdsap384)", v)
}

type loginFlags struct {
	domain          string
	privateKey      string
	cryptoAlgorithm CryptoAlgorithm
	registryURL     string
	token           string
}

func LoginCommand(args []string) error {
	if len(args) < 1 {
		return errors.New("authentication method required\n\nUsage: mcp-publisher login <method>\n\nMethods:\n  github        Interactive GitHub authentication\n  github-oidc   GitHub Actions OIDC authentication\n  dns           DNS-based authentication (requires --domain and --private-key)\n  http          HTTP-based authentication (requires --domain and --private-key)\n  none          Anonymous authentication (for testing)")
	}

	method := args[0]
	flags, err := parseLoginFlags(method, args[1:])
	if err != nil {
		return err
	}

	authProvider, err := createAuthProvider(method, flags)
	if err != nil {
		return err
	}

	return performLogin(authProvider, method, flags.registryURL)
}

func parseLoginFlags(method string, args []string) (*loginFlags, error) {
	flags := &loginFlags{
		cryptoAlgorithm: CryptoAlgorithm(auth.AlgorithmEd25519), // default
	}
	loginFlagSet := flag.NewFlagSet("login", flag.ExitOnError)

	loginFlagSet.StringVar(&flags.registryURL, "registry", DefaultRegistryURL, "Registry URL")

	if method == MethodGitHub {
		loginFlagSet.StringVar(&flags.token, "token", "", "GitHub Personal Access Token")
	}

	if method == MethodDNS || method == MethodHTTP {
		loginFlagSet.StringVar(&flags.domain, "domain", "", "Domain name")
		loginFlagSet.StringVar(&flags.privateKey, "private-key", "", "Private key (64-char hex)")
		loginFlagSet.Var(&flags.cryptoAlgorithm, "algorithm", "Cryptographic algorithm (ed25519, ecdsap384)")
	}

	if err := loginFlagSet.Parse(args); err != nil {
		return nil, err
	}

	return flags, nil
}

func createAuthProvider(method string, flags *loginFlags) (auth.Provider, error) {
	switch method {
	case MethodGitHub:
		return auth.NewGitHubATProvider(true, flags.registryURL, flags.token), nil
	case MethodGitHubOIDC:
		return auth.NewGitHubOIDCProvider(flags.registryURL), nil
	case MethodDNS:
		if flags.domain == "" || flags.privateKey == "" {
			return nil, errors.New("dns authentication requires --domain and --private-key")
		}
		return auth.NewDNSProvider(flags.registryURL, flags.domain, flags.privateKey, auth.CryptoAlgorithm(flags.cryptoAlgorithm)), nil
	case MethodHTTP:
		if flags.domain == "" || flags.privateKey == "" {
			return nil, errors.New("http authentication requires --domain and --private-key")
		}
		return auth.NewHTTPProvider(flags.registryURL, flags.domain, flags.privateKey, auth.CryptoAlgorithm(flags.cryptoAlgorithm)), nil
	case MethodNone:
		return auth.NewNoneProvider(flags.registryURL), nil
	default:
		return nil, fmt.Errorf("unknown authentication method: %s\nFor a list of available methods, run: mcp-publisher login", method)
	}
}

func performLogin(authProvider auth.Provider, method, registryURL string) error {
	ctx := context.Background()
	_, _ = fmt.Fprintf(os.Stdout, "Logging in with %s...\n", method)

	if err := authProvider.Login(ctx); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Get and save token
	token, err := authProvider.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	// Save token to file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	tokenPath := filepath.Join(homeDir, TokenFileName)
	tokenData := map[string]string{
		"token":    token,
		"method":   method,
		"registry": registryURL,
	}

	jsonData, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	if err := os.WriteFile(tokenPath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	_, _ = fmt.Fprintln(os.Stdout, "✓ Successfully logged in")
	return nil
}
