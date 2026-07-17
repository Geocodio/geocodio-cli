package api_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGeocodeResponseRoundTripsAllAPIFields guards against response fields
// silently disappearing from CLI output. The CLI decodes API responses into
// typed structs and re-encodes them for display, so any response key missing
// from the structs is dropped without error (this happened to match_type,
// address_lines, and formatted_street).
//
// The fixture is a real API response. To refresh it after an API change:
//
//	curl -s "https://api.geocod.io/v2/geocode?street=996C+East+St&street2=Unit+10&city=Walpole&state=MA&postal_code=02081&fields=census&api_key=$GEOCODIO_API_KEY" \
//	  | python3 -m json.tool > internal/api/testdata/geocode_full_response.json
func TestGeocodeResponseRoundTripsAllAPIFields(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "geocode_full_response.json"))
	require.NoError(t, err)

	var resp api.GeocodeResponse
	require.NoError(t, json.Unmarshal(raw, &resp))

	reencoded, err := json.Marshal(resp)
	require.NoError(t, err)

	var original, output interface{}
	require.NoError(t, json.Unmarshal(raw, &original))
	require.NoError(t, json.Unmarshal(reencoded, &output))

	originalPaths := map[string]bool{}
	outputPaths := map[string]bool{}
	collectKeyPaths(original, "", originalPaths)
	collectKeyPaths(output, "", outputPaths)

	var missing []string
	for path := range originalPaths {
		if !outputPaths[path] {
			missing = append(missing, path)
		}
	}
	sort.Strings(missing)

	assert.Empty(t, missing,
		"API response fields dropped by the structs in types.go; add the missing fields so CLI output includes them")
}

// collectKeyPaths records every object key path (e.g. "results[].match_type")
// that holds a non-empty value. Empty values (null, "", 0, false, [], {}) are
// skipped because the structs use omitempty, which legitimately omits them on
// re-encode. Keys prefixed with "_" (e.g. _warnings) are API metadata that the
// structs intentionally do not model.
func collectKeyPaths(value interface{}, prefix string, out map[string]bool) {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, child := range typed {
			if strings.HasPrefix(key, "_") {
				continue
			}
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			if isEmptyJSONValue(child) {
				continue
			}
			out[path] = true
			collectKeyPaths(child, path, out)
		}
	case []interface{}:
		for _, item := range typed {
			collectKeyPaths(item, prefix+"[]", out)
		}
	}
}

func isEmptyJSONValue(value interface{}) bool {
	switch typed := value.(type) {
	case nil:
		return true
	case string:
		return typed == ""
	case float64:
		return typed == 0
	case bool:
		return !typed
	case []interface{}:
		return len(typed) == 0
	case map[string]interface{}:
		return len(typed) == 0
	}
	return false
}
