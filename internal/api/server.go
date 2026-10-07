package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-sembako/config"
	"backend-sembako/internal/routes"
	"backend-sembako/pkg/logger"

	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	return &Server{
		config: cfg,
		db:     db,
	}
}

func (s *Server) Run() error {
	router := routes.SetupRouter(s.config, s.db)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.config.Server.Port),
		Handler:      router,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
	}

	// Channel to listen for errors from Server Goroutine
	serverErrors := make(chan error, 1)

	go func() {
		logger.Log.Info().
			Str("port", s.config.Server.Port).
			Str("env", s.config.App.Env).
			Msg("HTTP Server starting to listen")

		serverErrors <- srv.ListenAndServe()
	}()

	// Channel to listen for OS signals for Graceful Shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("error starting server: %w", err)
		}
	case sig := <-shutdown:
		logger.Log.Warn().
			Str("signal", sig.String()).
			Msg("Graceful shutdown signal received, shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Log.Error().Err(err).Msg("Server forced to shutdown")
			_ = srv.Close()
			return err
		}

		// Close Database Connections
		if s.db != nil {
			sqlDB, err := s.db.DB()
			if err == nil {
				if err := sqlDB.Close(); err != nil {
					logger.Log.Error().Err(err).Msg("Error closing database connection pool")
				} else {
					logger.Log.Info().Msg("PostgreSQL Connection pool closed successfully")
				}
			}
		}

		logger.Log.Info().Msg("Server exited gracefully")
	}

	return nil
}
