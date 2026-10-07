package middleware

import (
	"time"

	"backend-sembako/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		reqID, _ := c.Get(HeaderXRequestID)

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		entry := logger.Log.Info()
		if statusCode >= 400 && statusCode < 500 {
			entry = logger.Log.Warn()
		} else if statusCode >= 500 {
			entry = logger.Log.Error()
		}

		entry.
			Str("request_id", reqID.(string)).
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP Request")
	}
}
