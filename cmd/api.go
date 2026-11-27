package main

import (
	"context"
	"log"
	"time"

	"bearuang/internal/config"
	"bearuang/internal/database"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db, err := database.NewPostgresDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer db.Close()
}
