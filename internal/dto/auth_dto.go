package dto

import (
	"time"

	"backend-sembako/internal/constants"

	"github.com/google/uuid"
)

type RegisterCustomerDTO struct {
	FullName          string  `json:"full_name" binding:"required,min=3,max=100"`
	PhoneNumber       string  `json:"phone_number" binding:"required,phone_id"`
	Email             *string `json:"email" binding:"omitempty,email"`
	Password          string  `json:"password" binding:"required,min=8"`
	IsB2BRegistration bool    `json:"is_b2b_registration"`
	StoreName         string  `json:"store_name" binding:"required_if=IsB2BRegistration true"`
	StoreAddress      string  `json:"store_address" binding:"required_if=IsB2BRegistration true"`
	TaxID             *string `json:"tax_id_npwp" binding:"omitempty"`
}

type LoginDTO struct {
	PhoneNumber string `json:"phone_number" binding:"required,phone_id"`
	Password    string `json:"password" binding:"required"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ForgotPasswordDTO struct {
	PhoneNumber string `json:"phone_number" binding:"required,phone_id"`
}

type ResetPasswordDTO struct {
	ResetToken  string `json:"reset_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type VerifyEmailDTO struct {
	Token string `json:"token" binding:"required"`
}

type AuthResponseDTO struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	TokenType    string              `json:"token_type"`
	ExpiresIn    int64               `json:"expires_in"`
	User         UserProfileResponse `json:"user"`
}

type UserProfileResponse struct {
	ID              uuid.UUID          `json:"id"`
	FullName        string             `json:"full_name"`
	Email           *string            `json:"email,omitempty"`
	PhoneNumber     string             `json:"phone_number"`
	Role            constants.UserRole `json:"role"`
	IsActive        bool               `json:"is_active"`
	IsEmailVerified bool               `json:"is_email_verified"`
	AvatarURL       *string            `json:"avatar_url,omitempty"`
	LastLoginAt     *time.Time         `json:"last_login_at,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	B2BProfile      *B2BProfileDTO     `json:"b2b_profile,omitempty"`
}
