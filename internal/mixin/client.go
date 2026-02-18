package mixin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

func NewClient(baseURL, token string) *Client {
	if baseURL == "" {
		baseURL = "https://api.mixin.one"
	}
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTP:    &http.Client{Timeout: 15 * time.Second},
	}
}

type ListSnapshotsResponse struct {
	Data []Snapshot `json:"data"`
}

// ListSnapshots calls Mixin snapshots API.
// NOTE: Mixin auth is typically JWT-based; this uses a Bearer token placeholder.
// Zed needs to confirm the correct auth header format.
func (c *Client) ListSnapshots(ctx context.Context, limit int, offset string, assetID string) ([]Snapshot, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = "/snapshots"
	q := u.Query()
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset != "" {
		q.Set("offset", offset)
	}
	if assetID != "" {
		q.Set("asset", assetID)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("mixin snapshots status=%d body=%s", resp.StatusCode, string(b))
	}

	var out ListSnapshotsResponse
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
