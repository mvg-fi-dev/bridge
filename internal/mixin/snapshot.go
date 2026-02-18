package mixin

import (
	"encoding/json"
	"fmt"
	"time"
)

type Snapshot struct {
	SnapshotID string `json:"snapshot_id"`
	Type       string `json:"type"`
	AssetID    string `json:"asset_id"`
	Amount     string `json:"amount"`
	CreatedAt  string `json:"created_at"`
	Memo       string `json:"memo"`
	OpponentID string `json:"opponent_id"`
}

type SnapshotEnvelope struct {
	Data Snapshot `json:"data"`
}

func ParseSnapshot(body []byte) (*Snapshot, error) {
	var env SnapshotEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot envelope: %w", err)
	}
	return &env.Data, nil
}

func (s *Snapshot) CreatedAtTime() (*time.Time, error) {
	if s.CreatedAt == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
