package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateNotificationDTO struct {
	UserID  uuid.UUID `json:"user_id" binding:"required"`
	Title   string    `json:"title" binding:"required,min=3,max=200"`
	Message string    `json:"message" binding:"required"`
	Type    string    `json:"type" binding:"required"`
	Channel string    `json:"channel" binding:"required,oneof=EMAIL WHATSAPP SMS PUSH_NOTIFICATION IN_APP"`
}

type NotificationTemplateDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	EventType string    `json:"event_type"`
	Channel   string    `json:"channel"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	IsActive  bool      `json:"is_active"`
}

type CreateNotificationTemplateDTO struct {
	Name      string `json:"name" binding:"required,min=3,max=100"`
	EventType string `json:"event_type" binding:"required"`
	Channel   string `json:"channel" binding:"required,oneof=EMAIL WHATSAPP SMS PUSH_NOTIFICATION IN_APP"`
	Subject   string `json:"subject" binding:"required,min=3,max=200"`
	Body      string `json:"body" binding:"required"`
}

type NotificationLogResponseDTO struct {
	ID           uuid.UUID  `json:"id"`
	Channel      string     `json:"channel"`
	Status       string     `json:"status"`
	ErrorMessage string     `json:"error_message,omitempty"`
	RetryCount   int        `json:"retry_count"`
	SentAt       *time.Time `json:"sent_at,omitempty"`
}

type NotificationResponseDTO struct {
	ID        uuid.UUID                    `json:"id"`
	UserID    uuid.UUID                    `json:"user_id"`
	Title     string                       `json:"title"`
	Message   string                       `json:"message"`
	Type      string                       `json:"type"`
	Channel   string                       `json:"channel"`
	Status    string                       `json:"status"`
	IsRead    bool                         `json:"is_read"`
	ReadAt    *time.Time                   `json:"read_at,omitempty"`
	CreatedAt time.Time                    `json:"created_at"`
	Logs      []NotificationLogResponseDTO `json:"logs,omitempty"`
}

type EventDTO struct {
	ID          uuid.UUID `json:"id"`
	EventName   string    `json:"event_name"`
	Aggregate   string    `json:"aggregate"`
	AggregateID uuid.UUID `json:"aggregate_id"`
	Payload     string    `json:"payload"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type RetryNotificationDTO struct {
	NotificationID uuid.UUID `json:"notification_id" binding:"required"`
}
