package api

import (
	"net/url"
	"testing"
)

func TestAddDestinationParams(t *testing.T) {
	t.Run("no destinations", func(t *testing.T) {
		query := url.Values{}
		addDestinationParams(query, &DestinationParams{})

		if len(query) != 0 {
			t.Errorf("expected empty query, got %v", query)
		}
	})

	t.Run("single destination", func(t *testing.T) {
		query := url.Values{}
		addDestinationParams(query, &DestinationParams{
			Destinations: []string{"New York"},
		})

		dests := query["destinations[]"]
		if len(dests) != 1 || dests[0] != "New York" {
			t.Errorf("expected destinations[]=[New York], got %v", dests)
		}
	})

	t.Run("multiple destinations", func(t *testing.T) {
		query := url.Values{}
		addDestinationParams(query, &DestinationParams{
			Destinations: []string{"New York", "Boston", "Philadelphia"},
		})

		dests := query["destinations[]"]
		if len(dests) != 3 {
			t.Errorf("expected 3 destinations, got %d", len(dests))
		}
	})

	t.Run("all params", func(t *testing.T) {
		query := url.Values{}
		addDestinationParams(query, &DestinationParams{
			Destinations: []string{"New York"},
			Mode:         "driving",
			Units:        "km",
			MaxResults:   5,
			MaxDistance:  100.5,
			MaxDuration:  3600,
			MinDistance:  10.0,
			MinDuration:  60,
			OrderBy:      "distance",
			SortOrder:    "asc",
		})

		checks := map[string]string{
			"distance_mode":         "driving",
			"distance_units":        "km",
			"distance_max_results":  "5",
			"distance_max_distance": "100.5",
			"distance_max_duration": "3600",
			"distance_min_distance": "10",
			"distance_min_duration": "60",
			"distance_order_by":     "distance",
			"distance_sort_order":   "asc",
		}

		for key, want := range checks {
			got := query.Get(key)
			if got != want {
				t.Errorf("query[%q] = %q, want %q", key, got, want)
			}
		}
	})

	t.Run("zero values not added", func(t *testing.T) {
		query := url.Values{}
		addDestinationParams(query, &DestinationParams{
			Destinations: []string{"New York"},
			// All other fields are zero values
		})

		unwanted := []string{
			"distance_mode", "distance_units", "distance_max_results",
			"distance_max_distance", "distance_max_duration",
			"distance_min_distance", "distance_min_duration",
			"distance_order_by", "distance_sort_order",
		}

		for _, key := range unwanted {
			if query.Get(key) != "" {
				t.Errorf("query should not contain %q when zero value", key)
			}
		}
	})
}

func TestValidateDistanceParams(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		units   string
		wantErr bool
	}{
		{"empty params", "", "", false},
		{"valid driving", "driving", "", false},
		{"valid straightline", "straightline", "", false},
		{"valid miles", "", "miles", false},
		{"valid km", "", "km", false},
		{"valid both", "driving", "miles", false},
		{"invalid mode", "walking", "", true},
		{"invalid units", "", "meters", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDistanceParams(tt.mode, tt.units)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDistanceParams(%q, %q) error = %v, wantErr %v", tt.mode, tt.units, err, tt.wantErr)
			}
		})
	}
}

func TestAddDistanceFilters(t *testing.T) {
	t.Run("zero values not added", func(t *testing.T) {
		query := url.Values{}
		addDistanceFilters(query, &DistanceFilters{})

		if len(query) != 0 {
			t.Errorf("expected empty query, got %v", query)
		}
	})

	t.Run("all filters", func(t *testing.T) {
		query := url.Values{}
		addDistanceFilters(query, &DistanceFilters{
			MaxResults:  5,
			MaxDistance: 100.5,
			MinDistance: 10,
			MaxDuration: 3600,
			MinDuration: 60,
			OrderBy:     "duration",
			SortOrder:   "desc",
		})

		// The standalone distance endpoints take these unprefixed, unlike the
		// distance_-prefixed spelling used by geocode and reverse.
		checks := map[string]string{
			"max_results":  "5",
			"max_distance": "100.5",
			"min_distance": "10",
			"max_duration": "3600",
			"min_duration": "60",
			"order_by":     "duration",
			"sort_order":   "desc",
		}

		for key, want := range checks {
			if got := query.Get(key); got != want {
				t.Errorf("query[%q] = %q, want %q", key, got, want)
			}
		}
	})
}

func TestAddDistanceFiltersToBody(t *testing.T) {
	t.Run("zero values not added", func(t *testing.T) {
		body := map[string]interface{}{}
		addDistanceFiltersToBody(body, &DistanceFilters{})

		if len(body) != 0 {
			t.Errorf("expected empty body, got %v", body)
		}
	})

	t.Run("all filters", func(t *testing.T) {
		body := map[string]interface{}{}
		addDistanceFiltersToBody(body, &DistanceFilters{
			MaxResults:  5,
			MaxDistance: 100.5,
			MinDistance: 10,
			MaxDuration: 3600,
			MinDuration: 60,
			OrderBy:     "duration",
			SortOrder:   "desc",
		})

		checks := map[string]interface{}{
			"max_results":  5,
			"max_distance": 100.5,
			"min_distance": 10.0,
			"max_duration": 3600,
			"min_duration": 60,
			"order_by":     "duration",
			"sort_order":   "desc",
		}

		for key, want := range checks {
			if got := body[key]; got != want {
				t.Errorf("body[%q] = %v, want %v", key, got, want)
			}
		}
	})
}

func TestValidateDistanceFilters(t *testing.T) {
	tests := []struct {
		name    string
		filters DistanceFilters
		wantErr bool
	}{
		{"empty", DistanceFilters{}, false},
		{"radius only", DistanceFilters{MaxDistance: 25}, false},
		{"valid order by distance", DistanceFilters{OrderBy: "distance"}, false},
		{"valid order by duration", DistanceFilters{OrderBy: "duration"}, false},
		{"valid sort asc", DistanceFilters{SortOrder: "asc"}, false},
		{"valid sort desc", DistanceFilters{SortOrder: "desc"}, false},
		{"valid distance range", DistanceFilters{MinDistance: 5, MaxDistance: 25}, false},
		{"valid duration range", DistanceFilters{MinDuration: 60, MaxDuration: 600}, false},
		{"invalid order by", DistanceFilters{OrderBy: "price"}, true},
		{"invalid sort order", DistanceFilters{SortOrder: "descending"}, true},
		{"inverted distance range", DistanceFilters{MinDistance: 50, MaxDistance: 25}, true},
		{"inverted duration range", DistanceFilters{MinDuration: 600, MaxDuration: 60}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDistanceFilters(&tt.filters)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDistanceFilters(%+v) error = %v, wantErr %v", tt.filters, err, tt.wantErr)
			}
		})
	}
}
