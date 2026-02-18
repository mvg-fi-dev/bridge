package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/mvg-fi-dev/bridge/internal/api"
	"github.com/mvg-fi-dev/bridge/internal/config"
	"github.com/mvg-fi-dev/bridge/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.Open(cfg.SQLitePath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.SQL.Close()

	if err := db.Migrate(dbConn.SQL); err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	s := &api.Server{DB: dbConn.SQL, PayWindowSeconds: cfg.PayWindowSeconds, MixinBotUserID: cfg.MixinBotUserID, MixinWebhookSecret: cfg.MixinWebhookSecret}
	s.Register(r)

	log.Printf("bridge-api listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
