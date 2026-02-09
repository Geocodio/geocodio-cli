package config

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		apiKeyFlag  string
		baseURLFlag string
		debug       bool
		envAPIKey   string
		wantAPIKey  string
		wantBaseURL string
		wantDebug   bool
	}{
		{
			name:        "flag takes precedence over env",
			apiKeyFlag:  "flag-key",
			envAPIKey:   "env-key",
			wantAPIKey:  "flag-key",
			wantBaseURL: DefaultBaseURL,
		},
		{
			name:        "falls back to env when flag empty",
			apiKeyFlag:  "",
			envAPIKey:   "env-key",
			wantAPIKey:  "env-key",
			wantBaseURL: DefaultBaseURL,
		},
		{
			name:        "uses default base URL when not provided",
			apiKeyFlag:  "test-key",
			baseURLFlag: "",
			wantAPIKey:  "test-key",
			wantBaseURL: DefaultBaseURL,
		},
		{
			name:        "uses custom base URL when provided",
			apiKeyFlag:  "test-key",
			baseURLFlag: "https://custom.api.com",
			wantAPIKey:  "test-key",
			wantBaseURL: "https://custom.api.com",
		},
		{
			name:        "debug flag is passed through",
			apiKeyFlag:  "test-key",
			debug:       true,
			wantAPIKey:  "test-key",
			wantBaseURL: DefaultBaseURL,
			wantDebug:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.envAPIKey != "" {
				os.Setenv(EnvAPIKey, tt.envAPIKey)
				defer os.Unsetenv(EnvAPIKey)
			} else {
				os.Unsetenv(EnvAPIKey)
			}

			cfg := New(tt.apiKeyFlag, tt.baseURLFlag, tt.debug)

			if cfg.APIKey != tt.wantAPIKey {
				t.Errorf("APIKey = %q, want %q", cfg.APIKey, tt.wantAPIKey)
			}
			if cfg.BaseURL != tt.wantBaseURL {
				t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, tt.wantBaseURL)
			}
			if cfg.Debug != tt.wantDebug {
				t.Errorf("Debug = %v, want %v", cfg.Debug, tt.wantDebug)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid with API key",
			apiKey:  "test-key",
			wantErr: false,
		},
		{
			name:    "invalid without API key",
			apiKey:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				APIKey:  tt.apiKey,
				BaseURL: DefaultBaseURL,
			}

			err := cfg.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil {
				if _, ok := err.(*MissingAPIKeyError); !ok {
					t.Errorf("Validate() error type = %T, want *MissingAPIKeyError", err)
				}
			}
		})
	}
}

func TestMissingAPIKeyError_Error(t *testing.T) {
	err := &MissingAPIKeyError{}
	msg := err.Error()

	if msg == "" {
		t.Error("Error() returned empty string")
	}

	// Check that it mentions the expected ways to provide the key
	if !containsAll(msg, "API key", "GEOCODIO_API_KEY", "--api-key") {
		t.Errorf("Error() = %q, should mention API key, GEOCODIO_API_KEY, and --api-key", msg)
	}
}

func containsAll(s string, substrings ...string) bool {
	for _, sub := range substrings {
		if !contains(s, sub) {
			return false
		}
	}
	return true
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsHelper(s, sub))
}

func containsHelper(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
