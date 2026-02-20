package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mvg-fi-dev/bridge/internal/route"
)

func main() {
	acct := mustEnv("TEST_MIXIN_ACCOUNT_ID")
	mn := mustEnv("TEST_MIXIN_MNEMONIC")
	botPK := mustEnv("ROUTE_BOT_PUBLIC_KEY")

	inputMint := mustEnv("TEST_INPUT_MINT")
	outputMint := mustEnv("TEST_OUTPUT_MINT")
	source := getenv("TEST_SOURCE", "web3")
	amounts := strings.Fields(getenv("TEST_AMOUNTS", "10.0 100.0"))

	// Optional
	withdrawDest := os.Getenv("TEST_WITHDRAWAL_DEST")
	var wd *string
	if withdrawDest != "" {
		wd = &withdrawDest
	}

	ctx := context.Background()
	c := route.NewClient(getenv("ROUTE_BASE_URL", "https://api.route.mixin.one"))
	c.AccountID = acct
	c.Mnemonic = mn
	c.RouteBotPKB64 = botPK

	for _, amt := range amounts {
		fmt.Println("\n=== amount=", amt, "===")
		q, err := c.Quote(ctx, inputMint, outputMint, amt, source)
		if err != nil {
			fmt.Println("quote error:", err)
			continue
		}
		fmt.Printf("quote ok: in=%s out=%s slippage=%d payload_len=%d source=%s\n", q.InAmount, q.OutAmount, q.Slippage, len(q.Payload), q.Source)

		sreq := route.SwapRequest{
			Payer:                 acct,
			InputMint:             inputMint,
			InputAmount:           q.InAmount,
			OutputMint:            outputMint,
			Payload:               q.Payload,
			Source:                source,
			WithdrawalDestination: wd,
			Referral:              nil,
			WalletId:              nil,
		}
		resp, err := c.Swap(ctx, sreq)
		if err != nil {
			fmt.Println("swap error:", err)
			continue
		}
		fmt.Printf("swap ok: tx_present=%v depositDestination=%v displayUserId=%v out=%s\n",
			resp.Tx != nil && *resp.Tx != "", resp.DepositDestination, resp.DisplayUserId, resp.Quote.OutAmount)
		if resp.Tx != nil {
			fmt.Printf("tx_len=%d\n", len(*resp.Tx))
		}
	}
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
