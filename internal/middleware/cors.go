package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS allows the configured frontend origins. "*" (development only)
// reflects any origin. Credentials are not used — the API authenticates with
// a bearer token — so no Allow-Credentials header is sent.
func CORS(allowOrigins []string) gin.HandlerFunc {
	allowAll := false
	allowed := make(map[string]bool, len(allowOrigins))
	for _, o := range allowOrigins {
		if o == "*" {
			allowAll = true
		}
		allowed[o] = true
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" && (allowAll || allowed[origin]) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Add("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, Cache-Control, X-Requested-With, X-Request-ID, X-Device-Info")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Max-Age", "600")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
