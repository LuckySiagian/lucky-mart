package routes

import "backend-sembako/internal/constants"

// Role sets used to guard routes. Super Admin passes every check.
var (
	managerRoles  = []constants.UserRole{constants.RoleSuperAdmin, constants.RoleStoreManager}
	warehouseRole = []constants.UserRole{constants.RoleSuperAdmin, constants.RoleStoreManager, constants.RoleWarehouseStaff}
	deliveryRoles = []constants.UserRole{constants.RoleSuperAdmin, constants.RoleStoreManager, constants.RoleWarehouseStaff, constants.RoleDriver}
)
