package middleware

import (
	"context"
	"net/http"
	"strings"

	"golangAPI_construct/responses"
	"golangAPI_construct/security"
)

// JWTAuthMiddleware 統一的 JWT 驗證中間件，支援 Token 黑名單
func JWTAuthMiddleware(store security.TokenStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "MISSING_AUTH_HEADER", "Authorization header is required"))
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "INVALID_AUTH_FORMAT", "Authorization header must start with 'Bearer '"))
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "EMPTY_TOKEN", "Token cannot be empty"))
				return
			}

			claims, err := security.ValidateToken(tokenString)
			if err != nil {
				responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token"))
				return
			}

			if store != nil && claims.ID != "" {
				revoked, err := store.IsRevoked(claims.ID)
				if err != nil {
					responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "TOKEN_CHECK_FAILED", "Failed to verify token status"))
					return
				}
				if revoked {
					responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "TOKEN_REVOKED", "Token has been revoked"))
					return
				}
			}

			ctx := context.WithValue(r.Context(), "user", claims.Username)
			ctx = context.WithValue(ctx, "user_id", claims.Subject)
			ctx = context.WithValue(ctx, "roles", claims.Roles)
			ctx = context.WithValue(ctx, "claims", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalJWTAuthMiddleware 可選的 JWT 驗證中間件（用於不需要強制登入的端點）
func OptionalJWTAuthMiddleware(store security.TokenStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := security.ValidateToken(tokenString)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if store != nil && claims.ID != "" {
				revoked, err := store.IsRevoked(claims.ID)
				if err != nil || revoked {
					next.ServeHTTP(w, r)
					return
				}
			}

			ctx := context.WithValue(r.Context(), "user", claims.Username)
			ctx = context.WithValue(ctx, "user_id", claims.Subject)
			ctx = context.WithValue(ctx, "roles", claims.Roles)
			ctx = context.WithValue(ctx, "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole 基於角色的權限檢查中間件
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value("user") == nil {
				responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "NOT_AUTHENTICATED", "User not authenticated"))
				return
			}

			roles, _ := r.Context().Value("roles").([]string)
			if !containsRole(roles, requiredRole) {
				responses.Fail(w, r, responses.NewAppError(http.StatusForbidden, "FORBIDDEN", "Insufficient role"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func containsRole(roles []string, required string) bool {
	for _, role := range roles {
		if role == required {
			return true
		}
	}
	return false
}
