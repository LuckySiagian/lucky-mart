package services

import (
	"context"
	"errors"
	"strings"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
)

var ErrPaymentMethodNotFound = errors.New("metode pembayaran tidak ditemukan")

type PaymentMethodService interface {
	GetPaymentMethods(ctx context.Context, includeInactive bool) ([]dto.PaymentMethodDTO, error)
	CreatePaymentMethod(ctx context.Context, req dto.CreatePaymentMethodDTO) (*dto.PaymentMethodDTO, error)
	UpdatePaymentMethod(ctx context.Context, id uuid.UUID, req dto.UpdatePaymentMethodDTO) (*dto.PaymentMethodDTO, error)
	DeletePaymentMethod(ctx context.Context, id uuid.UUID) error
}

type paymentMethodService struct {
	methodRepo repositories.PaymentMethodRepository
}

func NewPaymentMethodService(methodRepo repositories.PaymentMethodRepository) PaymentMethodService {
	return &paymentMethodService{methodRepo: methodRepo}
}

func (s *paymentMethodService) GetPaymentMethods(ctx context.Context, includeInactive bool) ([]dto.PaymentMethodDTO, error) {
	var methods []models.PaymentMethod
	var err error
	if includeInactive {
		methods, err = s.methodRepo.FindAllIncludingInactive(ctx)
	} else {
		methods, err = s.methodRepo.FindAll(ctx)
	}
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.PaymentMethodDTO, len(methods))
	for i := range methods {
		dtos[i] = mapPaymentMethod(&methods[i])
	}
	return dtos, nil
}

func (s *paymentMethodService) CreatePaymentMethod(ctx context.Context, req dto.CreatePaymentMethodDTO) (*dto.PaymentMethodDTO, error) {
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	existing, err := s.methodRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("kode metode pembayaran sudah dipakai")
	}

	method := &models.PaymentMethod{
		Code:        code,
		Name:        strings.TrimSpace(req.Name),
		Type:        req.Type,
		Description: strings.TrimSpace(req.Description),
		IsActive:    true,
	}
	if err := s.methodRepo.Create(ctx, method); err != nil {
		return nil, err
	}
	res := mapPaymentMethod(method)
	return &res, nil
}

// UpdatePaymentMethod changes the label, instructions (e.g. the bank account
// customers transfer to) and availability. Code and type stay fixed because
// orders and payment flows depend on them.
func (s *paymentMethodService) UpdatePaymentMethod(ctx context.Context, id uuid.UUID, req dto.UpdatePaymentMethodDTO) (*dto.PaymentMethodDTO, error) {
	method, err := s.methodRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if method == nil {
		return nil, ErrPaymentMethodNotFound
	}

	if req.Name != nil {
		method.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		method.Description = strings.TrimSpace(*req.Description)
	}
	if req.IsActive != nil {
		method.IsActive = *req.IsActive
	}
	if err := s.methodRepo.Update(ctx, method); err != nil {
		return nil, err
	}
	res := mapPaymentMethod(method)
	return &res, nil
}

func (s *paymentMethodService) DeletePaymentMethod(ctx context.Context, id uuid.UUID) error {
	method, err := s.methodRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if method == nil {
		return ErrPaymentMethodNotFound
	}
	return s.methodRepo.Delete(ctx, id)
}

func mapPaymentMethod(m *models.PaymentMethod) dto.PaymentMethodDTO {
	return dto.PaymentMethodDTO{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Type:        m.Type,
		IsActive:    m.IsActive,
		Description: m.Description,
	}
}
