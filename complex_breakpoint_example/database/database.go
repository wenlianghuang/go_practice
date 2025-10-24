package database

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"complex_breakpoint_example/models"
)

// 模擬數據庫
type Database struct {
	users            map[int]*models.User
	transactions     map[int]*models.Transaction
	bankAccounts     map[int]*models.BankAccount
	loanApplications map[int]*models.LoanApplication
	mu               sync.RWMutex
	nextID           int
}

// 創建新的數據庫實例
func New() *Database {
	return &Database{
		users:            make(map[int]*models.User),
		transactions:     make(map[int]*models.Transaction),
		bankAccounts:     make(map[int]*models.BankAccount),
		loanApplications: make(map[int]*models.LoanApplication),
		nextID:           1,
	}
}

// 初始化數據庫
func (d *Database) Initialize() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 添加測試數據
	d.addTestData()
}

// 添加測試數據
func (d *Database) addTestData() {
	// 添加用戶
	users := []*models.User{
		{ID: 1, Username: "alice", Email: "alice@example.com", Balance: 1000.0, IsActive: true},
		{ID: 2, Username: "bob", Email: "bob@example.com", Balance: 500.0, IsActive: true},
		{ID: 3, Username: "charlie", Email: "charlie@example.com", Balance: 2000.0, IsActive: false},
	}

	for _, user := range users {
		d.users[user.ID] = user
	}

	// 添加銀行帳戶
	accounts := []*models.BankAccount{
		{ID: 1, UserID: 1, Balance: 1000.0, Currency: "USD", IsLocked: false},
		{ID: 2, UserID: 2, Balance: 500.0, Currency: "USD", IsLocked: false},
		{ID: 3, UserID: 3, Balance: 2000.0, Currency: "USD", IsLocked: true},
	}

	for _, account := range accounts {
		d.bankAccounts[account.ID] = account
	}

	// 添加一些交易記錄
	transactions := []*models.Transaction{
		{ID: 1, UserID: 1, Amount: 100.0, Type: "deposit", Description: "Initial deposit", Timestamp: time.Now().Add(-24 * time.Hour), Status: "completed"},
		{ID: 2, UserID: 2, Amount: 50.0, Type: "withdraw", Description: "ATM withdrawal", Timestamp: time.Now().Add(-12 * time.Hour), Status: "completed"},
		{ID: 3, UserID: 1, Amount: 200.0, Type: "transfer", Description: "Transfer to Bob", Timestamp: time.Now().Add(-6 * time.Hour), Status: "pending"},
	}

	for _, tx := range transactions {
		d.transactions[tx.ID] = tx
	}

	d.nextID = 4
}

// 創建用戶
func (d *Database) CreateUser(username, email string, initialBalance float64) (*models.User, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 檢查用戶名是否已存在
	for _, user := range d.users {
		if user.Username == username {
			return nil, fmt.Errorf("username already exists")
		}
	}

	user := &models.User{
		ID:       d.nextID,
		Username: username,
		Email:    email,
		Balance:  initialBalance,
		IsActive: true,
	}

	d.users[d.nextID] = user
	d.nextID++

	// 創建對應的銀行帳戶
	account := &models.BankAccount{
		ID:       d.nextID,
		UserID:   user.ID,
		Balance:  initialBalance,
		Currency: "USD",
		IsLocked: false,
	}

	d.bankAccounts[d.nextID] = account
	d.nextID++

	return user, nil
}

// 獲取用戶
func (d *Database) GetUser(id int) (*models.User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	user, exists := d.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// 獲取所有用戶
func (d *Database) GetAllUsers() ([]*models.User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var users []*models.User
	for _, user := range d.users {
		users = append(users, user)
	}

	return users, nil
}

// 獲取用戶的銀行帳戶
func (d *Database) GetUserAccount(userID int) (*models.BankAccount, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, account := range d.bankAccounts {
		if account.UserID == userID {
			return account, nil
		}
	}

	return nil, fmt.Errorf("account not found")
}

// 處理存款
func (d *Database) ProcessDeposit(userID int, amount float64, description string) (*models.Transaction, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 檢查用戶是否存在
	user, exists := d.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// 檢查用戶是否活躍
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// 獲取銀行帳戶
	var account *models.BankAccount
	for _, acc := range d.bankAccounts {
		if acc.UserID == userID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, fmt.Errorf("account not found")
	}

	// 檢查帳戶是否被鎖定
	if account.IsLocked {
		return nil, fmt.Errorf("account is locked")
	}

	// 檢查金額是否有效
	if amount <= 0 {
		return nil, fmt.Errorf("invalid amount")
	}

	// 模擬處理時間
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	// 創建交易記錄
	transaction := &models.Transaction{
		ID:          d.nextID,
		UserID:      userID,
		Amount:      amount,
		Type:        "deposit",
		Description: description,
		Timestamp:   time.Now(),
		Status:      "pending",
	}

	d.transactions[d.nextID] = transaction
	d.nextID++

	// 模擬隨機失敗
	if rand.Float64() < 0.1 { // 10% 失敗率
		transaction.Status = "failed"
		return transaction, fmt.Errorf("deposit processing failed")
	}

	// 更新餘額
	account.Balance += amount
	user.Balance += amount
	transaction.Status = "completed"

	return transaction, nil
}

// 處理提款
func (d *Database) ProcessWithdrawal(userID int, amount float64, description string) (*models.Transaction, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 檢查用戶是否存在
	user, exists := d.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// 檢查用戶是否活躍
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// 獲取銀行帳戶
	var account *models.BankAccount
	for _, acc := range d.bankAccounts {
		if acc.UserID == userID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, fmt.Errorf("account not found")
	}

	// 檢查帳戶是否被鎖定
	if account.IsLocked {
		return nil, fmt.Errorf("account is locked")
	}

	// 檢查金額是否有效
	if amount <= 0 {
		return nil, fmt.Errorf("invalid amount")
	}

	// 檢查餘額是否足夠
	if account.Balance < amount {
		return nil, fmt.Errorf("insufficient funds")
	}

	// 模擬處理時間
	time.Sleep(time.Duration(rand.Intn(150)) * time.Millisecond)

	// 創建交易記錄
	transaction := &models.Transaction{
		ID:          d.nextID,
		UserID:      userID,
		Amount:      amount,
		Type:        "withdraw",
		Description: description,
		Timestamp:   time.Now(),
		Status:      "pending",
	}

	d.transactions[d.nextID] = transaction
	d.nextID++

	// 模擬隨機失敗
	if rand.Float64() < 0.15 { // 15% 失敗率
		transaction.Status = "failed"
		return transaction, fmt.Errorf("withdrawal processing failed")
	}

	// 更新餘額
	account.Balance -= amount
	user.Balance -= amount
	transaction.Status = "completed"

	return transaction, nil
}

// 處理轉帳
func (d *Database) ProcessTransfer(fromUserID, toUserID int, amount float64, description string) (*models.Transaction, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 檢查發送方用戶
	fromUser, exists := d.users[fromUserID]
	if !exists {
		return nil, fmt.Errorf("sender user not found")
	}

	if !fromUser.IsActive {
		return nil, fmt.Errorf("sender account is inactive")
	}

	// 檢查接收方用戶
	toUser, exists := d.users[toUserID]
	if !exists {
		return nil, fmt.Errorf("recipient user not found")
	}

	if !toUser.IsActive {
		return nil, fmt.Errorf("recipient account is inactive")
	}

	// 獲取發送方帳戶
	var fromAccount *models.BankAccount
	for _, acc := range d.bankAccounts {
		if acc.UserID == fromUserID {
			fromAccount = acc
			break
		}
	}

	if fromAccount == nil {
		return nil, fmt.Errorf("sender account not found")
	}

	// 獲取接收方帳戶
	var toAccount *models.BankAccount
	for _, acc := range d.bankAccounts {
		if acc.UserID == toUserID {
			toAccount = acc
			break
		}
	}

	if toAccount == nil {
		return nil, fmt.Errorf("recipient account not found")
	}

	// 檢查帳戶是否被鎖定
	if fromAccount.IsLocked || toAccount.IsLocked {
		return nil, fmt.Errorf("one or both accounts are locked")
	}

	// 檢查金額是否有效
	if amount <= 0 {
		return nil, fmt.Errorf("invalid amount")
	}

	// 檢查餘額是否足夠
	if fromAccount.Balance < amount {
		return nil, fmt.Errorf("insufficient funds")
	}

	// 模擬處理時間
	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

	// 創建交易記錄
	transaction := &models.Transaction{
		ID:          d.nextID,
		UserID:      fromUserID,
		Amount:      amount,
		Type:        "transfer",
		Description: description,
		Timestamp:   time.Now(),
		Status:      "pending",
	}

	d.transactions[d.nextID] = transaction
	d.nextID++

	// 模擬隨機失敗
	if rand.Float64() < 0.2 { // 20% 失敗率
		transaction.Status = "failed"
		return transaction, fmt.Errorf("transfer processing failed")
	}

	// 更新餘額
	fromAccount.Balance -= amount
	toAccount.Balance += amount
	fromUser.Balance -= amount
	toUser.Balance += amount
	transaction.Status = "completed"

	return transaction, nil
}

// 申請貸款
func (d *Database) ApplyForLoan(userID int, amount float64, term int) (*models.LoanApplication, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 檢查用戶是否存在
	user, exists := d.users[userID]
	if !exists {
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

	// 模擬處理時間
	time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)

	// 創建貸款申請
	application := &models.LoanApplication{
		ID:           d.nextID,
		UserID:       userID,
		Amount:       amount,
		InterestRate: baseRate,
		Term:         term,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	d.loanApplications[d.nextID] = application
	d.nextID++

	// 模擬審核過程
	go d.processLoanApplication(application)

	return application, nil
}

// 處理貸款申請（並發處理）
func (d *Database) processLoanApplication(application *models.LoanApplication) {
	// 模擬審核時間
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)

	d.mu.Lock()
	defer d.mu.Unlock()

	// 獲取用戶信息
	user, exists := d.users[application.UserID]
	if !exists {
		application.Status = "rejected"
		return
	}

	// 簡單的審核邏輯
	if user.Balance > application.Amount*0.1 && application.InterestRate < 0.08 {
		application.Status = "approved"
		now := time.Now()
		application.ApprovedAt = &now
	} else {
		application.Status = "rejected"
	}
}

// 獲取用戶的所有交易
func (d *Database) GetUserTransactions(userID int) ([]*models.Transaction, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var transactions []*models.Transaction
	for _, tx := range d.transactions {
		if tx.UserID == userID {
			transactions = append(transactions, tx)
		}
	}

	return transactions, nil
}

// 獲取用戶的貸款申請
func (d *Database) GetUserLoanApplications(userID int) ([]*models.LoanApplication, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var applications []*models.LoanApplication
	for _, app := range d.loanApplications {
		if app.UserID == userID {
			applications = append(applications, app)
		}
	}

	return applications, nil
}
