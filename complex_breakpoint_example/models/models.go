package models

import (
	"time"

	"gorm.io/gorm"
)

// 用戶結構 - 兼容 GORM 和內存數據庫
type User struct {
	ID        uint           `gorm:"primarykey;column:id" json:"id"`
	Username  string         `gorm:"uniqueIndex;column:username" json:"username"`
	Email     string         `gorm:"uniqueIndex;column:email" json:"email"`
	Balance   float64        `gorm:"type:numeric(15,2);column:balance" json:"balance"`
	IsActive  bool           `gorm:"column:is_active" json:"is_active"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// 交易結構
type Transaction struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"index" json:"user_id"`
	ToUserID    uint           `gorm:"index" json:"to_user_id,omitempty"` // For transfer transactions
	Amount      float64        `gorm:"type:numeric(15,2)" json:"amount"`
	Type        string         `json:"type"` // "deposit", "withdraw", "transfer"
	Description string         `json:"description"`
	Timestamp   time.Time      `json:"timestamp"`
	Status      string         `json:"status"` // "pending", "completed", "failed"
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// 銀行帳戶結構
type BankAccount struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex" json:"user_id"`
	Balance   float64        `gorm:"type:numeric(15,2)" json:"balance"`
	Currency  string         `json:"currency"`
	IsLocked  bool           `json:"is_locked"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// 貸款申請結構
type LoanApplication struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `gorm:"index" json:"user_id"`
	Amount       float64        `gorm:"type:numeric(15,2)" json:"amount"`
	InterestRate float64        `gorm:"type:numeric(5,4)" json:"interest_rate"`
	Term         int            `json:"term"`   // 月數
	Status       string         `json:"status"` // "pending", "approved", "rejected"
	CreatedAt    time.Time      `json:"created_at"`
	ApprovedAt   *time.Time     `json:"approved_at,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// 並發測試請求結構
type ConcurrentTestRequest struct {
	UserID    int     `json:"user_id"`
	Amount    float64 `json:"amount"`
	Operation string  `json:"operation"` // "deposit", "withdraw", "transfer"
	Count     int     `json:"count"`     // 並發操作次數
}

// 並發測試響應結構
type ConcurrentTestResponse struct {
	SuccessCount int            `json:"success_count"`
	ErrorCount   int            `json:"error_count"`
	Results      []*Transaction `json:"results"`
	Errors       []string       `json:"errors"`
}

// Custom errors for testing
type ErrInsufficientFunds struct {
	Message string
}

func (e *ErrInsufficientFunds) Error() string {
	return "insufficient funds"
}

type ErrUserNotFound struct {
	UserID  int
	Message string
}

func (e *ErrUserNotFound) Error() string {
	return "user not found"
}
