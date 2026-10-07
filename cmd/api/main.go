package main

import (
	"io"
	"os"
	"path/filepath"

	"backend-sembako/config"
	"backend-sembako/internal/api"
	"backend-sembako/pkg/logger"
)

func main() {
	// 1. Load Configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// 2. Initialize Structured Logger
	logger.InitLogger(cfg.App.Env)
	logger.Log.Info().Msgf("Starting %s v%s...", cfg.App.Name, cfg.App.Version)

	// 3. Initialize Database Connection (PostgreSQL)
	db, err := config.InitDatabase(&cfg.Database, cfg.App.Env)
	if err != nil {
		logger.Log.Warn().Err(err).Msg("Database connection failed on startup. Server starting in degraded mode.")
	} else if db != nil {
		// Auto-migrate all 38 GORM entities
		if err := config.AutoMigrate(db); err != nil {
			logger.Log.Error().Err(err).Msg("Failed to auto-migrate database tables")
		} else {
			logger.Log.Info().Msg("PostgreSQL database tables auto-migrated successfully.")
			_ = config.SeedDatabase(db)
		}
	}

	// Sync Logo Binary Files to Public Directories
	syncLogoFiles()

	// 4. Initialize & Bootstrap HTTP Server Engine
	server := api.NewServer(cfg, db)
	if err := server.Run(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to run HTTP server")
	}
}

func syncLogoFiles() {
	copyFile := func(src, dst string) {
		in, err := os.Open(src)
		if err != nil {
			return
		}
		defer in.Close()

		_ = os.MkdirAll(filepath.Dir(dst), 0755)
		out, err := os.Create(dst)
		if err != nil {
			return
		}
		defer out.Close()

		_, _ = io.Copy(out, in)
	}

	copyFile("logo.png", "customer-web/public/logo.png")
	copyFile("logo.png", "admin-dashboard/public/logo.png")
}
