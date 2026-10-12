package main

import (
	"fmt"
	"os"

	"backend-sembako/config"
	"backend-sembako/internal/api"
	"backend-sembako/pkg/logger"
)

func main() {
	// 1. Load Configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load configuration:", err)
		os.Exit(1)
	}

	// 2. Initialize Structured Logger
	logger.InitLogger(cfg.App.Env)
	logger.Log.Info().Msgf("Starting %s v%s...", cfg.App.Name, cfg.App.Version)

	// 3. Connect to PostgreSQL. The API cannot work without its database,
	// so fail loudly instead of starting in a half-working state.
	db, err := config.InitDatabase(&cfg.Database, cfg.App.Env)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Database connection failed")
	}
	if err := config.AutoMigrate(db); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to auto-migrate database tables")
	}
	logger.Log.Info().Msg("PostgreSQL database tables auto-migrated successfully.")
	if err := config.SeedDatabase(db); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to seed database")
	}

	// 4. Start HTTP Server
	server := api.NewServer(cfg, db)
	if err := server.Run(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to run HTTP server")
	}
}
