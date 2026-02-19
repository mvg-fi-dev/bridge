package route

import (
	"encoding/json"
	"fmt"
)

// BotKeystoreForSessions is the minimal bot keystore fields required to call /sessions/fetch.
// We reuse bot-api-go-client/v2's SignAuthenticationToken under the hood.
type BotKeystoreForSessions struct {
	UserID     string `json:"user_id"`
	SessionID  string `json:"session_id"`
	PrivateKey string `json:"private_key"`
}

func ParseBotKeystoreForSessions(b []byte) (*BotKeystoreForSessions, error) {
	var ks BotKeystoreForSessions
	if err := json.Unmarshal(b, &ks); err != nil {
		return nil, err
	}
	if ks.UserID == "" || ks.SessionID == "" || ks.PrivateKey == "" {
		return nil, fmt.Errorf("keystore missing user_id/session_id/private_key")
	}
	return &ks, nil
}
