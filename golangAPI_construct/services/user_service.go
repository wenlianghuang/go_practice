package services

import (
	"strings"

	"golangAPI_construct/models"

	"gorm.io/gorm"
)

// UserService 提供與使用者相關的資料庫操作
type UserService struct {
	db *gorm.DB
}

// NewUserService 建立新的使用者服務實例
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// FindByUsername 依使用者名稱取得使用者
func (s *UserService) FindByUsername(username string) (*models.UserGORM, error) {
	var user models.UserGORM
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 依電子郵件取得使用者
func (s *UserService) FindByEmail(email string) (*models.UserGORM, error) {
	var user models.UserGORM
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 建立新的使用者
func (s *UserService) CreateUser(user models.UserGORM) (*models.UserGORM, error) {
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ParseRoles 解析資料庫中儲存的角色欄位
func ParseRoles(roles string) []string {
	if roles == "" {
		return []string{}
	}
	chunks := strings.Split(roles, ",")
	result := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		clean := strings.TrimSpace(chunk)
		if clean != "" {
			result = append(result, clean)
		}
	}
	return result
}
