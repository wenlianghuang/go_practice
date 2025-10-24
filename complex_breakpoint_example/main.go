package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

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

// 模擬數據庫
type Database struct {
	users            map[int]*User
	transactions     map[int]*Transaction
	bankAccounts     map[int]*BankAccount
	loanApplications map[int]*LoanApplication
	mu               sync.RWMutex
	nextID           int
}

// 全局數據庫實例
var db *Database

// 初始化數據庫
func initDatabase() {
	db = &Database{
		users:            make(map[int]*User),
		transactions:     make(map[int]*Transaction),
		bankAccounts:     make(map[int]*BankAccount),
		loanApplications: make(map[int]*LoanApplication),
		nextID:           1,
	}

	// 添加測試數據
	db.addTestData()
}

// 添加測試數據
func (d *Database) addTestData() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 添加用戶
	users := []*User{
		{ID: 1, Username: "alice", Email: "alice@example.com", Balance: 1000.0, IsActive: true},
		{ID: 2, Username: "bob", Email: "bob@example.com", Balance: 500.0, IsActive: true},
		{ID: 3, Username: "charlie", Email: "charlie@example.com", Balance: 2000.0, IsActive: false},
	}

	for _, user := range users {
		d.users[user.ID] = user
	}

	// 添加銀行帳戶
	accounts := []*BankAccount{
		{ID: 1, UserID: 1, Balance: 1000.0, Currency: "USD", IsLocked: false},
		{ID: 2, UserID: 2, Balance: 500.0, Currency: "USD", IsLocked: false},
		{ID: 3, UserID: 3, Balance: 2000.0, Currency: "USD", IsLocked: true},
	}

	for _, account := range accounts {
		d.bankAccounts[account.ID] = account
	}

	// 添加一些交易記錄
	transactions := []*Transaction{
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
func (d *Database) CreateUser(username, email string) (*User, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 檢查用戶名是否已存在
	for _, user := range d.users {
		if user.Username == username {
			return nil, fmt.Errorf("username already exists")
		}
	}

	user := &User{
		ID:       d.nextID,
		Username: username,
		Email:    email,
		Balance:  0.0,
		IsActive: true,
	}

	d.users[d.nextID] = user
	d.nextID++

	// 創建對應的銀行帳戶
	account := &BankAccount{
		ID:       d.nextID,
		UserID:   user.ID,
		Balance:  0.0,
		Currency: "USD",
		IsLocked: false,
	}

	d.bankAccounts[d.nextID] = account
	d.nextID++

	return user, nil
}

// 獲取用戶
func (d *Database) GetUser(id int) (*User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	user, exists := d.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// 獲取用戶的銀行帳戶
func (d *Database) GetUserAccount(userID int) (*BankAccount, error) {
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
func (d *Database) ProcessDeposit(userID int, amount float64, description string) (*Transaction, error) {
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
	var account *BankAccount
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
	transaction := &Transaction{
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
func (d *Database) ProcessWithdrawal(userID int, amount float64, description string) (*Transaction, error) {
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
	var account *BankAccount
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
	transaction := &Transaction{
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
func (d *Database) ProcessTransfer(fromUserID, toUserID int, amount float64, description string) (*Transaction, error) {
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
	var fromAccount *BankAccount
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
	var toAccount *BankAccount
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
	transaction := &Transaction{
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
func (d *Database) ApplyForLoan(userID int, amount float64, term int) (*LoanApplication, error) {
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
	application := &LoanApplication{
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
func (d *Database) processLoanApplication(application *LoanApplication) {
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
func (d *Database) GetUserTransactions(userID int) ([]*Transaction, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var transactions []*Transaction
	for _, tx := range d.transactions {
		if tx.UserID == userID {
			transactions = append(transactions, tx)
		}
	}

	return transactions, nil
}

// 獲取用戶的貸款申請
func (d *Database) GetUserLoanApplications(userID int) ([]*LoanApplication, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var applications []*LoanApplication
	for _, app := range d.loanApplications {
		if app.UserID == userID {
			applications = append(applications, app)
		}
	}

	return applications, nil
}

// HTTP 處理器
type Handler struct {
	db *Database
}

// 創建用戶
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.db.CreateUser(req.Username, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 獲取用戶
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 存款
func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	transaction, err := h.db.ProcessDeposit(userID, req.Amount, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// 提款
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	transaction, err := h.db.ProcessWithdrawal(userID, req.Amount, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// 轉帳
func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromUserID  int     `json:"from_user_id"`
		ToUserID    int     `json:"to_user_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	transaction, err := h.db.ProcessTransfer(req.FromUserID, req.ToUserID, req.Amount, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// 申請貸款
func (h *Handler) ApplyForLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
		Term   int     `json:"term"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	application, err := h.db.ApplyForLoan(userID, req.Amount, req.Term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(application)
}

// 獲取用戶交易記錄
func (h *Handler) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	transactions, err := h.db.GetUserTransactions(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// 獲取用戶貸款申請
func (h *Handler) GetUserLoanApplications(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	applications, err := h.db.GetUserLoanApplications(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applications)
}

// 獲取用戶帳戶信息
func (h *Handler) GetUserAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	account, err := h.db.GetUserAccount(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// 並發測試處理器
func (h *Handler) ConcurrentTest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    int     `json:"user_id"`
		Amount    float64 `json:"amount"`
		Operation string  `json:"operation"` // "deposit", "withdraw", "transfer"
		Count     int     `json:"count"`     // 並發操作次數
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Count <= 0 || req.Count > 100 {
		http.Error(w, "Count must be between 1 and 100", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []*Transaction
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

			var tx *Transaction
			var err error

			switch req.Operation {
			case "deposit":
				tx, err = h.db.ProcessDeposit(req.UserID, req.Amount, fmt.Sprintf("Concurrent deposit %d", index))
			case "withdraw":
				tx, err = h.db.ProcessWithdrawal(req.UserID, req.Amount, fmt.Sprintf("Concurrent withdrawal %d", index))
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

	response := struct {
		SuccessCount int            `json:"success_count"`
		ErrorCount   int            `json:"error_count"`
		Results      []*Transaction `json:"results"`
		Errors       []string       `json:"errors"`
	}{
		SuccessCount: len(results),
		ErrorCount:   len(errors),
		Results:      results,
		Errors:       make([]string, len(errors)),
	}

	for i, err := range errors {
		response.Errors[i] = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// 初始化數據庫
	initDatabase()

	// 創建處理器
	handler := &Handler{db: db}

	// 創建路由器
	router := mux.NewRouter()

	// API 路由
	api := router.PathPrefix("/api/v1").Subrouter()

	// 用戶相關路由
	api.HandleFunc("/users", handler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", handler.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}/account", handler.GetUserAccount).Methods("GET")
	api.HandleFunc("/users/{id}/transactions", handler.GetUserTransactions).Methods("GET")
	api.HandleFunc("/users/{id}/loans", handler.GetUserLoanApplications).Methods("GET")

	// 交易相關路由
	api.HandleFunc("/users/{id}/deposit", handler.Deposit).Methods("POST")
	api.HandleFunc("/users/{id}/withdraw", handler.Withdraw).Methods("POST")
	api.HandleFunc("/transfer", handler.Transfer).Methods("POST")
	api.HandleFunc("/users/{id}/apply-loan", handler.ApplyForLoan).Methods("POST")

	// 測試路由
	api.HandleFunc("/test/concurrent", handler.ConcurrentTest).Methods("POST")

	// 健康檢查
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	fmt.Println("🚀 Complex Breakpoint Example Server starting on :8080")
	fmt.Println("📚 Available endpoints:")
	fmt.Println("  POST /api/v1/users - Create user")
	fmt.Println("  GET  /api/v1/users/{id} - Get user")
	fmt.Println("  GET  /api/v1/users/{id}/account - Get user account")
	fmt.Println("  GET  /api/v1/users/{id}/transactions - Get user transactions")
	fmt.Println("  GET  /api/v1/users/{id}/loans - Get user loan applications")
	fmt.Println("  POST /api/v1/users/{id}/deposit - Deposit money")
	fmt.Println("  POST /api/v1/users/{id}/withdraw - Withdraw money")
	fmt.Println("  POST /api/v1/transfer - Transfer money")
	fmt.Println("  POST /api/v1/users/{id}/apply-loan - Apply for loan")
	fmt.Println("  POST /api/v1/test/concurrent - Concurrent operations test")
	fmt.Println("  GET  /health - Health check")

	log.Fatal(http.ListenAndServe(":8080", router))
}
