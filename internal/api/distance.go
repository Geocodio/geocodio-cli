package api

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

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

// Distance calculates driving distance and duration from a single origin to one or more destinations.
// The mode parameter can be "driving" (default) or "straightline".
// The units parameter can be "miles" (default) or "km".
func (c *Client) Distance(ctx context.Context, origin string, destinations []string, mode, units string) (*DistanceResponse, error) {
	if err := validateDistanceParams(mode, units); err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("origin", origin)
	for _, dest := range destinations {
		query.Add("destinations[]", dest)
	}

	if mode != "" {
		query.Set("mode", mode)
	}
	if units != "" {
		query.Set("units", units)
	}

	var resp DistanceResponse
	if err := c.get(ctx, "/distance", query, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DistanceMatrix calculates distances between multiple origins and destinations.
// Returns a matrix of distance/duration results for all origin-destination pairs.
// The mode parameter can be "driving" (default) or "straightline".
// The units parameter can be "miles" (default) or "km".
func (c *Client) DistanceMatrix(ctx context.Context, origins, destinations []string, mode, units string) (*DistanceMatrixResponse, error) {
	if err := validateDistanceParams(mode, units); err != nil {
		return nil, err
	}

	query := url.Values{}

	if mode != "" {
		query.Set("mode", mode)
	}
	if units != "" {
		query.Set("units", units)
	}

	body := map[string][]string{
		"origins":      origins,
		"destinations": destinations,
	}

	var resp DistanceMatrixResponse
	if err := c.post(ctx, "/distance-matrix", query, body, &resp); err != nil {
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
