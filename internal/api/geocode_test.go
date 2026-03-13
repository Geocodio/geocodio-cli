package api_test

import (
	"context"
	"testing"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeocode(t *testing.T) {
	client := newTestClient(t, "geocode_single")

	resp, err := client.Geocode(context.Background(), &api.GeocodeRequest{
		Address: "1600 Pennsylvania Ave NW, Washington DC",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Results)

	result := resp.Results[0]
	assert.Contains(t, result.FormattedAddress, "Pennsylvania")
	assert.InDelta(t, 38.8976, result.Location.Lat, 0.01)
	assert.InDelta(t, -77.0365, result.Location.Lng, 0.01)
}

func TestGeocodeWithFields(t *testing.T) {
	client := newTestClient(t, "geocode_fields")

	resp, err := client.Geocode(context.Background(), &api.GeocodeRequest{
		Address: "1600 Pennsylvania Ave NW, Washington DC",
		Fields:  []string{"timezone"},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Results)
}

func TestGeocodeWithLimit(t *testing.T) {
	client := newTestClient(t, "geocode_limit")

	resp, err := client.Geocode(context.Background(), &api.GeocodeRequest{
		Address: "1109 N Highland St, Arlington VA",
		Limit:   1,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.LessOrEqual(t, len(resp.Results), 1)
}

func TestGeocodeWithDestinations(t *testing.T) {
	client := newTestClient(t, "geocode_destinations")

	resp, err := client.Geocode(context.Background(), &api.GeocodeRequest{
		Address: "1600 Pennsylvania Ave NW, Washington DC",
		DestinationParams: api.DestinationParams{
			Destinations: []string{"New York"},
			Mode:         "straightline",
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Results)

	result := resp.Results[0]
	assert.NotEmpty(t, result.Destinations, "expected destinations in response")
	assert.Greater(t, result.Destinations[0].DistanceMiles, 0.0)
}

func TestGeocodeStableAddressKey(t *testing.T) {
	client := newTestClient(t, "geocode_single")

	resp, err := client.Geocode(context.Background(), &api.GeocodeRequest{
		Address: "1600 Pennsylvania Ave NW, Washington DC",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Results)

	result := resp.Results[0]
	assert.NotEmpty(t, result.StableAddressKey, "expected stable address key for rooftop result")
	assert.Contains(t, result.StableAddressKey, "gcod_")
}

func TestBatchGeocode(t *testing.T) {
	client := newTestClient(t, "geocode_batch")

	resp, err := client.BatchGeocode(context.Background(), &api.BatchGeocodeRequest{
		Addresses: []string{
			"1600 Pennsylvania Ave NW, Washington DC",
			"1 Infinite Loop, Cupertino CA",
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.Results, 2)
}
