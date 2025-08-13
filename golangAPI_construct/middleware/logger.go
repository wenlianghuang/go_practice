package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger outputs structured single-line log per request.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqID := c.GetString("request_id")
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		errCode := ""
		if len(c.Errors) > 0 {
			errCode = c.Errors[0].Err.Error()
		}
		if rawQuery != "" {
			path += "?" + rawQuery
		}

		log.Printf("[REQ] id=%s method=%s path=%s status=%d latency=%s err=%s",
			reqID, c.Request.Method, path, status, latency, errCode)
	}
}
