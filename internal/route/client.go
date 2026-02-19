package route

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client

	// Auth inputs
	AccountID    string
	Mnemonic     string
	RouteBotPKB64 string
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://api.route.mixin.one"
	}
	return &Client{BaseURL: baseURL, HTTP: &http.Client{Timeout: 15 * time.Second}}
}

type QuoteResult struct {
	InputMint  string `json:"inputMint"`
	InAmount   string `json:"inAmount"`
	OutputMint string `json:"outputMint"`
	OutAmount  string `json:"outAmount"`
	Slippage   int    `json:"slippage"`
	Source     string `json:"source"`
	Payload    string `json:"payload"`
}

type SwapRequest struct {
	Payer                 string  `json:"payer"`
	InputMint             string  `json:"inputMint"`
	InputAmount           string  `json:"inputAmount"`
	OutputMint            string  `json:"outputMint"`
	Payload               string  `json:"payload"`
	Source                string  `json:"source"`
	WithdrawalDestination *string `json:"withdrawalDestination"`
	Referral              *string `json:"referral"`
	WalletId              *string `json:"walletId"`
}

type SwapResponse struct {
	Tx                 *string     `json:"tx"`
	Source             string      `json:"source"`
	DisplayUserId      *string     `json:"displayUserId"`
	DepositDestination *string     `json:"depositDestination"`
	Quote              QuoteResult `json:"quote"`
}

func (c *Client) Quote(ctx context.Context, inputMint, outputMint, amount, source string) (*QuoteResult, error) {
	u, _ := url.Parse(c.BaseURL)
	u.Path = "/web3/quote"
	q := u.Query()
	q.Set("inputMint", inputMint)
	q.Set("outputMint", outputMint)
	q.Set("amount", amount) // amount format to be confirmed (likely base units)
	if source != "" {
		q.Set("source", source)
	}
	u.RawQuery = q.Encode()

	pathForSig := u.Path + "?" + u.RawQuery
	ts := time.Now().UTC().Unix()
	sign, err := ComputeMRAccessSign(c.AccountID, c.Mnemonic, c.RouteBotPKB64, ts, http.MethodGet, pathForSig, "")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("MR-ACCESS-TIMESTAMP", strconv.FormatInt(ts, 10))
	req.Header.Set("MR-ACCESS-SIGN", sign)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("route quote status=%d body=%s", resp.StatusCode, string(b))
	}

	var out struct {
		Data  QuoteResult `json:"data"`
		Error any        `json:"error"`
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) Swap(ctx context.Context, reqBody SwapRequest) (*SwapResponse, error) {
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(c.BaseURL)
	u.Path = "/web3/swap"
	pathForSig := u.Path
	ts := time.Now().UTC().Unix()
	sign, err := ComputeMRAccessSign(c.AccountID, c.Mnemonic, c.RouteBotPKB64, ts, http.MethodPost, pathForSig, string(b))
	if err != nil {
		return nil, err
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	hreq.Header.Set("Content-Type", "application/json")
	hreq.Header.Set("MR-ACCESS-TIMESTAMP", strconv.FormatInt(ts, 10))
	hreq.Header.Set("MR-ACCESS-SIGN", sign)

	resp, err := c.HTTP.Do(hreq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	outBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("route swap status=%d body=%s", resp.StatusCode, string(outBytes))
	}

	var out struct {
		Data  SwapResponse `json:"data"`
		Error any         `json:"error"`
	}
	if err := json.Unmarshal(outBytes, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}
