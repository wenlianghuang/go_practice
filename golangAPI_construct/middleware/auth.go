package middleware

import (
	"net/http"
	"strings"

	"golangAPI_construct/responses"
	"golangAPI_construct/security"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 統一的 JWT 驗證中間件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 獲取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(responses.NewAppError(http.StatusUnauthorized, "MISSING_AUTH_HEADER", "Authorization header is required"))
			c.Abort()
			return
		}

		// 檢查 Bearer token 格式
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Error(responses.NewAppError(http.StatusUnauthorized, "INVALID_AUTH_FORMAT", "Authorization header must start with 'Bearer '"))
			c.Abort()
			return
		}

		// 提取 token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.Error(responses.NewAppError(http.StatusUnauthorized, "EMPTY_TOKEN", "Token cannot be empty"))
			c.Abort()
			return
		}

		// 驗證 token
		claims, err := security.ValidateToken(tokenString)
		if err != nil {
			c.Error(responses.NewAppError(http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token"))
			c.Abort()
			return
		}

		// 將用戶信息存儲到 context 中
		c.Set("user", claims.Username)
		c.Set("user_id", claims.Subject)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalJWTAuthMiddleware 可選的 JWT 驗證中間件（用於不需要強制登入的端點）
func OptionalJWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.Next()
			return
		}

		claims, err := security.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user", claims.Username)
		c.Set("user_id", claims.Subject)
		c.Set("claims", claims)
		c.Next()
	}
}

// RequireRole 基於角色的權限檢查中間件（為未來擴展預留）
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先確保用戶已通過 JWT 驗證
		user, exists := c.Get("user")
		if !exists {
			c.Error(responses.NewAppError(http.StatusUnauthorized, "NOT_AUTHENTICATED", "User not authenticated"))
			c.Abort()
			return
		}

		// 這裡可以添加角色檢查邏輯
		// 目前只是示例，實際實現需要從資料庫或 token 中獲取用戶角色
		_ = user
		_ = requiredRole

		c.Next()
	}
}
