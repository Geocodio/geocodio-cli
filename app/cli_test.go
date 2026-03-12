package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoArgs(t *testing.T) {
	err, output := RunAppForTesting([]string{})

	assert.Nil(t, err, "running with no args should show help without error")
	assert.Contains(t, output, "Geocodio - Geocode lists using the Geocodio API", "Output should contain expected string")
}

func TestWithDummyAPIKey(t *testing.T) {
	err, output := RunAppForTesting([]string{"--apikey=DEMO"})

	assert.Nil(t, err)
	assert.Contains(t, output, "Geocodio - Geocode lists using the Geocodio API", "Output should contain expected string")
}

