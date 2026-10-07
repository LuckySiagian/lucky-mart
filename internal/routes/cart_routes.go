package routes

import (
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(r *gin.RouterGroup, cartHandler *handler.CartHandler, jwtService jwt.JWTService) {
	cart := r.Group("/cart")
	cart.Use(middleware.JWTAuth(jwtService))
	{
		cart.GET("", cartHandler.GetCart)
		cart.POST("/items", cartHandler.AddToCart)
		cart.PUT("/items/:id", cartHandler.UpdateCartItem)
		cart.DELETE("/items/:id", cartHandler.RemoveCartItem)
		cart.DELETE("", cartHandler.ClearCart)
		cart.POST("/checkout", cartHandler.GetCheckoutSummary)
	}
}
