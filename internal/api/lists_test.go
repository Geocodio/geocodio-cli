package api_test

import (
	"context"
	"testing"

	"github.com/geocodio/geocodio-cli/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadList(t *testing.T) {
	client := newTestClient(t, "list_upload")

	csvData := []byte("address\n1600 Pennsylvania Ave NW, Washington DC\n")

	resp, err := client.UploadList(context.Background(), &api.ListUploadRequest{
		Filename:  "test.csv",
		Data:      csvData,
		Direction: "forward",
		Format:    "{{A}}",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Greater(t, resp.ID, 0)
}

func TestListLists(t *testing.T) {
	client := newTestClient(t, "list_list")

	resp, err := client.ListLists(context.Background())

	require.NoError(t, err)
	require.NotNil(t, resp)
	// May be empty, that's ok
}

func TestGetList(t *testing.T) {
	client := newTestClient(t, "list_get")

	// First upload a list to get
	csvData := []byte("address\n1600 Pennsylvania Ave NW, Washington DC\n")
	uploadResp, err := client.UploadList(context.Background(), &api.ListUploadRequest{
		Filename:  "test.csv",
		Data:      csvData,
		Direction: "forward",
		Format:    "{{A}}",
	})
	require.NoError(t, err)

	// Then get its status
	resp, err := client.GetList(context.Background(), uploadResp.ID)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uploadResp.ID, resp.ID)
}

func TestDeleteList(t *testing.T) {
	client := newTestClient(t, "list_delete")

	// First upload a list to delete
	csvData := []byte("address\n1600 Pennsylvania Ave NW, Washington DC\n")
	uploadResp, err := client.UploadList(context.Background(), &api.ListUploadRequest{
		Filename:  "test.csv",
		Data:      csvData,
		Direction: "forward",
		Format:    "{{A}}",
	})
	require.NoError(t, err)

	// Delete it
	err = client.DeleteList(context.Background(), uploadResp.ID)
	require.NoError(t, err)
}
