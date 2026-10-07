package models

import (
	"backend-sembako/internal/constants"

	"github.com/google/uuid"
)

type B2BProfile struct {
	BaseModel
	UserID              uuid.UUID                       `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	StoreName           string                          `gorm:"type:varchar(150);not null" json:"store_name"`
	StoreAddress        string                          `gorm:"type:text;not null" json:"store_address"`
	TaxID               *string                         `gorm:"type:varchar(30);nullable" json:"tax_id_npwp,omitempty"`
	VerificationStatus  constants.B2BVerificationStatus `gorm:"type:varchar(20);not null;default:'PENDING'" json:"verification_status"`
	ApprovedCreditLimit float64                         `gorm:"type:numeric(15,2);not null;default:0.00" json:"approved_credit_limit"`
}

func (B2BProfile) TableName() string {
	return "b2b_profiles"
}
