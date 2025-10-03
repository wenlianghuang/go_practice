package middleware

import (
	"net/http"
	"time"

	"golangAPI_construct/logging"
)

// Logger outputs one structured line per request to the unified logger.
func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID := r.Context().Value("request_id")
			if reqID == nil {
				reqID = "unknown"
			}
			path := r.URL.Path
			rawQuery := r.URL.RawQuery

			// Create a response writer wrapper to capture status
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			if rawQuery != "" {
				path += "?" + rawQuery
			}
			latency := time.Since(start)
			status := wrapped.statusCode

			logging.Logger.Printf("[REQ] id=%s method=%s path=%s status=%d latency=%s",
				reqID, r.Method, path, status, latency)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
