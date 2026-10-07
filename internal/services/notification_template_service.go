package services

import (
	"context"
	"errors"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
)

type NotificationTemplateService interface {
	GetTemplates(ctx context.Context) ([]dto.NotificationTemplateDTO, error)
	GetTemplateByID(ctx context.Context, id uuid.UUID) (*dto.NotificationTemplateDTO, error)
	CreateTemplate(ctx context.Context, req dto.CreateNotificationTemplateDTO) (*dto.NotificationTemplateDTO, error)
	UpdateTemplate(ctx context.Context, id uuid.UUID, req dto.CreateNotificationTemplateDTO) (*dto.NotificationTemplateDTO, error)
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
}

type notificationTemplateService struct {
	tmplRepo repositories.NotificationTemplateRepository
}

func NewNotificationTemplateService(tmplRepo repositories.NotificationTemplateRepository) NotificationTemplateService {
	return &notificationTemplateService{tmplRepo: tmplRepo}
}

func (s *notificationTemplateService) GetTemplates(ctx context.Context) ([]dto.NotificationTemplateDTO, error) {
	templates, err := s.tmplRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.NotificationTemplateDTO, len(templates))
	for i, t := range templates {
		dtos[i] = mapTemplateToDTO(&t)
	}
	return dtos, nil
}

func (s *notificationTemplateService) GetTemplateByID(ctx context.Context, id uuid.UUID) (*dto.NotificationTemplateDTO, error) {
	tmpl, err := s.tmplRepo.FindByID(ctx, id)
	if err != nil || tmpl == nil {
		return nil, errors.New("notification template not found")
	}
	res := mapTemplateToDTO(tmpl)
	return &res, nil
}

func (s *notificationTemplateService) CreateTemplate(ctx context.Context, req dto.CreateNotificationTemplateDTO) (*dto.NotificationTemplateDTO, error) {
	existing, err := s.tmplRepo.FindByEventTypeAndChannel(ctx, req.EventType, req.Channel)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("notification template for this EventType and Channel already exists")
	}

	tmpl := &models.NotificationTemplate{
		Name:      req.Name,
		EventType: req.EventType,
		Channel:   req.Channel,
		Subject:   req.Subject,
		Body:      req.Body,
		IsActive:  true,
	}

	if err := s.tmplRepo.Create(ctx, tmpl); err != nil {
		return nil, err
	}

	return s.GetTemplateByID(ctx, tmpl.ID)
}

func (s *notificationTemplateService) UpdateTemplate(ctx context.Context, id uuid.UUID, req dto.CreateNotificationTemplateDTO) (*dto.NotificationTemplateDTO, error) {
	tmpl, err := s.tmplRepo.FindByID(ctx, id)
	if err != nil || tmpl == nil {
		return nil, errors.New("notification template not found")
	}

	tmpl.Name = req.Name
	tmpl.EventType = req.EventType
	tmpl.Channel = req.Channel
	tmpl.Subject = req.Subject
	tmpl.Body = req.Body

	if err := s.tmplRepo.Update(ctx, tmpl); err != nil {
		return nil, err
	}

	return s.GetTemplateByID(ctx, id)
}

func (s *notificationTemplateService) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	return s.tmplRepo.Delete(ctx, id)
}

func mapTemplateToDTO(t *models.NotificationTemplate) dto.NotificationTemplateDTO {
	return dto.NotificationTemplateDTO{
		ID:        t.ID,
		Name:      t.Name,
		EventType: t.EventType,
		Channel:   t.Channel,
		Subject:   t.Subject,
		Body:      t.Body,
		IsActive:  t.IsActive,
	}
}
