package middleware

import (
	"fmt"

	"backend-sembako/helpers"
	"backend-sembako/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				reqID, _ := c.Get(HeaderXRequestID)

				logger.Log.Error().
					Str("request_id", fmt.Sprintf("%v", reqID)).
					Interface("panic_error", err).
					Msg("Unhandled panic recovered")

				helpers.SendInternalServerError(c, "An unexpected internal server error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
