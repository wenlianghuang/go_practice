# 🎯 Go-Chi Breakpoint 調試練習指南

## 📋 項目升級說明

項目已成功從 `gorilla/mux` 升級到 `go-chi`，這帶來了以下改進：

### 🚀 go-chi 的優勢

1. **更好的性能**: 更輕量級，路由匹配更快
2. **更現代的中間件支持**: 內建中間件和更好的中間件鏈
3. **更清晰的 API**: 更符合 Go 慣用法
4. **更好的錯誤處理**: 內建 panic 恢復和錯誤處理
5. **更活躍的維護**: 更頻繁的更新和更好的文檔

## 🎯 新的調試練習案例

### 🟢 初級練習 - Handler 層調試（go-chi 版本）

#### 案例 1: 用戶創建 Handler 調試
**目標**: 學習 go-chi 的 HTTP 請求處理調試

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

#### 案例 2: URL 參數解析調試（go-chi 特色）
**目標**: 學習 go-chi 的 URL 參數解析

**斷點設置**:
```
handlers/handlers.go:50  - 檢查 URL 參數提取
handlers/handlers.go:52  - 檢查參數轉換
```

**對應代碼**:
```go
// 行 50 - 檢查 URL 參數提取（go-chi 方式）
userIDStr := chi.URLParam(r, "id")

// 行 52 - 檢查參數轉換
userID, err := strconv.Atoi(userIDStr)
```

**Watch 表達式**:
```
userIDStr
userID
err
chi.URLParam(r, "id")
```

**調試重點**:
- 觀察 `chi.URLParam(r, "id")` 如何提取 URL 參數
- 比較與 `mux.Vars(r)["id"]` 的差異
- 理解 go-chi 的參數提取機制

**測試步驟**:
```
GET http://localhost:9090/api/v1/users/1
```

---

### 🟡 中級練習 - 中間件和服務層調試

#### 案例 3: 中間件鏈執行調試（修正版）
**目標**: 理解中間件如何包裝和執行，學習中間件鏈的工作原理

**斷點設置**:
```
middleware/middleware.go:12  - 中間件開始執行
middleware/middleware.go:18  - 調用下一個中間件/處理器
middleware/middleware.go:20  - 中間件結束執行
```

**對應代碼**:
```go
// middleware/middleware.go:12 - 中間件開始執行
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()  // 斷點設置在這裡
        
        log.Printf("🚀 Request started: %s %s", r.Method, r.URL.Path)
        
        // middleware/middleware.go:18 - 調用下一個中間件/處理器
        next.ServeHTTP(w, r)  // 斷點設置在這裡（關鍵！）
        
        // middleware/middleware.go:20 - 中間件結束執行
        log.Printf("✅ Request completed: %s %s in %v", 
            r.Method, r.URL.Path, time.Since(start))  // 斷點設置在這裡
    })
}
```

**調試按鍵使用**:
- **F10**: 執行 `start := time.Now()`，觀察 start 的值
- **F11**: 執行 `next.ServeHTTP(w, r)`，進入下一個中間件或 handler
- **F10**: 執行日誌記錄，觀察執行時間

**Watch 表達式**:
```
start
time.Since(start)
r.Method
r.URL.Path
w.Header().Get("Content-Type")
```

**學習重點**:
1. **中間件包裝**: 觀察 `next.ServeHTTP(w, r)` 如何調用下一個處理器
2. **執行順序**: 理解中間件的執行順序（註冊順序 = 執行順序）
3. **請求/響應流**: 觀察請求如何流經中間件鏈
4. **性能監控**: 理解如何測量請求處理時間

**測試步驟**:
```
GET http://localhost:9090/api/v1/users/1
```

---

#### 案例 4: 服務層調試（修正版）
**目標**: 理解 Handler → Service → Database 的調用鏈，學習分層架構的調試

**斷點設置**:
```
handlers/handlers.go:57  - Handler 調用 Service
services/services.go:25  - Service 調用 Database
database/database.go:45  - Database 實際操作
```

**對應代碼**:
```go
// handlers/handlers.go:57 - Handler 調用 Service
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    userIDStr := chi.URLParam(r, "id")
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    // 斷點設置在這裡
    user, err := h.userService.GetUser(userID)  // Handler 調用 Service
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// services/services.go:25 - Service 調用 Database
func (s *UserService) GetUser(userID int) (*models.User, error) {
    // 斷點設置在這裡
    user, err := s.db.GetUser(userID)  // Service 調用 Database
    
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

// database/database.go:45 - Database 實際操作
func (d *Database) GetUser(userID int) (*User, error) {
    d.mu.RLock()  // 斷點設置在這裡
    defer d.mu.RUnlock()
    
    user, exists := d.users[userID]
    if !exists {
        return nil, fmt.Errorf("user not found")
    }
    
    return user, nil
}
```

**調試按鍵使用**:
- **F11**: 在 Handler 中按 F11 進入 Service 層
- **F11**: 在 Service 中按 F11 進入 Database 層
- **F10**: 在 Database 中按 F10 執行鎖操作
- **Shift+F11**: 快速跳出不需要的函數

**Watch 表達式**:
```
userID
user
err
err != nil
d.users[userID]
exists
```

**學習重點**:
1. **分層調用**: 理解 Handler → Service → Database 的調用鏈
2. **錯誤傳播**: 觀察錯誤如何從 Database 層傳播到 Handler 層
3. **數據流**: 理解數據如何在各層之間傳遞
4. **並發安全**: 觀察 Database 層的鎖機制

**測試步驟**:
```
GET http://localhost:9090/api/v1/users/1
GET http://localhost:9090/api/v1/users/999  # 測試錯誤情況
```

---

### 🔴 高級練習 - 並發和錯誤處理調試

#### 案例 5: 並發調試（修正版）
**目標**: 學習並發環境下的調試技巧，理解 Goroutine 和鎖機制

**斷點設置**:
```
handlers/handlers.go:180  - 並發測試開始
handlers/handlers.go:185  - Goroutine 創建
handlers/handlers.go:195  - WaitGroup 等待
database/database.go:45   - 鎖競爭點
```

**對應代碼**:
```go
// handlers/handlers.go:180 - 並發測試開始
func (h *TransactionHandler) ConcurrentTest(w http.ResponseWriter, r *http.Request) {
    var wg sync.WaitGroup
    results := make(chan string, 10)
    
    // 斷點設置在這裡
    for i := 0; i < 5; i++ {  // 斷點設置在這裡
        wg.Add(1)
        
        // handlers/handlers.go:185 - Goroutine 創建
        go func(id int) {  // 斷點設置在這裡
            defer wg.Done()
            
            // 模擬並發操作
            user, err := h.transactionService.GetUser(id + 1)
            if err != nil {
                results <- fmt.Sprintf("Goroutine %d: Error - %v", id, err)
                return
            }
            
            results <- fmt.Sprintf("Goroutine %d: User %s found", id, user.Username)
        }(i)
    }
    
    // handlers/handlers.go:195 - WaitGroup 等待
    go func() {  // 斷點設置在這裡
        wg.Wait()
        close(results)
    }()
    
    var response []string
    for result := range results {
        response = append(response, result)
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// database/database.go:45 - 鎖競爭點
func (d *Database) GetUser(userID int) (*User, error) {
    d.mu.RLock()  // 斷點設置在這裡（關鍵！）
    defer d.mu.RUnlock()
    
    user, exists := d.users[userID]
    if !exists {
        return nil, fmt.Errorf("user not found")
    }
    
    return user, nil
}
```

**調試按鍵使用**:
- **F11**: 進入 Goroutine 函數，觀察並發執行
- **F10**: 執行鎖操作，觀察鎖競爭
- **F5**: 在並發環境中快速跳轉
- **Shift+F11**: 跳出 Goroutine

**Watch 表達式**:
```
wg
len(results)
id
userID
d.mu
r.Method
r.URL.Path
```

**學習重點**:
1. **並發調試**: 理解 Goroutine 的執行順序
2. **鎖機制**: 觀察讀寫鎖的競爭
3. **WaitGroup**: 理解同步等待機制
4. **Channel**: 觀察數據在 Goroutine 間的傳遞

**測試步驟**:
```
POST http://localhost:9090/api/v1/test/concurrent
```

---

#### 案例 6: 錯誤處理和 Panic 恢復調試（修正版）
**目標**: 學習錯誤處理機制和 Panic 恢復的調試

**斷點設置**:
```
middleware/middleware.go:45  - 自定義錯誤處理中間件
handlers/handlers.go:300     - 測試 Panic 端點
routes/routes.go:20          - Recoverer 中間件註冊
```

**對應代碼**:
```go
// middleware/middleware.go:45 - 自定義錯誤處理中間件
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {  // 斷點設置在這裡
            if err := recover(); err != nil {
                log.Printf("❌ Panic recovered: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

// handlers/handlers.go:300 - 測試 Panic 端點
func (h *UserHandler) TestPanic(w http.ResponseWriter, r *http.Request) {
    log.Println("🚨 About to panic...")
    
    // 斷點設置在這裡
    panic("Test panic for debugging")  // 斷點設置在這裡
}

// routes/routes.go:20 - Recoverer 中間件註冊
router.Use(chimiddleware.Recoverer)  // 斷點設置在這裡
```

**調試按鍵使用**:
- **F10**: 執行 panic 語句，觀察 panic 觸發
- **F11**: 進入錯誤處理中間件
- **F5**: 在錯誤處理流程中跳轉

**Watch 表達式**:
```
err
r.Method
r.URL.Path
w.Header().Get("Content-Type")
```

**學習重點**:
1. **Panic 機制**: 理解 panic 如何觸發
2. **Recover 機制**: 觀察 panic 如何被恢復
3. **錯誤傳播**: 理解錯誤如何在各層傳播
4. **中間件保護**: 學習如何用中間件保護應用

**測試步驟**:
```
# 需要先添加測試端點到路由
GET http://localhost:9090/api/v1/test/panic
```

---

## 🎮 Breakpoint 調試按鍵指南

### 🔍 基本調試按鍵

#### **F5 (Continue)**
- **用途**: 繼續執行到下一個斷點
- **使用時機**: 當你已經檢查完當前斷點的變量，想要跳到下一個斷點
- **示例**: 在 `handlers.go:50` 檢查完 `userIDStr` 後，按 F5 跳到下一個斷點

#### **F10 (Step Over)**
- **用途**: 執行當前行，但不進入函數內部
- **使用時機**: 當你不想深入函數內部，只想看當前行的執行結果
- **示例**: 在 `strconv.Atoi(userIDStr)` 這行，按 F10 會執行轉換但不進入 `Atoi` 函數內部

#### **F11 (Step Into)**
- **用途**: 進入函數內部進行詳細調試
- **使用時機**: 當你想要深入了解函數的內部邏輯
- **示例**: 在 `h.userService.GetUser(userID)` 這行，按 F11 會進入 `GetUser` 函數內部

#### **Shift+F11 (Step Out)**
- **用途**: 跳出當前函數，回到調用者
- **使用時機**: 當你在函數內部調試完畢，想要回到調用者
- **示例**: 在 `GetUser` 函數內部調試完後，按 Shift+F11 回到 handler

### 🎯 針對不同調試場景的按鍵選擇

#### **場景 1: Handler 層調試**
```go
// handlers/handlers.go:50
userIDStr := chi.URLParam(r, "id")  // 斷點設置在這裡

// 按鍵選擇:
// F10 - 執行這行，觀察 userIDStr 的值
// F11 - 進入 chi.URLParam 函數內部（通常不需要）

// handlers/handlers.go:52  
userID, err := strconv.Atoi(userIDStr)  // 斷點設置在這裡

// 按鍵選擇:
// F10 - 執行轉換，觀察 userID 和 err 的值
// F11 - 進入 strconv.Atoi 內部（通常不需要）
```

#### **場景 2: 中間件調試**
```go
// middleware/middleware.go:12
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()  // 斷點設置在這裡
        
        // 按鍵選擇:
        // F10 - 執行這行，觀察 start 的值
        // F11 - 進入 time.Now() 內部（通常不需要）
        
        log.Printf("🚀 Request started: %s %s", r.Method, r.URL.Path)
        
        // 斷點設置在這裡
        next.ServeHTTP(w, r)  // 這是關鍵！
        
        // 按鍵選擇:
        // F11 - 進入下一個中間件或 handler（推薦！）
        // F10 - 跳過，不進入內部（會錯過重要調試機會）
        
        log.Printf("✅ Request completed: %s %s in %v", 
            r.Method, r.URL.Path, time.Since(start))
    })
}
```

#### **場景 3: 服務層調試**
```go
// services/services.go
func (s *UserService) GetUser(userID int) (*models.User, error) {
    // 斷點設置在這裡
    user, err := s.db.GetUser(userID)
    
    // 按鍵選擇:
    // F11 - 進入 database 層（推薦！了解數據庫操作）
    // F10 - 跳過，只觀察結果（如果已經了解數據庫邏輯）
    
    if err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### 🎯 推薦的調試流程

#### **初級練習 - 使用 F10**
```go
// 1. 設置斷點在 handlers/handlers.go:50
userIDStr := chi.URLParam(r, "id")

// 2. 按 F10 執行這行
// 3. 在 Watch 中觀察 userIDStr 的值
// 4. 按 F5 跳到下一個斷點

// 5. 設置斷點在 handlers/handlers.go:52
userID, err := strconv.Atoi(userIDStr)

// 6. 按 F10 執行這行
// 7. 在 Watch 中觀察 userID 和 err 的值
```

#### **中級練習 - 混合使用 F10 和 F11**
```go
// 1. 設置斷點在 middleware/middleware.go:12
start := time.Now()

// 2. 按 F10 執行這行
// 3. 觀察 start 的值

// 4. 設置斷點在 middleware/middleware.go:18
next.ServeHTTP(w, r)

// 5. 按 F11 進入下一個中間件或 handler
// 6. 觀察中間件鏈的執行順序
```

#### **高級練習 - 深度調試使用 F11**
```go
// 1. 設置斷點在 handlers/handlers.go:57
user, err := h.userService.GetUser(userID)

// 2. 按 F11 進入 UserService.GetUser
// 3. 在服務層設置斷點
// 4. 按 F11 進入 Database.GetUser
// 5. 觀察完整的調用鏈
```

### 🎯 實際調試建議

#### **對於初學者**:
- **主要使用 F10**: 專注於觀察變量值，不深入函數內部
- **偶爾使用 F11**: 只在關鍵函數調用時使用
- **使用 F5**: 在檢查完變量後跳到下一個斷點

#### **對於中級開發者**:
- **混合使用 F10 和 F11**: 根據需要選擇是否深入函數
- **使用 Shift+F11**: 快速跳出不需要的函數
- **使用 F5**: 在熟悉代碼後快速跳轉

#### **對於高級開發者**:
- **主要使用 F11**: 深入理解代碼執行流程
- **使用 Call Stack**: 觀察完整的調用鏈
- **使用 F5**: 快速定位問題

### 🎯 調試技巧

#### **1. 設置多個斷點**
```go
// 在關鍵位置設置斷點
userIDStr := chi.URLParam(r, "id")     // 斷點 1
userID, err := strconv.Atoi(userIDStr) // 斷點 2
user, err := h.userService.GetUser(userID) // 斷點 3
```

#### **2. 使用 Watch 表達式**
```
userIDStr
userID
err
err != nil
r.URL.Path
r.Method
```

#### **3. 觀察 Call Stack**
- 按 F11 進入函數後，觀察 Call Stack 面板
- 了解函數的調用層次
- 理解代碼的執行流程

---

## 🎯 go-chi 特有的調試技巧

### 1. **URL 參數調試**
```go
// 使用 chi.URLParam 調試
userID := chi.URLParam(r, "id")
fmt.Printf("Extracted userID: %s\n", userID)
```

### 2. **路由調試**
```go
// 檢查路由匹配
fmt.Printf("Request path: %s\n", r.URL.Path)
fmt.Printf("Request method: %s\n", r.Method)
```

### 3. **中間件調試**
```go
// 在中間件中添加調試信息
log.Printf("Middleware executing: %s", r.URL.Path)
```

### 4. **上下文調試**
```go
// 檢查請求上下文
ctx := r.Context()
if deadline, ok := ctx.Deadline(); ok {
    fmt.Printf("Request deadline: %v\n", deadline)
}
```

## 🔧 調試工具和技巧

### 1. **內建中間件監控**
- **Logger**: 自動記錄請求日誌
- **RequestID**: 為每個請求生成唯一 ID
- **RealIP**: 獲取真實客戶端 IP
- **Timeout**: 設置請求超時

### 2. **路由調試**
```bash
# 測試不同路由
curl http://localhost:9090/api/v1/users
curl http://localhost:9090/api/v1/users/1
curl http://localhost:9090/api/v1/users/1/account
```

### 3. **中間件調試**
```bash
# 觀察中間件日誌
curl -H "X-Request-ID: test-123" http://localhost:9090/api/v1/users
```

### 4. **性能調試**
```bash
# 測試並發請求
for i in {1..10}; do
  curl http://localhost:9090/api/v1/users &
done
wait
```

## 📊 調試檢查清單

### 🟢 初級檢查項
- [ ] URL 參數正確提取 (`chi.URLParam(r, "id")`)
- [ ] 參數類型轉換正確 (`strconv.Atoi`)
- [ ] 錯誤處理適當 (`err != nil`)
- [ ] JSON 響應格式正確
- [ ] HTTP 狀態碼正確

### 🟡 中級檢查項
- [ ] 中間件鏈執行順序正確
- [ ] Handler → Service → Database 調用鏈正常
- [ ] 錯誤從底層正確傳播到頂層
- [ ] 數據在各層之間正確傳遞
- [ ] 鎖機制正常工作 (讀寫鎖)

### 🔴 高級檢查項
- [ ] Goroutine 並發執行正常
- [ ] WaitGroup 同步機制正確
- [ ] Channel 數據傳遞正常
- [ ] Panic 恢復機制正常
- [ ] 內建中間件功能正常 (Logger, Recoverer, RequestID, RealIP, Timeout)

### 🎯 go-chi 特有檢查項
- [ ] 路由嵌套結構正確
- [ ] URL 參數提取機制正常
- [ ] 內建中間件集成正確
- [ ] 自定義中間件正確執行
- [ ] 錯誤處理中間件正常工作

## 🎮 實際操作步驟

### 1. **URL 參數調試**
```bash
# 設置斷點在 handlers/handlers.go:50
# 發送請求: GET http://localhost:9090/api/v1/users/1
# 觀察 chi.URLParam(r, "id") 的執行
# 使用 F10 執行，觀察 userIDStr 的值
```

### 2. **中間件調試**
```bash
# 設置斷點在 middleware/middleware.go:12 和 18
# 發送請求並觀察執行順序
# 使用 F11 進入 next.ServeHTTP(w, r)
# 檢查內建中間件的功能
```

### 3. **服務層調試**
```bash
# 設置斷點在 handlers/handlers.go:57
# 使用 F11 進入 Service 層
# 繼續使用 F11 進入 Database 層
# 觀察完整的調用鏈
```

### 4. **並發調試**
```bash
# 設置斷點在 handlers/handlers.go:180, 185, 195
# 設置斷點在 database/database.go:45
# 發送並發請求: POST http://localhost:9090/api/v1/test/concurrent
# 觀察 Goroutine 和鎖的競爭
```

### 5. **錯誤處理調試**
```bash
# 設置斷點在 middleware/middleware.go:45
# 創建會 panic 的端點
# 測試 Recoverer 中間件
# 觀察錯誤處理流程
```

## 🎯 調試技巧總結

### **按鍵使用建議**:
- **初學者**: 主要使用 F10，偶爾使用 F11
- **中級**: 混合使用 F10 和 F11，使用 Shift+F11 跳出
- **高級**: 主要使用 F11，觀察 Call Stack

### **斷點設置策略**:
- **關鍵變量**: 在變量賦值後設置斷點
- **函數調用**: 在函數調用前設置斷點
- **錯誤處理**: 在錯誤檢查處設置斷點
- **並發點**: 在 Goroutine 創建和鎖操作處設置斷點

### **Watch 表達式建議**:
- **變量值**: 直接變量名
- **條件判斷**: `err != nil`, `len(slice) > 0`
- **請求信息**: `r.Method`, `r.URL.Path`
- **響應信息**: `w.Header().Get("Content-Type")`

---

**記住**: go-chi 提供了更現代、更高效的 HTTP 路由解決方案，學習其調試技巧對成為優秀的 Go 開發者非常重要！通過這些練習，你將掌握：

1. **基礎調試**: 變量觀察、錯誤處理
2. **中級調試**: 中間件鏈、分層架構
3. **高級調試**: 並發處理、錯誤恢復
4. **go-chi 特色**: 路由嵌套、內建中間件

**開始你的調試之旅吧！** 🚀
