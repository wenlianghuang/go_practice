# 🎯 重構後的 Go Breakpoint 調試練習指南

## 📋 項目結構

重構後的項目採用了分層架構，更符合實際項目的最佳實踐：

```
complex_breakpoint_example/
├── main.go                 # 程序入口
├── config/
│   └── config.go          # 配置管理
├── models/
│   └── models.go          # 數據模型定義
├── database/
│   └── database.go        # 數據庫操作
├── handlers/
│   └── handlers.go        # HTTP 處理器
├── services/
│   └── services.go        # 業務邏輯服務
├── middleware/
│   └── middleware.go      # 中間件
├── routes/
│   └── routes.go          # 路由配置
└── tests/                 # 測試文件
```

## 🎯 分層調試練習

### 🟢 初級練習 - Handler 層調試

#### 案例 1: 用戶創建 Handler 調試
**目標**: 學習 HTTP 請求處理的調試

**斷點設置**:
```
handlers/handlers.go:25  - 檢查請求解析
handlers/handlers.go:32  - 檢查服務調用
handlers/handlers.go:39  - 檢查響應構建
```

**對應代碼**:
```go
// 行 25 - 檢查請求解析
var req struct {
    Username string  `json:"username"`
    Email    string  `json:"email"`
    Balance  float64 `json:"balance"`
}

// 行 32 - 檢查服務調用
user, err := h.userService.CreateUser(req.Username, req.Email, req.Balance)

// 行 39 - 檢查響應構建
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(user)
```

**Watch 表達式**:
```
req.Username
req.Email
req.Balance
user
err
err != nil
```

**測試步驟**:
```
POST http://localhost:9090/api/v1/users
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "balance": 500.0
}
```

---

#### 案例 2: 存款 Handler 調試
**目標**: 理解複雜 HTTP 處理的調試

**斷點設置**:
```
handlers/handlers.go:85  - 檢查用戶 ID 解析
handlers/handlers.go:95  - 檢查請求體解析
handlers/handlers.go:105 - 檢查服務調用
handlers/handlers.go:112 - 檢查響應構建
```

**Watch 表達式**:
```
userID
req.Amount
req.Description
transaction
err
```

**測試步驟**:
```
POST http://localhost:9090/api/v1/users/1/deposit
Content-Type: application/json

{
  "amount": 100.0,
  "description": "Test deposit"
}
```

---

### 🟡 中級練習 - Service 層調試

#### 案例 3: 用戶服務調試
**目標**: 學習業務邏輯層的調試

**斷點設置**:
```
services/services.go:20  - 檢查服務方法調用
services/services.go:25  - 檢查數據庫調用
```

**對應代碼**:
```go
// 行 20 - 檢查服務方法調用
func (s *UserService) CreateUser(username, email string, initialBalance float64) (*models.User, error) {
    return s.db.CreateUser(username, email, initialBalance)
}

// 行 25 - 檢查數據庫調用
func (s *UserService) GetUser(id int) (*models.User, error) {
    return s.db.GetUser(id)
}
```

**Watch 表達式**:
```
username
email
initialBalance
s.db
```

**調試重點**:
- 觀察服務層如何調用數據庫層
- 理解業務邏輯的執行流程
- 檢查參數傳遞的正確性

---

#### 案例 4: 交易服務調試
**目標**: 理解複雜業務邏輯的調試

**斷點設置**:
```
services/services.go:45  - 檢查存款服務調用
services/services.go:50  - 檢查提款服務調用
services/services.go:55  - 檢查轉帳服務調用
```

**Watch 表達式**:
```
userID
amount
description
fromUserID
toUserID
```

**測試步驟**:
```
POST http://localhost:9090/api/v1/users/1/deposit
POST http://localhost:9090/api/v1/users/1/withdraw
POST http://localhost:9090/api/v1/transfer
```

---

### 🔴 高級練習 - Database 層調試

#### 案例 5: 數據庫操作調試
**目標**: 學習數據層的調試

**斷點設置**:
```
database/database.go:45  - 檢查用戶創建
database/database.go:65  - 檢查用戶名驗證
database/database.go:75  - 檢查用戶對象創建
database/database.go:85  - 檢查銀行帳戶創建
```

**對應代碼**:
```go
// 行 45 - 檢查用戶創建
func (d *Database) CreateUser(username, email string, initialBalance float64) (*models.User, error) {
    d.mu.Lock()
    defer d.mu.Unlock()

// 行 65 - 檢查用戶名驗證
for _, user := range d.users {
    if user.Username == username {
        return nil, fmt.Errorf("username already exists")
    }
}

// 行 75 - 檢查用戶對象創建
user := &models.User{
    ID:       d.nextID,
    Username: username,
    Email:    email,
    Balance:  initialBalance,
    IsActive: true,
}

// 行 85 - 檢查銀行帳戶創建
account := &models.BankAccount{
    ID:       d.nextID,
    UserID:   user.ID,
    Balance:  initialBalance,
    Currency: "USD",
    IsLocked: false,
}
```

**Watch 表達式**:
```
username
email
initialBalance
d.nextID
d.users
d.bankAccounts
user
account
```

**調試重點**:
- 觀察數據庫鎖的使用
- 理解數據一致性的維護
- 檢查 ID 的分配和遞增

---

#### 案例 6: 並發操作調試
**目標**: 學習並發數據操作的調試

**斷點設置**:
```
database/database.go:120 - 檢查存款處理
database/database.go:140 - 檢查餘額更新
database/database.go:200 - 檢查轉帳處理
database/database.go:250 - 檢查轉帳餘額更新
```

**Watch 表達式**:
```
d.mu
user.Balance
account.Balance
fromAccount.Balance
toAccount.Balance
transaction.Status
```

**調試重點**:
- 觀察互斥鎖的使用
- 理解原子操作的執行
- 檢查數據競爭的避免

---

### 🔧 中間件調試練習

#### 案例 7: 日誌中間件調試
**目標**: 學習中間件的調試

**斷點設置**:
```
middleware/middleware.go:12 - 檢查請求開始
middleware/middleware.go:18 - 檢查請求完成
```

**對應代碼**:
```go
// 行 12 - 檢查請求開始
log.Printf("🚀 Started %s %s", r.Method, r.URL.Path)

// 行 18 - 檢查請求完成
log.Printf("✅ Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
```

**Watch 表達式**:
```
r.Method
r.URL.Path
start
time.Since(start)
```

**調試重點**:
- 觀察中間件的執行順序
- 理解請求生命週期
- 檢查性能監控

---

#### 案例 8: CORS 中間件調試
**目標**: 學習跨域處理的調試

**斷點設置**:
```
middleware/middleware.go:25 - 檢查 CORS 頭設置
middleware/middleware.go:35 - 檢查預檢請求處理
```

**Watch 表達式**:
```
w.Header()
r.Method
r.Header.Get("Origin")
```

---

## 🎯 調試技巧總結

### 1. **分層調試策略**
- **Handler 層**: 關注 HTTP 請求處理
- **Service 層**: 關注業務邏輯執行
- **Database 層**: 關注數據操作和一致性
- **Middleware 層**: 關注橫切關注點

### 2. **斷點設置原則**
- 在函數入口設置斷點
- 在關鍵業務邏輯處設置斷點
- 在錯誤處理處設置斷點
- 在數據轉換處設置斷點

### 3. **Watch 表達式設計**
- 監視輸入參數
- 監視中間變量
- 監視狀態變化
- 監視錯誤條件

### 4. **調試流程**
1. 設置斷點和 Watch 表達式
2. 啟動調試模式
3. 發送測試請求
4. 觀察變量變化
5. 單步執行分析
6. 檢查錯誤處理

## 🎮 實際操作步驟

### 1. **Handler 層調試**
```bash
# 設置斷點在 handlers/handlers.go
# 發送 HTTP 請求
# 觀察請求解析和響應構建
```

### 2. **Service 層調試**
```bash
# 設置斷點在 services/services.go
# 觀察業務邏輯執行
# 檢查服務間的調用
```

### 3. **Database 層調試**
```bash
# 設置斷點在 database/database.go
# 觀察數據操作
# 檢查鎖的使用
```

### 4. **Middleware 層調試**
```bash
# 設置斷點在 middleware/middleware.go
# 觀察中間件執行
# 檢查橫切關注點
```

## 📊 調試檢查清單

### Handler 層
- [ ] 請求解析正確
- [ ] 參數驗證通過
- [ ] 服務調用成功
- [ ] 響應構建正確
- [ ] 錯誤處理適當

### Service 層
- [ ] 業務邏輯正確
- [ ] 參數傳遞正確
- [ ] 數據庫調用成功
- [ ] 錯誤處理適當
- [ ] 返回值正確

### Database 層
- [ ] 鎖使用正確
- [ ] 數據一致性
- [ ] 事務處理正確
- [ ] 錯誤處理適當
- [ ] 性能良好

### Middleware 層
- [ ] 執行順序正確
- [ ] 功能實現正確
- [ ] 性能影響最小
- [ ] 錯誤處理適當
- [ ] 日誌記錄完整

---

**記住**: 分層調試讓你更深入地理解系統架構，這對成為優秀的 Go 開發者非常重要！
