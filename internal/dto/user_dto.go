package dto

import (
	"time"

	"backend-sembako/internal/constants"

	"github.com/google/uuid"
)

type UserResponseDTO struct {
	ID          uuid.UUID          `json:"id"`
	FullName    string             `json:"full_name"`
	Email       *string            `json:"email,omitempty"`
	PhoneNumber string             `json:"phone_number"`
	Role        constants.UserRole `json:"role"`
	IsActive    bool               `json:"is_active"`
	CreatedAt   time.Time          `json:"created_at"`
	B2BProfile  *B2BProfileDTO     `json:"b2b_profile,omitempty"`
}

type CreateUserDTO struct {
	FullName    string             `json:"full_name" binding:"required,min=3,max=100"`
	Email       *string            `json:"email" binding:"omitempty,email"`
	PhoneNumber string             `json:"phone_number" binding:"required,phone_id"`
	Password    string             `json:"password" binding:"required,min=8"`
	Role        constants.UserRole `json:"role" binding:"required"`
}
