package repositories

import (
	"context"

	"backend-sembako/internal/models"

	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) GetRolePermissions(ctx context.Context, role string) ([]string, error) {
	var codes []string
	err := r.db.WithContext(ctx).
		Table("role_permissions").
		Select("permissions.code").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role = ?", role).
		Pluck("permissions.code", &codes).Error

	return codes, err
}

func (r *permissionRepository) HasPermission(ctx context.Context, role string, permissionCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.RolePermission{}).
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role = ? AND permissions.code = ?", role, permissionCode).
		Count(&count).Error

	return count > 0, err
}
