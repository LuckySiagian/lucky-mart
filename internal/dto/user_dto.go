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

type UpdateUserDTO struct {
	FullName *string             `json:"full_name" binding:"omitempty,min=3,max=100"`
	Email    *string             `json:"email" binding:"omitempty,email"`
	Role     *constants.UserRole `json:"role" binding:"omitempty"`
	IsActive *bool               `json:"is_active" binding:"omitempty"`
}

type AdminResetPasswordDTO struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type VerifyB2BDTO struct {
	Status      constants.B2BVerificationStatus `json:"status" binding:"required,oneof=APPROVED REJECTED PENDING"`
	CreditLimit *float64                        `json:"credit_limit" binding:"omitempty,gte=0"`
}

type UpdateProfileDTO struct {
	FullName *string `json:"full_name" binding:"omitempty,min=3,max=100"`
	Email    *string `json:"email" binding:"omitempty,email"`
}
