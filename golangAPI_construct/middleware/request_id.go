package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// RequestID sets a unique X-Request-ID header & context value.
func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b := make([]byte, 8)
			if _, err := rand.Read(b); err != nil {
				copy(b, []byte("fallback1"))
			}
			id := hex.EncodeToString(b)
			ctx := context.WithValue(r.Context(), "request_id", id)
			w.Header().Set("X-Request-ID", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
