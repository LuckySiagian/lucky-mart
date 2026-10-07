package models

type Category struct {
	BaseModel
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Slug        string `gorm:"type:varchar(120);not null;uniqueIndex" json:"slug"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	IsActive    bool   `gorm:"not null;default:true" json:"is_active"`

	// Relationships
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
