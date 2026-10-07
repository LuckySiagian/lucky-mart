package models

type Courier struct {
	BaseModel
	Code     string  `gorm:"type:varchar(30);not null;uniqueIndex" json:"code"` // KURIR_TOKO, JNE, TIKI, POS
	Name     string  `gorm:"type:varchar(100);not null" json:"name"`
	Phone    *string `gorm:"type:varchar(20)" json:"phone,omitempty"`
	IsActive bool    `gorm:"not null;default:true" json:"is_active"`
}

func (Courier) TableName() string {
	return "couriers"
}
