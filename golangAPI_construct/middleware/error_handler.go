package middleware

import (
	"net/http"
	"runtime/debug"

	"golangAPI_construct/logging"
	"golangAPI_construct/responses"

	"github.com/gin-gonic/gin"
)

// ErrorHandler centralizes errors and recovers panics into unified JSON.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logging.Logger.Printf("[PANIC] req_id=%s panic=%v\n%s",
					c.GetString("request_id"), rec, string(debug.Stack()))
				responses.Fail(c, responses.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
				c.Abort()
			}
		}()

		c.Next()

		if c.IsAborted() || len(c.Errors) == 0 {
			return
		}

		gErr := c.Errors[0]
		if appErr, ok := gErr.Err.(*responses.AppError); ok {
			responses.Fail(c, appErr)
			return
		}
		responses.Fail(c, responses.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
	}
}
