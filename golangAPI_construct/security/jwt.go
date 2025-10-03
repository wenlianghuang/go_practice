package security

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 統一的 Claims 結構，包含用戶名和標準 JWT 聲明
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 獲取 JWT 密鑰
func secret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "dev-insecure-secret-change"
	}
	return []byte(s)
}

// GenerateToken 生成 JWT token
func GenerateToken(username string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username, // 使用 username 作為 subject
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret())
}

// ValidateToken 驗證 JWT token 並返回 claims
func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// 驗證簽名方法
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
