package routes

import (
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterAddressRoutes(r *gin.RouterGroup, addressHandler *handler.AddressHandler, jwtService jwt.JWTService) {
	addresses := r.Group("/addresses")
	addresses.Use(middleware.JWTAuth(jwtService))
	{
		addresses.POST("", addressHandler.CreateAddress)
		addresses.GET("", addressHandler.GetAddresses)
		addresses.GET("/:id", addressHandler.GetAddressByID)
		addresses.PUT("/:id", addressHandler.UpdateAddress)
		addresses.DELETE("/:id", addressHandler.DeleteAddress)
		addresses.PATCH("/:id/default", addressHandler.SetDefaultAddress)
	}
}
