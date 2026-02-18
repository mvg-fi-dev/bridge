package ids

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

// NewToken returns a URL/memo-friendly token.
// Uses crypto/rand. Output is base32 lowercase without padding.
func NewToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	return strings.ToLower(s), nil
}
