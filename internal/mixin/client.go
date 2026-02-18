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
	// Auth via bot keystore (preferred)
	UID        string
	SID        string
	PrivateKey string
	Scope      string

	// Legacy/simple token mode (optional)
	Token string

	HTTP *http.Client
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://api.mixin.one"
	}
	return &Client{
		BaseURL: baseURL,
		Scope:   "FULL",
		HTTP:    &http.Client{Timeout: 15 * time.Second},
	}
}

type ListSnapshotsResponse struct {
	Data []Snapshot `json:"data"`
}

// ListSnapshots calls Mixin snapshots API.
// Auth: if UID/SID/PrivateKey are set, signs an EdDSA JWT per Mixin docs.
// Otherwise, if Token is set, uses Bearer Token.
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

	uriForSig := u.Path
	if u.RawQuery != "" {
		uriForSig = uriForSig + "?" + u.RawQuery
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Auth header
	if c.UID != "" && c.SID != "" && c.PrivateKey != "" {
		tok, err := SignAuthenticationToken(c.UID, c.SID, c.PrivateKey, http.MethodGet, uriForSig, "", c.Scope)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+tok)
	} else if c.Token != "" {
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
