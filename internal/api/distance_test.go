package api_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDistance(t *testing.T) {
	client := newTestClient(t, "distance_single")

	resp, err := client.Distance(context.Background(),
		"Washington DC",
		[]string{"New York"},
		"driving",
		"miles",
	)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Destinations)

	dest := resp.Destinations[0]
	assert.Greater(t, dest.DistanceMiles, 200.0)
	assert.Greater(t, dest.DistanceKm, 300.0)
	assert.NotNil(t, dest.DurationSeconds)
}

func TestDistanceMultipleDestinations(t *testing.T) {
	client := newTestClient(t, "distance_multiple")

	resp, err := client.Distance(context.Background(),
		"Washington DC",
		[]string{"New York", "Boston"},
		"driving",
		"miles",
	)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.Destinations, 2)
}

func TestDistanceStraightline(t *testing.T) {
	client := newTestClient(t, "distance_straightline")

	resp, err := client.Distance(context.Background(),
		"Washington DC",
		[]string{"New York"},
		"straightline",
		"miles",
	)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Destinations)

	dest := resp.Destinations[0]
	assert.Greater(t, dest.DistanceMiles, 0.0)
	assert.Nil(t, dest.DurationSeconds)
}

func TestDistanceMatrix(t *testing.T) {
	client := newTestClient(t, "distance_matrix")

	resp, err := client.DistanceMatrix(context.Background(),
		[]string{"Washington DC", "New York"},
		[]string{"Boston", "Philadelphia"},
		"driving",
		"miles",
	)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.Results, 2)

	for _, r := range resp.Results {
		assert.NotNil(t, r.Origin)
		assert.Len(t, r.Destinations, 2)
	}
}
