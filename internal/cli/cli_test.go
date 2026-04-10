package cli

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestParseCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLat float64
		wantLng float64
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid coordinates",
			input:   "38.8977,-77.0365",
			wantLat: 38.8977,
			wantLng: -77.0365,
			wantErr: false,
		},
		{
			name:    "valid with spaces",
			input:   " 38.8977 , -77.0365 ",
			wantLat: 38.8977,
			wantLng: -77.0365,
			wantErr: false,
		},
		{
			name:    "valid edge case - equator prime meridian",
			input:   "0,0",
			wantLat: 0,
			wantLng: 0,
			wantErr: false,
		},
		{
			name:    "valid edge case - max bounds",
			input:   "90,180",
			wantLat: 90,
			wantLng: 180,
			wantErr: false,
		},
		{
			name:    "valid edge case - min bounds",
			input:   "-90,-180",
			wantLat: -90,
			wantLng: -180,
			wantErr: false,
		},
		{
			name:    "invalid format - missing comma",
			input:   "38.8977 -77.0365",
			wantErr: true,
			errMsg:  "invalid coordinate format",
		},
		{
			name:    "invalid format - extra comma",
			input:   "38.8977,-77.0365,0",
			wantErr: true,
			errMsg:  "invalid coordinate format",
		},
		{
			name:    "invalid format - empty string",
			input:   "",
			wantErr: true,
			errMsg:  "invalid coordinate format",
		},
		{
			name:    "invalid latitude - not a number",
			input:   "abc,-77.0365",
			wantErr: true,
			errMsg:  "invalid latitude",
		},
		{
			name:    "invalid longitude - not a number",
			input:   "38.8977,xyz",
			wantErr: true,
			errMsg:  "invalid longitude",
		},
		{
			name:    "out of bounds - latitude too high",
			input:   "91,0",
			wantErr: true,
			errMsg:  "latitude must be between",
		},
		{
			name:    "out of bounds - latitude too low",
			input:   "-91,0",
			wantErr: true,
			errMsg:  "latitude must be between",
		},
		{
			name:    "out of bounds - longitude too high",
			input:   "0,181",
			wantErr: true,
			errMsg:  "longitude must be between",
		},
		{
			name:    "out of bounds - longitude too low",
			input:   "0,-181",
			wantErr: true,
			errMsg:  "longitude must be between",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lng, err := parseCoordinates(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseCoordinates() expected error, got nil")
					return
				}
				if tt.errMsg != "" && !containsString(err.Error(), tt.errMsg) {
					t.Errorf("parseCoordinates() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("parseCoordinates() unexpected error: %v", err)
				return
			}

			if lat != tt.wantLat {
				t.Errorf("parseCoordinates() lat = %v, want %v", lat, tt.wantLat)
			}
			if lng != tt.wantLng {
				t.Errorf("parseCoordinates() lng = %v, want %v", lng, tt.wantLng)
			}
		})
	}
}

func TestReadLines(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantLines []string
		wantErr   bool
	}{
		{
			name:      "normal file",
			content:   "line1\nline2\nline3",
			wantLines: []string{"line1", "line2", "line3"},
			wantErr:   false,
		},
		{
			name:      "file with blank lines",
			content:   "line1\n\nline2\n   \nline3",
			wantLines: []string{"line1", "line2", "line3"},
			wantErr:   false,
		},
		{
			name:      "file with trailing newline",
			content:   "line1\nline2\n",
			wantLines: []string{"line1", "line2"},
			wantErr:   false,
		},
		{
			name:      "empty file",
			content:   "",
			wantLines: nil,
			wantErr:   false,
		},
		{
			name:      "whitespace only file",
			content:   "   \n\t\n  \n",
			wantLines: nil,
			wantErr:   false,
		},
		{
			name:      "lines with whitespace trimmed",
			content:   "  line1  \n\tline2\t\n",
			wantLines: []string{"line1", "line2"},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.txt")
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0600); err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}

			lines, err := readLines(tmpFile)

			if tt.wantErr {
				if err == nil {
					t.Errorf("readLines() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("readLines() unexpected error: %v", err)
				return
			}

			if len(lines) != len(tt.wantLines) {
				t.Errorf("readLines() returned %d lines, want %d", len(lines), len(tt.wantLines))
				return
			}

			for i, want := range tt.wantLines {
				if lines[i] != want {
					t.Errorf("readLines()[%d] = %q, want %q", i, lines[i], want)
				}
			}
		})
	}

	t.Run("non-existent file", func(t *testing.T) {
		_, err := readLines("/nonexistent/path/to/file.txt")
		if err == nil {
			t.Error("readLines() expected error for non-existent file, got nil")
		}
	})
}

func TestGeocodeWithCommaDestination(t *testing.T) {
	wantDestinations := []string{
		"1600 Pennsylvania Ave NW, Washington DC",
		"350 Fifth Ave, New York, NY",
	}

	var sawRequest bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawRequest = true

		if r.URL.Path != "/geocode" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/geocode")
		}

		destinations := r.URL.Query()["destinations[]"]
		if len(destinations) != len(wantDestinations) {
			t.Errorf("expected %d destinations, got %d: %v", len(wantDestinations), len(destinations), destinations)
		}
		for i, want := range wantDestinations {
			if i < len(destinations) && destinations[i] != want {
				t.Errorf("destination[%d] = %q, want %q", i, destinations[i], want)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"results": [{
				"formatted_address": "1 Main St, Washington, DC 20001",
				"location": {"lat": 38.9000, "lng": -77.0000},
				"accuracy": 1,
				"accuracy_type": "rooftop",
				"destinations": [{
					"query": "1600 Pennsylvania Ave NW, Washington DC",
					"distance_miles": 1.2
				}]
			}]
		}`)
	}))
	defer server.Close()

	err := Run(context.Background(), []string{
		"geocodio",
		"--api-key", "test-api-key",
		"--base-url", server.URL,
		"geocode", "1 Main St",
		"-d", wantDestinations[0],
		"-d", wantDestinations[1],
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !sawRequest {
		t.Fatal("expected geocode request")
	}
}

func TestNewApp(t *testing.T) {
	app := NewApp()

	if app == nil {
		t.Fatal("NewApp() returned nil")
	}

	if app.Name != "geocodio" {
		t.Errorf("app.Name = %q, want %q", app.Name, "geocodio")
	}

	if app.Version == "" {
		t.Error("app.Version is empty")
	}

	// Verify expected commands exist
	expectedCmds := []string{"geocode", "reverse", "distance", "distance-matrix", "distance-jobs", "lists"}
	cmdNames := make(map[string]bool)
	for _, cmd := range app.Commands {
		cmdNames[cmd.Name] = true
	}

	for _, expected := range expectedCmds {
		if !cmdNames[expected] {
			t.Errorf("missing command %q", expected)
		}
	}

	// Verify expected flags exist
	expectedFlags := []string{"api-key", "base-url", "json", "debug"}
	flagNames := make(map[string]bool)
	for _, flag := range app.Flags {
		for _, name := range flag.Names() {
			flagNames[name] = true
		}
	}

	for _, expected := range expectedFlags {
		if !flagNames[expected] {
			t.Errorf("missing flag %q", expected)
		}
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
