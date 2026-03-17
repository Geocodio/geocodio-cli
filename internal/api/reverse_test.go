package api_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReverseGeocode(t *testing.T) {
	client := newTestClient(t, "reverse_single")

	resp, err := client.ReverseGeocode(context.Background(), &api.ReverseGeocodeRequest{
		Lat: 38.8976763,
		Lng: -77.0365298,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Results)

	result := resp.Results[0]
	assert.Contains(t, result.FormattedAddress, "Washington")
}

func TestReverseGeocodeSkipGeocoding(t *testing.T) {
	client := newTestClient(t, "reverse_skip_geocoding")

	resp, err := client.ReverseGeocode(context.Background(), &api.ReverseGeocodeRequest{
		Lat:           38.8976763,
		Lng:           -77.0365298,
		Fields:        []string{"timezone"},
		SkipGeocoding: true,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Results)

	result := resp.Results[0]
	assert.Empty(t, result.FormattedAddress, "expected empty address when skip geocoding")
}

func TestReverseGeocodeWithDestinations(t *testing.T) {
	client := newTestClient(t, "reverse_destinations")

	resp, err := client.ReverseGeocode(context.Background(), &api.ReverseGeocodeRequest{
		Lat:   38.8976763,
		Lng:   -77.0365298,
		Limit: 1,
		DestinationParams: api.DestinationParams{
			Destinations: []string{"New York"},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Results)

	result := resp.Results[0]
	assert.NotEmpty(t, result.Destinations, "expected destinations in response")
}

func TestBatchReverseGeocode(t *testing.T) {
	client := newTestClient(t, "reverse_batch")

	resp, err := client.BatchReverseGeocode(context.Background(), &api.BatchReverseGeocodeRequest{
		Coordinates: []api.Location{
			{Lat: 38.8976763, Lng: -77.0365298},
			{Lat: 37.3318, Lng: -122.0312},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.Results, 2)
}

func TestBatchReverseGeocodeWithDestinations_AddsQueryParams(t *testing.T) {
	transport := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		q := r.URL.Query()

		require.Contains(t, q["destinations[]"], "New York")
		require.Equal(t, "straightline", q.Get("distance_mode"))

		return &http.Response{
			StatusCode: 200,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"results":[]}`)),
			Request:    r,
		}, nil
	})

	client := api.NewClient(
		"https://api.geocod.io/v1.9",
		"test-api-key",
		api.WithHTTPClient(&http.Client{Transport: transport}),
	)

	_, err := client.BatchReverseGeocode(context.Background(), &api.BatchReverseGeocodeRequest{
		Coordinates: []api.Location{
			{Lat: 38.8976763, Lng: -77.0365298},
		},
		DestinationParams: api.DestinationParams{
			Destinations: []string{"New York"},
			Mode:         "straightline",
		},
	})
	require.NoError(t, err)
}
