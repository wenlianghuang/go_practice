package middleware

import (
	"context"
	"net/http"
	"strings"

	"golangAPI_construct/responses"
	"golangAPI_construct/security"
)

// JWTAuthMiddleware 統一的 JWT 驗證中間件
func JWTAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 獲取 Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "MISSING_AUTH_HEADER", "Authorization header is required"))
				return
			}

			// 檢查 Bearer token 格式
			if !strings.HasPrefix(authHeader, "Bearer ") {
				responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "INVALID_AUTH_FORMAT", "Authorization header must start with 'Bearer '"))
				return
			}

			// 提取 token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "EMPTY_TOKEN", "Token cannot be empty"))
				return
			}

			// 驗證 token
			claims, err := security.ValidateToken(tokenString)
			if err != nil {
				responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token"))
				return
			}

			// 將用戶信息存儲到 context 中
			ctx := context.WithValue(r.Context(), "user", claims.Username)
			ctx = context.WithValue(ctx, "user_id", claims.Subject)
			ctx = context.WithValue(ctx, "claims", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalJWTAuthMiddleware 可選的 JWT 驗證中間件（用於不需要強制登入的端點）
func OptionalJWTAuthMiddleware() func(http.Handler) http.Handler {
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

			ctx := context.WithValue(r.Context(), "user", claims.Username)
			ctx = context.WithValue(ctx, "user_id", claims.Subject)
			ctx = context.WithValue(ctx, "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole 基於角色的權限檢查中間件（為未來擴展預留）
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 首先確保用戶已通過 JWT 驗證
			user := r.Context().Value("user")
			if user == nil {
				responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "NOT_AUTHENTICATED", "User not authenticated"))
				return
			}

			// 這裡可以添加角色檢查邏輯
			// 目前只是示例，實際實現需要從資料庫或 token 中獲取用戶角色
			_ = user
			_ = requiredRole

			next.ServeHTTP(w, r)
		})
	}
}
