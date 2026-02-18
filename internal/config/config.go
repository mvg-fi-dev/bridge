package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port string
	SQLitePath string

	PayWindowSeconds int64

	// Mixin-first MVP payment
	MixinBotUserID string
	MixinWebhookSecret string
}

func Load() (*Config, error) {
	c := &Config{}
	c.Port = getenv("PORT", "8080")
	c.SQLitePath = getenv("SQLITE_PATH", "bridge.db")

	c.MixinBotUserID = os.Getenv("MIXIN_BOT_USER_ID")
	c.MixinWebhookSecret = os.Getenv("MIXIN_WEBHOOK_SECRET")

	pws := getenv("PAY_WINDOW_SECONDS", "900")
	v, err := strconv.ParseInt(pws, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid PAY_WINDOW_SECONDS: %w", err)
	}
	c.PayWindowSeconds = v
	return c, nil
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
