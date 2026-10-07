package models

type Brand struct {
	BaseModel
	Name    string  `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Slug    string  `gorm:"type:varchar(120);not null;uniqueIndex" json:"slug"`
	LogoURL *string `gorm:"type:varchar(500)" json:"logo_url,omitempty"`
}

func (Brand) TableName() string {
	return "brands"
}
