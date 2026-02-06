// Package config handles configuration resolution for the Geocodio CLI.
package config

import "os"

const (
	// EnvAPIKey is the environment variable for the Geocodio API key.
	EnvAPIKey = "GEOCODIO_API_KEY"

	// DefaultBaseURL is the default Geocodio API base URL.
	DefaultBaseURL = "https://api.geocod.io/v1.9"
)

// Config holds the CLI configuration.
type Config struct {
	APIKey  string
	BaseURL string
	Debug   bool
}

// New creates a Config by resolving values from flags and environment.
// Flag values take precedence over environment variables.
func New(apiKeyFlag, baseURLFlag string, debug bool) *Config {
	cfg := &Config{
		APIKey:  apiKeyFlag,
		BaseURL: baseURLFlag,
		Debug:   debug,
	}

	// Fall back to environment variable if flag not provided
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv(EnvAPIKey)
	}

	// Use default base URL if not provided
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}

	return cfg
}

// Validate checks that required configuration is present.
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return &MissingAPIKeyError{}
	}
	return nil
}

// MissingAPIKeyError indicates the API key was not provided.
type MissingAPIKeyError struct{}

func (e *MissingAPIKeyError) Error() string {
	return "API key required: set GEOCODIO_API_KEY environment variable or use --api-key flag"
}
