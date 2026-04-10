package main_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	geocli "github.com/geocodio/geocodio-cli/internal/cli"
)

// runCLI runs the CLI with the given args against a test server.
// Returns stdout output and any error.
func runCLI(t *testing.T, server *httptest.Server, args ...string) (string, error) {
	t.Helper()

	fullArgs := []string{"geocodio", "--base-url", server.URL, "--api-key", "test-key", "--no-color"}
	fullArgs = append(fullArgs, args...)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	app := geocli.NewApp()
	err := app.Run(context.Background(), fullArgs)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	return buf.String(), err
}

// loadFixture reads a JSON fixture file from testdata/.
func loadFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("failed to load fixture %s: %v", name, err)
	}
	return string(data)
}

// newFixtureServer creates a test server that serves a fixture for a given path.
func newFixtureServer(t *testing.T, pathToFixture map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for path, fixture := range pathToFixture {
			if strings.Contains(r.URL.Path, path) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				_, _ = w.Write([]byte(fixture))
				return
			}
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.String())
		w.WriteHeader(404)
	}))
}

// README Example: geocodio geocode "1600 Pennsylvania Ave NW, Washington DC"
func TestIntegration_GeocodeSimple(t *testing.T) {
	server := newFixtureServer(t, map[string]string{
		"/geocode": loadFixture(t, "geocode_simple.json"),
	})
	defer server.Close()

	output, err := runCLI(t, server, "geocode", "1600 Pennsylvania Ave NW, Washington DC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, output, "1600 Pennsylvania Ave NW, Washington, DC 20500")
	assertContains(t, output, "38.8976750")
	assertContains(t, output, "rooftop")
	assertContains(t, output, "Statewide (City of Washington)")
}

// README Example: geocodio geocode "30 Rockefeller Plaza, New York NY" --fields timezone,cd
func TestIntegration_GeocodeWithFields(t *testing.T) {
	server := newFixtureServer(t, map[string]string{
		"/geocode": loadFixture(t, "geocode_fields.json"),
	})
	defer server.Close()

	output, err := runCLI(t, server, "geocode", "30 Rockefeller Plaza, New York NY", "--fields", "timezone,cd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, output, "30 Rockefeller Plz, New York, NY 10112")
	assertContains(t, output, "Timezone")
	assertContains(t, output, "America/New_York")
	assertContains(t, output, "Congressional District")
	assertContains(t, output, "Congressional District 12")
}

// README Example: geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --destinations "New York" --destinations "Boston"
func TestIntegration_GeocodeWithDestinations(t *testing.T) {
	server := newFixtureServer(t, map[string]string{
		"/geocode": loadFixture(t, "geocode_destinations.json"),
	})
	defer server.Close()

	output, err := runCLI(t, server, "geocode", "1600 Pennsylvania Ave NW, Washington DC", "-d", "New York", "-d", "Boston")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, output, "Distances:")
	assertContains(t, output, "New York, NY 10001")
	assertContains(t, output, "Boston, MA 02106")
	assertContains(t, output, "206.2 miles")
	assertContains(t, output, "393.7 miles")
}

// README Example: geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --destinations "New York" --distance-mode driving
func TestIntegration_GeocodeWithDrivingDistance(t *testing.T) {
	server := newFixtureServer(t, map[string]string{
		"/geocode": loadFixture(t, "geocode_driving.json"),
	})
	defer server.Close()

	output, err := runCLI(t, server, "geocode", "1600 Pennsylvania Ave NW, Washington DC", "-d", "New York", "--distance-mode", "driving")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, output, "227.2 miles")
	assertContains(t, output, "288 minutes")
}

// README Example: geocodio geocode "1600 Pennsylvania Ave NW, Washington DC" --show-address-key
func TestIntegration_GeocodeShowAddressKey(t *testing.T) {
	server := newFixtureServer(t, map[string]string{
		"/geocode": loadFixture(t, "geocode_address_key.json"),
	})
	defer server.Close()

	output, err := runCLI(t, server, "geocode", "1600 Pennsylvania Ave NW, Washington DC", "--show-address-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, output, "Address Key")
	assertContains(t, output, "gcod_usnqu2qtyk7ac38dlpzanec598p2c")
}

// README Example: geocodio reverse "38.8976,-77.0365"
func TestIntegration_ReverseGeocode(t *testing.T) {
	server := newFixtureServer(t, map[string]string{
		"/reverse": loadFixture(t, "reverse_simple.json"),
	})
	defer server.Close()

	output, err := runCLI(t, server, "reverse", "38.8976,-77.0365")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, output, "1600 Pennsylvania Ave NW, Washington, DC 20500")
	assertContains(t, output, "rooftop")
}

// README Example: geocodio reverse "40.7588,-73.9788" --skip-geocoding --fields timezone,cd
func TestIntegration_ReverseSkipGeocoding(t *testing.T) {
	server := newFixtureServer(t, map[string]string{
		"/reverse": loadFixture(t, "reverse_skip_fields.json"),
	})
	defer server.Close()

	output, err := runCLI(t, server, "reverse", "40.7588,-73.9788", "--skip-geocoding", "--fields", "timezone,cd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, output, "Timezone")
	assertContains(t, output, "America/New_York")
	assertContains(t, output, "Congressional District")
	assertContains(t, output, "Congressional District 12")
}

// README Example: geocodio distance "Washington DC" "New York"
func TestIntegration_Distance(t *testing.T) {
	server := newFixtureServer(t, map[string]string{
		"/distance": loadFixture(t, "distance_simple.json"),
	})
	defer server.Close()

	output, err := runCLI(t, server, "distance", "Washington DC", "New York")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, output, "Washington, DC 20001")
	assertContains(t, output, "New York, NY 10001")
	assertContains(t, output, "226.6 miles")
	assertContains(t, output, "284 minutes")
}

// README Example: geocodio distance "Washington DC" "New York" "Boston" "Philadelphia"
func TestIntegration_DistanceMultiple(t *testing.T) {
	server := newFixtureServer(t, map[string]string{
		"/distance": loadFixture(t, "distance_multiple.json"),
	})
	defer server.Close()

	output, err := runCLI(t, server, "distance", "Washington DC", "New York", "Boston", "Philadelphia")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, output, "New York, NY 10001")
	assertContains(t, output, "Boston, MA 02106")
	assertContains(t, output, "Philadelphia, PA 19101")
}

// README Example: geocodio distance "Washington DC" "New York" --mode driving --units km
func TestIntegration_DistanceWithUnitsKm(t *testing.T) {
	server := newFixtureServer(t, map[string]string{
		"/distance": loadFixture(t, "distance_km.json"),
	})
	defer server.Close()

	output, err := runCLI(t, server, "distance", "Washington DC", "New York", "--mode", "driving", "--units", "km")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertContains(t, output, "364.7 km (226.6 miles)")
}

func assertContains(t *testing.T, output, want string) {
	t.Helper()
	if !strings.Contains(output, want) {
		t.Errorf("output missing %q:\n%s", want, output)
	}
}
