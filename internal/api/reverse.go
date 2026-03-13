package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ReverseGeocode reverse geocodes a single coordinate pair.
func (c *Client) ReverseGeocode(ctx context.Context, req *ReverseGeocodeRequest) (*GeocodeResponse, error) {
	query := url.Values{}
	query.Set("q", fmt.Sprintf("%f,%f", req.Lat, req.Lng))

	if len(req.Fields) > 0 {
		query.Set("fields", strings.Join(req.Fields, ","))
	}
	if req.Limit > 0 {
		query.Set("limit", strconv.Itoa(req.Limit))
	}
	if req.SkipGeocoding {
		query.Set("skipGeocoding", "true")
	}

	addDestinationParams(query, &req.DestinationParams)

	var resp GeocodeResponse
	if err := c.get(ctx, "/reverse", query, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// BatchReverseGeocode reverse geocodes multiple coordinate pairs in a single request.
func (c *Client) BatchReverseGeocode(ctx context.Context, req *BatchReverseGeocodeRequest) (*BatchReverseGeocodeResponse, error) {
	query := url.Values{}

	if len(req.Fields) > 0 {
		query.Set("fields", strings.Join(req.Fields, ","))
	}
	if req.Limit > 0 {
		query.Set("limit", strconv.Itoa(req.Limit))
	}

	// Convert coordinates to string format expected by API
	coords := make([]string, len(req.Coordinates))
	for i, coord := range req.Coordinates {
		coords[i] = fmt.Sprintf("%f,%f", coord.Lat, coord.Lng)
	}

	var resp BatchReverseGeocodeResponse
	if err := c.post(ctx, "/reverse", query, coords, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
