package routes

import (
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup, userHandler *handler.UserHandler, jwtService jwt.JWTService) {
	users := r.Group("/users")
	users.Use(middleware.JWTAuth(jwtService), middleware.RequireRoles(managerRoles...))
	{
		// Which accounts a manager may change is checked in the service:
		// only a Super Admin can manage admin-level accounts.
		users.GET("", userHandler.GetUsers)
		users.GET("/:id", userHandler.GetUserByID)
		users.POST("", userHandler.CreateUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.POST("/:id/reset-password", userHandler.ResetPassword)
		users.PATCH("/:id/b2b", userHandler.VerifyB2B)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}
