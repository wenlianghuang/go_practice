package models

import "time"

// 用戶結構
type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
	IsActive bool    `json:"is_active"`
}

// 交易結構
type Transaction struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"` // "deposit", "withdraw", "transfer"
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"` // "pending", "completed", "failed"
}

// 銀行帳戶結構
type BankAccount struct {
	ID       int     `json:"id"`
	UserID   int     `json:"user_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	IsLocked bool    `json:"is_locked"`
}

// 貸款申請結構
type LoanApplication struct {
	ID           int        `json:"id"`
	UserID       int        `json:"user_id"`
	Amount       float64    `json:"amount"`
	InterestRate float64    `json:"interest_rate"`
	Term         int        `json:"term"`   // 月數
	Status       string     `json:"status"` // "pending", "approved", "rejected"
	CreatedAt    time.Time  `json:"created_at"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
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
