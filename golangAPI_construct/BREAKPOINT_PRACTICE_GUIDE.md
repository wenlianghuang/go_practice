# Go API 斷點調試練習指南

## 🎯 目標
通過實際的斷點調試練習，深入理解 Go API 的請求處理流程、錯誤處理機制和數據流動。

## 📋 目錄
- [環境準備](#環境準備)
- [基礎調試設置](#基礎調試設置)
- [練習案例](#練習案例)
  - [案例 1: JWT 登入流程](#案例-1-jwt-登入流程)
  - [案例 2: 獲取書籍列表](#案例-2-獲取書籍列表)
  - [案例 3: 錯誤處理調試](#案例-3-錯誤處理調試)
- [調試技巧](#調試技巧)
- [常見問題](#常見問題)

## 🔧 環境準備

### 1. 安裝 Delve 調試器
```bash
# 安裝 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 驗證安裝
dlv version
```

### 2. 配置 launch.json
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Go API (Debug)",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/golangAPI_construct",
            "env": {
                "PORT": "8081",
                "USE_DB": "true",
                "USE_GORM": "true",
                "DB_DSN": "file:books.db",
                "JWT_SECRET": "your-secret-key"
            },
            "args": []
        }
    ]
}
```

### 3. 準備測試工具
- **Postman** 或 **curl** 用於發送 HTTP 請求
- **Cursor** 調試功能

## 🚀 基礎調試設置

### 調試控制按鈕
- **`F5` (Continue)**: 繼續執行到下一個斷點
- **`F10` (Step Over)**: 單步執行，不進入函數內部
- **`F11` (Step Into)**: 單步執行，會進入函數內部
- **`Shift+F11` (Step Out)**: 跳出當前函數
- **`Shift+F5` (Stop)**: 停止調試

### 調試面板使用
- **Variables 面板**: 查看當前作用域的所有變量
- **Watch 面板**: 監視特定表達式的值
- **Call Stack 面板**: 查看函數調用順序
- **Debug Console**: 執行 Go 表達式

## 📚 練習案例

### 案例 1: JWT 登入流程

#### 🎯 目標
理解用戶登入、JWT token 生成和驗證的完整流程

#### 📍 斷點設置
```
handlers/auth.go:35    - 檢查登入請求數據
handlers/auth.go:41    - 檢查用戶驗證
handlers/auth.go:48    - 檢查 token 生成
security/jwt.go:28     - 檢查 JWT 生成過程
security/jwt.go:38     - 檢查 token 簽名
```

#### 🧪 測試步驟

**步驟 1: 設置斷點**
在指定行號左側點擊設置斷點（紅色圓點）

**步驟 2: 啟動調試**
- 按 `Cmd+Shift+D` 打開調試面板
- 選擇 "Launch Go API (Debug)" 配置
- 按 `F5` 開始調試

**步驟 3: 發送登入請求**
**Postman 設置**:
- **URL**: `POST http://localhost:8081/api/v1/auth/login`
- **Headers**: 
  ```
  Content-Type: application/json
  ```
- **Body**:
  ```json
  {
    "username": "Matt",
    "password": "password"
  }
  ```

#### 🔍 調試觀察點

**斷點 1 (auth.go:35)**:
- 檢查 `req.Username` 和 `req.Password`
- 預期: `req.Username = "Matt"`, `req.Password = "password"`

**斷點 2 (auth.go:41)**:
- 檢查 `demoUser.Username` 和 `demoUser.PasswordHash`
- 檢查 `security.CheckPassword()` 的結果

**斷點 3 (auth.go:48)**:
- 檢查 `req.Username` 和 `ttl` 值
- 預期: `username = "Matt"`, `ttl = 2h0m0s`

**斷點 4 (jwt.go:28)**:
- 檢查 `username` 參數
- 檢查 `ttl` 參數

**斷點 5 (jwt.go:38)**:
- 檢查 `claims` 結構
- 檢查 `secret()` 返回的密鑰

#### 📊 預期結果
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": 1703123456,
    "user": "Matt"
  }
}
```

### 案例 2: 獲取書籍列表

#### 🎯 目標
理解數據查詢、業務邏輯處理和響應構建的流程

#### 📍 斷點設置
```
handlers/book.go:28    - 檢查查詢參數
handlers/book.go:30    - 檢查結果變量
handlers/book.go:32    - 檢查作者查詢分支
handlers/book.go:34    - 檢查全部查詢分支
handlers/book.go:37    - 檢查響應構建
services/book_gorm.go:48  - 檢查 GetAllBooks
services/book_gorm.go:49  - 檢查數據庫查詢
```

#### 🧪 測試步驟

**測試 1: 獲取所有書籍**
- **URL**: `GET http://localhost:8081/api/v1/books`
- **Headers**: 
  ```
  Authorization: Bearer <你的JWT token>
  Content-Type: application/json
  ```

**測試 2: 按作者查詢**
- **URL**: `GET http://localhost:8081/api/v1/books?author=張三`
- **Headers**: 
  ```
  Authorization: Bearer <你的JWT token>
  Content-Type: application/json
  ```

#### 🔍 調試觀察點

**獲取所有書籍流程**:
1. **斷點 1**: `author = ""` (空字串)
2. **斷點 2**: `result = []` (空切片)
3. **斷點 3**: 進入 `GetAllBooks()` 分支
4. **斷點 4**: 檢查數據庫查詢結果
5. **斷點 5**: 檢查最終響應數據

**按作者查詢流程**:
1. **斷點 1**: `author = "張三"`
2. **斷點 2**: 進入 `GetBooksByAuthor()` 分支
3. **斷點 3**: 檢查 SQL 查詢: `WHERE author ILIKE '%張三%'`
4. **斷點 4**: 檢查查詢結果

### 案例 3: 錯誤處理調試

#### 🎯 目標
理解錯誤檢測、錯誤處理和失敗響應的完整流程

#### 📍 斷點設置
```
handlers/book.go:96    - 錯誤處理和失敗響應
middleware/auth.go:18  - 檢查 Authorization header
middleware/auth.go:37  - 檢查 token 驗證
middleware/validation.go:52  - JSON 解析錯誤
middleware/validation.go:65  - 數據驗證錯誤
```

#### 🧪 測試案例

**案例 3.1: 無效的書籍 ID**
- **URL**: `GET http://localhost:8081/api/v1/books/99999`
- **Headers**: 
  ```
  Authorization: Bearer <有效token>
  Content-Type: application/json
  ```

**調試過程**:
1. 程序停在 `handlers/book.go:96`
2. **Variables 面板**:
   - `id = "99999"`
   - `book = nil`
   - `err = "book not found"`
3. **預期響應**:
   ```json
   {
     "success": false,
     "error": {
       "code": "NOT_FOUND",
       "message": "book not found",
       "status": 404
     }
   }
   ```

**案例 3.2: 無效的 JWT Token**
- **URL**: `GET http://localhost:8081/api/v1/books`
- **Headers**: 
  ```
  Authorization: Bearer invalid_token_here
  Content-Type: application/json
  ```

**調試過程**:
1. 程序停在 `middleware/auth.go:37`
2. **Variables 面板**:
   - `authHeader = "Bearer invalid_token_here"`
   - `tokenString = "invalid_token_here"`
   - `err != nil` (token 驗證失敗)
3. **預期響應**:
   ```json
   {
     "success": false,
     "error": {
       "code": "INVALID_TOKEN",
       "message": "Invalid or expired token",
       "status": 401
     }
   }
   ```

**案例 3.3: 驗證失敗**
- **URL**: `POST http://localhost:8081/api/v1/books`
- **Headers**: 
  ```
  Authorization: Bearer <有效token>
  Content-Type: application/json
  ```
- **Body**:
  ```json
  {
    "title": "",
    "author": "",
    "price": -100
  }
  ```

**調試過程**:
1. 程序停在 `middleware/validation.go:65`
2. **Variables 面板**:
   - `requestData` 包含無效數據
   - `err != nil` (驗證失敗)
3. **預期響應**:
   ```json
   {
     "success": false,
     "error": {
       "code": "VALIDATION_ERROR",
       "message": "Field 'title' is required",
       "status": 400
     }
   }
   ```

## 🔧 調試技巧

### 1. 斷點管理
- **右鍵斷點** → "Disable Breakpoint" 暫時禁用
- **右鍵斷點** → "Edit Breakpoint" 設置條件
- **條件斷點示例**: `req.Username == "Matt"`

### 2. 變量監視
**重要變量**:
- `req.Username`, `req.Password` - 登入數據
- `token`, `claims` - JWT 相關
- `authHeader`, `tokenString` - 認證相關
- `requestData` - 請求數據
- `err` - 錯誤對象

### 3. Watch 表達式
```
len(token)                    - token 長度
claims.Username              - 用戶名
r.Header.Get("Authorization") - 完整 header
len(result)                  - 結果數量
err.Error()                  - 錯誤信息
```

### 4. 調試控制台
在調試控制台中執行 Go 表達式：
```go
fmt.Printf("Username: %s\n", req.Username)
fmt.Printf("Token length: %d\n", len(token))
```

## ❓ 常見問題

### Q: 調試器啟動後立即退出
**A**: 檢查環境變量設置，確保 `USE_DB=true` 和 `USE_GORM=true`

### Q: 斷點沒有觸發
**A**: 
1. 確認斷點設置在可執行代碼行
2. 檢查請求 URL 和端口是否正確
3. 確認服務器正在運行

### Q: Variables 面板為空
**A**: 
1. 確認程序停在斷點處
2. 檢查變量是否在當前作用域內
3. 使用 `locals` 命令查看局部變量

### Q: 端口衝突
**A**: 
- 調試模式使用 8081 端口
- 普通運行使用 8080 端口
- 確保沒有其他程序占用相同端口

## 🎉 成功標誌

調試成功時，你應該能夠：
- **看到** 請求數據在每個斷點處的狀態
- **理解** 中間件如何處理請求
- **觀察** JWT token 的生成和驗證過程
- **掌握** 錯誤處理的完整流程
- **熟悉** 數據如何在組件間傳遞

## 📝 練習記錄

### 練習日期: ___________

### 完成的案例:
- [ ] 案例 1: JWT 登入流程
- [ ] 案例 2: 獲取書籍列表  
- [ ] 案例 3: 錯誤處理調試

### 學到的重點:
1. ________________________________
2. ________________________________
3. ________________________________

### 遇到的問題:
1. ________________________________
2. ________________________________
3. ________________________________

### 解決方案:
1. ________________________________
2. ________________________________
3. ________________________________

---

**記住**: 調試是一個實踐的過程，多練習才能熟練掌握！
