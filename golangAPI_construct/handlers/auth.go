package handlers

import (
	"net/http"
	"time"

	"golangAPI_construct/responses"
	"golangAPI_construct/security"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      string `json:"user"`
}

// DEMO 用：硬編碼一個 bcrypt 雜湊（密碼為 "password"）
var demoUser = struct {
	Username     string
	PasswordHash string
}{
	Username:     "Matt",
	PasswordHash: "$2a$10$AQuMpFYbHBfGx2F2bS0.x.Nm.YTFzwjHaznp9uUCN9V5t3sweZ4w6", // 請用 security.HashPassword("password") 產生
}

// Login 處理用戶登入
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(responses.NewAppError(http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body"))
		return
	}

	// 驗證用戶名和密碼 - 修復函數名
	if req.Username != demoUser.Username || !security.CheckPassword(demoUser.PasswordHash, req.Password) {
		c.Error(responses.NewAppError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid username or password"))
		return
	}

	// 生成 JWT token
	ttl := 2 * time.Hour
	token, err := security.GenerateToken(req.Username, ttl)
	if err != nil {
		c.Error(responses.NewAppError(http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate token"))
		return
	}

	// 返回登入響應
	responses.Success(c, http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(ttl).Unix(),
		User:      req.Username,
	})
}

// Logout 處理用戶登出（目前只是示例，實際實現需要 token 黑名單）
func Logout(c *gin.Context) {
	// 在實際應用中，這裡應該將 token 加入黑名單
	// 或者使用 Redis 來管理 token 狀態

	responses.Success(c, http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

// RefreshToken 刷新 JWT token
func RefreshToken(c *gin.Context) {
	// 從 context 中獲取當前用戶信息
	user, exists := c.Get("user")
	if !exists {
		c.Error(responses.NewAppError(http.StatusUnauthorized, "NOT_AUTHENTICATED", "User not authenticated"))
		return
	}

	username, ok := user.(string)
	if !ok {
		c.Error(responses.NewAppError(http.StatusInternalServerError, "INVALID_USER_DATA", "Invalid user data"))
		return
	}

	// 生成新的 token
	ttl := 2 * time.Hour
	token, err := security.GenerateToken(username, ttl)
	if err != nil {
		c.Error(responses.NewAppError(http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate token"))
		return
	}

	responses.Success(c, http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(ttl).Unix(),
		User:      username,
	})
}
