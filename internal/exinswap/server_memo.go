package exinswap

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// ServerMemo represents ExinSwap server transfer memo.
// Spec (V2 doc): BASE64("RESULT|TRACE|SOURCE|TYPE")
// RESULT: 0 success / 1 fail
// TRACE: user's payment trace id (we use order.ID)
// SOURCE: e.g. SW
// TYPE: for swap: RL (release) or RF (refund)
// Other combos exist for liquidity actions.
type ServerMemo struct {
	Result string
	Trace  string
	Source string
	Type   string
}

func ParseServerMemo(b64 string) (*ServerMemo, error) {
	if b64 == "" {
		return nil, fmt.Errorf("empty memo")
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		// Some clients may use raw url encoding; try that too.
		if raw2, err2 := base64.RawURLEncoding.DecodeString(b64); err2 == nil {
			raw = raw2
		} else {
			return nil, err
		}
	}
	parts := strings.Split(string(raw), "|")
	if len(parts) != 4 {
		return nil, fmt.Errorf("bad memo parts=%d raw=%q", len(parts), string(raw))
	}
	m := &ServerMemo{Result: parts[0], Trace: parts[1], Source: parts[2], Type: parts[3]}
	if m.Trace == "" {
		return nil, fmt.Errorf("missing trace")
	}
	return m, nil
}
