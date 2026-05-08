package main

import (
	"os"

	"operra/api/internal/app"
	"operra/api/internal/platform/config"
	"operra/api/internal/platform/database"
	"operra/api/internal/platform/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv)

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Printf("database connection unavailable: %v", err)
	}

	server := app.New(cfg, db, log)

	addr := cfg.ListenAddr
	if err := server.Listen(addr); err != nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}
