package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// captureOutput temporarily redirects os.Stdout and os.Stderr while fn runs,
// so callers can assert on what a Run() invocation actually printed to each
// stream (newApp reads os.Stdout/os.Stderr at call time).
func captureOutput(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()

	origOut, origErr := os.Stdout, os.Stderr
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = outW
	os.Stderr = errW
	defer func() {
		os.Stdout = origOut
		os.Stderr = origErr
	}()

	fn()

	_ = outW.Close()
	_ = errW.Close()

	outData, _ := io.ReadAll(outR)
	errData, _ := io.ReadAll(errR)
	return string(outData), string(errData)
}

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

func TestDistanceRadiusFlags(t *testing.T) {
	var gotQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"origin": {"query": "Washington DC"}, "destinations": []}`)
	}))
	defer server.Close()

	err := Run(context.Background(), []string{
		"geocodio",
		"--api-key", "test-api-key",
		"--base-url", server.URL,
		"distance", "Washington DC", "New York", "Boston",
		"--radius", "150.5",
		"--min-distance", "10",
		"--max-duration", "7200",
		"--min-duration", "60",
		"--max-results", "3",
		"--order-by", "duration",
		"--sort-order", "desc",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	want := map[string]string{
		"max_distance": "150.5",
		"min_distance": "10",
		"max_duration": "7200",
		"min_duration": "60",
		"max_results":  "3",
		"order_by":     "duration",
		"sort_order":   "desc",
	}
	for key, wantValue := range want {
		if got := gotQuery.Get(key); got != wantValue {
			t.Errorf("query[%q] = %q, want %q", key, got, wantValue)
		}
	}
}

func TestDistanceRejectsInvalidRadiusRange(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("expected no request to be made")
	}))
	defer server.Close()

	err := Run(context.Background(), []string{
		"geocodio",
		"--api-key", "test-api-key",
		"--base-url", server.URL,
		"distance", "Washington DC", "New York",
		"--radius", "10",
		"--min-distance", "50",
	})
	if err == nil {
		t.Fatal("expected an error for min-distance greater than max-distance")
	}
}

func TestDistanceMatrixRadiusFlags(t *testing.T) {
	dir := t.TempDir()
	originsFile := filepath.Join(dir, "origins.txt")
	destinationsFile := filepath.Join(dir, "destinations.txt")
	if err := os.WriteFile(originsFile, []byte("Washington DC\n"), 0o600); err != nil {
		t.Fatalf("writing origins: %v", err)
	}
	if err := os.WriteFile(destinationsFile, []byte("New York\nBoston\n"), 0o600); err != nil {
		t.Fatalf("writing destinations: %v", err)
	}

	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"results": []}`)
	}))
	defer server.Close()

	err := Run(context.Background(), []string{
		"geocodio",
		"--api-key", "test-api-key",
		"--base-url", server.URL,
		"distance-matrix",
		"--origins", originsFile,
		"--destinations", destinationsFile,
		"--radius", "25",
		"--max-results", "1",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got := gotBody["max_distance"]; got != 25.0 {
		t.Errorf("body[\"max_distance\"] = %v, want 25", got)
	}
	if got := gotBody["max_results"]; got != 1.0 {
		t.Errorf("body[\"max_results\"] = %v, want 1", got)
	}
}

func TestDistanceDefaultsToStraightlineMode(t *testing.T) {
	var gotQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"origin": {"query": "Washington DC"}, "destinations": []}`)
	}))
	defer server.Close()

	err := Run(context.Background(), []string{
		"geocodio",
		"--api-key", "test-api-key",
		"--base-url", server.URL,
		"distance", "Washington DC", "New York",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got := gotQuery.Get("mode"); got != "straightline" {
		t.Errorf("mode = %q, want %q (straightline is the API default and costs half as much)", got, "straightline")
	}
}

func TestDistanceMatrixDefaultsToStraightlineMode(t *testing.T) {
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&gotBody)

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"results": []}`)
	}))
	defer server.Close()

	dir := t.TempDir()
	originsFile := filepath.Join(dir, "origins.txt")
	destsFile := filepath.Join(dir, "destinations.txt")
	if err := os.WriteFile(originsFile, []byte("Washington DC\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destsFile, []byte("New York\n"), 0600); err != nil {
		t.Fatal(err)
	}

	err := Run(context.Background(), []string{
		"geocodio",
		"--api-key", "test-api-key",
		"--base-url", server.URL,
		"distance-matrix", "--origins", originsFile, "--destinations", destsFile,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got := gotBody["mode"]; got != "straightline" {
		t.Errorf("body mode = %v, want %q", got, "straightline")
	}
}

func TestDistanceJobsCreateDefaultsToStraightlineMode(t *testing.T) {
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&gotBody)

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"identifier": "abc123", "status": "PROCESSING"}`)
	}))
	defer server.Close()

	dir := t.TempDir()
	originsFile := filepath.Join(dir, "origins.txt")
	destsFile := filepath.Join(dir, "destinations.txt")
	if err := os.WriteFile(originsFile, []byte("Washington DC\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destsFile, []byte("New York\n"), 0600); err != nil {
		t.Fatal(err)
	}

	err := Run(context.Background(), []string{
		"geocodio",
		"--api-key", "test-api-key",
		"--base-url", server.URL,
		"--json",
		"distance-jobs", "create", "--name", "test", "--origins", originsFile, "--destinations", destsFile,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got := gotBody["distance_mode"]; got != "straightline" {
		t.Errorf("body distance_mode = %v, want %q", got, "straightline")
	}
}

func TestBatchGeocodeOverLimitPointsAtLists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("expected no request to be made")
	}))
	defer server.Close()

	batchFile := filepath.Join(t.TempDir(), "addresses.txt")
	var sb strings.Builder
	for i := 0; i < 10001; i++ {
		sb.WriteString("Washington DC\n")
	}
	if err := os.WriteFile(batchFile, []byte(sb.String()), 0600); err != nil {
		t.Fatal(err)
	}

	err := Run(context.Background(), []string{
		"geocodio",
		"--api-key", "test-api-key",
		"--base-url", server.URL,
		"geocode", "--batch", batchFile,
	})
	if err == nil {
		t.Fatal("expected an error for an oversized batch file")
	}

	for _, want := range []string{"10,000", "10001", "appends", "lists upload"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error message missing %q; got: %v", want, err)
		}
	}
}

func TestListsStatusEnqueuedShowsConcurrencyHint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id": 42, "status": {"state": "ENQUEUED", "progress": 0}}`)
	}))
	defer server.Close()

	var runErr error
	stdout, stderr := captureOutput(t, func() {
		runErr = Run(context.Background(), []string{
			"geocodio",
			"--api-key", "test-api-key",
			"--base-url", server.URL,
			"lists", "status", "42",
		})
	})
	if runErr != nil {
		t.Fatalf("Run() error = %v", runErr)
	}

	for _, want := range []string{"queued", "depends on your plan", "spreadsheet-concurrency"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("stderr missing %q; got: %s", want, stderr)
		}
	}

	if strings.Contains(stdout, "queued") || strings.Contains(stdout, "spreadsheet-concurrency") {
		t.Errorf("concurrency hint leaked into stdout: %s", stdout)
	}
}

func TestListsStatusProcessingOmitsConcurrencyHint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id": 42, "status": {"state": "PROCESSING", "progress": 40}}`)
	}))
	defer server.Close()

	var runErr error
	_, stderr := captureOutput(t, func() {
		runErr = Run(context.Background(), []string{
			"geocodio",
			"--api-key", "test-api-key",
			"--base-url", server.URL,
			"lists", "status", "42",
		})
	})
	if runErr != nil {
		t.Fatalf("Run() error = %v", runErr)
	}

	if strings.Contains(stderr, "queued") || strings.Contains(stderr, "Geocodio limits") {
		t.Errorf("expected no concurrency hint for PROCESSING status; got: %s", stderr)
	}
}

func TestListsStatusCompletedOmitsConcurrencyHint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id": 42, "status": {"state": "COMPLETED", "progress": 100}}`)
	}))
	defer server.Close()

	var runErr error
	_, stderr := captureOutput(t, func() {
		runErr = Run(context.Background(), []string{
			"geocodio",
			"--api-key", "test-api-key",
			"--base-url", server.URL,
			"lists", "status", "42",
		})
	})
	if runErr != nil {
		t.Fatalf("Run() error = %v", runErr)
	}

	if strings.Contains(stderr, "queued") || strings.Contains(stderr, "Geocodio limits") {
		t.Errorf("expected no concurrency hint for COMPLETED status; got: %s", stderr)
	}
}

func TestListsStatusJSONNotCorruptedByConcurrencyHint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id": 42, "status": {"state": "ENQUEUED", "progress": 0}}`)
	}))
	defer server.Close()

	var runErr error
	stdout, stderr := captureOutput(t, func() {
		runErr = Run(context.Background(), []string{
			"geocodio",
			"--api-key", "test-api-key",
			"--base-url", server.URL,
			"--json",
			"lists", "status", "42",
		})
	})
	if runErr != nil {
		t.Fatalf("Run() error = %v", runErr)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(stdout), &parsed); err != nil {
		t.Fatalf("stdout is not valid JSON: %v; got: %s", err, stdout)
	}

	// The hint is informational and belongs on stderr only, so it must not
	// appear in the machine-readable stdout payload.
	if strings.Contains(stdout, "Geocodio limits") {
		t.Errorf("concurrency hint leaked into JSON stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "queued") {
		t.Errorf("expected concurrency hint on stderr; got: %s", stderr)
	}
}

func TestShouldShowEnqueuedHint(t *testing.T) {
	tests := []struct {
		name    string
		state   string
		elapsed time.Duration
		want    bool
	}{
		{"enqueued at threshold", "ENQUEUED", enqueuedHintDelay, true},
		{"enqueued past threshold", "ENQUEUED", enqueuedHintDelay + time.Second, true},
		{"enqueued before threshold", "ENQUEUED", enqueuedHintDelay - time.Second, false},
		{"processing past threshold", "PROCESSING", enqueuedHintDelay + time.Second, false},
		{"completed past threshold", "COMPLETED", enqueuedHintDelay + time.Second, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldShowEnqueuedHint(tt.state, tt.elapsed); got != tt.want {
				t.Errorf("shouldShowEnqueuedHint(%q, %v) = %v, want %v", tt.state, tt.elapsed, got, tt.want)
			}
		})
	}
}
