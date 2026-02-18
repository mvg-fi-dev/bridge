package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvg-fi-dev/bridge/internal/webhooks"
)

type Server struct {
	DB *sql.DB

	PayWindowSeconds int64
	MixinBotUserID   string
	MixinWebhookSecret string
}

func (s *Server) Register(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// Orders (Mixin-first MVP)
	r.POST("/v1/orders", s.handleCreateOrder)
	r.GET("/v1/orders/:public_id", s.handleGetOrder)

	// Optional webhook ingestion (can be replaced by polling or blaze).
	mw := &webhooks.MixinWebhookHandler{Secret: s.MixinWebhookSecret, DB: s.DB, MixinBotUserID: s.MixinBotUserID}
	r.POST("/v1/webhooks/mixin", mw.Handle)
}
