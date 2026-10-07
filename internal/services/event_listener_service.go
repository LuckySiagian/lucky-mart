package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"
	"backend-sembako/pkg/eventbus"
	"backend-sembako/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventListenerService interface {
	RegisterSubscribers(bus eventbus.EventBus)
	GetEvents(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (*dto.EventDTO, error)
}

type eventListenerService struct {
	db               *gorm.DB
	eventRepo        repositories.EventLogRepository
	tmplRepo         repositories.NotificationTemplateRepository
	notificationRepo repositories.NotificationRepository
	notifLogRepo     repositories.NotificationLogRepository
}

func NewEventListenerService(
	db *gorm.DB,
	eventRepo repositories.EventLogRepository,
	tmplRepo repositories.NotificationTemplateRepository,
	notificationRepo repositories.NotificationRepository,
	notifLogRepo repositories.NotificationLogRepository,
) EventListenerService {
	return &eventListenerService{
		db:               db,
		eventRepo:        eventRepo,
		tmplRepo:         tmplRepo,
		notificationRepo: notificationRepo,
		notifLogRepo:     notifLogRepo,
	}
}

func (s *eventListenerService) RegisterSubscribers(bus eventbus.EventBus) {
	events := []string{
		"USER_REGISTERED", "USER_LOGIN",
		"PRODUCT_CREATED", "LOW_STOCK",
		"ORDER_CREATED", "ORDER_CANCELLED", "ORDER_COMPLETED",
		"PAYMENT_PENDING", "PAYMENT_PAID", "PAYMENT_REFUNDED",
		"SHIPMENT_CREATED", "SHIPMENT_PICKED_UP", "SHIPMENT_DELIVERED",
		"PASSWORD_RESET", "EMAIL_VERIFICATION",
	}

	for _, evt := range events {
		bus.Subscribe(evt, s.handleDomainEvent)
	}
}

func (s *eventListenerService) handleDomainEvent(ctx context.Context, evt eventbus.DomainEvent) error {
	logger.Log.Info().
		Str("event_name", evt.Name).
		Str("aggregate", evt.Aggregate).
		Str("aggregate_id", evt.AggregateID.String()).
		Msg("Processing Domain Event Subscriber")

	// 1. Audit Log into EventLog Table
	eventLog := &models.EventLog{
		EventName:   evt.Name,
		Aggregate:   evt.Aggregate,
		AggregateID: evt.AggregateID,
		Payload:     evt.ToJSON(),
		Status:      "PROCESSED",
	}
	_ = s.eventRepo.Create(ctx, nil, eventLog)

	// Extract Target User ID from payload
	userIDVal, exists := evt.Payload["user_id"]
	if !exists {
		return nil
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		return nil
	}
	targetUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil
	}

	channels := []string{"IN_APP", "EMAIL"}
	now := time.Now()

	for _, ch := range channels {
		// Find Notification Template for event + channel
		tmpl, err := s.tmplRepo.FindByEventTypeAndChannel(ctx, evt.Name, ch)
		title := fmt.Sprintf("Notifikasi %s", evt.Name)
		body := fmt.Sprintf("Event %s telah diproses untuk ID %s", evt.Name, evt.AggregateID)

		if err == nil && tmpl != nil {
			title = tmpl.Subject
			body = renderTemplateBody(tmpl.Body, evt.Payload)
		}

		notif := &models.Notification{
			UserID:  targetUserID,
			Title:   title,
			Message: body,
			Type:    evt.Name,
			Channel: ch,
			Status:  "SENT",
			IsRead:  false,
		}

		_ = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.notificationRepo.Create(ctx, tx, notif); err != nil {
				return err
			}
			return s.notifLogRepo.Create(ctx, tx, &models.NotificationLog{
				NotificationID: notif.ID,
				Channel:        ch,
				Status:         "SENT",
				RetryCount:     0,
				SentAt:         &now,
			})
		})
	}

	return nil
}

func (s *eventListenerService) GetEvents(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	events, total, err := s.eventRepo.FindAll(ctx, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	dtos := make([]dto.EventDTO, len(events))
	for i, e := range events {
		dtos[i] = dto.EventDTO{
			ID:          e.ID,
			EventName:   e.EventName,
			Aggregate:   e.Aggregate,
			AggregateID: e.AggregateID,
			Payload:     e.Payload,
			Status:      e.Status,
			CreatedAt:   e.CreatedAt,
		}
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(dtos, req.Page, req.Limit, total), nil
}

func (s *eventListenerService) GetEventByID(ctx context.Context, id uuid.UUID) (*dto.EventDTO, error) {
	e, err := s.eventRepo.FindByID(ctx, id)
	if err != nil || e == nil {
		return nil, fmt.Errorf("event log not found")
	}
	return &dto.EventDTO{
		ID:          e.ID,
		EventName:   e.EventName,
		Aggregate:   e.Aggregate,
		AggregateID: e.AggregateID,
		Payload:     e.Payload,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
	}, nil
}

func renderTemplateBody(body string, payload map[string]interface{}) string {
	res := body
	for k, v := range payload {
		placeholder := fmt.Sprintf("{{%s}}", k)
		valStr := fmt.Sprintf("%v", v)
		res = strings.ReplaceAll(res, placeholder, valStr)
	}
	return res
}
