package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"golangAPI_construct/responses"

	"github.com/gin-gonic/gin"
)

// ErrorHandler centralizes error → unified JSON envelope.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC] req_id=%s panic=%v\n%s", c.GetString("request_id"), rec, string(debug.Stack()))
				appErr := responses.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
				responses.Fail(c, appErr)
				c.Abort()
			}
		}()

		c.Next()

		// If response already written or no errors → do nothing.
		if c.IsAborted() || len(c.Errors) == 0 {
			return
		}

		// Use first error; could be extended to map types.
		gErr := c.Errors[0]
		if appErr, ok := gErr.Err.(*responses.AppError); ok {
			responses.Fail(c, appErr)
			return
		}

		// Fallback generic
		responses.Fail(c, responses.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
	}
}
