package middleware

import (
	"backend-sembako/helpers"
	"backend-sembako/internal/constants"

	"github.com/gin-gonic/gin"
)

func RequireRoles(allowedRoles ...constants.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(ContextUserRole)
		if !exists {
			helpers.SendUnauthorized(c, "User role missing in request context")
			c.Abort()
			return
		}

		userRole, ok := roleVal.(constants.UserRole)
		if !ok {
			helpers.SendUnauthorized(c, "Invalid user role type")
			c.Abort()
			return
		}

		// Super Admin bypasses role restrictions
		if userRole == constants.RoleSuperAdmin {
			c.Next()
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		helpers.SendForbidden(c, "You do not have permission to access this resource")
		c.Abort()
	}
}
