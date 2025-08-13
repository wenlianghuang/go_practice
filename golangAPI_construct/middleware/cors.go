package middleware

import (
	"net/http"

	"golangAPI_construct/config"

	"github.com/gin-gonic/gin"
)

// CORS now uses config.LoadCORSOrigins() to build whitelist.
func CORS() gin.HandlerFunc {
	allowOrigins := config.LoadCORSOrigins()
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := false
		if len(allowOrigins) > 0 {
			for _, o := range allowOrigins {
				if o == origin {
					allowed = true
					break
				}
			}
		}
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Max-Age", "600")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
