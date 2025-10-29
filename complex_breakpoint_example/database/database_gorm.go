package database

import (
	"fmt"
	"time"

	"complex_breakpoint_example/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 使用 GORM 的數據庫
type GormDatabase struct {
	db *gorm.DB
}

// 創建新的 GORM 數據庫實例
func NewGormDatabase(databaseURL string) (*GormDatabase, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	gdb := &GormDatabase{db: db}

	// 自動遷移
	// 為了讓資料庫表結構自動與 model 結構保持同步，減少人為遺漏或 migrations 步驟過程中的錯誤
	err = gdb.AutoMigrate()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return gdb, nil
}

// 自動遷移表結構
func (d *GormDatabase) AutoMigrate() error {
	return d.db.AutoMigrate(
		&models.User{},
		&models.Transaction{},
		&models.BankAccount{},
		&models.LoanApplication{},
	)
}

// 初始化數據庫
func (d *GormDatabase) Initialize() error {
	// 檢查是否已有數據
	var count int64
	d.db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return nil // 已有數據，不重複初始化
	}

	// 添加測試數據
	return d.addTestData()
}

// 添加測試數據
func (d *GormDatabase) addTestData() error {
	// 創建用戶
	users := []models.User{
		{Username: "alice", Email: "alice@example.com", Balance: 1000.0, IsActive: true},
		{Username: "bob", Email: "bob@example.com", Balance: 500.0, IsActive: true},
		{Username: "charlie", Email: "charlie@example.com", Balance: 2000.0, IsActive: false},
	}

	for i := range users {
		if err := d.db.Create(&users[i]).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// 創建對應的銀行帳戶
		account := models.BankAccount{
			UserID:   users[i].ID,
			Balance:  users[i].Balance,
			Currency: "USD",
			IsLocked: false,
		}
		if err := d.db.Create(&account).Error; err != nil {
			return fmt.Errorf("failed to create bank account: %w", err)
		}

		// 添加一些交易記錄
		if i < 3 {
			tx := models.Transaction{
				UserID:      users[i].ID,
				Amount:      100.0,
				Type:        "deposit",
				Description: "Initial deposit",
				Timestamp:   time.Now().Add(-24 * time.Hour),
				Status:      "completed",
			}
			if err := d.db.Create(&tx).Error; err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}
		}
	}

	return nil
}

// 創建用戶
func (d *GormDatabase) CreateUser(username, email string, initialBalance float64) (*models.User, error) {
	// 檢查用戶名是否已存在
	var existingUser models.User
	result := d.db.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		return nil, fmt.Errorf("username already exists")
	} else if result.Error != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", result.Error)
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Balance:  initialBalance,
		IsActive: true,
	}

	if err := d.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 創建對應的銀行帳戶
	account := &models.BankAccount{
		UserID:   user.ID,
		Balance:  initialBalance,
		Currency: "USD",
		IsLocked: false,
	}

	if err := d.db.Create(account).Error; err != nil {
		return nil, fmt.Errorf("failed to create bank account: %w", err)
	}

	return user, nil
}

// 獲取用戶
func (d *GormDatabase) GetUser(id uint) (*models.User, error) {
	var user models.User
	result := d.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}
	return &user, nil
}

// 獲取所有用戶
func (d *GormDatabase) GetAllUsers() ([]*models.User, error) {
	var users []models.User
	if err := d.db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	result := make([]*models.User, len(users))
	for i := range users {
		result[i] = &users[i]
	}
	return result, nil
}

// 獲取用戶的銀行帳戶
func (d *GormDatabase) GetUserAccount(userID uint) (*models.BankAccount, error) {
	var account models.BankAccount
	result := d.db.Where("user_id = ?", userID).First(&account)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get bank account: %w", result.Error)
	}
	return &account, nil
}

// 處理存款
func (d *GormDatabase) ProcessDeposit(userID uint, amount float64, description string) (*models.Transaction, error) {
	// 開始事務
	tx := d.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 檢查用戶是否存在
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("user not found")
	}

	// 檢查用戶是否活躍
	if !user.IsActive {
		tx.Rollback()
		return nil, fmt.Errorf("user account is inactive")
	}

	// 獲取銀行帳戶
	var account models.BankAccount
	if err := tx.Where("user_id = ?", userID).First(&account).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("account not found")
	}

	// 檢查帳戶是否被鎖定
	if account.IsLocked {
		tx.Rollback()
		return nil, fmt.Errorf("account is locked")
	}

	// 檢查金額是否有效
	if amount <= 0 {
		tx.Rollback()
		return nil, fmt.Errorf("invalid amount")
	}

	// 創建交易記錄
	transaction := &models.Transaction{
		UserID:      userID,
		Amount:      amount,
		Type:        "deposit",
		Description: description,
		Timestamp:   time.Now(),
		Status:      "completed",
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 更新餘額
	account.Balance += amount
	user.Balance += amount

	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// 處理提款
func (d *GormDatabase) ProcessWithdrawal(userID uint, amount float64, description string) (*models.Transaction, error) {
	// 開始事務
	tx := d.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 檢查用戶是否存在
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("user not found")
	}

	// 檢查用戶是否活躍
	if !user.IsActive {
		tx.Rollback()
		return nil, fmt.Errorf("user account is inactive")
	}

	// 獲取銀行帳戶
	var account models.BankAccount
	if err := tx.Where("user_id = ?", userID).First(&account).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("account not found")
	}

	// 檢查帳戶是否被鎖定
	if account.IsLocked {
		tx.Rollback()
		return nil, fmt.Errorf("account is locked")
	}

	// 檢查金額是否有效
	if amount <= 0 {
		tx.Rollback()
		return nil, fmt.Errorf("invalid amount")
	}

	// 檢查餘額是否足夠
	if account.Balance < amount {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient funds")
	}

	// 創建交易記錄
	transaction := &models.Transaction{
		UserID:      userID,
		Amount:      amount,
		Type:        "withdraw",
		Description: description,
		Timestamp:   time.Now(),
		Status:      "completed",
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 更新餘額
	account.Balance -= amount
	user.Balance -= amount

	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// 處理轉帳
func (d *GormDatabase) ProcessTransfer(fromUserID, toUserID uint, amount float64, description string) (*models.Transaction, error) {
	// 開始事務
	tx := d.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 檢查發送方用戶
	var fromUser models.User
	if err := tx.First(&fromUser, fromUserID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("sender user not found")
	}

	if !fromUser.IsActive {
		tx.Rollback()
		return nil, fmt.Errorf("sender account is inactive")
	}

	// 檢查接收方用戶
	var toUser models.User
	if err := tx.First(&toUser, toUserID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("recipient user not found")
	}

	if !toUser.IsActive {
		tx.Rollback()
		return nil, fmt.Errorf("recipient account is inactive")
	}

	// 獲取發送方帳戶
	var fromAccount models.BankAccount
	if err := tx.Where("user_id = ?", fromUserID).First(&fromAccount).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("sender account not found")
	}

	// 獲取接收方帳戶
	var toAccount models.BankAccount
	if err := tx.Where("user_id = ?", toUserID).First(&toAccount).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("recipient account not found")
	}

	// 檢查帳戶是否被鎖定
	if fromAccount.IsLocked || toAccount.IsLocked {
		tx.Rollback()
		return nil, fmt.Errorf("one or both accounts are locked")
	}

	// 檢查金額是否有效
	if amount <= 0 {
		tx.Rollback()
		return nil, fmt.Errorf("invalid amount")
	}

	// 檢查餘額是否足夠
	if fromAccount.Balance < amount {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient funds")
	}

	// 創建交易記錄
	transaction := &models.Transaction{
		UserID:      fromUserID,
		ToUserID:    toUserID,
		Amount:      amount,
		Type:        "transfer",
		Description: description,
		Timestamp:   time.Now(),
		Status:      "completed",
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 更新餘額
	fromAccount.Balance -= amount
	toAccount.Balance += amount
	fromUser.Balance -= amount
	toUser.Balance += amount

	if err := tx.Save(&fromAccount).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update sender account: %w", err)
	}

	if err := tx.Save(&toAccount).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update recipient account: %w", err)
	}

	if err := tx.Save(&fromUser).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update sender balance: %w", err)
	}

	if err := tx.Save(&toUser).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update recipient balance: %w", err)
	}

	// 提交事務
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// 申請貸款
func (d *GormDatabase) ApplyForLoan(userID uint, amount float64, term int) (*models.LoanApplication, error) {
	// 檢查用戶是否存在
	var user models.User
	if err := d.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 檢查用戶是否活躍
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// 檢查金額是否有效
	if amount <= 0 {
		return nil, fmt.Errorf("invalid loan amount")
	}

	// 檢查期限是否有效
	if term <= 0 || term > 60 {
		return nil, fmt.Errorf("invalid loan term")
	}

	// 計算利率（基於用戶餘額和期限）
	baseRate := 0.05 // 5% 基礎利率
	if user.Balance > 1000 {
		baseRate -= 0.01 // 高餘額用戶享受優惠
	}
	if term > 24 {
		baseRate += 0.02 // 長期貸款利率更高
	}

	// 創建貸款申請
	application := &models.LoanApplication{
		UserID:       userID,
		Amount:       amount,
		InterestRate: baseRate,
		Term:         term,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	if err := d.db.Create(application).Error; err != nil {
		return nil, fmt.Errorf("failed to create loan application: %w", err)
	}

	// 在後台處理貸款審核
	go d.processLoanApplication(application)

	return application, nil
}

// 處理貸款申請（並發處理）
func (d *GormDatabase) processLoanApplication(application *models.LoanApplication) {
	// 模擬審核時間
	time.Sleep(2 * time.Second)

	var user models.User
	if err := d.db.First(&user, application.UserID).Error; err != nil {
		application.Status = "rejected"
		d.db.Save(application)
		return
	}

	// 簡單的審核邏輯
	if user.Balance > application.Amount*0.1 && application.InterestRate < 0.08 {
		now := time.Now()
		application.Status = "approved"
		application.ApprovedAt = &now
	} else {
		application.Status = "rejected"
	}

	d.db.Save(application)
}

// 獲取用戶的所有交易
func (d *GormDatabase) GetUserTransactions(userID uint) ([]*models.Transaction, error) {
	var transactions []models.Transaction
	result := d.db.Where("user_id = ?", userID).Find(&transactions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", result.Error)
	}

	txns := make([]*models.Transaction, len(transactions))
	for i := range transactions {
		txns[i] = &transactions[i]
	}
	return txns, nil
}

// 獲取用戶的貸款申請
func (d *GormDatabase) GetUserLoanApplications(userID uint) ([]*models.LoanApplication, error) {
	var applications []models.LoanApplication
	result := d.db.Where("user_id = ?", userID).Find(&applications)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get loan applications: %w", result.Error)
	}

	apps := make([]*models.LoanApplication, len(applications))
	for i := range applications {
		apps[i] = &applications[i]
	}
	return apps, nil
}
