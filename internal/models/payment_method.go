package models

type PaymentMethod struct {
	BaseModel
	Code        string `gorm:"type:varchar(30);not null;uniqueIndex" json:"code"` // CASH, COD, QRIS, MIDTRANS, B2B_CREDIT
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Type        string `gorm:"type:varchar(30);not null;index" json:"type"` // CASH, COD, MANUAL_TRANSFER, GATEWAY, B2B_CREDIT
	IsActive    bool   `gorm:"not null;default:true" json:"is_active"`
	Description string `gorm:"type:text" json:"description,omitempty"`
}

func (PaymentMethod) TableName() string {
	return "payment_methods"
}
