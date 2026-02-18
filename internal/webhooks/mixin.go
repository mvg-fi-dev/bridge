package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvg-fi-dev/bridge/internal/mixin"
)

type MixinWebhookHandler struct {
	Secret string
}

// Verify is a simple HMAC-SHA256 verifier.
// Header name/signature format varies by provider; adjust once you confirm Mixin's webhook spec.
func (h *MixinWebhookHandler) verify(body []byte, providedHex string) bool {
	if h.Secret == "" || providedHex == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(h.Secret))
	mac.Write(body)
	expected := mac.Sum(nil)
	provided, err := hex.DecodeString(providedHex)
	if err != nil {
		return false
	}
	return hmac.Equal(expected, provided)
}

func (h *MixinWebhookHandler) Handle(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "read body"})
		return
	}

	// TODO: confirm header name with actual Mixin webhook docs.
	sig := c.GetHeader("X-Signature")
	if h.Secret != "" {
		if !h.verify(body, sig) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}
	}

	_, err = mixin.ParseSnapshot(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid snapshot"})
		return
	}

	// TODO: persist snapshot + dispatch to order matching.
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
