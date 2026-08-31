package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// addDestinationParams adds inline destination parameters to a query string.
// Used by geocode and reverse endpoints to calculate distances as part of the request.
func addDestinationParams(query url.Values, p *DestinationParams) {
	if len(p.Destinations) == 0 {
		return
	}
	for _, d := range p.Destinations {
		query.Add("destinations[]", d)
	}
	if p.Mode != "" {
		query.Set("distance_mode", p.Mode)
	}
	if p.Units != "" {
		query.Set("distance_units", p.Units)
	}
	if p.MaxResults > 0 {
		query.Set("distance_max_results", strconv.Itoa(p.MaxResults))
	}
	if p.MaxDistance > 0 {
		query.Set("distance_max_distance", strconv.FormatFloat(p.MaxDistance, 'f', -1, 64))
	}
	if p.MaxDuration > 0 {
		query.Set("distance_max_duration", strconv.Itoa(p.MaxDuration))
	}
	if p.MinDistance > 0 {
		query.Set("distance_min_distance", strconv.FormatFloat(p.MinDistance, 'f', -1, 64))
	}
	if p.MinDuration > 0 {
		query.Set("distance_min_duration", strconv.Itoa(p.MinDuration))
	}
	if p.OrderBy != "" {
		query.Set("distance_order_by", p.OrderBy)
	}
	if p.SortOrder != "" {
		query.Set("distance_sort_order", p.SortOrder)
	}
}

// validateDistanceParams validates mode and units parameters for distance calculations.
func validateDistanceParams(mode, units string) error {
	if mode != "" && mode != "driving" && mode != "straightline" {
		return fmt.Errorf("invalid mode %q: must be \"driving\" or \"straightline\"", mode)
	}
	if units != "" && units != "miles" && units != "km" {
		return fmt.Errorf("invalid units %q: must be \"miles\" or \"km\"", units)
	}
	return nil
}

// validateDistanceFilters validates the filtering and sorting parameters shared by
// the /distance and /distance-matrix endpoints.
func validateDistanceFilters(f *DistanceFilters) error {
	if f.OrderBy != "" && f.OrderBy != "distance" && f.OrderBy != "duration" {
		return fmt.Errorf("invalid order_by %q: must be \"distance\" or \"duration\"", f.OrderBy)
	}
	if f.SortOrder != "" && f.SortOrder != "asc" && f.SortOrder != "desc" {
		return fmt.Errorf("invalid sort_order %q: must be \"asc\" or \"desc\"", f.SortOrder)
	}
	if f.MaxDistance > 0 && f.MinDistance > f.MaxDistance {
		return fmt.Errorf("min_distance (%g) must not be greater than max_distance (%g)", f.MinDistance, f.MaxDistance)
	}
	if f.MaxDuration > 0 && f.MinDuration > f.MaxDuration {
		return fmt.Errorf("min_duration (%d) must not be greater than max_duration (%d)", f.MinDuration, f.MaxDuration)
	}
	return nil
}

// addDistanceFilters adds the filtering and sorting parameters to a query string.
// Zero values are skipped so the API keeps its own defaults.
func addDistanceFilters(query url.Values, f *DistanceFilters) {
	if f.MaxResults > 0 {
		query.Set("max_results", strconv.Itoa(f.MaxResults))
	}
	if f.MaxDistance > 0 {
		query.Set("max_distance", strconv.FormatFloat(f.MaxDistance, 'f', -1, 64))
	}
	if f.MinDistance > 0 {
		query.Set("min_distance", strconv.FormatFloat(f.MinDistance, 'f', -1, 64))
	}
	if f.MaxDuration > 0 {
		query.Set("max_duration", strconv.Itoa(f.MaxDuration))
	}
	if f.MinDuration > 0 {
		query.Set("min_duration", strconv.Itoa(f.MinDuration))
	}
	if f.OrderBy != "" {
		query.Set("order_by", f.OrderBy)
	}
	if f.SortOrder != "" {
		query.Set("sort_order", f.SortOrder)
	}
}

// addDistanceFiltersToBody adds the filtering and sorting parameters to a JSON
// request body, for the endpoints that take their parameters as POST bodies.
func addDistanceFiltersToBody(body map[string]interface{}, f *DistanceFilters) {
	if f.MaxResults > 0 {
		body["max_results"] = f.MaxResults
	}
	if f.MaxDistance > 0 {
		body["max_distance"] = f.MaxDistance
	}
	if f.MinDistance > 0 {
		body["min_distance"] = f.MinDistance
	}
	if f.MaxDuration > 0 {
		body["max_duration"] = f.MaxDuration
	}
	if f.MinDuration > 0 {
		body["min_duration"] = f.MinDuration
	}
	if f.OrderBy != "" {
		body["order_by"] = f.OrderBy
	}
	if f.SortOrder != "" {
		body["sort_order"] = f.SortOrder
	}
}

// Distance calculates driving distance and duration from a single origin to one or more destinations.
// Req.Mode can be "driving" (default) or "straightline"; Req.Units can be "miles" (default) or "km".
// The embedded DistanceFilters optionally limit the destinations to a radius, a
// travel time, or the N nearest.
func (c *Client) Distance(ctx context.Context, req *DistanceRequest) (*DistanceResponse, error) {
	if err := validateDistanceParams(req.Mode, req.Units); err != nil {
		return nil, err
	}
	if err := validateDistanceFilters(&req.DistanceFilters); err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("origin", req.Origin)
	for _, dest := range req.Destinations {
		query.Add("destinations[]", dest)
	}

	if req.Mode != "" {
		query.Set("mode", req.Mode)
	}
	if req.Units != "" {
		query.Set("units", req.Units)
	}
	addDistanceFilters(query, &req.DistanceFilters)

	var resp DistanceResponse
	if err := c.get(ctx, "/distance", query, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DistanceMatrix calculates distances between multiple origins and destinations.
// Returns a matrix of distance/duration results for all origin-destination pairs.
// Req.Mode can be "driving" (default) or "straightline"; Req.Units can be "miles" (default) or "km".
// The embedded DistanceFilters optionally limit each origin's destinations to a
// radius, a travel time, or the N nearest.
func (c *Client) DistanceMatrix(ctx context.Context, req *DistanceMatrixRequest) (*DistanceMatrixResponse, error) {
	if err := validateDistanceParams(req.Mode, req.Units); err != nil {
		return nil, err
	}
	if err := validateDistanceFilters(&req.DistanceFilters); err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"origins":      req.Origins,
		"destinations": req.Destinations,
	}
	if req.Mode != "" {
		body["mode"] = req.Mode
	}
	if req.Units != "" {
		body["units"] = req.Units
	}
	addDistanceFiltersToBody(body, &req.DistanceFilters)

	var resp DistanceMatrixResponse
	if err := c.post(ctx, "/distance-matrix", nil, body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateDistanceJob creates an asynchronous distance matrix calculation job.
// Use this for large-scale distance calculations that would exceed API limits.
func (c *Client) CreateDistanceJob(ctx context.Context, req *DistanceJobCreateRequest) (*DistanceJobResponse, error) {
	body := map[string]interface{}{
		"name":         req.Name,
		"origins":      req.Origins,
		"destinations": req.Destinations,
	}

	if req.Mode != "" {
		body["distance_mode"] = req.Mode
	}
	if req.Units != "" {
		body["units"] = req.Units
	}

	var job DistanceJob
	if err := c.post(ctx, "/distance-jobs", nil, body, &job); err != nil {
		return nil, err
	}

	return &DistanceJobResponse{Data: &job}, nil
}

// ListDistanceJobs returns all distance jobs for the account.
func (c *Client) ListDistanceJobs(ctx context.Context) (*DistanceJobListResponse, error) {
	var resp DistanceJobListResponse
	if err := c.get(ctx, "/distance-jobs", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDistanceJob retrieves the status and details of a specific distance job.
func (c *Client) GetDistanceJob(ctx context.Context, identifier string) (*DistanceJobResponse, error) {
	var resp DistanceJobResponse
	if err := c.get(ctx, fmt.Sprintf("/distance-jobs/%s", identifier), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DownloadDistanceJob downloads the results of a completed distance job as CSV.
func (c *Client) DownloadDistanceJob(ctx context.Context, identifier string) ([]byte, error) {
	return c.doRaw(ctx, "GET", fmt.Sprintf("/distance-jobs/%s/download", identifier), nil)
}

// DeleteDistanceJob deletes a distance job and its results.
func (c *Client) DeleteDistanceJob(ctx context.Context, identifier string) error {
	return c.delete(ctx, fmt.Sprintf("/distance-jobs/%s", identifier), nil)
}

const pollInterval = 2 * time.Second

// PollDistanceJob polls a distance job until it completes or fails.
// The optional callback is invoked on each poll with the current job status.
func (c *Client) PollDistanceJob(ctx context.Context, identifier string, callback func(*DistanceJobResponse)) (*DistanceJobResponse, error) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		job, err := c.GetDistanceJob(ctx, identifier)
		if err != nil {
			return nil, err
		}

		if callback != nil {
			callback(job)
		}

		status := ""
		if job.Data != nil {
			status = job.Data.Status
		}
		if status == "COMPLETED" || status == "FAILED" {
			return job, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}
