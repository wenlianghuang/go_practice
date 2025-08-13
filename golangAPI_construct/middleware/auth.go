package middleware

import (
	"net/http"
	"strings"

	"golangAPI_construct/responses"
	"golangAPI_construct/security"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.Error(responses.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "missing bearer token"))
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		claims, err := security.ValidateToken(tokenStr)
		if err != nil {
			c.Error(responses.NewAppError(http.StatusUnauthorized, "INVALID_TOKEN", "invalid token"))
			c.Abort()
			return
		}
		c.Set("user", claims.Subject)
		c.Next()
	}
}
