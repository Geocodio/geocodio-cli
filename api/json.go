package api

import (
	"encoding/json"
	"github.com/geocodio/geocodio-cli/output"
)

func ParseJson(body []byte, job interface{}) error {
	jsonErr := json.Unmarshal(body, &job)
	if jsonErr != nil {
		return output.ErrorStringAndExit("Could not parse JSON from the Geocodio API")
	}
	return nil
}
