package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"complex_breakpoint_example/database"
	"complex_breakpoint_example/models"
)

// 用戶服務
type UserService struct {
	db database.DatabaseInterface
}

func NewUserService(db interface{}) *UserService {
	// 类型断言
	var dbInterface database.DatabaseInterface
	switch v := db.(type) {
	case *database.Database:
		dbInterface = v
	case *database.GormAdapter:
		dbInterface = v
	default:
		panic("unsupported database type")
	}

	return &UserService{
		db: dbInterface,
	}
}

// 創建用戶
func (s *UserService) CreateUser(username, email string, initialBalance float64) (*models.User, error) {
	return s.db.CreateUser(username, email, initialBalance)
}

// 獲取用戶
func (s *UserService) GetUser(id int) (*models.User, error) {
	return s.db.GetUser(uint(id))
}

// 獲取所有用戶
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.db.GetAllUsers()
}

// 獲取用戶帳戶
func (s *UserService) GetUserAccount(userID int) (*models.BankAccount, error) {
	return s.db.GetUserAccount(uint(userID))
}

// 交易服務
type TransactionService struct {
	db database.DatabaseInterface
}

func NewTransactionService(db interface{}) *TransactionService {
	var dbInterface database.DatabaseInterface
	switch v := db.(type) {
	case *database.Database:
		dbInterface = v
	case *database.GormAdapter:
		dbInterface = v
	default:
		panic("unsupported database type")
	}

	return &TransactionService{
		db: dbInterface,
	}
}

// 存款
func (s *TransactionService) Deposit(userID int, amount float64, description string) (*models.Transaction, error) {
	return s.db.ProcessDeposit(uint(userID), amount, description)
}

// 提款
func (s *TransactionService) Withdraw(userID int, amount float64, description string) (*models.Transaction, error) {
	return s.db.ProcessWithdrawal(uint(userID), amount, description)
}

// 轉帳
func (s *TransactionService) Transfer(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
	return s.db.ProcessTransfer(uint(fromUserID), uint(toUserID), amount, description)
}

// 獲取用戶交易記錄
func (s *TransactionService) GetUserTransactions(userID int) ([]*models.Transaction, error) {
	return s.db.GetUserTransactions(uint(userID))
}

// 並發測試
func (s *TransactionService) ConcurrentTest(req *models.ConcurrentTestRequest) (*models.ConcurrentTestResponse, error) {
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
type LoanService struct {
	db database.DatabaseInterface
}

func NewLoanService(db interface{}) *LoanService {
	var dbInterface database.DatabaseInterface
	switch v := db.(type) {
	case *database.Database:
		dbInterface = v
	case *database.GormAdapter:
		dbInterface = v
	default:
		panic("unsupported database type")
	}

	return &LoanService{
		db: dbInterface,
	}
}

// 申請貸款
func (s *LoanService) ApplyForLoan(userID int, amount float64, term int) (*models.LoanApplication, error) {
	return s.db.ApplyForLoan(uint(userID), amount, term)
}

// 獲取用戶貸款申請
func (s *LoanService) GetUserLoanApplications(userID int) ([]*models.LoanApplication, error) {
	return s.db.GetUserLoanApplications(uint(userID))
}
