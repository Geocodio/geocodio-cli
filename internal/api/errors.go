// Package api provides a client for the Geocodio API.
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error response from the Geocodio API.
type APIError struct {
	StatusCode int
	Message    string `json:"error"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("geocodio API error (%d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("geocodio API error (%d): %s", e.StatusCode, http.StatusText(e.StatusCode))
}

// parseAPIError attempts to parse an error response from the API.
func parseAPIError(resp *http.Response) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiErr
	}

	// Try to parse JSON error message
	var errResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil {
		if errResp.Error != "" {
			apiErr.Message = errResp.Error
		} else if errResp.Message != "" {
			apiErr.Message = errResp.Message
		} else {
			apiErr.Message = string(body)
		}
	} else {
		apiErr.Message = string(body)
	}

	return apiErr
}

// IsNotFound returns true if the error is a 404 Not Found error.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsRateLimited returns true if the error is a 429 rate limit error.
func IsRateLimited(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// IsUnauthorized returns true if the error is a 401/403 authentication error.
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized || apiErr.StatusCode == http.StatusForbidden
	}
	return false
}
