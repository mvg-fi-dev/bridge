package mixin

import (
	"context"
	"fmt"

	bot "github.com/MixinNetwork/bot-api-go-client/v2"
)

// Transfer sends an internal Mixin transfer to a user (opponent) with optional memo.
// This uses safe transaction signing (SpendPrivateKey required).
func (c *SDKClient) Transfer(ctx context.Context, assetID, opponentUserID, amount, memo, traceID string) (*bot.SequencerTransactionRequest, error) {
	ks := c.Keystore
	if ks == nil {
		return nil, fmt.Errorf("missing keystore")
	}
	u := ks.ToSafeUser()
	if u.SpendPrivateKey == "" {
		return nil, fmt.Errorf("keystore missing spend_private_key")
	}
	extra := []byte(nil)
	if memo != "" {
		extra = []byte(memo)
	}
	recipients := []*bot.TransactionRecipient{{
		MixAddress: bot.NewUUIDMixAddress([]string{opponentUserID}, 1).String(),
		Amount:     amount,
	}}
	return bot.SendTransaction(ctx, assetID, recipients, traceID, extra, nil, u)
}
