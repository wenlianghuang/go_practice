package middleware

import (
	"net/http"
	"runtime/debug"

	"golangAPI_construct/logging"
	"golangAPI_construct/responses"
)

// ErrorHandler centralizes errors and recovers panics into unified JSON.
func ErrorHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					reqID := r.Context().Value("request_id")
					if reqID == nil {
						reqID = "unknown"
					}
					logging.Logger.Printf("[PANIC] req_id=%s panic=%v\n%s",
						reqID, rec, string(debug.Stack()))
					responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error"))
					return
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
