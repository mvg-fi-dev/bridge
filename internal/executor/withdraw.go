package executor

import (
	"context"
	"fmt"
	"log"

	"github.com/mvg-fi-dev/bridge/internal/db"
	"github.com/mvg-fi-dev/bridge/internal/ids"
	"github.com/mvg-fi-dev/bridge/internal/mixin"
	"github.com/mvg-fi-dev/bridge/internal/models"
)

type WithdrawExecutor struct {
	Orders *db.OrdersRepo
	Mixin  *mixin.SDKClient
}

func NewWithdrawExecutor(orders *db.OrdersRepo, mixinClient *mixin.SDKClient) *WithdrawExecutor {
	return &WithdrawExecutor{Orders: orders, Mixin: mixinClient}
}

// ExecuteWithdrawing submits a safe withdrawal to the target chain address.
// MVP assumption (per user): single chain support, no tag.
func (e *WithdrawExecutor) ExecuteWithdrawing(ctx context.Context, o *models.Order) error {
	if o.Status != models.StatusWithdrawing {
		return nil
	}
	if o.FinalOut == nil || *o.FinalOut == "" {
		return fmt.Errorf("missing final_out")
	}
	if o.TargetAsset == "" {
		return fmt.Errorf("missing target_asset")
	}
	if o.TargetAddress == "" {
		return fmt.Errorf("missing target_address")
	}

	traceID := ids.DeterministicUUID(o.ID + ":withdraw")
	log.Printf("withdraw order=%s asset=%s amount=%s dest=%s", o.PublicID, o.TargetAsset, *o.FinalOut, o.TargetAddress)
	resp, err := e.Mixin.Withdraw(ctx, o.TargetAsset, o.TargetAddress, "", *o.FinalOut, traceID)
	if err != nil {
		return err
	}
	withdrawRef := resp.RequestID
	if withdrawRef == "" {
		withdrawRef = resp.SnapshotID
	}
	if withdrawRef == "" {
		withdrawRef = traceID
	}
	return e.Orders.MarkCompleted(ctx, o.ID, withdrawRef)
}
