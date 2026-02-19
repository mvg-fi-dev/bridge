package route

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/curve25519"

	"github.com/tyler-smith/go-bip39"
)

// DeriveEd25519SeedFromMixinMnemonic derives a 32-byte Ed25519 seed from a Mixin Messenger mnemonic.
//
// Android client path (observed):
//   BIP39 seed (passphrase="")
//   -> BIP32 master private key (HMAC-SHA512 key="Bitcoin seed")
//   -> take master private key (32 bytes) as Ed25519 seed
//
// Mixin mnemonics may include a checksum word (13/25 words); Android drops the last word in those cases.
func DeriveEd25519SeedFromMixinMnemonic(mnemonic string) ([]byte, error) {
	mnemonic = strings.TrimSpace(mnemonic)
	words := strings.Fields(mnemonic)
	if len(words) == 13 {
		words = words[:12]
	} else if len(words) == 25 {
		words = words[:24]
	}
	mn := strings.Join(words, " ")
	if !bip39.IsMnemonicValid(mn) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	seed := bip39.NewSeed(mn, "")
	// BIP32 master key: I = HMAC-SHA512(key="Bitcoin seed", data=seed)
	mac := hmac.New(sha512.New, []byte("Bitcoin seed"))
	mac.Write(seed)
	I := mac.Sum(nil)
	IL := I[:32]

	// master private key = parse256(IL) mod n (must be 1..n-1)
	k := new(big.Int).SetBytes(IL)
	k.Mod(k, secp256k1N())
	if k.Sign() == 0 {
		return nil, errors.New("derived master key is zero")
	}
	out := make([]byte, 32)
	kb := k.Bytes()
	copy(out[32-len(kb):], kb)
	return out, nil
}

func secp256k1N() *big.Int {
	// 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141
	n, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	return n
}

// ed25519SeedToX25519Scalar converts Ed25519 seed -> X25519 scalar as in Android's privateKeyToCurve25519.
func ed25519SeedToX25519Scalar(seed32 []byte) ([]byte, error) {
	if len(seed32) != 32 {
		return nil, fmt.Errorf("seed must be 32 bytes")
	}
	priv := ed25519.NewKeyFromSeed(seed32)
	h := sha512.Sum512(priv[:32])
	s := h[:32]
	// clamp
	s[0] &= 248
	s[31] &= 127
	s[31] |= 64
	out := make([]byte, 32)
	copy(out, s)
	return out, nil
}

// ComputeMRAccessSign computes MR-ACCESS-SIGN for Route API.
//
// Based on Android implementation:
// - botPublicKey is base64url decoded and used directly as Curve25519 public key (32 bytes).
// - sharedKey = X25519(userScalar, botCurve25519Pub)
// - content = ts + method + path + body
// - signature = base64url( accountId_bytes || HMAC_SHA256(sharedKey, content) )
func ComputeMRAccessSign(accountID string, userMnemonic string, routeBotCurve25519PublicKeyBase64URL string, tsSeconds int64, method, path, body string) (string, error) {
	seed, err := DeriveEd25519SeedFromMixinMnemonic(userMnemonic)
	if err != nil {
		return "", err
	}
	userScalar, err := ed25519SeedToX25519Scalar(seed)
	if err != nil {
		return "", err
	}

	botPubCurve, err := base64.RawURLEncoding.DecodeString(routeBotCurve25519PublicKeyBase64URL)
	if err != nil {
		return "", fmt.Errorf("decode bot public key: %w", err)
	}
	if len(botPubCurve) != 32 {
		return "", fmt.Errorf("bot public key must be 32 bytes, got %d", len(botPubCurve))
	}

	shared, err := curve25519.X25519(userScalar, botPubCurve)
	if err != nil {
		return "", err
	}

	content := fmt.Sprintf("%d%s%s", tsSeconds, method, path)
	if body != "" {
		content += body
	}

	mac := hmac.New(sha256.New, shared)
	mac.Write([]byte(content))
	h := mac.Sum(nil)

	payload := append([]byte(accountID), h...)
	return base64.RawURLEncoding.EncodeToString(payload), nil
}
