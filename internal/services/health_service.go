package services

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type HealthService interface {
	CheckLiveness() bool
	CheckReadiness(ctx context.Context) (bool, string)
}

type healthService struct {
	db *gorm.DB
}

func NewHealthService(db *gorm.DB) HealthService {
	return &healthService{db: db}
}

func (s *healthService) CheckLiveness() bool {
	return true
}

func (s *healthService) CheckReadiness(ctx context.Context) (bool, string) {
	if s.db == nil {
		return false, "Database instance is null"
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		return false, fmt.Sprintf("Failed to get SQL DB: %v", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return false, fmt.Sprintf("Database ping failed: %v", err)
	}

	return true, "Database connected and healthy"
}
