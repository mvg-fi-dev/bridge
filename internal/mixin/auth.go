package mixin

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// SignAuthenticationToken signs an EdDSA JWT for Mixin API calls.
//
// Per Mixin docs, sig = sha256(method + uri + body), where uri is the path without hostname.
// It's not 100% clear if query string must be included; we include it as part of uri.
func SignAuthenticationToken(uid, sid, privateKeyBase64URL, method, uri, body string, scp string) (string, error) {
	expire := time.Now().UTC().Add(90 * 24 * time.Hour)
	sum := sha256.Sum256([]byte(method + uri + body))

	claims := jwt.MapClaims{
		"uid": uid,
		"sid": sid,
		"iat": time.Now().UTC().Unix(),
		"exp": expire.Unix(),
		"jti": uuid.NewString(),
		"sig": hex.EncodeToString(sum[:]),
		"scp": scp,
	}

	privRaw, err := base64.RawURLEncoding.DecodeString(privateKeyBase64URL)
	if err != nil {
		return "", err
	}
	if len(privRaw) != 64 {
		return "", fmt.Errorf("bad ed25519 private key length=%d", len(privRaw))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(ed25519.PrivateKey(privRaw))
}
