package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/geocodio/geocodio-cli/internal/api"
)

func TestJSON_FormatGeocode(t *testing.T) {
	tests := []struct {
		name string
		resp *api.GeocodeResponse
	}{
		{
			name: "single result",
			resp: &api.GeocodeResponse{
				Results: []api.GeocodeResult{
					{
						FormattedAddress: "1600 Pennsylvania Ave NW, Washington, DC",
						Location:         api.Location{Lat: 38.8977, Lng: -77.0365},
						Accuracy:         1.0,
						AccuracyType:     "rooftop",
					},
				},
			},
		},
		{
			name: "multiple results",
			resp: &api.GeocodeResponse{
				Results: []api.GeocodeResult{
					{FormattedAddress: "Address 1", Location: api.Location{Lat: 1, Lng: 2}},
					{FormattedAddress: "Address 2", Location: api.Location{Lat: 3, Lng: 4}},
				},
			},
		},
		{
			name: "empty results",
			resp: &api.GeocodeResponse{Results: []api.GeocodeResult{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			j := NewJSON(&buf)

			err := j.FormatGeocode(tt.resp)
			if err != nil {
				t.Fatalf("FormatGeocode() error = %v", err)
			}

			// Verify valid JSON output
			var decoded api.GeocodeResponse
			if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
				t.Errorf("output is not valid JSON: %v", err)
			}

			if len(decoded.Results) != len(tt.resp.Results) {
				t.Errorf("decoded results count = %d, want %d", len(decoded.Results), len(tt.resp.Results))
			}
		})
	}
}

func TestJSON_FormatBatchGeocode(t *testing.T) {
	tests := []struct {
		name string
		resp *api.BatchGeocodeResponse
	}{
		{
			name: "multiple queries",
			resp: &api.BatchGeocodeResponse{
				Results: []api.BatchGeocodeResult{
					{Query: "query1", Response: api.GeocodeResponse{Results: []api.GeocodeResult{{FormattedAddress: "addr1"}}}},
					{Query: "query2", Response: api.GeocodeResponse{Results: []api.GeocodeResult{{FormattedAddress: "addr2"}}}},
				},
			},
		},
		{
			name: "empty responses",
			resp: &api.BatchGeocodeResponse{
				Results: []api.BatchGeocodeResult{
					{Query: "query1", Response: api.GeocodeResponse{}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			j := NewJSON(&buf)

			err := j.FormatBatchGeocode(tt.resp)
			if err != nil {
				t.Fatalf("FormatBatchGeocode() error = %v", err)
			}

			var decoded api.BatchGeocodeResponse
			if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
				t.Errorf("output is not valid JSON: %v", err)
			}
		})
	}
}

func TestJSON_FormatDistance(t *testing.T) {
	durationSecs := 900
	tests := []struct {
		name string
		resp *api.DistanceResponse
	}{
		{
			name: "with duration",
			resp: &api.DistanceResponse{
				Origin: &api.DistanceLocation{Query: "Origin"},
				Destinations: []api.DistanceDestination{
					{
						Query:           "Dest",
						DistanceMiles:   10.5,
						DistanceKm:      16.9,
						DurationSeconds: &durationSecs,
					},
				},
			},
		},
		{
			name: "without duration (straightline)",
			resp: &api.DistanceResponse{
				Origin: &api.DistanceLocation{Query: "Origin"},
				Destinations: []api.DistanceDestination{
					{Query: "Dest", DistanceMiles: 5.1, DistanceKm: 8.2},
				},
			},
		},
		{
			name: "empty results",
			resp: &api.DistanceResponse{Destinations: []api.DistanceDestination{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			j := NewJSON(&buf)

			err := j.FormatDistance(tt.resp)
			if err != nil {
				t.Fatalf("FormatDistance() error = %v", err)
			}

			var decoded api.DistanceResponse
			if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
				t.Errorf("output is not valid JSON: %v", err)
			}
		})
	}
}

func TestJSON_FormatDistanceMatrix(t *testing.T) {
	durationSecs := 900
	tests := []struct {
		name string
		resp *api.DistanceMatrixResponse
	}{
		{
			name: "multiple origins",
			resp: &api.DistanceMatrixResponse{
				Mode: "driving",
				Results: []api.DistanceMatrixResult{
					{
						Origin: &api.DistanceLocation{Query: "Origin 1"},
						Destinations: []api.DistanceDestination{
							{Query: "Dest 1", DistanceMiles: 10.5, DistanceKm: 16.9, DurationSeconds: &durationSecs},
						},
					},
					{
						Origin: &api.DistanceLocation{Query: "Origin 2"},
						Destinations: []api.DistanceDestination{
							{Query: "Dest 1", DistanceMiles: 5.2, DistanceKm: 8.4},
						},
					},
				},
			},
		},
		{
			name: "empty results",
			resp: &api.DistanceMatrixResponse{Results: []api.DistanceMatrixResult{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			j := NewJSON(&buf)

			err := j.FormatDistanceMatrix(tt.resp)
			if err != nil {
				t.Fatalf("FormatDistanceMatrix() error = %v", err)
			}

			var decoded api.DistanceMatrixResponse
			if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
				t.Errorf("output is not valid JSON: %v", err)
			}
		})
	}
}

func TestJSON_FormatDistanceJob(t *testing.T) {
	var buf bytes.Buffer
	j := NewJSON(&buf)

	resp := &api.DistanceJobResponse{
		Data: &api.DistanceJob{
			Identifier: "abc123",
			Status:     "COMPLETED",
			Progress:   100,
		},
	}

	err := j.FormatDistanceJob(resp)
	if err != nil {
		t.Fatalf("FormatDistanceJob() error = %v", err)
	}

	var decoded api.DistanceJobResponse
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}

	if decoded.Data == nil || decoded.Data.Identifier != "abc123" {
		t.Errorf("decoded identifier = %v, want abc123", decoded.Data)
	}
}

func TestJSON_FormatDistanceJobList(t *testing.T) {
	var buf bytes.Buffer
	j := NewJSON(&buf)

	resp := &api.DistanceJobListResponse{
		Jobs: []api.DistanceJob{
			{Identifier: "abc123", Status: "COMPLETED"},
			{Identifier: "def456", Status: "PROCESSING"},
		},
	}

	err := j.FormatDistanceJobList(resp)
	if err != nil {
		t.Fatalf("FormatDistanceJobList() error = %v", err)
	}

	var decoded api.DistanceJobListResponse
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}

	if len(decoded.Jobs) != 2 {
		t.Errorf("decoded jobs count = %d, want 2", len(decoded.Jobs))
	}
}

func TestJSON_FormatList(t *testing.T) {
	var buf bytes.Buffer
	j := NewJSON(&buf)

	resp := &api.ListResponse{
		ID:     456,
		Status: &api.ListStatus{State: "COMPLETED", Progress: 100},
	}

	err := j.FormatList(resp)
	if err != nil {
		t.Fatalf("FormatList() error = %v", err)
	}

	var decoded api.ListResponse
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}

	if decoded.ID != 456 {
		t.Errorf("decoded ID = %d, want 456", decoded.ID)
	}
}

func TestJSON_FormatListList(t *testing.T) {
	var buf bytes.Buffer
	j := NewJSON(&buf)

	resp := &api.ListListResponse{
		Lists: []api.ListResponse{
			{ID: 1},
			{ID: 2},
		},
	}

	err := j.FormatListList(resp)
	if err != nil {
		t.Fatalf("FormatListList() error = %v", err)
	}

	var decoded api.ListListResponse
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}

	if len(decoded.Lists) != 2 {
		t.Errorf("decoded lists count = %d, want 2", len(decoded.Lists))
	}
}

func TestJSON_FormatError(t *testing.T) {
	var buf bytes.Buffer
	j := NewJSON(&buf)

	testErr := errors.New("test error message")
	err := j.FormatError(testErr)
	if err != nil {
		t.Fatalf("FormatError() error = %v", err)
	}

	var decoded map[string]string
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}

	if decoded["error"] != "test error message" {
		t.Errorf("decoded error = %q, want %q", decoded["error"], "test error message")
	}
}

func TestJSON_FormatMessage(t *testing.T) {
	var buf bytes.Buffer
	j := NewJSON(&buf)

	err := j.FormatMessage("test message")
	if err != nil {
		t.Fatalf("FormatMessage() error = %v", err)
	}

	var decoded map[string]string
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}

	if decoded["message"] != "test message" {
		t.Errorf("decoded message = %q, want %q", decoded["message"], "test message")
	}
}
