package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateAddressDTO struct {
	Title         string `json:"title" binding:"required,min=2,max=50"`
	RecipientName string `json:"recipient_name" binding:"required,min=3,max=100"`
	Phone         string `json:"phone" binding:"required,phone_id"`
	FullAddress   string `json:"full_address" binding:"required"`
	City          string `json:"city" binding:"required"`
	Province      string `json:"province" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	IsDefault     bool   `json:"is_default"`
}

type AddressResponseDTO struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Title         string    `json:"title"`
	RecipientName string    `json:"recipient_name"`
	Phone         string    `json:"phone"`
	FullAddress   string    `json:"full_address"`
	City          string    `json:"city"`
	Province      string    `json:"province"`
	PostalCode    string    `json:"postal_code"`
	IsDefault     bool      `json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
}
