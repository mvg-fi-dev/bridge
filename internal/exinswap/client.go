package exinswap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: "https://app.exinswap.com/api/v2",
		HTTP:    &http.Client{Timeout: 15 * time.Second},
	}
}

type APIResponse[T any] struct {
	Code        int    `json:"code"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	Data        T      `json:"data"`
	TimestampMs int64  `json:"timestampMs"`
}

type Asset struct {
	UUID     string `json:"uuid"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	IconURL  string `json:"iconUrl"`
	RouteID  string `json:"routeId"`
	PriceUSDT string `json:"priceUsdt"`
	ChainAsset *struct {
		UUID    string `json:"uuid"`
		Symbol  string `json:"symbol"`
		Name    string `json:"name"`
		IconURL string `json:"iconUrl"`
		RouteID string `json:"routeId"`
	} `json:"chainAsset"`
}

type Pair struct {
	Asset0UUID string `json:"asset0Uuid"`
	Asset1UUID string `json:"asset1Uuid"`
	LPAssetUUID string `json:"lpAssetUuid"`
	Asset0Balance string `json:"asset0Balance"`
	Asset1Balance string `json:"asset1Balance"`
	LPAssetSupply string `json:"lpAssetSupply"`
	TradeType string `json:"tradeType"`
	CurveAmplifier string `json:"curveAmplifier"`
	CreatedAt int64 `json:"createdAt"`
	UpdatedAt int64 `json:"updatedAt"`
}

func (c *Client) GetAssets(ctx context.Context) ([]Asset, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/assets", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("exinswap assets status=%d body=%s", resp.StatusCode, string(b))
	}
	var out APIResponse[[]Asset]
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	if !out.Success || out.Code != 0 {
		return nil, fmt.Errorf("exinswap assets code=%d msg=%s", out.Code, out.Message)
	}
	return out.Data, nil
}

func (c *Client) GetPairs(ctx context.Context) ([]Pair, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/pairs", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("exinswap pairs status=%d body=%s", resp.StatusCode, string(b))
	}
	var out APIResponse[[]Pair]
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	if !out.Success || out.Code != 0 {
		return nil, fmt.Errorf("exinswap pairs code=%d msg=%s", out.Code, out.Message)
	}
	return out.Data, nil
}
