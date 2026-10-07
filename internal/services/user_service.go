package services

import (
	"context"
	"errors"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/models"
	"backend-sembako/internal/repositories"

	"github.com/google/uuid"
)

type UserService interface {
	GetUsers(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserProfileResponse, error)
	CreateUser(ctx context.Context, req dto.CreateUserDTO) (*dto.UserProfileResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, fullName string, email *string, isActive bool) (*dto.UserProfileResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUsers(ctx context.Context, req dto.PaginationRequest) (dto.PaginatedResponse, error) {
	users, total, err := s.userRepo.FindAll(ctx, req)
	if err != nil {
		return dto.PaginatedResponse{}, err
	}

	userDtos := make([]dto.UserProfileResponse, len(users))
	for i, u := range users {
		userDtos[i] = dto.UserProfileResponse{
			ID:              u.ID,
			FullName:        u.FullName,
			Email:           u.Email,
			PhoneNumber:     u.PhoneNumber,
			Role:            u.Role,
			IsActive:        u.IsActive,
			IsEmailVerified: u.IsEmailVerified,
			AvatarURL:       u.AvatarURL,
			LastLoginAt:     u.LastLoginAt,
		}
	}

	req.SetDefaults()
	return dto.NewPaginatedResponse(userDtos, req.Page, req.Limit, total), nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserProfileResponse{
		ID:              user.ID,
		FullName:        user.FullName,
		Email:           user.Email,
		PhoneNumber:     user.PhoneNumber,
		Role:            user.Role,
		IsActive:        user.IsActive,
		IsEmailVerified: user.IsEmailVerified,
		AvatarURL:       user.AvatarURL,
		LastLoginAt:     user.LastLoginAt,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserDTO) (*dto.UserProfileResponse, error) {
	if err := helpers.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	existing, err := s.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("phone number is already registered")
	}

	hashed, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: hashed,
		Role:         req.Role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserProfileResponse{
		ID:          user.ID,
		FullName:    user.FullName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		IsActive:    user.IsActive,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, fullName string, email *string, isActive bool) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	user.FullName = fullName
	user.Email = email
	user.IsActive = isActive

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserProfileResponse{
		ID:          user.ID,
		FullName:    user.FullName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		IsActive:    user.IsActive,
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
