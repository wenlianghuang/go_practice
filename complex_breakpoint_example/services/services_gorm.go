package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"complex_breakpoint_example/models"
)

// Database interface for both in-memory and GORM implementations
type DatabaseInterface interface {
	CreateUser(username, email string, initialBalance float64) (*models.User, error)
	GetUser(id uint) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetUserAccount(userID uint) (*models.BankAccount, error)
	ProcessDeposit(userID uint, amount float64, description string) (*models.Transaction, error)
	ProcessWithdrawal(userID uint, amount float64, description string) (*models.Transaction, error)
	ProcessTransfer(fromUserID, toUserID uint, amount float64, description string) (*models.Transaction, error)
	GetUserTransactions(userID uint) ([]*models.Transaction, error)
	ApplyForLoan(userID uint, amount float64, term int) (*models.LoanApplication, error)
	GetUserLoanApplications(userID uint) ([]*models.LoanApplication, error)
}

// 用戶服務
type UserServiceGorm struct {
	db DatabaseInterface
}

func NewUserServiceGorm(db DatabaseInterface) *UserServiceGorm {
	return &UserServiceGorm{
		db: db,
	}
}

// 創建用戶
func (s *UserServiceGorm) CreateUser(username, email string, initialBalance float64) (*models.User, error) {
	return s.db.CreateUser(username, email, initialBalance)
}

// 獲取用戶
func (s *UserServiceGorm) GetUser(id int) (*models.User, error) {
	return s.db.GetUser(uint(id))
}

// 獲取所有用戶
func (s *UserServiceGorm) GetAllUsers() ([]*models.User, error) {
	return s.db.GetAllUsers()
}

// 獲取用戶帳戶
func (s *UserServiceGorm) GetUserAccount(userID int) (*models.BankAccount, error) {
	return s.db.GetUserAccount(uint(userID))
}

// 交易服務
type TransactionServiceGorm struct {
	db DatabaseInterface
}

func NewTransactionServiceGorm(db DatabaseInterface) *TransactionServiceGorm {
	return &TransactionServiceGorm{
		db: db,
	}
}

// 存款
func (s *TransactionServiceGorm) Deposit(userID int, amount float64, description string) (*models.Transaction, error) {
	return s.db.ProcessDeposit(uint(userID), amount, description)
}

// 提款
func (s *TransactionServiceGorm) Withdraw(userID int, amount float64, description string) (*models.Transaction, error) {
	return s.db.ProcessWithdrawal(uint(userID), amount, description)
}

// 轉帳
func (s *TransactionServiceGorm) Transfer(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
	return s.db.ProcessTransfer(uint(fromUserID), uint(toUserID), amount, description)
}

// 獲取用戶交易記錄
func (s *TransactionServiceGorm) GetUserTransactions(userID int) ([]*models.Transaction, error) {
	return s.db.GetUserTransactions(uint(userID))
}

// 並發測試
func (s *TransactionServiceGorm) ConcurrentTest(req *models.ConcurrentTestRequest) (*models.ConcurrentTestResponse, error) {
	if req.Count <= 0 || req.Count > 100 {
		return nil, fmt.Errorf("count must be between 1 and 100")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []*models.Transaction
	var errors []error

	// 創建上下文用於超時控制
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i := 0; i < req.Count; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				mu.Lock()
				errors = append(errors, fmt.Errorf("operation %d timed out", index))
				mu.Unlock()
				return
			default:
			}

			var tx *models.Transaction
			var err error

			switch req.Operation {
			case "deposit":
				tx, err = s.db.ProcessDeposit(uint(req.UserID), req.Amount, fmt.Sprintf("Concurrent deposit %d", index))
			case "withdraw":
				tx, err = s.db.ProcessWithdrawal(uint(req.UserID), req.Amount, fmt.Sprintf("Concurrent withdrawal %d", index))
			default:
				err = fmt.Errorf("unsupported operation: %s", req.Operation)
			}

			mu.Lock()
			if err != nil {
				errors = append(errors, err)
			} else {
				results = append(results, tx)
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	response := &models.ConcurrentTestResponse{
		SuccessCount: len(results),
		ErrorCount:   len(errors),
		Results:      results,
		Errors:       make([]string, len(errors)),
	}

	for i, err := range errors {
		response.Errors[i] = err.Error()
	}

	return response, nil
}

// 貸款服務
type LoanServiceGorm struct {
	db DatabaseInterface
}

func NewLoanServiceGorm(db DatabaseInterface) *LoanServiceGorm {
	return &LoanServiceGorm{
		db: db,
	}
}

// 申請貸款
func (s *LoanServiceGorm) ApplyForLoan(userID int, amount float64, term int) (*models.LoanApplication, error) {
	return s.db.ApplyForLoan(uint(userID), amount, term)
}

// 獲取用戶貸款申請
func (s *LoanServiceGorm) GetUserLoanApplications(userID int) ([]*models.LoanApplication, error) {
	return s.db.GetUserLoanApplications(uint(userID))
}
