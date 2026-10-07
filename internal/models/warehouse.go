package models

type Warehouse struct {
	BaseModel
	Code     string `gorm:"type:varchar(30);not null;uniqueIndex" json:"code"`
	Name     string `gorm:"type:varchar(150);not null" json:"name"`
	Address  string `gorm:"type:text;not null" json:"address"`
	IsActive bool   `gorm:"not null;default:true" json:"is_active"`
}

func (Warehouse) TableName() string {
	return "warehouses"
}
