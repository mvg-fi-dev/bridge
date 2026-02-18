package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvg-fi-dev/bridge/internal/webhooks"
)

type Server struct {
	MixinWebhookSecret string
}

func (s *Server) Register(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	mw := &webhooks.MixinWebhookHandler{Secret: s.MixinWebhookSecret}
	r.POST("/v1/webhooks/mixin", mw.Handle)
}
