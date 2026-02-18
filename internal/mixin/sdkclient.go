package mixin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	bot "github.com/MixinNetwork/bot-api-go-client/v2"
)

// SafeKeystore is a minimal subset of the bot keystore JSON from Mixin dashboard.
// Field names may vary; adjust if your keystore uses different keys.
//
// NOTE: For safe withdrawals, spend_private_key is required.
type SafeKeystore struct {
	UserID          string `json:"user_id"`
	SessionID       string `json:"session_id"`
	PrivateKey      string `json:"private_key"` // session private key (base64 raw url)
	ServerPublicKey string `json:"server_public_key"`
	SpendPrivateKey string `json:"spend_private_key"`
}

func (ks *SafeKeystore) Validate() error {
	if ks.UserID == "" || ks.SessionID == "" || ks.PrivateKey == "" {
		return fmt.Errorf("keystore missing user_id/session_id/private_key")
	}
	return nil
}

func ParseSafeKeystore(b []byte) (*SafeKeystore, error) {
	var ks SafeKeystore
	if err := json.Unmarshal(b, &ks); err != nil {
		return nil, err
	}
	if err := ks.Validate(); err != nil {
		return nil, err
	}
	return &ks, nil
}

func (ks *SafeKeystore) ToSafeUser() *bot.SafeUser {
	return &bot.SafeUser{
		UserId:            ks.UserID,
		SessionId:         ks.SessionID,
		SessionPrivateKey: ks.PrivateKey,
		ServerPublicKey:   ks.ServerPublicKey,
		SpendPrivateKey:   ks.SpendPrivateKey,
	}
}

type SDKClient struct {
	Keystore *SafeKeystore
}

func NewSDKClient(ks *SafeKeystore) *SDKClient {
	return &SDKClient{Keystore: ks}
}

// ListSafeSnapshots polls /safe/snapshots via the official Go SDK.
func (c *SDKClient) ListSafeSnapshots(ctx context.Context, limit int, offset string) ([]*bot.SafeSnapshot, error) {
	ks := c.Keystore
	if ks == nil {
		return nil, fmt.Errorf("missing keystore")
	}
	return bot.SafeSnapshots(ctx, limit, "", "", "", offset, ks.UserID, ks.SessionID, ks.PrivateKey)
}

// Withdraw uses safe withdrawal (no PIN). Tag is chain-specific memo/tag (can be empty).
func (c *SDKClient) Withdraw(ctx context.Context, assetID, destination, tag, amount, traceID string) (*bot.SequencerTransactionRequest, error) {
	ks := c.Keystore
	if ks == nil {
		return nil, fmt.Errorf("missing keystore")
	}
	u := ks.ToSafeUser()
	if u.SpendPrivateKey == "" {
		return nil, fmt.Errorf("keystore missing spend_private_key")
	}
	return bot.SendWithdrawal(ctx, assetID, destination, tag, amount, traceID, u)
}

func SafeSnapshotToInternal(s *bot.SafeSnapshot) *Snapshot {
	if s == nil {
		return nil
	}
	return &Snapshot{
		SnapshotID: s.SnapshotID,
		Type:       s.Type,
		AssetID:    s.AssetID,
		Amount:     s.Amount,
		CreatedAt:  s.CreatedAt.UTC().Format(time.RFC3339Nano),
		Memo:       s.Memo,
		OpponentID: s.OpponentID,
	}
}
