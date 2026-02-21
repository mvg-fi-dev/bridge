package models

import "time"

type OrderStatus string

const (
	StatusQuoteCreated       OrderStatus = "quote_created"
	StatusAwaitingDeposit    OrderStatus = "awaiting_deposit"
	StatusDepositDetected    OrderStatus = "deposit_tx_detected"
	StatusDepositPendingMixin OrderStatus = "deposit_pending_mixin"
	StatusDepositCredited    OrderStatus = "deposit_credited"
	StatusExecutingSwap      OrderStatus = "executing_swap"
	StatusWithdrawing        OrderStatus = "withdrawing"
	StatusCompleted          OrderStatus = "completed"
	StatusRefunding          OrderStatus = "refunding"
	StatusRefunded           OrderStatus = "refunded"
	StatusFailedManual       OrderStatus = "failed_manual_review"
)

type Order struct {
	ID        string
	PublicID  string
	Status    OrderStatus
	CreatedAt time.Time
	UpdatedAt time.Time

	// Requested swap
	SourceChain  string
	SourceAsset  string
	AmountIn     string
	TargetChain  string
	TargetAsset  string
	TargetAddress string

	// Quote
	EstimatedOut string
	MinOut       string
	QuoteExpiryAt *time.Time

	// Timing
	PayWindowSeconds int64

	// Mixin payment UX (for Mixin-first MVP)
	MixinOpponentID   string
	MixinAssetID      string
	MixinPayMemo      string
	MixinPayURL       string

	// Deposit tracking
	DepositTxID          *string
	DepositTxDetectedAt  *time.Time
	DepositCreditedAt    *time.Time
	AmountCredited       *string
	RefundToAddress      *string

	// Execution
	FinalOut        *string
	SwapRef         *string
	ExinSwapTraceID *string
	WithdrawTxID    *string
	RefundTxID      *string
	RefundAssetID   *string
	RefundAmount    *string
	RefundReceivedSnapshotID *string
}
