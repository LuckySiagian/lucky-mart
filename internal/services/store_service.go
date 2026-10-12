package services

import (
	"context"
	"errors"
	"strings"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"

	"gorm.io/gorm"
)

const storeSettingID = 1

type StoreService interface {
	Get(ctx context.Context) (*models.StoreSetting, error)
	Update(ctx context.Context, req dto.UpdateStoreSettingDTO) (*models.StoreSetting, error)
}

type storeService struct {
	db *gorm.DB
}

func NewStoreService(db *gorm.DB) StoreService {
	return &storeService{db: db}
}

// Get returns the store details, creating the row with defaults on first use.
func (s *storeService) Get(ctx context.Context) (*models.StoreSetting, error) {
	setting := models.StoreSetting{ID: storeSettingID}
	err := s.db.WithContext(ctx).
		Attrs(models.StoreSetting{Name: "Lucky Mart", Tagline: "Toko sembako keluarga, grosir & eceran", ReceiptNote: "Terima kasih sudah berbelanja!"}).
		FirstOrCreate(&setting, models.StoreSetting{ID: storeSettingID}).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (s *storeService) Update(ctx context.Context, req dto.UpdateStoreSettingDTO) (*models.StoreSetting, error) {
	wa := strings.TrimSpace(req.WhatsApp)
	if wa != "" && !phonePattern.MatchString(wa) {
		return nil, errors.New("nomor WhatsApp tidak valid (contoh: 081234567890)")
	}
	maps := strings.TrimSpace(req.MapsURL)
	if maps != "" && !strings.HasPrefix(maps, "https://") {
		return nil, errors.New("link Google Maps harus diawali https://")
	}

	setting, err := s.Get(ctx)
	if err != nil {
		return nil, err
	}
	setting.Name = strings.TrimSpace(req.Name)
	setting.Tagline = strings.TrimSpace(req.Tagline)
	setting.Address = strings.TrimSpace(req.Address)
	setting.WhatsApp = wa
	setting.OpenHours = strings.TrimSpace(req.OpenHours)
	setting.MapsURL = maps
	setting.ReceiptNote = strings.TrimSpace(req.ReceiptNote)
	setting.Announcement = strings.TrimSpace(req.Announcement)
	if err := s.db.WithContext(ctx).Save(setting).Error; err != nil {
		return nil, err
	}
	return setting, nil
}
