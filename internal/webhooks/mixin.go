package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mvg-fi-dev/bridge/internal/db"
	"github.com/mvg-fi-dev/bridge/internal/mixin"
)

type MixinWebhookHandler struct {
	Secret string
	DB     *sql.DB
	MixinBotUserID string
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

	snap, err := mixin.ParseSnapshot(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid snapshot"})
		return
	}

	// Minimal matching: memo maps to order; opponent_id must be our bot; asset_id must match.
	// Snapshot itself is credited, so we can mark the order as deposit_credited.
	createdAt, err := snap.CreatedAtTime()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid created_at"})
		return
	}
	creditedAt := time.Now().UTC()
	if createdAt != nil {
		creditedAt = createdAt.UTC()
	}

	if h.DB != nil && snap.Memo != "" {
		repo := db.NewOrdersRepo(h.DB)
		_, _ = repo.SetDepositCreditedByMemo(c.Request.Context(), snap.Memo, snap.SnapshotID, creditedAt, snap.Amount, snap.AssetID, snap.OpponentID)
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
