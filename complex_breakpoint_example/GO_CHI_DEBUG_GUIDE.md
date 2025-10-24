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

### 🟡 中級練習 - 中間件調試（go-chi 特色）

#### 案例 3: 內建中間件調試
**目標**: 學習 go-chi 內建中間件的調試

**斷點設置**:
```
routes/routes.go:19  - 檢查 Logger 中間件
routes/routes.go:20  - 檢查 Recoverer 中間件
routes/routes.go:21  - 檢查 RequestID 中間件
routes/routes.go:22  - 檢查 RealIP 中間件
routes/routes.go:23  - 檢查 Timeout 中間件
```

**對應代碼**:
```go
// 行 19 - Logger 中間件
router.Use(chimiddleware.Logger)

// 行 20 - Recoverer 中間件
router.Use(chimiddleware.Recoverer)

// 行 21 - RequestID 中間件
router.Use(chimiddleware.RequestID)

// 行 22 - RealIP 中間件
router.Use(chimiddleware.RealIP)

// 行 23 - Timeout 中間件
router.Use(chimiddleware.Timeout(60))
```

**Watch 表達式**:
```
r.Header.Get("X-Request-ID")
r.RemoteAddr
r.Context().Deadline()
```

**調試重點**:
- 觀察內建中間件的執行順序
- 理解 RequestID 的生成
- 檢查 RealIP 的處理
- 觀察 Timeout 的設置

---

#### 案例 4: 路由嵌套調試（go-chi 特色）
**目標**: 學習 go-chi 的路由嵌套機制

**斷點設置**:
```
routes/routes.go:38  - 檢查 API 路由設置
routes/routes.go:40  - 檢查用戶路由嵌套
routes/routes.go:41  - 檢查具體路由註冊
```

**對應代碼**:
```go
// 行 38 - API 路由設置
router.Route("/api/v1", func(r chi.Router) {
    // 行 40 - 用戶路由嵌套
    r.Route("/users", func(r chi.Router) {
        // 行 41 - 具體路由註冊
        r.Post("/", userHandler.CreateUser)
        r.Get("/", userHandler.GetUsers)
        r.Get("/{id}", userHandler.GetUser)
        r.Get("/{id}/account", userHandler.GetUserAccount)
        r.Get("/{id}/transactions", transactionHandler.GetUserTransactions)
        r.Get("/{id}/loans", loanHandler.GetUserLoanApplications)
        r.Post("/{id}/deposit", transactionHandler.Deposit)
        r.Post("/{id}/withdraw", transactionHandler.Withdraw)
        r.Post("/{id}/apply-loan", loanHandler.ApplyForLoan)
    })
})
```

**Watch 表達式**:
```
r
router
```

**調試重點**:
- 觀察路由嵌套的結構
- 理解 `chi.Router` 的層次關係
- 檢查路由註冊的順序

---

### 🔴 高級練習 - 並發和性能調試

#### 案例 5: 中間件鏈調試
**目標**: 學習中間件鏈的執行順序

**斷點設置**:
```
middleware/middleware.go:12  - 檢查日誌中間件開始
middleware/middleware.go:18  - 檢查日誌中間件結束
routes/routes.go:19-23       - 檢查內建中間件
routes/routes.go:26-30       - 檢查自定義中間件
```

**Watch 表達式**:
```
r.Method
r.URL.Path
start
time.Since(start)
r.Header.Get("X-Request-ID")
```

**調試重點**:
- 觀察中間件的執行順序
- 理解中間件鏈的構建
- 檢查性能影響

---

#### 案例 6: 錯誤處理調試
**目標**: 學習 go-chi 的錯誤處理機制

**斷點設置**:
```
routes/routes.go:20  - 檢查 Recoverer 中間件
middleware/middleware.go:45  - 檢查自定義錯誤處理
```

**對應代碼**:
```go
// routes/routes.go:20 - Recoverer 中間件
router.Use(chimiddleware.Recoverer)

// middleware/middleware.go:45 - 自定義錯誤處理
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("❌ Panic recovered: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}
```

**測試步驟**:
創建一個會 panic 的處理器來測試錯誤處理：
```go
// 在 handlers 中添加測試端點
func (h *UserHandler) TestPanic(w http.ResponseWriter, r *http.Request) {
    panic("Test panic for debugging")
}
```

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

### go-chi 特有檢查項
- [ ] URL 參數正確提取
- [ ] 路由嵌套正確
- [ ] 中間件鏈執行順序
- [ ] 內建中間件正常工作
- [ ] 自定義中間件正確集成
- [ ] 錯誤處理機制正常
- [ ] 請求 ID 生成正確
- [ ] 超時處理正常

### 通用檢查項
- [ ] Handler 層調試正常
- [ ] Service 層調試正常
- [ ] Database 層調試正常
- [ ] 錯誤處理適當
- [ ] 性能良好

## 🎮 實際操作步驟

### 1. **URL 參數調試**
```bash
# 設置斷點在 handlers/handlers.go:50
# 發送請求: GET http://localhost:9090/api/v1/users/1
# 觀察 chi.URLParam(r, "id") 的執行
```

### 2. **中間件調試**
```bash
# 設置斷點在各個中間件
# 發送請求並觀察執行順序
# 檢查內建中間件的功能
```

### 3. **路由調試**
```bash
# 設置斷點在 routes/routes.go
# 觀察路由註冊過程
# 檢查路由嵌套結構
```

### 4. **錯誤處理調試**
```bash
# 創建會 panic 的端點
# 測試 Recoverer 中間件
# 觀察錯誤處理流程
```

---

**記住**: go-chi 提供了更現代、更高效的 HTTP 路由解決方案，學習其調試技巧對成為優秀的 Go 開發者非常重要！
