package executor

import (
	"context"
	"fmt"
	"log"

	bot "github.com/MixinNetwork/bot-api-go-client/v2"
	"github.com/mvg-fi-dev/bridge/internal/db"
	"github.com/mvg-fi-dev/bridge/internal/ids"
	"github.com/mvg-fi-dev/bridge/internal/mixin"
	"github.com/mvg-fi-dev/bridge/internal/models"
)

type RefundExecutor struct {
	Orders *db.OrdersRepo
	Mixin  *mixin.SDKClient
}

func NewRefundExecutor(orders *db.OrdersRepo, mixinClient *mixin.SDKClient) *RefundExecutor {
	return &RefundExecutor{Orders: orders, Mixin: mixinClient}
}

// ExecuteRefunding refunds funds to the original sender (refund_to_address).
// MVP: Mixin-internal refund transfer.
func (e *RefundExecutor) ExecuteRefunding(ctx context.Context, o *models.Order) error {
	if o.Status != models.StatusRefunding {
		return nil
	}
	if o.RefundToAddress == nil || *o.RefundToAddress == "" {
		return fmt.Errorf("missing refund_to_address")
	}

	assetID := o.SourceAsset
	amount := ""
	if o.RefundAssetID != nil && *o.RefundAssetID != "" {
		assetID = *o.RefundAssetID
	}
	if o.RefundAmount != nil && *o.RefundAmount != "" {
		amount = *o.RefundAmount
	} else if o.AmountCredited != nil {
		amount = *o.AmountCredited
	}
	if amount == "" {
		return fmt.Errorf("missing refund amount")
	}

	traceID := ids.DeterministicUUID(o.ID + ":refund")
	memo := "" // optional; could include reason
	log.Printf("refund order=%s asset=%s amount=%s to=%s", o.PublicID, assetID, amount, *o.RefundToAddress)

	resp, err := e.Mixin.Transfer(ctx, assetID, *o.RefundToAddress, amount, memo, traceID)
	if err != nil {
		return err
	}
	refundRef := resp.RequestID
	if refundRef == "" {
		refundRef = resp.SnapshotID
	}
	if refundRef == "" {
		refundRef = traceID
	}
	return e.Orders.MarkRefunded(ctx, o.ID, refundRef)
}

var _ = bot.SequencerTransactionRequest{} // keep import stable if needed
