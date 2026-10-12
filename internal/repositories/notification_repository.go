package repositories

import (
	"context"
	"errors"
	"time"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, tx *gorm.DB, notification *models.Notification) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepository) Update(ctx context.Context, tx *gorm.DB, notification *models.Notification) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Save(notification).Error
}

func (r *notificationRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	var n models.Notification
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Logs").
		Where("id = ?", id).
		First(&n).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID uuid.UUID, req dto.PaginationRequest) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *notificationRepository) FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Notification{}).Preload("User")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Notification{}, id).Error
}

type notificationTemplateRepository struct {
	db *gorm.DB
}

func NewNotificationTemplateRepository(db *gorm.DB) NotificationTemplateRepository {
	return &notificationTemplateRepository{db: db}
}

func (r *notificationTemplateRepository) FindAll(ctx context.Context) ([]models.NotificationTemplate, error) {
	var templates []models.NotificationTemplate
	err := r.db.WithContext(ctx).Find(&templates).Error
	return templates, err
}

func (r *notificationTemplateRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.NotificationTemplate, error) {
	var tmpl models.NotificationTemplate
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tmpl).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tmpl, nil
}

func (r *notificationTemplateRepository) FindByEventTypeAndChannel(ctx context.Context, eventType, channel string) (*models.NotificationTemplate, error) {
	var tmpl models.NotificationTemplate
	err := r.db.WithContext(ctx).
		Where("event_type = ? AND channel = ? AND is_active = true", eventType, channel).
		First(&tmpl).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tmpl, nil
}

func (r *notificationTemplateRepository) Create(ctx context.Context, tmpl *models.NotificationTemplate) error {
	return r.db.WithContext(ctx).Create(tmpl).Error
}

func (r *notificationTemplateRepository) Update(ctx context.Context, tmpl *models.NotificationTemplate) error {
	return r.db.WithContext(ctx).Save(tmpl).Error
}

func (r *notificationTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.NotificationTemplate{}, id).Error
}

type notificationLogRepository struct {
	db *gorm.DB
}

func NewNotificationLogRepository(db *gorm.DB) NotificationLogRepository {
	return &notificationLogRepository{db: db}
}

func (r *notificationLogRepository) Create(ctx context.Context, tx *gorm.DB, log *models.NotificationLog) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(log).Error
}

type eventLogRepository struct {
	db *gorm.DB
}

func NewEventLogRepository(db *gorm.DB) EventLogRepository {
	return &eventLogRepository{db: db}
}

func (r *eventLogRepository) Create(ctx context.Context, tx *gorm.DB, event *models.EventLog) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(event).Error
}

func (r *eventLogRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.EventLog, error) {
	var event models.EventLog
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&event).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *eventLogRepository) FindAll(ctx context.Context, req dto.PaginationRequest) ([]models.EventLog, int64, error) {
	var events []models.EventLog
	var total int64

	db := r.db.WithContext(ctx).Model(&models.EventLog{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = utils.ApplyPagination(db, req)
	if err := db.Order("created_at DESC").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *notificationRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).Count(&count).Error
	return count, err
}
