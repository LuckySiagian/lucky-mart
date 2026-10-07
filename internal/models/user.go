package models

import (
	"time"

	"backend-sembako/internal/constants"
)

type User struct {
	BaseModel
	FullName     string             `gorm:"type:varchar(100);not null" json:"full_name"`
	Email        *string            `gorm:"type:varchar(150);uniqueIndex;nullable" json:"email,omitempty"`
	PhoneNumber  string             `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone_number"`
	PasswordHash string             `gorm:"type:varchar(255);not null" json:"-"`
	Role         constants.UserRole `gorm:"type:varchar(30);not null;default:'B2C_CUSTOMER'" json:"role"`
	IsActive     bool               `gorm:"not null;default:true" json:"is_active"`
	AvatarURL    *string            `gorm:"type:varchar(500);nullable" json:"avatar_url,omitempty"`

	// Email Verification Fields
	IsEmailVerified            bool       `gorm:"not null;default:false" json:"is_email_verified"`
	EmailVerificationToken     *string    `gorm:"type:varchar(255);nullable;index" json:"-"`
	EmailVerificationExpiresAt *time.Time `gorm:"nullable" json:"-"`

	// Password Reset Fields
	PasswordResetToken     *string    `gorm:"type:varchar(255);nullable;index" json:"-"`
	PasswordResetExpiresAt *time.Time `gorm:"nullable" json:"-"`

	// Security & Lockout Fields
	FailedLoginAttempts int        `gorm:"not null;default:0" json:"-"`
	LockoutUntil        *time.Time `gorm:"nullable" json:"-"`
	LastLoginAt         *time.Time `gorm:"nullable" json:"last_login_at,omitempty"`

	// Relationships
	B2BProfile *B2BProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"b2b_profile,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsLockedOut() bool {
	if u.LockoutUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockoutUntil)
}
