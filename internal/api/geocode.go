package api

import (
	"context"
	"net/url"
	"strconv"
	"strings"
)

// Geocode geocodes a single address.
func (c *Client) Geocode(ctx context.Context, req *GeocodeRequest) (*GeocodeResponse, error) {
	query := url.Values{}
	query.Set("q", req.Address)

	if len(req.Fields) > 0 {
		query.Set("fields", strings.Join(req.Fields, ","))
	}
	if req.Limit > 0 {
		query.Set("limit", strconv.Itoa(req.Limit))
	}
	if req.Country != "" {
		query.Set("country", req.Country)
	}

	addDestinationParams(query, &req.DestinationParams)

	var resp GeocodeResponse
	if err := c.get(ctx, "/geocode", query, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// BatchGeocode geocodes multiple addresses in a single request.
func (c *Client) BatchGeocode(ctx context.Context, req *BatchGeocodeRequest) (*BatchGeocodeResponse, error) {
	query := url.Values{}

	if len(req.Fields) > 0 {
		query.Set("fields", strings.Join(req.Fields, ","))
	}
	if req.Limit > 0 {
		query.Set("limit", strconv.Itoa(req.Limit))
	}

	var resp BatchGeocodeResponse
	if err := c.post(ctx, "/geocode", query, req.Addresses, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
