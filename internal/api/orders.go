package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mvg-fi-dev/bridge/internal/db"
	"github.com/mvg-fi-dev/bridge/internal/ids"
	"github.com/mvg-fi-dev/bridge/internal/models"
)

type CreateOrderRequest struct {
	// For now: Mixin-first MVP. These map to Mixin fields directly.
	// Later we can add chain/asset mapping.
	MixinAssetID string `json:"mixin_asset_id" binding:"required"`
	AmountIn     string `json:"amount_in" binding:"required"`
	
	TargetChain   string `json:"target_chain" binding:"required"`
	TargetAsset   string `json:"target_asset" binding:"required"`
	TargetAddress string `json:"target_address" binding:"required"`

	// Pricing inputs (MVP: caller provides; later: computed from pricing module)
	EstimatedOut string `json:"estimated_out" binding:"required"`
	MinOut       string `json:"min_out" binding:"required"`
}

type CreateOrderResponse struct {
	PublicID string            `json:"public_id"`
	Status   models.OrderStatus `json:"status"`
	PayWindowSeconds int64     `json:"pay_window_seconds"`
	MixinPayment struct {
		OpponentID string `json:"opponent_id"`
		AssetID    string `json:"asset_id"`
		Amount     string `json:"amount"`
		Memo       string `json:"memo"`
	} `json:"mixin_payment"`
	Quote struct {
		EstimatedOut string `json:"estimated_out"`
		MinOut       string `json:"min_out"`
	} `json:"quote"`
	Terms map[string]string `json:"terms"`
}

func (s *Server) handleCreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memo, err := ids.NewToken(10) // ~16 chars base32
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token"})
		return
	}

	now := time.Now().UTC()
	o := &models.Order{
		ID:        ids.NewUUID(),
		PublicID:  ids.NewPublicID("BRG"),
		Status:    models.StatusAwaitingDeposit,
		CreatedAt: now,
		UpdatedAt: now,

		SourceChain:  "MIXIN",
		SourceAsset:  req.MixinAssetID,
		AmountIn:     req.AmountIn,
		TargetChain:  req.TargetChain,
		TargetAsset:  req.TargetAsset,
		TargetAddress: req.TargetAddress,

		EstimatedOut: req.EstimatedOut,
		MinOut:       req.MinOut,
		PayWindowSeconds: s.PayWindowSeconds,

		MixinOpponentID: s.MixinBotUserID,
		MixinAssetID:    req.MixinAssetID,
		MixinPayMemo:    memo,
	}

	repo := db.NewOrdersRepo(s.DB)
	if err := repo.Insert(c.Request.Context(), o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}

	var resp CreateOrderResponse
	resp.PublicID = o.PublicID
	resp.Status = o.Status
	resp.PayWindowSeconds = o.PayWindowSeconds
	resp.MixinPayment.OpponentID = o.MixinOpponentID
	resp.MixinPayment.AssetID = o.MixinAssetID
	resp.MixinPayment.Amount = o.AmountIn
	resp.MixinPayment.Memo = o.MixinPayMemo
	resp.Quote.EstimatedOut = o.EstimatedOut
	resp.Quote.MinOut = o.MinOut
	resp.Terms = map[string]string{
		"late_deposit": "auto_refund",
		"below_min_out": "auto_refund",
		"refund_fee": "paid_by_user",
		"refund_to": "original_address",
		"paid_definition": "mixin_credited",
		"late_cutoff": "first_detected_time",
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) handleGetOrder(c *gin.Context) {
	pid := c.Param("public_id")
	repo := db.NewOrdersRepo(s.DB)
	o, err := repo.GetByPublicID(c.Request.Context(), pid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, o)
}
