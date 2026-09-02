package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDistanceSurfacesBillingHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Billable-Lookups-Count", "4")
		w.Header().Set("X-Billable-Distance-Calculations", "8")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"origin": {"query": "Washington DC"}, "destinations": []}`))
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-api-key")

	resp, err := client.Distance(context.Background(), &api.DistanceRequest{
		Origin:       "Washington DC",
		Destinations: []string{"New York"},
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Billing)
	require.NotNil(t, resp.Billing.Lookups)
	require.NotNil(t, resp.Billing.DistanceCalculations)
	assert.Equal(t, 4, *resp.Billing.Lookups)
	assert.Equal(t, 8, *resp.Billing.DistanceCalculations)
}

func TestGeocodeBillingOmittedWhenHeadersAbsent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results": []}`))
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-api-key")

	resp, err := client.Geocode(context.Background(), &api.GeocodeRequest{Address: "Washington DC"})
	require.NoError(t, err)
	assert.Nil(t, resp.Billing)
}

func TestGeocodeBillingIgnoresUnparseableHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Billable-Lookups-Count", "not-a-number")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results": []}`))
	}))
	defer server.Close()

	client := api.NewClient(server.URL, "test-api-key")

	resp, err := client.Geocode(context.Background(), &api.GeocodeRequest{Address: "Washington DC"})
	require.NoError(t, err)
	assert.Nil(t, resp.Billing)
}
