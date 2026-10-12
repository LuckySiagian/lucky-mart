package dto

// UpdateStoreSettingDTO replaces the store's public details. Empty optional
// fields clear the value.
type UpdateStoreSettingDTO struct {
	Name         string `json:"name" binding:"required,min=2,max=100"`
	Tagline      string `json:"tagline" binding:"max=150"`
	Address      string `json:"address" binding:"max=500"`
	WhatsApp     string `json:"whatsapp" binding:"max=20"`
	OpenHours    string `json:"open_hours" binding:"max=100"`
	MapsURL      string `json:"maps_url" binding:"max=500"`
	ReceiptNote  string `json:"receipt_note" binding:"max=255"`
	Announcement string `json:"announcement" binding:"max=255"`
}
