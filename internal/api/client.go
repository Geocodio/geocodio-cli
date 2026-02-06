package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the Geocodio API client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	userAgent  string
	debug      bool
	debugOut   io.Writer
}

// ClientOption configures the Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = c
	}
}

// WithDebug enables debug logging to the given writer.
func WithDebug(w io.Writer) ClientOption {
	return func(client *Client) {
		client.debug = true
		client.debugOut = w
	}
}

// WithUserAgent sets a custom User-Agent header for API requests.
func WithUserAgent(ua string) ClientOption {
	return func(client *Client) {
		client.userAgent = ua
	}
}

// NewClient creates a new Geocodio API client.
func NewClient(baseURL, apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL:   strings.TrimSuffix(baseURL, "/"),
		apiKey:    apiKey,
		userAgent: "geocodio-cli/dev",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// do executes an HTTP request and decodes the JSON response.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body interface{}, result interface{}) error {
	// Build URL
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Add API key to query
	if query == nil {
		query = url.Values{}
	}
	query.Set("api_key", c.apiKey)
	u.RawQuery = query.Encode()

	// Encode body if present
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encoding request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	// Debug logging
	if c.debug && c.debugOut != nil {
		fmt.Fprintf(c.debugOut, "DEBUG: %s %s\n", method, c.redactURL(u))
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	// Debug response
	if c.debug && c.debugOut != nil {
		fmt.Fprintf(c.debugOut, "DEBUG: Response status: %d\n", resp.StatusCode)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}

	// Decode response
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

// doRaw executes an HTTP request and returns the raw response body.
func (c *Client) doRaw(ctx context.Context, method, path string, query url.Values) ([]byte, error) {
	// Build URL
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Add API key to query
	if query == nil {
		query = url.Values{}
	}
	query.Set("api_key", c.apiKey)
	u.RawQuery = query.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	// Debug logging
	if c.debug && c.debugOut != nil {
		fmt.Fprintf(c.debugOut, "DEBUG: %s %s\n", method, c.redactURL(u))
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	// Debug response
	if c.debug && c.debugOut != nil {
		fmt.Fprintf(c.debugOut, "DEBUG: Response status: %d\n", resp.StatusCode)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}

	return io.ReadAll(resp.Body)
}

// redactURL returns a URL string with the API key redacted.
func (c *Client) redactURL(u *url.URL) string {
	q := u.Query()
	if q.Get("api_key") != "" {
		q.Set("api_key", "[REDACTED]")
	}
	u2 := *u
	u2.RawQuery = q.Encode()
	return u2.String()
}

// get performs a GET request.
func (c *Client) get(ctx context.Context, path string, query url.Values, result interface{}) error {
	return c.do(ctx, http.MethodGet, path, query, nil, result)
}

// post performs a POST request.
func (c *Client) post(ctx context.Context, path string, query url.Values, body interface{}, result interface{}) error {
	return c.do(ctx, http.MethodPost, path, query, body, result)
}

// delete performs a DELETE request.
func (c *Client) delete(ctx context.Context, path string, query url.Values) error {
	return c.do(ctx, http.MethodDelete, path, query, nil, nil)
}
