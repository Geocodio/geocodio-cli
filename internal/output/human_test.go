package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/geocodio/geocodio-cli/internal/api"
)

func TestHuman_FormatGeocode(t *testing.T) {
	tests := []struct {
		name         string
		resp         *api.GeocodeResponse
		wantContains []string
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
			wantContains: []string{"1600 Pennsylvania Ave NW", "38.8977", "rooftop"},
		},
		{
			name: "multiple results",
			resp: &api.GeocodeResponse{
				Results: []api.GeocodeResult{
					{FormattedAddress: "Address 1", Location: api.Location{Lat: 1, Lng: 2}, AccuracyType: "range_interpolation"},
					{FormattedAddress: "Address 2", Location: api.Location{Lat: 3, Lng: 4}, AccuracyType: "place"},
				},
			},
			wantContains: []string{"Result 1 of 2", "Result 2 of 2", "Address 1", "Address 2"},
		},
		{
			name:         "empty results",
			resp:         &api.GeocodeResponse{Results: []api.GeocodeResult{}},
			wantContains: []string{"No results found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := NewHuman(&buf, false, Options{})

			err := h.FormatGeocode(tt.resp)
			if err != nil {
				t.Fatalf("FormatGeocode() error = %v", err)
			}

			output := buf.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q:\n%s", want, output)
				}
			}
		})
	}
}

func TestHuman_FormatBatchGeocode(t *testing.T) {
	tests := []struct {
		name         string
		resp         *api.BatchGeocodeResponse
		wantContains []string
	}{
		{
			name: "multiple queries",
			resp: &api.BatchGeocodeResponse{
				Results: []api.BatchGeocodeResult{
					{Query: "query one", Response: api.GeocodeResponse{Results: []api.GeocodeResult{{FormattedAddress: "addr1", AccuracyType: "rooftop"}}}},
					{Query: "query two", Response: api.GeocodeResponse{Results: []api.GeocodeResult{{FormattedAddress: "addr2", AccuracyType: "place"}}}},
				},
			},
			wantContains: []string{"Query:", "query one", "Query:", "query two", "addr1", "addr2"},
		},
		{
			name: "empty response for query",
			resp: &api.BatchGeocodeResponse{
				Results: []api.BatchGeocodeResult{
					{Query: "nowhere", Response: api.GeocodeResponse{}},
				},
			},
			wantContains: []string{"Query:", "nowhere", "No results found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := NewHuman(&buf, false, Options{})

			err := h.FormatBatchGeocode(tt.resp)
			if err != nil {
				t.Fatalf("FormatBatchGeocode() error = %v", err)
			}

			output := buf.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q:\n%s", want, output)
				}
			}
		})
	}
}

func TestHuman_FormatDistance(t *testing.T) {
	durationSecs := 900
	tests := []struct {
		name         string
		resp         *api.DistanceResponse
		wantContains []string
		wantMissing  []string
	}{
		{
			name: "with duration",
			resp: &api.DistanceResponse{
				Origin: &api.DistanceLocation{
					Query:   "Origin Addr",
					Geocode: &api.GeocodeResult{FormattedAddress: "Origin Addr"},
				},
				Destinations: []api.DistanceDestination{
					{
						Query:           "Dest Addr",
						Geocode:         &api.GeocodeResult{FormattedAddress: "Dest Addr"},
						DistanceMiles:   10.5,
						DistanceKm:      16.9,
						DurationSeconds: &durationSecs,
					},
				},
			},
			wantContains: []string{"Origin Addr", "Dest Addr", "10.5 miles", "15 minutes"},
		},
		{
			name: "without duration (straightline)",
			resp: &api.DistanceResponse{
				Origin: &api.DistanceLocation{Query: "Origin"},
				Destinations: []api.DistanceDestination{
					{
						Query:         "Dest",
						DistanceMiles: 5.1,
						DistanceKm:    8.2,
					},
				},
			},
			wantContains: []string{"8.2 km"},
			wantMissing:  []string{"Duration"},
		},
		{
			name:         "empty results",
			resp:         &api.DistanceResponse{Destinations: []api.DistanceDestination{}},
			wantContains: []string{"No results found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := NewHuman(&buf, false, Options{})

			err := h.FormatDistance(tt.resp)
			if err != nil {
				t.Fatalf("FormatDistance() error = %v", err)
			}

			output := buf.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q:\n%s", want, output)
				}
			}
			for _, notWant := range tt.wantMissing {
				if strings.Contains(output, notWant) {
					t.Errorf("output should not contain %q:\n%s", notWant, output)
				}
			}
		})
	}
}

func TestHuman_FormatDistanceMatrix(t *testing.T) {
	durationSecs := 900
	tests := []struct {
		name         string
		resp         *api.DistanceMatrixResponse
		wantContains []string
	}{
		{
			name: "multiple origins",
			resp: &api.DistanceMatrixResponse{
				Mode: "driving",
				Results: []api.DistanceMatrixResult{
					{
						Origin: &api.DistanceLocation{Query: "Washington DC"},
						Destinations: []api.DistanceDestination{
							{Query: "New York", DistanceMiles: 225.5, DistanceKm: 363.0, DurationSeconds: &durationSecs},
						},
					},
				},
			},
			wantContains: []string{"Origin:", "Washington DC", "New York", "225.5 miles", "15 minutes"},
		},
		{
			name:         "empty results",
			resp:         &api.DistanceMatrixResponse{Results: []api.DistanceMatrixResult{}},
			wantContains: []string{"No results found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := NewHuman(&buf, false, Options{})

			err := h.FormatDistanceMatrix(tt.resp)
			if err != nil {
				t.Fatalf("FormatDistanceMatrix() error = %v", err)
			}

			output := buf.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q:\n%s", want, output)
				}
			}
		})
	}
}

func TestHuman_FormatDistanceJob(t *testing.T) {
	var buf bytes.Buffer
	h := NewHuman(&buf, false, Options{})

	resp := &api.DistanceJobResponse{
		Data: &api.DistanceJob{
			Identifier: "abc123xyz",
			Name:       "test job",
			Status:     "COMPLETED",
			Progress:   100,
			CreatedAt:  "2024-01-15T10:00:00Z",
		},
	}

	err := h.FormatDistanceJob(resp)
	if err != nil {
		t.Fatalf("FormatDistanceJob() error = %v", err)
	}

	output := buf.String()
	wantContains := []string{"abc123xyz", "test job", "COMPLETED", "100%", "2024-01-15"}
	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}

func TestHuman_FormatDistanceJobList(t *testing.T) {
	tests := []struct {
		name         string
		resp         *api.DistanceJobListResponse
		wantContains []string
	}{
		{
			name: "with jobs",
			resp: &api.DistanceJobListResponse{
				Jobs: []api.DistanceJob{
					{Identifier: "abc123", Name: "job1", Status: "COMPLETED", Progress: 100},
					{Identifier: "def456", Name: "job2", Status: "PROCESSING", Progress: 50},
				},
			},
			wantContains: []string{"ID", "Status", "abc123", "def456", "COMPLETED", "PROCESSING"},
		},
		{
			name:         "empty list",
			resp:         &api.DistanceJobListResponse{Jobs: []api.DistanceJob{}},
			wantContains: []string{"No distance jobs found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := NewHuman(&buf, false, Options{})

			err := h.FormatDistanceJobList(tt.resp)
			if err != nil {
				t.Fatalf("FormatDistanceJobList() error = %v", err)
			}

			output := buf.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q:\n%s", want, output)
				}
			}
		})
	}
}

func TestHuman_FormatList(t *testing.T) {
	var buf bytes.Buffer
	h := NewHuman(&buf, false, Options{})

	resp := &api.ListResponse{
		ID:   456,
		File: &api.ListFile{Filename: "test.csv"},
		Status: &api.ListStatus{
			State:    "COMPLETED",
			Progress: 100,
		},
		Rows: &api.ListRowCounts{
			Total:     1000,
			Processed: 998,
			Failed:    2,
		},
	}

	err := h.FormatList(resp)
	if err != nil {
		t.Fatalf("FormatList() error = %v", err)
	}

	output := buf.String()
	wantContains := []string{"456", "test.csv", "COMPLETED", "100", "1000", "998", "2"}
	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}

func TestHuman_FormatListList(t *testing.T) {
	tests := []struct {
		name         string
		resp         *api.ListListResponse
		wantContains []string
	}{
		{
			name: "with lists",
			resp: &api.ListListResponse{
				Lists: []api.ListResponse{
					{ID: 1, File: &api.ListFile{Filename: "file1.csv"}, Status: &api.ListStatus{State: "COMPLETED"}},
					{ID: 2, File: &api.ListFile{Filename: "file2.csv"}, Status: &api.ListStatus{State: "PROCESSING"}},
				},
			},
			wantContains: []string{"ID", "Filename", "file1.csv", "file2.csv", "COMPLETED", "PROCESSING"},
		},
		{
			name:         "empty list",
			resp:         &api.ListListResponse{Lists: []api.ListResponse{}},
			wantContains: []string{"No lists found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := NewHuman(&buf, false, Options{})

			err := h.FormatListList(tt.resp)
			if err != nil {
				t.Fatalf("FormatListList() error = %v", err)
			}

			output := buf.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q:\n%s", want, output)
				}
			}
		})
	}
}

func TestHuman_FormatError(t *testing.T) {
	var buf bytes.Buffer
	h := NewHuman(&buf, false, Options{})

	testErr := errors.New("something went wrong")
	err := h.FormatError(testErr)
	if err != nil {
		t.Fatalf("FormatError() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("output missing 'Error:' prefix:\n%s", output)
	}
	if !strings.Contains(output, "something went wrong") {
		t.Errorf("output missing error message:\n%s", output)
	}
}

func TestHuman_FormatMessage(t *testing.T) {
	var buf bytes.Buffer
	h := NewHuman(&buf, false, Options{})

	err := h.FormatMessage("operation successful")
	if err != nil {
		t.Fatalf("FormatMessage() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "operation successful") {
		t.Errorf("output missing message:\n%s", output)
	}
}

func TestNew(t *testing.T) {
	var buf bytes.Buffer

	jsonFormatter := New(&buf, OutputModeJSON, false)
	if _, ok := jsonFormatter.(*JSON); !ok {
		t.Errorf("New(_, OutputModeJSON, _) returned %T, want *JSON", jsonFormatter)
	}

	humanFormatter := New(&buf, OutputModeHuman, false)
	if _, ok := humanFormatter.(*Human); !ok {
		t.Errorf("New(_, OutputModeHuman, _) returned %T, want *Human", humanFormatter)
	}

	agentFormatter := New(&buf, OutputModeAgent, false)
	if _, ok := agentFormatter.(*Agent); !ok {
		t.Errorf("New(_, OutputModeAgent, _) returned %T, want *Agent", agentFormatter)
	}
}
