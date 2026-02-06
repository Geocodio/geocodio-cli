package api_test

import (
	"context"
	"testing"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDistanceJob(t *testing.T) {
	client := newTestClient(t, "distance_job_create")

	resp, err := client.CreateDistanceJob(context.Background(), &api.DistanceJobCreateRequest{
		Name:         "test job",
		Origins:      []string{"Washington DC"},
		Destinations: []string{"New York"},
		Mode:         "driving",
		Units:        "miles",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	assert.NotEmpty(t, resp.Data.Identifier)
	assert.NotEmpty(t, resp.Data.Status)
}

func TestListDistanceJobs(t *testing.T) {
	client := newTestClient(t, "distance_job_list")

	resp, err := client.ListDistanceJobs(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	// May be empty if no jobs exist, that's ok
}

func TestGetDistanceJob(t *testing.T) {
	client := newTestClient(t, "distance_job_get")

	// First create a job to get
	createResp, err := client.CreateDistanceJob(context.Background(), &api.DistanceJobCreateRequest{
		Name:         "test job for get",
		Origins:      []string{"Washington DC"},
		Destinations: []string{"New York"},
	})
	require.NoError(t, err)
	require.NotNil(t, createResp.Data)

	// Then get its status
	resp, err := client.GetDistanceJob(context.Background(), createResp.Data.Identifier)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	assert.Equal(t, createResp.Data.Identifier, resp.Data.Identifier)
}

func TestDeleteDistanceJob(t *testing.T) {
	client := newTestClient(t, "distance_job_delete")

	// First create a job to delete
	createResp, err := client.CreateDistanceJob(context.Background(), &api.DistanceJobCreateRequest{
		Name:         "test job for delete",
		Origins:      []string{"Washington DC"},
		Destinations: []string{"New York"},
	})
	require.NoError(t, err)
	require.NotNil(t, createResp.Data)

	// Delete it - may return 404 if already processed/removed
	err = client.DeleteDistanceJob(context.Background(), createResp.Data.Identifier)
	if err != nil && !api.IsNotFound(err) {
		t.Fatalf("unexpected error: %v", err)
	}
}
