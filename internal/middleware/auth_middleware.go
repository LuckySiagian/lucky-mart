package middleware

import (
	"strings"

	"backend-sembako/helpers"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey = "user_id"
	ContextUserRole  = "user_role"
)

func JWTAuth(jwtService jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.SendUnauthorized(c, "Authorization header is missing")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			helpers.SendUnauthorized(c, "Invalid Authorization header format. Expected Bearer <token>")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := jwtService.ValidateAccessToken(tokenStr)
		if err != nil {
			helpers.SendUnauthorized(c, "Invalid or expired access token: "+err.Error())
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUserRole, claims.Role)

		c.Next()
	}
}
