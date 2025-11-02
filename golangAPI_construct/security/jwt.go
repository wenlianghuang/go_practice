package security

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 定義系統使用的 JWT claims
type Claims struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// GenerateTokenWithClaims 生成包含角色與使用者代號的 JWT
func GenerateTokenWithClaims(username string, userID uint, roles []string, ttl time.Duration) (string, error) {
	now := time.Now()
	jti, err := generateJTI()
	if err != nil {
		return "", fmt.Errorf("failed to generate token id: %w", err)
	}

	subject := strconv.FormatUint(uint64(userID), 10)
	if subject == "0" {
		subject = username
	}

	claims := Claims{
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			ID:        jti,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret())
}

// GenerateToken 與舊版相容的函式
func GenerateToken(username string, ttl time.Duration) (string, error) {
	return GenerateTokenWithClaims(username, 0, nil, ttl)
}

// ValidateToken 驗證 JWT token 並返回 claims
func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret(), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// GetUsernameFromToken 從 token 中提取用戶名（便利函數）
func GetUsernameFromToken(tokenStr string) (string, error) {
	claims, err := ValidateToken(tokenStr)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}

// secret 取得 JWT 秘鑰
func secret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "dev-insecure-secret-change"
	}
	return []byte(s)
}

func generateJTI() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
