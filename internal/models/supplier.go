package models

type Supplier struct {
	BaseModel
	Name        string  `gorm:"type:varchar(150);not null" json:"name"`
	ContactName *string `gorm:"type:varchar(100)" json:"contact_name,omitempty"`
	Phone       string  `gorm:"type:varchar(20);not null" json:"phone"`
	Email       *string `gorm:"type:varchar(150)" json:"email,omitempty"`
	Address     string  `gorm:"type:text;not null" json:"address"`
	IsActive    bool    `gorm:"not null;default:true" json:"is_active"`
}

func (Supplier) TableName() string {
	return "suppliers"
}
