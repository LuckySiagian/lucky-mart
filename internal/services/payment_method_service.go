package services

import (
	"context"
	"errors"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
)

type PaymentMethodService interface {
	GetPaymentMethods(ctx context.Context) ([]dto.PaymentMethodDTO, error)
	CreatePaymentMethod(ctx context.Context, req dto.CreatePaymentMethodDTO) (*dto.PaymentMethodDTO, error)
	UpdatePaymentMethod(ctx context.Context, id uuid.UUID, req dto.CreatePaymentMethodDTO) (*dto.PaymentMethodDTO, error)
	DeletePaymentMethod(ctx context.Context, id uuid.UUID) error
}

type paymentMethodService struct {
	methodRepo repositories.PaymentMethodRepository
}

func NewPaymentMethodService(methodRepo repositories.PaymentMethodRepository) PaymentMethodService {
	return &paymentMethodService{methodRepo: methodRepo}
}

func (s *paymentMethodService) GetPaymentMethods(ctx context.Context) ([]dto.PaymentMethodDTO, error) {
	methods, err := s.methodRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.PaymentMethodDTO, len(methods))
	for i, m := range methods {
		dtos[i] = dto.PaymentMethodDTO{
			ID:          m.ID,
			Code:        m.Code,
			Name:        m.Name,
			Type:        m.Type,
			IsActive:    m.IsActive,
			Description: m.Description,
		}
	}
	return dtos, nil
}

func (s *paymentMethodService) CreatePaymentMethod(ctx context.Context, req dto.CreatePaymentMethodDTO) (*dto.PaymentMethodDTO, error) {
	existing, err := s.methodRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("payment method code already exists")
	}

	method := &models.PaymentMethod{
		Code:        req.Code,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.methodRepo.Create(ctx, method); err != nil {
		return nil, err
	}

	return &dto.PaymentMethodDTO{
		ID:          method.ID,
		Code:        method.Code,
		Name:        method.Name,
		Type:        method.Type,
		IsActive:    method.IsActive,
		Description: method.Description,
	}, nil
}

func (s *paymentMethodService) UpdatePaymentMethod(ctx context.Context, id uuid.UUID, req dto.CreatePaymentMethodDTO) (*dto.PaymentMethodDTO, error) {
	method, err := s.methodRepo.FindByID(ctx, id)
	if err != nil || method == nil {
		return nil, errors.New("payment method not found")
	}

	method.Name = req.Name
	method.Type = req.Type
	method.Description = req.Description

	if err := s.methodRepo.Update(ctx, method); err != nil {
		return nil, err
	}

	return &dto.PaymentMethodDTO{
		ID:          method.ID,
		Code:        method.Code,
		Name:        method.Name,
		Type:        method.Type,
		IsActive:    method.IsActive,
		Description: method.Description,
	}, nil
}

func (s *paymentMethodService) DeletePaymentMethod(ctx context.Context, id uuid.UUID) error {
	return s.methodRepo.Delete(ctx, id)
}
