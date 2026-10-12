package services

import (
	"context"
	"fmt"
	"regexp"
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

// Events that are recorded in the audit log and may create notifications.
var domainEvents = []string{
	"USER_REGISTERED", "B2B_REGISTERED", "B2B_VERIFIED", "PASSWORD_RESET_REQUESTED",
	"ORDER_CREATED", "ORDER_CANCELLED", "ORDER_STATUS_CHANGED", "ORDER_COMPLETED",
	"PAYMENT_PROOF_UPLOADED", "PAYMENT_PAID", "PAYMENT_REJECTED",
	"SHIPMENT_OUT_FOR_DELIVERY", "CREDIT_SETTLED",
}

func (s *eventListenerService) RegisterSubscribers(bus eventbus.EventBus) {
	for _, evt := range domainEvents {
		bus.Subscribe(evt, s.handleDomainEvent)
	}
}

// handleDomainEvent writes the event to the audit log, then creates in-app
// notifications: one for payload["user_id"] (the customer) and, when
// payload["notify_staff"] is true, one for every active admin/manager.
// The title and message come from a notification template for the event
// when one exists, otherwise from payload["title"] / payload["message"].
func (s *eventListenerService) handleDomainEvent(ctx context.Context, evt eventbus.DomainEvent) error {
	logger.Log.Info().
		Str("event_name", evt.Name).
		Str("aggregate", evt.Aggregate).
		Str("aggregate_id", evt.AggregateID.String()).
		Msg("Processing domain event")

	_ = s.eventRepo.Create(ctx, nil, &models.EventLog{
		EventName:   evt.Name,
		Aggregate:   evt.Aggregate,
		AggregateID: evt.AggregateID,
		Payload:     evt.ToJSON(),
		Status:      "PROCESSED",
	})

	if userIDStr, ok := evt.Payload["user_id"].(string); ok {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			title, body := payloadString(evt.Payload, "title"), payloadString(evt.Payload, "message")
			if tmpl, err := s.tmplRepo.FindByEventTypeAndChannel(ctx, evt.Name, "IN_APP"); err == nil && tmpl != nil && tmpl.IsActive {
				title = renderTemplateBody(tmpl.Subject, evt.Payload)
				body = renderTemplateBody(tmpl.Body, evt.Payload)
			}
			if title != "" {
				s.notify(ctx, userID, evt.Name, title, body)
			}
		}
	}

	if notifyStaff, _ := evt.Payload["notify_staff"].(bool); notifyStaff {
		title := payloadString(evt.Payload, "staff_title")
		body := payloadString(evt.Payload, "staff_message")
		if title != "" {
			var staff []models.User
			if err := s.db.WithContext(ctx).
				Where("role IN ? AND is_active = true", []string{"SUPER_ADMIN", "STORE_MANAGER"}).
				Find(&staff).Error; err == nil {
				for _, u := range staff {
					s.notify(ctx, u.ID, evt.Name, title, body)
				}
			}
		}
	}

	return nil
}

func (s *eventListenerService) notify(ctx context.Context, userID uuid.UUID, eventName, title, body string) {
	now := time.Now()
	notif := &models.Notification{
		UserID:  userID,
		Title:   title,
		Message: body,
		Type:    eventName,
		Channel: "IN_APP",
		Status:  "SENT",
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.notificationRepo.Create(ctx, tx, notif); err != nil {
			return err
		}
		return s.notifLogRepo.Create(ctx, tx, &models.NotificationLog{
			NotificationID: notif.ID,
			Channel:        "IN_APP",
			Status:         "SENT",
			SentAt:         &now,
		})
	})
	if err != nil {
		logger.Log.Error().Err(err).Str("event_name", eventName).Msg("Failed to create notification")
	}
}

func payloadString(payload map[string]interface{}, key string) string {
	if v, ok := payload[key].(string); ok {
		return v
	}
	return ""
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
	// Never show customers a raw {{placeholder}} that had no value.
	return unknownPlaceholder.ReplaceAllString(res, "")
}

var unknownPlaceholder = regexp.MustCompile(`\{\{[a-zA-Z0-9_]+\}\}`)
