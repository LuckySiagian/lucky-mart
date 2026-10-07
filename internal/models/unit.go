package models

type Unit struct {
	BaseModel
	Name      string `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`       // Karung, Dus, Kg, Pcs
	ShortName string `gorm:"type:varchar(15);not null;uniqueIndex" json:"short_name"` // krg, dus, kg, pcs
}

func (Unit) TableName() string {
	return "units"
}
