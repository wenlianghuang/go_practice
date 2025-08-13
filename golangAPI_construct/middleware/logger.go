package middleware

import (
	"time"

	"golangAPI_construct/logging"

	"github.com/gin-gonic/gin"
)

// Logger outputs one structured line per request to the unified logger.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqID := c.GetString("request_id")
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		if rawQuery != "" {
			path += "?" + rawQuery
		}
		latency := time.Since(start)
		status := c.Writer.Status()
		errCode := ""
		if len(c.Errors) > 0 {
			// 如果是 AppError 會在 error_handler 已格式化回應；這裡只記錄第一個
			errCode = c.Errors[0].Error()
		}

		logging.Logger.Printf("[REQ] id=%s method=%s path=%s status=%d latency=%s err=%s",
			reqID, c.Request.Method, path, status, latency, errCode)
	}
}
