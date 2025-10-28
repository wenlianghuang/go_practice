package database

import (
	"complex_breakpoint_example/models"
)

// DatabaseInterface 定義數據庫操作的通用接口
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
