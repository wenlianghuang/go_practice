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

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(responses.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid request body"))
		return
	}
	// DEMO 用：請改為真正的驗證 (DB / LDAP ...)
	if req.Username != "admin" || req.Password != "password" {
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
