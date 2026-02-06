package api_test

import (
	"context"
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
