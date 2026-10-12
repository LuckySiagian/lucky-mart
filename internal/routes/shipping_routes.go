package routes

import (
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

		staff := shipping.Group("")
		staff.Use(middleware.JWTAuth(jwtService), middleware.RequireRoles(deliveryRoles...))
		{
			staff.GET("/shipments", shippingHandler.GetShipments)
			staff.POST("/shipments", shippingHandler.CreateShipment)
			staff.PATCH("/shipments/:id/status", shippingHandler.UpdateShipmentStatus)
		}
	}
}
