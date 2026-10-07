package routes

import (
	"backend-sembako/internal/constants"
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup, userHandler *handler.UserHandler, jwtService jwt.JWTService) {
	users := r.Group("/users")
	users.Use(middleware.JWTAuth(jwtService))
	users.Use(middleware.RequireRoles(constants.RoleSuperAdmin, constants.RoleStoreManager))
	{
		users.GET("", userHandler.GetUsers)
		users.GET("/:id", userHandler.GetUserByID)
		users.POST("", userHandler.CreateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}
