package main

import (
	"context"
	"fmt"
	"log"
	"os"

	bot "github.com/MixinNetwork/bot-api-go-client/v2"
	"github.com/mvg-fi-dev/bridge/internal/route"
)

func main() {
	mn := os.Getenv("TEST_MIXIN_MNEMONIC")
	acct := os.Getenv("TEST_MIXIN_ACCOUNT_ID")
	if mn == "" || acct == "" {
		log.Fatal("set TEST_MIXIN_MNEMONIC and TEST_MIXIN_ACCOUNT_ID")
	}

	// Fetch route bot session public key via Mixin API
	ksPath := os.Getenv("MIXIN_KEYSTORE_PATH")
	if ksPath == "" {
		log.Fatal("set MIXIN_KEYSTORE_PATH to fetch route bot session")
	}
	ksBytes, err := os.ReadFile(ksPath)
	if err != nil {
		log.Fatal(err)
	}
	ks, err := route.ParseBotKeystoreForSessions(ksBytes)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	// route bot user id from android source constants
	routeBotUserID := "61cb8dd4-16b1-4744-ba0c-7b2d2e52fc59"
	sessions, err := bot.FetchUserSession(ctx, []string{routeBotUserID}, ks.UserID, ks.SessionID, ks.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	if len(sessions) == 0 {
		log.Fatal("no session")
	}
	botPK := sessions[0].PublicKey
	fmt.Println("route bot public key (base64url):", botPK)

	c := route.NewClient("https://api.route.mixin.one")
	c.AccountID = acct
	c.Mnemonic = mn
	c.RouteBotPKB64 = botPK

	// Example mints: TODO set real values
	inputMint := os.Getenv("TEST_INPUT_MINT")
	outputMint := os.Getenv("TEST_OUTPUT_MINT")
	amount := os.Getenv("TEST_AMOUNT")
	if inputMint == "" || outputMint == "" || amount == "" {
		log.Fatal("set TEST_INPUT_MINT, TEST_OUTPUT_MINT, TEST_AMOUNT")
	}

	q, err := c.Quote(ctx, inputMint, outputMint, amount, "web3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("quote outAmount=%s payload_len=%d\n", q.OutAmount, len(q.Payload))
}
