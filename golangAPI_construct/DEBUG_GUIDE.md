# Cursor 斷點調試指南

## 🎯 調試目標
測試 `golangAPI_construct` 項目的請求驗證中間件和書籍創建功能

## 🔴 已設置的斷點位置

### 1. middleware/validation.go
- **第 53 行**: JSON 解析前檢查
- **第 65 行**: 驗證邏輯執行前檢查
- **第 78 行**: 驗證成功後數據檢查

### 2. handlers/book.go
- **第 50 行**: 從中間件獲取驗證數據

## 🚀 調試步驟

### 步驟 1: 啟動調試
1. 在 Cursor 中按 `Cmd+Shift+D` (Mac) 或 `Ctrl+Shift+D` (Windows) 打開調試面板
2. 選擇 "Launch Go API" 配置
3. 按 `F5` 或點擊綠色播放按鈕開始調試

### 步驟 2: 設置斷點
在以下行號左側點擊設置斷點（紅色圓點）：
- `middleware/validation.go:53`
- `middleware/validation.go:65`
- `middleware/validation.go:78`
- `handlers/book.go:50`

### 步驟 3: 發送測試請求

#### 測試 1: 有效請求
```bash
curl -X POST http://localhost:8081/api/v1/books \
  -H "Content-Type: application/json" \
  -d @test_book.json
```

#### 測試 2: 無效請求（觸發驗證錯誤）
```bash
curl -X POST http://localhost:8081/api/v1/books \
  -H "Content-Type: application/json" \
  -d @test_invalid_book.json
```

### 步驟 4: 調試過程

當程序停在斷點時，你可以：

1. **查看變量**:
   - 在 Variables 面板查看 `requestData`, `err`, `validatedData` 等
   - 將滑鼠懸停在變量上查看值

2. **控制執行**:
   - `F10`: 單步執行（不進入函數）
   - `F11`: 單步執行（進入函數）
   - `F5`: 繼續執行到下一個斷點

3. **檢查數據**:
   - 斷點 1: 檢查 `r.Body` 和請求方法
   - 斷點 2: 檢查 `requestData` 的內容
   - 斷點 3: 檢查 `err` 是否為 nil
   - 斷點 4: 檢查 `validatedData` 和 `ok` 的值

## 📊 預期結果

### 有效請求流程:
1. 斷點 1: `requestData` 為空 map
2. 斷點 2: `requestData` 包含完整的書籍數據
3. 斷點 3: `err` 為 nil，驗證通過
4. 斷點 4: `validatedData` 包含驗證後的數據，`ok` 為 true

### 無效請求流程:
1. 斷點 1: `requestData` 為空 map
2. 斷點 2: `requestData` 包含無效數據（空標題、負價格）
3. 斷點 3: `err` 不為 nil，包含驗證錯誤信息
4. 程序不會到達斷點 4，因為驗證失敗

## 🔧 調試技巧

1. **使用 Watch 面板**: 添加 `requestData["title"]` 等表達式進行監視
2. **檢查調用棧**: 在 Call Stack 面板查看函數調用順序
3. **查看控制台**: 在 Debug Console 中執行 Go 表達式
4. **條件斷點**: 右鍵斷點可以設置條件，如 `len(requestData) > 0`

## 🎉 成功標誌

調試成功時，你應該能夠：
- 看到請求數據在每個斷點處的狀態
- 理解中間件如何驗證數據
- 觀察驗證錯誤如何被處理
- 看到驗證後的數據如何傳遞給 handler

