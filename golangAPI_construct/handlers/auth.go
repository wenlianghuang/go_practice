package handlers

import (
	"net/http"
	"time"

	"golangAPI_construct/responses"
	"golangAPI_construct/security"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DEMO 用：硬編碼一個 bcrypt 雜湊（密碼為 "password"）
var demoUser = struct {
	Username     string
	PasswordHash string
}{
	Username:     "Matt",
	PasswordHash: "$2a$10$AQuMpFYbHBfGx2F2bS0.x.Nm.YTFzwjHaznp9uUCN9V5t3sweZ4w6", // 請用 security.HashPassword("password") 產生
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(responses.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid request body"))
		return
	}
	if req.Username != demoUser.Username || !security.CheckPassword(demoUser.PasswordHash, req.Password) {
		c.Error(responses.NewAppError(http.StatusUnauthorized, "BAD_CREDENTIALS", "invalid username or password"))
		return
	}
	ttl := 2 * time.Hour
	token, err := security.GenerateToken(req.Username, ttl)
	if err != nil {
		c.Error(responses.NewAppError(http.StatusInternalServerError, "TOKEN_ISSUE", "cannot generate token"))
		return
	}
	responses.Success(c, http.StatusOK, gin.H{
		"token": token,
		"exp":   time.Now().Add(ttl).Unix(),
	})
}
