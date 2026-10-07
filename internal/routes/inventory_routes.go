package routes

import (
	"backend-sembako/internal/constants"
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterInventoryRoutes(r *gin.RouterGroup, inventoryHandler *handler.InventoryHandler, jwtService jwt.JWTService) {
	inventory := r.Group("/inventory")
	inventory.Use(middleware.JWTAuth(jwtService))
	inventory.Use(middleware.RequireRoles(constants.RoleSuperAdmin, constants.RoleStoreManager, constants.RoleWarehouseStaff))
	{
		inventory.POST("/adjust", inventoryHandler.AdjustStock)
		inventory.GET("/low-stock", inventoryHandler.GetLowStockAlerts)
		inventory.GET("/:id/movements", inventoryHandler.GetStockHistory)
	}
}
