package database

import (
	"complex_breakpoint_example/models"
)

// GormAdapter 适配方 GORM 数据库接口到标准 Database 接口
type GormAdapter struct {
	gormDB *GormDatabase
}

// NewGormAdapter 创建 GORM 适配器
func NewGormAdapter(gormDB *GormDatabase) *GormAdapter {
	return &GormAdapter{
		gormDB: gormDB,
	}
}

// 实现 Database 接口
func (a *GormAdapter) CreateUser(username, email string, initialBalance float64) (*models.User, error) {
	return a.gormDB.CreateUser(username, email, initialBalance)
}

func (a *GormAdapter) GetUser(id uint) (*models.User, error) {
	return a.gormDB.GetUser(id)
}

func (a *GormAdapter) GetAllUsers() ([]*models.User, error) {
	return a.gormDB.GetAllUsers()
}

func (a *GormAdapter) GetUserAccount(userID uint) (*models.BankAccount, error) {
	return a.gormDB.GetUserAccount(userID)
}

func (a *GormAdapter) ProcessDeposit(userID uint, amount float64, description string) (*models.Transaction, error) {
	return a.gormDB.ProcessDeposit(userID, amount, description)
}

func (a *GormAdapter) ProcessWithdrawal(userID uint, amount float64, description string) (*models.Transaction, error) {
	return a.gormDB.ProcessWithdrawal(userID, amount, description)
}

func (a *GormAdapter) ProcessTransfer(fromUserID, toUserID uint, amount float64, description string) (*models.Transaction, error) {
	return a.gormDB.ProcessTransfer(fromUserID, toUserID, amount, description)
}

func (a *GormAdapter) GetUserTransactions(userID uint) ([]*models.Transaction, error) {
	return a.gormDB.GetUserTransactions(userID)
}

func (a *GormAdapter) ApplyForLoan(userID uint, amount float64, term int) (*models.LoanApplication, error) {
	return a.gormDB.ApplyForLoan(userID, amount, term)
}

func (a *GormAdapter) GetUserLoanApplications(userID uint) ([]*models.LoanApplication, error) {
	return a.gormDB.GetUserLoanApplications(userID)
}
