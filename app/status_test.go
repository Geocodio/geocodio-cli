package app

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestStatusWithoutArgument(t *testing.T) {
	err, output := RunAppForTesting([]string{"--apikey=" + os.Getenv("GEOCODIO_TEST_API_KEY"), "status"})

	assert.Contains(t, err.Error(), "Invalid spreadsheet job id specified")
	assert.Equal(t, "", output, "No standard output should be present")
}

func TestStatusWithInvalidSpreadsheetId(t *testing.T) {
	err, output := RunAppForTesting([]string{"--apikey=" + os.Getenv("GEOCODIO_TEST_API_KEY"), "status", "1"})

	assert.Contains(t, err.Error(), " No spreadsheet job with that id found")
	assert.Equal(t, "", output, "No standard output should be present")
}

func TestStatusWithValidSpreadsheetId(t *testing.T) {
	err, output := RunAppForTesting([]string{"--apikey=" + os.Getenv("GEOCODIO_TEST_API_KEY"), "status", "11471563"})

	assert.Nil(t, err)
	assert.Contains(t, output, "State: COMPLETED")
}