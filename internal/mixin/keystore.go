package mixin

import (
	"encoding/json"
	"fmt"
	"os"
)

// Keystore is a minimal subset of the keystore-xxxx.json downloaded from Mixin dashboard.
// Fields may vary by version; we only need uid/sid/private_key for JWT signing.
type Keystore struct {
	UserID     string `json:"user_id"`
	SessionID  string `json:"session_id"`
	PrivateKey string `json:"private_key"`
}

func LoadKeystore(path string) (*Keystore, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var ks Keystore
	if err := json.Unmarshal(b, &ks); err != nil {
		return nil, err
	}
	if ks.UserID == "" || ks.SessionID == "" || ks.PrivateKey == "" {
		return nil, fmt.Errorf("keystore missing user_id/session_id/private_key")
	}
	return &ks, nil
}
