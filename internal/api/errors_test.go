package api_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name       string
		err        *api.APIError
		wantSubstr string
	}{
		{
			name:       "with message",
			err:        &api.APIError{StatusCode: 400, Message: "Bad request"},
			wantSubstr: "Bad request",
		},
		{
			name:       "without message",
			err:        &api.APIError{StatusCode: 404},
			wantSubstr: "Not Found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Contains(t, tt.err.Error(), tt.wantSubstr)
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "404 error",
			err:  &api.APIError{StatusCode: http.StatusNotFound},
			want: true,
		},
		{
			name: "other status",
			err:  &api.APIError{StatusCode: http.StatusBadRequest},
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, api.IsNotFound(tt.err))
		})
	}
}

func TestIsRateLimited(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "429 error",
			err:  &api.APIError{StatusCode: http.StatusTooManyRequests},
			want: true,
		},
		{
			name: "other status",
			err:  &api.APIError{StatusCode: http.StatusBadRequest},
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, api.IsRateLimited(tt.err))
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "401 error",
			err:  &api.APIError{StatusCode: http.StatusUnauthorized},
			want: true,
		},
		{
			name: "403 error",
			err:  &api.APIError{StatusCode: http.StatusForbidden},
			want: true,
		},
		{
			name: "other status",
			err:  &api.APIError{StatusCode: http.StatusBadRequest},
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, api.IsUnauthorized(tt.err))
		})
	}
}
