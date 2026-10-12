package routes

import (
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.RouterGroup, orderHandler *handler.OrderHandler, jwtService jwt.JWTService) {
	orders := r.Group("/orders")
	orders.Use(middleware.JWTAuth(jwtService))
	{
		// Customers see only their own orders; staff see all (enforced in the service).
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("", orderHandler.GetOrders)
		orders.GET("/:id", orderHandler.GetOrderByID)
		orders.PUT("/:id/cancel", orderHandler.CancelOrder)
		orders.GET("/:id/history", orderHandler.GetOrderHistory)

		staff := orders.Group("")
		staff.Use(middleware.RequireRoles(deliveryRoles...))
		{
			staff.PUT("/:id/status", orderHandler.UpdateOrderStatus)
			staff.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
		}
	}

	invoices := r.Group("/invoices")
	invoices.Use(middleware.JWTAuth(jwtService))
	{
		invoices.GET("/:id", orderHandler.GetInvoiceByID)
	}
}
