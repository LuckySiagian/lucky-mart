package models

type NotificationTemplate struct {
	BaseModel
	Name      string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	EventType string `gorm:"type:varchar(50);not null;index" json:"event_type"` // ORDER_CREATED, PAYMENT_PAID, etc.
	Channel   string `gorm:"type:varchar(30);not null;index" json:"channel"`    // EMAIL, WHATSAPP, IN_APP
	Subject   string `gorm:"type:varchar(200);not null" json:"subject"`
	Body      string `gorm:"type:text;not null" json:"body"`
	IsActive  bool   `gorm:"not null;default:true" json:"is_active"`
}

func (NotificationTemplate) TableName() string {
	return "notification_templates"
}
