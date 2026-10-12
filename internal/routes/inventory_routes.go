package routes

import (
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterInventoryRoutes(r *gin.RouterGroup, inventoryHandler *handler.InventoryHandler, jwtService jwt.JWTService) {
	inventory := r.Group("/inventory")
	inventory.Use(middleware.JWTAuth(jwtService), middleware.RequireRoles(warehouseRole...))
	{
		inventory.GET("", inventoryHandler.GetInventories)
		inventory.GET("/low-stock", inventoryHandler.GetLowStockAlerts)
		inventory.GET("/movements", inventoryHandler.GetMovements)
		inventory.POST("/adjust", inventoryHandler.AdjustStock)
		inventory.GET("/:id/movements", inventoryHandler.GetStockHistory)
		inventory.PATCH("/:id/min-stock", inventoryHandler.UpdateMinStockAlert)
	}
}
