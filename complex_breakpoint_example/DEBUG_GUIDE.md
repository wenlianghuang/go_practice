# 🎯 複雜 Go Breakpoint 調試練習指南

## 📋 項目概述

這是一個複雜的銀行系統模擬器，包含多個業務場景，專門設計用於練習 Go 語言的 breakpoint 調試技巧。系統包含：

- **用戶管理**: 創建、查詢用戶
- **銀行帳戶**: 存款、提款、轉帳
- **貸款系統**: 申請、審核貸款
- **並發處理**: Goroutine 和 Channel 調試
- **錯誤處理**: 多種錯誤場景
- **數據競爭**: Race condition 調試

## 🔧 環境準備

### 1. 安裝依賴
```bash
cd complex_breakpoint_example
go mod tidy
```

### 2. 安裝 Delve 調試器
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

### 3. VS Code/Cursor 配置
創建 `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Complex Breakpoint Example",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/complex_breakpoint_example",
            "args": [],
            "showLog": true
        }
    ]
}
```

## 🚀 練習案例

### 案例 1: 用戶創建流程調試

#### 🎯 目標
理解用戶創建、數據驗證和錯誤處理的完整流程

#### 📍 斷點設置
```
main.go:45    - 檢查請求數據解析
main.go:50    - 檢查用戶名驗證
main.go:55    - 檢查用戶創建
main.go:65    - 檢查銀行帳戶創建
main.go:75    - 檢查響應構建
```

#### 🧪 測試步驟

**步驟 1: 設置斷點**
在指定行號左側點擊設置斷點

**步驟 2: 啟動調試**
- 按 `Cmd+Shift+D` 打開調試面板
- 選擇 "Launch Complex Breakpoint Example"
- 按 `F5` 開始調試

**步驟 3: 發送創建用戶請求**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com"
  }'
```

#### 🔍 調試觀察點

**斷點 1 (main.go:45)**:
- 檢查 `req.Username` 和 `req.Email`
- 預期: `req.Username = "testuser"`, `req.Email = "test@example.com"`

**斷點 2 (main.go:50)**:
- 檢查用戶名重複驗證邏輯
- 觀察 `for` 循環中的用戶名比較

**斷點 3 (main.go:55)**:
- 檢查新用戶對象的創建
- 觀察 `d.nextID` 的遞增

**斷點 4 (main.go:65)**:
- 檢查銀行帳戶的創建
- 觀察 `UserID` 的關聯

### 案例 2: 存款處理調試

#### 🎯 目標
理解存款處理、餘額更新和交易記錄的流程

#### 📍 斷點設置
```
main.go:120   - 檢查用戶存在性驗證
main.go:125   - 檢查用戶活躍狀態
main.go:135   - 檢查帳戶鎖定狀態
main.go:145   - 檢查金額有效性
main.go:155   - 檢查交易記錄創建
main.go:165   - 檢查餘額更新
```

#### 🧪 測試步驟

**測試 1: 正常存款**
```bash
curl -X POST http://localhost:8080/api/v1/users/1/deposit \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.0,
    "description": "Test deposit"
  }'
```

**測試 2: 無效金額**
```bash
curl -X POST http://localhost:8080/api/v1/users/1/deposit \
  -H "Content-Type: application/json" \
  -d '{
    "amount": -50.0,
    "description": "Invalid deposit"
  }'
```

#### 🔍 調試觀察點

**正常存款流程**:
1. **斷點 1**: `user` 對象不為 nil
2. **斷點 2**: `user.IsActive = true`
3. **斷點 3**: `account.IsLocked = false`
4. **斷點 4**: `amount > 0`
5. **斷點 5**: `transaction.Status = "pending"`
6. **斷點 6**: `account.Balance` 增加

**錯誤處理流程**:
- 觀察錯誤如何被檢測和返回
- 檢查交易狀態如何設置為 "failed"

### 案例 3: 轉帳處理調試

#### 🎯 目標
理解複雜業務邏輯、多用戶驗證和數據一致性

#### 📍 斷點設置
```
main.go:180   - 檢查發送方用戶
main.go:190   - 檢查接收方用戶
main.go:200   - 檢查發送方帳戶
main.go:210   - 檢查接收方帳戶
main.go:220   - 檢查帳戶鎖定狀態
main.go:230   - 檢查餘額充足性
main.go:250   - 檢查餘額更新
```

#### 🧪 測試步驟

**測試 1: 正常轉帳**
```bash
curl -X POST http://localhost:8080/api/v1/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "from_user_id": 1,
    "to_user_id": 2,
    "amount": 50.0,
    "description": "Transfer test"
  }'
```

**測試 2: 餘額不足**
```bash
curl -X POST http://localhost:8080/api/v1/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "from_user_id": 2,
    "to_user_id": 1,
    "amount": 1000.0,
    "description": "Insufficient funds test"
  }'
```

#### 🔍 調試觀察點

**轉帳成功流程**:
1. 兩個用戶都被找到且活躍
2. 兩個帳戶都存在且未鎖定
3. 發送方餘額充足
4. 兩個帳戶餘額同時更新
5. 交易狀態設為 "completed"

### 案例 4: 貸款申請調試

#### 🎯 目標
理解異步處理、Goroutine 調試和並發安全

#### 📍 斷點設置
```
main.go:280   - 檢查貸款申請創建
main.go:290   - 檢查利率計算
main.go:300   - 檢查 Goroutine 啟動
main.go:310   - 檢查審核邏輯
main.go:320   - 檢查審核結果
```

#### 🧪 測試步驟

**測試 1: 申請貸款**
```bash
curl -X POST http://localhost:8080/api/v1/users/1/apply-loan \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5000.0,
    "term": 24
  }'
```

**測試 2: 查詢貸款狀態**
```bash
curl http://localhost:8080/api/v1/users/1/loans
```

#### 🔍 調試觀察點

**貸款申請流程**:
1. **斷點 1**: 檢查申請對象創建
2. **斷點 2**: 觀察利率計算邏輯
3. **斷點 3**: 檢查 Goroutine 啟動
4. **斷點 4**: 觀察審核邏輯
5. **斷點 5**: 檢查最終狀態

**Goroutine 調試技巧**:
- 使用 `dlv goroutines` 查看所有 Goroutine
- 使用 `dlv goroutine <id>` 切換到特定 Goroutine
- 觀察並發執行的順序

### 案例 5: 並發操作調試

#### 🎯 目標
理解並發處理、Race Condition 和同步機制

#### 📍 斷點設置
```
main.go:350   - 檢查並發請求解析
main.go:360   - 檢查 Goroutine 創建
main.go:370   - 檢查 WaitGroup 使用
main.go:380   - 檢查互斥鎖使用
main.go:390   - 檢查結果收集
```

#### 🧪 測試步驟

**測試: 並發存款操作**
```bash
curl -X POST http://localhost:8080/api/v1/test/concurrent \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "amount": 10.0,
    "operation": "deposit",
    "count": 10
  }'
```

#### 🔍 調試觀察點

**並發調試技巧**:
1. **斷點 1**: 檢查請求參數
2. **斷點 2**: 觀察多個 Goroutine 創建
3. **斷點 3**: 檢查 WaitGroup.Add() 和 WaitGroup.Done()
4. **斷點 4**: 觀察互斥鎖的獲取和釋放
5. **斷點 5**: 檢查結果的正確性

**Race Condition 檢測**:
- 使用 `go run -race main.go` 檢測數據競爭
- 觀察沒有鎖保護的共享數據訪問
- 理解為什麼需要互斥鎖

## 🔧 高級調試技巧

### 1. 條件斷點
設置條件斷點來只在特定情況下停止：
```
userID == 1
amount > 100
operation == "deposit"
```

### 2. 日誌斷點
在斷點處添加日誌而不停止執行：
```
fmt.Printf("Processing user %d with amount %.2f\n", userID, amount)
```

### 3. 調試表達式
在調試控制台中執行表達式：
```go
// 查看所有用戶
db.users

// 檢查特定用戶的餘額
db.users[1].Balance

// 查看所有交易
len(db.transactions)

// 檢查 Goroutine 數量
runtime.NumGoroutine()
```

### 4. 內存調試
```go
// 查看內存使用
runtime.MemStats
runtime.ReadMemStats(&m)
```

### 5. 並發調試
```go
// 查看所有 Goroutine
runtime.NumGoroutine()

// 查看 Goroutine 堆棧
runtime.Stack(buf, true)
```

## 🐛 常見調試場景

### 場景 1: 數據競爭
**問題**: 多個 Goroutine 同時修改共享數據
**調試**: 使用 `-race` 標誌運行程序
**解決**: 添加適當的鎖保護

### 場景 2: 死鎖
**問題**: Goroutine 互相等待導致程序掛起
**調試**: 使用 `dlv goroutines` 查看 Goroutine 狀態
**解決**: 檢查鎖的獲取順序

### 場景 3: 內存洩漏
**問題**: 內存使用持續增長
**調試**: 使用 `runtime.MemStats` 監控內存
**解決**: 確保資源正確釋放

### 場景 4: 性能瓶頸
**問題**: 程序執行緩慢
**調試**: 使用 `pprof` 進行性能分析
**解決**: 優化熱點代碼

## 📊 調試檢查清單

### 基礎調試
- [ ] 設置適當的斷點
- [ ] 檢查變量值
- [ ] 觀察程序流程
- [ ] 理解錯誤信息

### 並發調試
- [ ] 檢查 Goroutine 狀態
- [ ] 觀察鎖的使用
- [ ] 檢測數據競爭
- [ ] 驗證同步機制

### 錯誤處理
- [ ] 理解錯誤來源
- [ ] 檢查錯誤處理邏輯
- [ ] 驗證錯誤恢復
- [ ] 測試邊界條件

### 性能調試
- [ ] 監控內存使用
- [ ] 檢查 CPU 使用率
- [ ] 分析瓶頸
- [ ] 優化熱點代碼

## 🎯 練習目標

完成所有練習後，你應該能夠：

1. **熟練使用斷點**: 設置條件斷點、日誌斷點
2. **調試並發程序**: 理解 Goroutine 調試
3. **檢測數據競爭**: 使用 race detector
4. **分析性能問題**: 使用 profiling 工具
5. **理解錯誤處理**: 調試複雜錯誤場景
6. **掌握同步機制**: 調試鎖和 Channel

