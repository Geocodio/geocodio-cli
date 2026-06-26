package cli

import (
	"testing"
)

func TestAppendCountry(t *testing.T) {
	tests := []struct {
		name    string
		address string
		country string
		want    string
	}{
		{
			name:    "appends Canada",
			address: "Ottawa ON",
			country: "Canada",
			want:    "Ottawa ON, Canada",
		},
		{
			name:    "appends value as-is without normalization",
			address: "Ottawa ON",
			country: "canada",
			want:    "Ottawa ON, canada",
		},
		{
			name:    "appends USA",
			address: "Springfield IL",
			country: "USA",
			want:    "Springfield IL, USA",
		},
		{
			name:    "appends United Kingdom",
			address: "10 Downing St, London",
			country: "United Kingdom",
			want:    "10 Downing St, London, United Kingdom",
		},
		{
			name:    "appends arbitrary country code (no validation)",
			address: "10 Downing St, London",
			country: "GB",
			want:    "10 Downing St, London, GB",
		},
		{
			name:    "appends abbreviated country (no validation)",
			address: "Berlin",
			country: "DE",
			want:    "Berlin, DE",
		},
		{
			name:    "no country flag",
			address: "Ottawa ON",
			country: "",
			want:    "Ottawa ON",
		},
		{
			name:    "address already contains country",
			address: "Ottawa Ontario Canada",
			country: "Canada",
			want:    "Ottawa Ontario Canada",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := appendCountry(tt.address, tt.country)
			if got != tt.want {
				t.Errorf("appendCountry(%q, %q) = %q, want %q", tt.address, tt.country, got, tt.want)
			}
		})
	}
}

func TestAppendCountryToSlice(t *testing.T) {
	addresses := []string{"Ottawa ON", "Toronto ON"}
	result := appendCountryToAll(addresses, "Canada")

	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0] != "Ottawa ON, Canada" {
		t.Errorf("result[0] = %q, want %q", result[0], "Ottawa ON, Canada")
	}
	if result[1] != "Toronto ON, Canada" {
		t.Errorf("result[1] = %q, want %q", result[1], "Toronto ON, Canada")
	}
}

func TestDestinationFlags(t *testing.T) {
	flags := destinationFlags()

	expectedNames := []string{
		"destinations",
		"distance-mode",
		"distance-units",
		"distance-max-results",
		"distance-max-distance",
		"distance-max-duration",
		"distance-min-distance",
		"distance-min-duration",
		"distance-order-by",
		"distance-sort-order",
	}

	flagNames := make(map[string]bool)
	for _, f := range flags {
		for _, name := range f.Names() {
			flagNames[name] = true
		}
	}

	for _, expected := range expectedNames {
		if !flagNames[expected] {
			t.Errorf("missing destination flag %q", expected)
		}
	}

	// Verify aliases
	expectedAliases := map[string]string{
		"d": "destinations",
		"m": "distance-mode",
		"u": "distance-units",
	}
	for alias := range expectedAliases {
		if !flagNames[alias] {
			t.Errorf("missing alias %q", alias)
		}
	}
}
