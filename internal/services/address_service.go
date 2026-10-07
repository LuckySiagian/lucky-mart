package services

import (
	"context"
	"errors"

	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
)

var (
	ErrAddressNotFound    = errors.New("address not found")
	ErrUnauthorizedAddress = errors.New("you do not have permission to access or modify this address")
)

type AddressService interface {
	CreateAddress(ctx context.Context, userID uuid.UUID, req dto.CreateAddressDTO) (*dto.AddressResponseDTO, error)
	GetAddresses(ctx context.Context, userID uuid.UUID) ([]dto.AddressResponseDTO, error)
	GetAddressByID(ctx context.Context, userID, addressID uuid.UUID) (*dto.AddressResponseDTO, error)
	UpdateAddress(ctx context.Context, userID, addressID uuid.UUID, req dto.CreateAddressDTO) (*dto.AddressResponseDTO, error)
	DeleteAddress(ctx context.Context, userID, addressID uuid.UUID) error
	SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID) (*dto.AddressResponseDTO, error)
}

type addressService struct {
	addressRepo repositories.AddressRepository
}

func NewAddressService(addressRepo repositories.AddressRepository) AddressService {
	return &addressService{addressRepo: addressRepo}
}

func (s *addressService) CreateAddress(ctx context.Context, userID uuid.UUID, req dto.CreateAddressDTO) (*dto.AddressResponseDTO, error) {
	if req.IsDefault {
		_ = s.addressRepo.UnsetDefault(ctx, userID)
	}

	addr := &models.Address{
		UserID:        userID,
		Title:         req.Title,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		FullAddress:   req.FullAddress,
		City:          req.City,
		Province:      req.Province,
		PostalCode:    req.PostalCode,
		IsDefault:     req.IsDefault,
	}

	if err := s.addressRepo.Create(ctx, addr); err != nil {
		return nil, err
	}

	return s.GetAddressByID(ctx, userID, addr.ID)
}

func (s *addressService) GetAddresses(ctx context.Context, userID uuid.UUID) ([]dto.AddressResponseDTO, error) {
	addresses, err := s.addressRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.AddressResponseDTO, len(addresses))
	for i, a := range addresses {
		dtos[i] = dto.AddressResponseDTO{
			ID:            a.ID,
			UserID:        a.UserID,
			Title:         a.Title,
			RecipientName: a.RecipientName,
			Phone:         a.Phone,
			FullAddress:   a.FullAddress,
			City:          a.City,
			Province:      a.Province,
			PostalCode:    a.PostalCode,
			IsDefault:     a.IsDefault,
			CreatedAt:     a.CreatedAt,
		}
	}

	return dtos, nil
}

func (s *addressService) GetAddressByID(ctx context.Context, userID, addressID uuid.UUID) (*dto.AddressResponseDTO, error) {
	addr, err := s.addressRepo.FindByID(ctx, addressID)
	if err != nil || addr == nil {
		return nil, ErrAddressNotFound
	}

	if addr.UserID != userID {
		return nil, ErrUnauthorizedAddress
	}

	return &dto.AddressResponseDTO{
		ID:            addr.ID,
		UserID:        addr.UserID,
		Title:         addr.Title,
		RecipientName: addr.RecipientName,
		Phone:         addr.Phone,
		FullAddress:   addr.FullAddress,
		City:          addr.City,
		Province:      addr.Province,
		PostalCode:    addr.PostalCode,
		IsDefault:     addr.IsDefault,
		CreatedAt:     addr.CreatedAt,
	}, nil
}

func (s *addressService) UpdateAddress(ctx context.Context, userID, addressID uuid.UUID, req dto.CreateAddressDTO) (*dto.AddressResponseDTO, error) {
	addr, err := s.addressRepo.FindByID(ctx, addressID)
	if err != nil || addr == nil {
		return nil, ErrAddressNotFound
	}

	if addr.UserID != userID {
		return nil, ErrUnauthorizedAddress
	}

	if req.IsDefault && !addr.IsDefault {
		_ = s.addressRepo.UnsetDefault(ctx, userID)
	}

	addr.Title = req.Title
	addr.RecipientName = req.RecipientName
	addr.Phone = req.Phone
	addr.FullAddress = req.FullAddress
	addr.City = req.City
	addr.Province = req.Province
	addr.PostalCode = req.PostalCode
	addr.IsDefault = req.IsDefault

	if err := s.addressRepo.Update(ctx, addr); err != nil {
		return nil, err
	}

	return s.GetAddressByID(ctx, userID, addressID)
}

func (s *addressService) DeleteAddress(ctx context.Context, userID, addressID uuid.UUID) error {
	addr, err := s.addressRepo.FindByID(ctx, addressID)
	if err != nil || addr == nil {
		return ErrAddressNotFound
	}

	if addr.UserID != userID {
		return ErrUnauthorizedAddress
	}

	return s.addressRepo.Delete(ctx, addressID)
}

func (s *addressService) SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID) (*dto.AddressResponseDTO, error) {
	addr, err := s.addressRepo.FindByID(ctx, addressID)
	if err != nil || addr == nil {
		return nil, ErrAddressNotFound
	}

	if addr.UserID != userID {
		return nil, ErrUnauthorizedAddress
	}

	if err := s.addressRepo.UnsetDefault(ctx, userID); err != nil {
		return nil, err
	}

	addr.IsDefault = true
	if err := s.addressRepo.Update(ctx, addr); err != nil {
		return nil, err
	}

	return s.GetAddressByID(ctx, userID, addressID)
}
