package executor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mvg-fi-dev/bridge/internal/db"
	"github.com/mvg-fi-dev/bridge/internal/exinswap"
	"github.com/mvg-fi-dev/bridge/internal/mixin"
	"github.com/mvg-fi-dev/bridge/internal/models"
)

const (
	// ExinSwap bot (from ExinSwap V2 docs)
	ExinSwapBotUserID = "29f23576-4651-47ff-8c16-6c8a5d76985e"
)

type ExinSwapExecutor struct {
	Orders *db.OrdersRepo
	Mixin  *mixin.SDKClient

	// Policy knobs (env configurable later)
	SwapTimeoutSeconds int64
}

func NewExinSwapExecutor(orders *db.OrdersRepo, mixinClient *mixin.SDKClient) *ExinSwapExecutor {
	return &ExinSwapExecutor{
		Orders:             orders,
		Mixin:              mixinClient,
		SwapTimeoutSeconds: 120,
	}
}

// ExecuteDepositCredited tries to execute one order:
// - enforce min_out using ExinSwap memo (min_out + latest_exec_time); ExinSwap will refund on failure.
// - transfer SourceAsset to ExinSwap bot with memo
// - TODO: wait for result credit + withdraw (next iteration)
func (e *ExinSwapExecutor) ExecuteDepositCredited(ctx context.Context, o *models.Order) error {
	if o.Status != models.StatusDepositCredited {
		return nil
	}
	if o.AmountCredited == nil {
		return fmt.Errorf("missing amount_credited")
	}

	ok, err := e.Orders.TryMarkExecutingSwap(ctx, o.ID)
	if err != nil {
		return err
	}
	if !ok {
		return nil // lost race
	}

	latest := time.Now().UTC().Add(time.Duration(e.SwapTimeoutSeconds) * time.Second)
	memo, err := exinswap.TradeMemoV2(o.TargetAsset, o.MinOut, &latest, "")
	if err != nil {
		_ = e.Orders.MarkRefunding(ctx, o.ID, "memo_build_failed")
		return err
	}

	traceID := o.ID // idempotency
	log.Printf("exinswap execute order=%s transfer asset=%s amt=%s target=%s minOut=%s latest=%d", o.PublicID, o.SourceAsset, *o.AmountCredited, o.TargetAsset, o.MinOut, latest.Unix())
	_, err = e.Mixin.Transfer(ctx, o.SourceAsset, ExinSwapBotUserID, *o.AmountCredited, memo, traceID)
	if err != nil {
		// Put it back to deposit_credited for retry? For now mark manual.
		_ = e.Orders.MarkRefunding(ctx, o.ID, "transfer_failed")
		return err
	}

	// For now we stop here; next iteration will:
	// 1) watch snapshots for credit from ExinSwap bot in o.TargetAsset and correlate by traceID in memo response
	// 2) mark withdrawing + withdraw to o.TargetAddress
	return nil
}
