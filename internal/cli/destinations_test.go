package cli

import (
	"testing"
)

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
