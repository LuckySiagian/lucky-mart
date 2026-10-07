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
		orders.POST("", orderHandler.CreateOrder)
		orders.GET("", orderHandler.GetOrders)
		orders.GET("/:id", orderHandler.GetOrderByID)
		orders.PUT("/:id/cancel", orderHandler.CancelOrder)
		orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		orders.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
		orders.GET("/:id/history", orderHandler.GetOrderHistory)
	}

	invoices := r.Group("/invoices")
	invoices.Use(middleware.JWTAuth(jwtService))
	{
		invoices.GET("/:id", orderHandler.GetInvoiceByID)
	}
}
