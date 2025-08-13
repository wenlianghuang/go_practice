package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// RequestID sets a unique X-Request-ID header & context value.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			copy(b, []byte("fallback1"))
		}
		id := hex.EncodeToString(b)
		c.Set("request_id", id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}
