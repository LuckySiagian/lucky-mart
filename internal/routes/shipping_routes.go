package routes

import (
	"backend-sembako/internal/constants"
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterShippingRoutes(r *gin.RouterGroup, shippingHandler *handler.ShippingHandler, jwtService jwt.JWTService) {
	shipping := r.Group("/shipping")
	{
		shipping.GET("/methods", shippingHandler.GetShippingMethods)
		shipping.POST("/calculate", shippingHandler.CalculateShippingRate)

		protected := shipping.Group("")
		protected.Use(middleware.JWTAuth(jwtService))
		{
			protected.GET("/orders/:order_id", shippingHandler.GetShipmentByOrderID)
		}

		adminOnly := shipping.Group("")
		adminOnly.Use(middleware.JWTAuth(jwtService))
		adminOnly.Use(middleware.RequireRoles(constants.RoleSuperAdmin, constants.RoleStoreManager, constants.RoleWarehouseStaff, constants.RoleDriver))
		{
			adminOnly.POST("/shipments", shippingHandler.CreateShipment)
			adminOnly.PATCH("/shipments/:id/status", shippingHandler.UpdateShipmentStatus)
		}
	}
}
