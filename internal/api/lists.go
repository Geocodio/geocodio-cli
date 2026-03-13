package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// UploadList uploads a spreadsheet file for batch geocoding.
// The file is processed asynchronously; use PollList to monitor progress.
func (c *Client) UploadList(ctx context.Context, req *ListUploadRequest) (*ListResponse, error) {
	u, err := url.Parse(c.baseURL + "/lists")
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	query := url.Values{}
	query.Set("api_key", c.apiKey)
	query.Set("direction", req.Direction)
	query.Set("format", req.Format)
	if req.Callback != "" {
		query.Set("callback", req.Callback)
	}
	if len(req.Fields) > 0 {
		query.Set("fields", strings.Join(req.Fields, ","))
	}
	u.RawQuery = query.Encode()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(req.Filename))
	if err != nil {
		return nil, fmt.Errorf("creating form file: %w", err)
	}

	if _, err := part.Write(req.Data); err != nil {
		return nil, fmt.Errorf("writing file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("closing multipart writer: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", c.userAgent)

	if c.debug && c.debugOut != nil {
		fmt.Fprintf(c.debugOut, "DEBUG: POST %s\n", c.redactURL(u))
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if c.debug && c.debugOut != nil {
		fmt.Fprintf(c.debugOut, "DEBUG: Response status: %d\n", resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}

	var listResp ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &listResp, nil
}

// ListLists returns all uploaded spreadsheet lists for the account.
func (c *Client) ListLists(ctx context.Context) (*ListListResponse, error) {
	var resp ListListResponse
	if err := c.get(ctx, "/lists", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetList retrieves the status and details of a specific list.
func (c *Client) GetList(ctx context.Context, id int) (*ListResponse, error) {
	var resp ListResponse
	if err := c.get(ctx, fmt.Sprintf("/lists/%d", id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DownloadList downloads the results of a completed list as a spreadsheet.
func (c *Client) DownloadList(ctx context.Context, id int) ([]byte, error) {
	return c.doRaw(ctx, "GET", fmt.Sprintf("/lists/%d/download", id), nil)
}

// DeleteList deletes a list and its results.
func (c *Client) DeleteList(ctx context.Context, id int) error {
	return c.delete(ctx, fmt.Sprintf("/lists/%d", id), nil)
}

// PollList polls a list until processing completes or fails.
// The optional callback is invoked on each poll with the current list status.
func (c *Client) PollList(ctx context.Context, id int, callback func(*ListResponse)) (*ListResponse, error) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		list, err := c.GetList(ctx, id)
		if err != nil {
			return nil, err
		}

		if callback != nil {
			callback(list)
		}

		if list.Status != nil && (list.Status.State == "COMPLETED" || list.Status.State == "FAILED") {
			return list, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}

// DownloadListToWriter downloads list results directly to a writer.
func (c *Client) DownloadListToWriter(ctx context.Context, id int, w io.Writer) error {
	data, err := c.DownloadList(ctx, id)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
