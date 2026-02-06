package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/geocodio/geocodio-cli/internal/api"
)

func TestAgent_FormatGeocode(t *testing.T) {
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
			wantContains: []string{
				"## Geocode Result",
				"| Field | Value |",
				"| Matched Address | 1600 Pennsylvania Ave NW",
				"| Coordinates | 38.8977",
				"| Accuracy | rooftop",
			},
		},
		{
			name: "multiple results",
			resp: &api.GeocodeResponse{
				Results: []api.GeocodeResult{
					{FormattedAddress: "Address 1", Location: api.Location{Lat: 1, Lng: 2}, AccuracyType: "rooftop"},
					{FormattedAddress: "Address 2", Location: api.Location{Lat: 3, Lng: 4}, AccuracyType: "place"},
				},
			},
			wantContains: []string{
				"### Result 1 of 2",
				"### Result 2 of 2",
				"Address 1",
				"Address 2",
			},
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
			a := NewAgent(&buf)

			err := a.FormatGeocode(tt.resp)
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

func TestAgent_FormatBatchGeocode(t *testing.T) {
	var buf bytes.Buffer
	a := NewAgent(&buf)

	resp := &api.BatchGeocodeResponse{
		Results: []api.BatchGeocodeResult{
			{
				Query: "test address",
				Response: api.GeocodeResponse{
					Results: []api.GeocodeResult{
						{FormattedAddress: "Formatted Address", AccuracyType: "rooftop"},
					},
				},
			},
		},
	}

	err := a.FormatBatchGeocode(resp)
	if err != nil {
		t.Fatalf("FormatBatchGeocode() error = %v", err)
	}

	output := buf.String()
	wantContains := []string{
		"## Batch Geocode Results",
		"### Query: test address",
		"Formatted Address",
	}
	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}

func TestAgent_FormatDistance(t *testing.T) {
	durationSecs := 900
	var buf bytes.Buffer
	a := NewAgent(&buf)

	resp := &api.DistanceResponse{
		Origin: &api.DistanceLocation{
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
	}

	err := a.FormatDistance(resp)
	if err != nil {
		t.Fatalf("FormatDistance() error = %v", err)
	}

	output := buf.String()
	wantContains := []string{
		"## Distance Results",
		"| From | To | Distance | Duration |",
		"Origin Addr",
		"Dest Addr",
		"10.5 mi",
		"15 min",
	}
	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}

func TestAgent_FormatDistanceJob(t *testing.T) {
	var buf bytes.Buffer
	a := NewAgent(&buf)

	resp := &api.DistanceJobResponse{
		Data: &api.DistanceJob{
			Identifier: "abc123",
			Name:       "test job",
			Status:     "COMPLETED",
			Progress:   100,
			CreatedAt:  "2024-01-15T10:00:00Z",
		},
	}

	err := a.FormatDistanceJob(resp)
	if err != nil {
		t.Fatalf("FormatDistanceJob() error = %v", err)
	}

	output := buf.String()
	wantContains := []string{
		"## Distance Job",
		"| Field | Value |",
		"| Job ID | abc123 |",
		"| Name | test job |",
		"| Status | COMPLETED |",
	}
	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}

func TestAgent_FormatDistanceJobList(t *testing.T) {
	var buf bytes.Buffer
	a := NewAgent(&buf)

	resp := &api.DistanceJobListResponse{
		Jobs: []api.DistanceJob{
			{Identifier: "abc123", Name: "job1", Status: "COMPLETED", Progress: 100},
			{Identifier: "def456", Name: "job2", Status: "PROCESSING", Progress: 50},
		},
	}

	err := a.FormatDistanceJobList(resp)
	if err != nil {
		t.Fatalf("FormatDistanceJobList() error = %v", err)
	}

	output := buf.String()
	wantContains := []string{
		"## Distance Jobs",
		"| ID | Name | Status | Progress | Created |",
		"abc123",
		"job1",
		"COMPLETED",
	}
	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}

func TestAgent_FormatList(t *testing.T) {
	var buf bytes.Buffer
	a := NewAgent(&buf)

	resp := &api.ListResponse{
		ID:   456,
		File: &api.ListFile{Filename: "test.csv"},
		Status: &api.ListStatus{
			State:    "COMPLETED",
			Progress: 100,
		},
	}

	err := a.FormatList(resp)
	if err != nil {
		t.Fatalf("FormatList() error = %v", err)
	}

	output := buf.String()
	wantContains := []string{
		"## List",
		"| Field | Value |",
		"| List ID | 456 |",
		"| Filename | test.csv |",
		"| Status | COMPLETED |",
	}
	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}

func TestAgent_FormatListList(t *testing.T) {
	var buf bytes.Buffer
	a := NewAgent(&buf)

	resp := &api.ListListResponse{
		Lists: []api.ListResponse{
			{ID: 1, File: &api.ListFile{Filename: "file1.csv"}, Status: &api.ListStatus{State: "COMPLETED"}},
			{ID: 2, File: &api.ListFile{Filename: "file2.csv"}, Status: &api.ListStatus{State: "PROCESSING"}},
		},
	}

	err := a.FormatListList(resp)
	if err != nil {
		t.Fatalf("FormatListList() error = %v", err)
	}

	output := buf.String()
	wantContains := []string{
		"## Lists",
		"| ID | Filename | Status | Progress |",
		"file1.csv",
		"file2.csv",
		"COMPLETED",
		"PROCESSING",
	}
	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}

func TestAgent_FormatError(t *testing.T) {
	var buf bytes.Buffer
	a := NewAgent(&buf)

	testErr := errors.New("something went wrong")
	err := a.FormatError(testErr)
	if err != nil {
		t.Fatalf("FormatError() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "**Error:**") {
		t.Errorf("output missing '**Error:**' prefix:\n%s", output)
	}
	if !strings.Contains(output, "something went wrong") {
		t.Errorf("output missing error message:\n%s", output)
	}
}

func TestAgent_FormatMessage(t *testing.T) {
	var buf bytes.Buffer
	a := NewAgent(&buf)

	err := a.FormatMessage("operation successful")
	if err != nil {
		t.Fatalf("FormatMessage() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "operation successful") {
		t.Errorf("output missing message:\n%s", output)
	}
}
