package models

import (
	"backend-sembako/internal/constants"

	"github.com/google/uuid"
)

type Permission struct {
	BaseModel
	Code        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"code"`
	Name        string `gorm:"type:varchar(150);not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}

type RolePermission struct {
	BaseModel
	Role         constants.UserRole `gorm:"type:varchar(30);not null;index" json:"role"`
	PermissionID uuid.UUID          `gorm:"type:uuid;not null;index" json:"permission_id"`

	// Relationship
	Permission *Permission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE" json:"permission,omitempty"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
