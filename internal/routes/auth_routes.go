package routes

import (
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup, authHandler *handler.AuthHandler, jwtService jwt.JWTService) {
	auth := r.Group("/auth")
	{
		// Public Endpoints
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/verify-email", authHandler.VerifyEmail)

		// Protected Endpoints (Requires JWT Access Token)
		protected := auth.Group("")
		protected.Use(middleware.JWTAuth(jwtService))
		{
			protected.GET("/me", authHandler.GetMe)
			protected.PUT("/me", authHandler.UpdateProfile)
			protected.POST("/logout", authHandler.Logout)
			protected.POST("/change-password", authHandler.ChangePassword)
		}
	}
}
