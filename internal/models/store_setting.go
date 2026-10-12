package models

import "time"

// StoreSetting holds the store's public details (one row, ID 1): shown in the
// customer footer, on the pickup page and on printed receipts.
type StoreSetting struct {
	ID           uint   `gorm:"primaryKey" json:"-"`
	Name         string `gorm:"type:varchar(100);not null;default:'Lucky Mart'" json:"name"`
	Tagline      string `gorm:"type:varchar(150)" json:"tagline"`
	Address      string `gorm:"type:text" json:"address"`
	WhatsApp     string `gorm:"type:varchar(20)" json:"whatsapp"`
	OpenHours    string `gorm:"type:varchar(100)" json:"open_hours"`
	MapsURL      string `gorm:"type:varchar(500)" json:"maps_url"`
	ReceiptNote  string `gorm:"type:varchar(255)" json:"receipt_note"`
	Announcement string `gorm:"type:varchar(255)" json:"announcement"`
	// Orders whose goods total reaches this amount ship free (0 = never).
	FreeShippingMin float64   `gorm:"type:numeric(12,2);not null;default:300000" json:"free_shipping_min"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (StoreSetting) TableName() string {
	return "store_settings"
}
