package services

import (
	"context"
	"errors"
	"time"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotificationNotFound     = errors.New("notification not found")
	ErrUnauthorizedNotification = errors.New("you do not have permission to access this notification")
)

type NotificationService interface {
	GetNotifications(ctx context.Context, userID uuid.UUID, userRole string, req dto.PaginationRequest) (dto.PaginatedResponse, error)
	GetNotificationByID(ctx context.Context, userID uuid.UUID, userRole string, id uuid.UUID) (*dto.NotificationResponseDTO, error)
	MarkAsRead(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.NotificationResponseDTO, error)
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	DeleteNotification(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	RetryNotification(ctx context.Context, notificationID uuid.UUID) (*dto.NotificationResponseDTO, error)
}

type notificationService struct {
	db               *gorm.DB
	notificationRepo repositories.NotificationRepository
	logRepo          repositories.NotificationLogRepository
}

func NewNotificationService(
	db *gorm.DB,
	notificationRepo repositories.NotificationRepository,
	logRepo repositories.NotificationLogRepository,
) NotificationService {
	return &notificationService{
		db:               db,
		notificationRepo: notificationRepo,
		logRepo:          logRepo,
	}
}

func (s *notificationService) GetNotifications(ctx context.Context, userID uuid.UUID, userRole string, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	// Everyone, staff included, sees their own inbox. Store alerts reach
	// staff as their own notifications (see EventListenerService).
	notifications, total, err := s.notificationRepo.FindByUserID(ctx, userID, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.NotificationResponseDTO, len(notifications))
	for i, n := range notifications {
		dtos[i] = mapNotificationToDTO(&n)
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *notificationService) GetNotificationByID(ctx context.Context, userID uuid.UUID, userRole string, id uuid.UUID) (*dto.NotificationResponseDTO, error) {
	n, err := s.notificationRepo.FindByID(ctx, id)
	if err != nil || n == nil {
		return nil, ErrNotificationNotFound
	}

	if userRole != "SUPER_ADMIN" && userRole != "STORE_MANAGER" {
		if n.UserID != userID {
			return nil, ErrUnauthorizedNotification
		}
	}

	res := mapNotificationToDTO(n)
	return &res, nil
}

func (s *notificationService) MarkAsRead(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.NotificationResponseDTO, error) {
	n, err := s.notificationRepo.FindByID(ctx, id)
	if err != nil || n == nil {
		return nil, ErrNotificationNotFound
	}

	if n.UserID != userID {
		return nil, ErrUnauthorizedNotification
	}

	now := time.Now()
	n.IsRead = true
	n.ReadAt = &now

	if err := s.notificationRepo.Update(ctx, nil, n); err != nil {
		return nil, err
	}

	res := mapNotificationToDTO(n)
	return &res, nil
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

func (s *notificationService) DeleteNotification(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	n, err := s.notificationRepo.FindByID(ctx, id)
	if err != nil || n == nil {
		return ErrNotificationNotFound
	}

	if n.UserID != userID {
		return ErrUnauthorizedNotification
	}

	return s.notificationRepo.Delete(ctx, id)
}

func (s *notificationService) RetryNotification(ctx context.Context, notificationID uuid.UUID) (*dto.NotificationResponseDTO, error) {
	n, err := s.notificationRepo.FindByID(ctx, notificationID)
	if err != nil || n == nil {
		return nil, ErrNotificationNotFound
	}

	now := time.Now()
	n.Status = "SENT"

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.notificationRepo.Update(ctx, tx, n); err != nil {
			return err
		}

		return s.logRepo.Create(ctx, tx, &models.NotificationLog{
			NotificationID: n.ID,
			Channel:        n.Channel,
			Status:         "SENT",
			RetryCount:     1,
			SentAt:         &now,
		})
	})

	if err != nil {
		return nil, err
	}

	return s.GetNotificationByID(ctx, n.UserID, "SUPER_ADMIN", notificationID)
}

func mapNotificationToDTO(n *models.Notification) dto.NotificationResponseDTO {
	logs := make([]dto.NotificationLogResponseDTO, len(n.Logs))
	for i, l := range n.Logs {
		logs[i] = dto.NotificationLogResponseDTO{
			ID:           l.ID,
			Channel:      l.Channel,
			Status:       l.Status,
			ErrorMessage: l.ErrorMessage,
			RetryCount:   l.RetryCount,
			SentAt:       l.SentAt,
		}
	}

	return dto.NotificationResponseDTO{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Message:   n.Message,
		Type:      n.Type,
		Channel:   n.Channel,
		Status:    n.Status,
		IsRead:    n.IsRead,
		ReadAt:    n.ReadAt,
		CreatedAt: n.CreatedAt,
		Logs:      logs,
	}
}

func (s *notificationService) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.notificationRepo.CountUnread(ctx, userID)
}
