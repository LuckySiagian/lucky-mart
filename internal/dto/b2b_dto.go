package dto

import (
	"time"

	"backend-sembako/internal/constants"

	"github.com/google/uuid"
)

type B2BProfileDTO struct {
	ID                  uuid.UUID                       `json:"id"`
	UserID              uuid.UUID                       `json:"user_id"`
	StoreName           string                          `json:"store_name"`
	StoreAddress        string                          `json:"store_address"`
	TaxID               *string                         `json:"tax_id_npwp,omitempty"`
	VerificationStatus  constants.B2BVerificationStatus `json:"verification_status"`
	ApprovedCreditLimit float64                         `json:"approved_credit_limit"`
	CreatedAt           time.Time                       `json:"created_at"`
}

type CreateB2BProfileDTO struct {
	StoreName    string  `json:"store_name" binding:"required,min=3,max=150"`
	StoreAddress string  `json:"store_address" binding:"required"`
	TaxID        *string `json:"tax_id_npwp" binding:"omitempty"`
}

type UpdateB2BCreditLimitDTO struct {
	ApprovedCreditLimit float64                         `json:"approved_credit_limit" binding:"required,gte=0"`
	VerificationStatus  constants.B2BVerificationStatus `json:"verification_status" binding:"required"`
}
