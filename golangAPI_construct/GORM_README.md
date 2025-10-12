# Go API Construct - GORM 整合功能

## 🎯 什麼是 GORM？

**GORM** 是 Go 語言最受歡迎的 ORM（Object-Relational Mapping）框架，它讓你可以用 Go 結構體來操作數據庫，而不需要寫複雜的 SQL 語句。

### 🤔 為什麼需要 GORM？

想像一下，如果沒有 GORM，你需要這樣寫代碼：
```go
// 沒有 GORM 的寫法 - 複雜且容易出錯
rows, err := db.Query("SELECT id, title, author, price FROM books WHERE author = ?", "George Orwell")
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

var books []Book
for rows.Next() {
    var book Book
    err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price)
    if err != nil {
        log.Fatal(err)
    }
    books = append(books, book)
}
```

有了 GORM，你只需要：
```go
// 使用 GORM 的寫法 - 簡潔且類型安全
var books []BookGORM
db.Where("author = ?", "George Orwell").Find(&books)
```

## 🚀 GORM 的核心優勢

### 1. **代碼簡潔性**
- **減少 70% 的數據庫代碼**：從 20 行代碼縮減到 2 行
- **自動處理 SQL 查詢**：不需要手動寫 SELECT、INSERT、UPDATE 語句
- **類型安全**：編譯時就能發現數據類型錯誤

### 2. **智能功能**
- **自動遷移**：數據庫表結構會自動根據 Go 結構體創建和更新
- **關聯查詢**：自動處理表之間的關係（一對一、一對多、多對多）
- **軟刪除**：刪除數據時不會真正從數據庫移除，只是標記為已刪除
- **時間戳**：自動管理 `CreatedAt`、`UpdatedAt` 欄位

#### 🧠 智能功能詳細說明

**🔄 自動遷移 (Auto Migration)**

GORM 會自動比較你的 Go 結構體和數據庫表結構，自動創建、更新表結構。

```go
// ❌ 傳統方式：需要手動寫 SQL
func createTable() {
    db.Exec(`
        CREATE TABLE IF NOT EXISTS books (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            author VARCHAR(255) NOT NULL,
            price DECIMAL(10,2) NOT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW()
        );
    `)
}

// ✅ GORM 方式：自動處理
type BookGORM struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `gorm:"size:255;not null" json:"title"`
    Author    string    `gorm:"size:255;not null" json:"author"`
    Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// 只需要這一行，GORM 會自動創建表
db.AutoMigrate(&BookGORM{})
```

**智能之處：**
- 如果你在結構體中**新增欄位**，GORM 會自動在數據庫中添加對應的列
- 如果你**修改欄位類型**，GORM 會自動更新數據庫結構
- 如果你**刪除欄位**，GORM 會自動從數據庫中移除對應的列

**🔗 關聯查詢 (Association Queries)**

GORM 能自動處理表之間的關係，不需要手動寫 JOIN 查詢。

```go
// 定義關聯關係
type BookGORM struct {
    ID         uint           `gorm:"primaryKey" json:"id"`
    Title      string         `json:"title"`
    Author     string         `json:"author"`
    CategoryID *uint         `json:"category_id"`
    CategoryRef *CategoryGORM `gorm:"foreignKey:CategoryID" json:"category_ref"`
    Reviews    []ReviewGORM   `gorm:"foreignKey:BookID" json:"reviews"`
}

type CategoryGORM struct {
    ID    uint        `gorm:"primaryKey" json:"id"`
    Name  string      `json:"name"`
    Books []BookGORM  `gorm:"foreignKey:CategoryID" json:"books"`
}

type ReviewGORM struct {
    ID     uint      `gorm:"primaryKey" json:"id"`
    BookID uint      `json:"book_id"`
    Rating int       `json:"rating"`
    Comment string   `json:"comment"`
    Book   BookGORM  `gorm:"foreignKey:BookID" json:"book"`
}
```

**智能查詢範例：**

```go
// ❌ 傳統方式：需要手動寫複雜的 JOIN
func getBookWithCategoryAndReviews(bookID uint) {
    rows, err := db.Query(`
        SELECT b.id, b.title, b.author, c.name as category_name,
               r.id as review_id, r.rating, r.comment
        FROM books_gorm b
        LEFT JOIN categories_gorm c ON b.category_id = c.id
        LEFT JOIN reviews_gorm r ON b.id = r.book_id
        WHERE b.id = ?
    `, bookID)
    // ... 複雜的結果處理
}

// ✅ GORM 方式：自動處理關聯
func getBookWithCategoryAndReviews(bookID uint) {
    var book BookGORM
    db.Preload("CategoryRef").Preload("Reviews").First(&book, bookID)
    // 自動載入關聯數據，book.CategoryRef 和 book.Reviews 都已經填充
}
```

**智能之處：**
- **一對一關聯**：自動處理主鍵-外鍵關係
- **一對多關聯**：自動載入相關的多條記錄
- **多對多關聯**：自動處理中間表
- **預載入**：避免 N+1 查詢問題

**🗑️ 軟刪除 (Soft Delete)**

數據不會真正從數據庫中刪除，只是標記為已刪除，可以恢復。

```go
// ❌ 傳統方式：數據真的被刪除
func deleteBook(id uint) {
    db.Exec("DELETE FROM books WHERE id = ?", id)
    // 數據永遠消失了！
}

// ✅ GORM 方式：軟刪除
type BookGORM struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Title     string         `json:"title"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 軟刪除欄位
}

func deleteBook(id uint) {
    db.Delete(&BookGORM{}, id)
    // 數據還在，只是 deleted_at 欄位被設置了時間戳
}
```

**智能之處：**
- **自動過濾**：查詢時自動過濾已刪除的記錄
- **恢復功能**：可以恢復"已刪除"的數據
- **審計追蹤**：保留刪除歷史記錄
- **數據安全**：避免意外刪除重要數據

```go
// 查詢時自動過濾已刪除的記錄
db.Find(&books) // 只會返回未刪除的書籍

// 包含已刪除的記錄
db.Unscoped().Find(&books) // 返回所有記錄，包括已刪除的

// 恢復已刪除的記錄
db.Unscoped().Model(&book).Update("deleted_at", nil)
```

**⏰ 時間戳 (Timestamps)**

GORM 會自動管理 `CreatedAt` 和 `UpdatedAt` 欄位。

```go
// ❌ 傳統方式：需要手動管理時間
func createBook(book *Book) {
    now := time.Now()
    book.CreatedAt = now
    book.UpdatedAt = now
    
    db.Exec(`
        INSERT INTO books (title, author, price, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
    `, book.Title, book.Author, book.Price, book.CreatedAt, book.UpdatedAt)
}

func updateBook(book *Book) {
    book.UpdatedAt = time.Now()
    
    db.Exec(`
        UPDATE books 
        SET title = ?, author = ?, price = ?, updated_at = ?
        WHERE id = ?
    `, book.Title, book.Author, book.Price, book.UpdatedAt, book.ID)
}

// ✅ GORM 方式：自動管理時間戳
type BookGORM struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `json:"title"`
    Author    string    `json:"author"`
    Price     float64   `json:"price"`
    CreatedAt time.Time `json:"created_at"` // 自動設置
    UpdatedAt time.Time `json:"updated_at"` // 自動更新
}

func createBook(book *BookGORM) {
    db.Create(book) // CreatedAt 和 UpdatedAt 自動設置為當前時間
}

func updateBook(book *BookGORM) {
    db.Save(book) // UpdatedAt 自動更新為當前時間
}
```

**智能之處：**
- **創建時**：`CreatedAt` 和 `UpdatedAt` 自動設置為當前時間
- **更新時**：`UpdatedAt` 自動更新為當前時間
- **一致性**：所有記錄都有統一的時間格式
- **審計追蹤**：可以追蹤數據的創建和修改時間

#### 🎯 為什麼這些功能是"智能"的？

1. **自動化程度高**
   - 不需要手動寫 SQL
   - 不需要手動管理時間戳
   - 不需要手動處理關聯關係

2. **錯誤減少**
   - 類型安全，編譯時檢查
   - 自動處理 SQL 注入防護
   - 自動處理數據類型轉換

3. **開發效率**
   - 代碼量減少 70%
   - 維護成本降低
   - 學習曲線平緩

4. **功能豐富**
   - 內建分頁、排序、篩選
   - 自動索引創建
   - 查詢優化

這些"智能功能"讓 GORM 成為一個真正智能的 ORM 框架，大大提升了 Go 開發者的生產力和代碼質量！

### 3. **開發者友好**
- **鏈式查詢**：可以像 `db.Where().Order().Limit().Find()` 這樣鏈式調用
- **預載入關聯**：一次查詢就能獲取相關聯的所有數據
- **批量操作**：一次處理多條記錄
- **分頁查詢**：內建分頁功能

## 📋 快速開始指南

### 步驟 1：環境準備

首先確保你的系統已經安裝了：
- **Go 1.19+**
- **PostgreSQL 12+**（或其他支持的數據庫）
- **Git**

### 步驟 2：數據庫設置

1. **安裝 PostgreSQL**（如果還沒有）：
   ```bash
   # macOS
   brew install postgresql
   brew services start postgresql
   
   # Ubuntu/Debian
   sudo apt-get install postgresql postgresql-contrib
   sudo systemctl start postgresql
   
   # Windows
   # 下載並安裝 PostgreSQL 官方安裝包
   ```

2. **創建數據庫**：
   ```bash
   # 連接到 PostgreSQL
   psql -U postgres
   
   # 創建數據庫
   CREATE DATABASE books_db;
   
   # 創建用戶（可選）
   CREATE USER books_user WITH PASSWORD 'your_password';
   GRANT ALL PRIVILEGES ON DATABASE books_db TO books_user;
   ```

### 步驟 3：環境變數配置

在項目根目錄創建 `.env` 文件：

```bash
# 啟用數據庫模式
USE_DB=true

# 啟用 GORM（而不是原生 SQL）
USE_GORM=true

# 數據庫連接配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=books_db
DB_SSLMODE=disable

# GORM 日誌級別（可選）
GORM_LOG_LEVEL=info  # silent, error, warn, info
```

### 步驟 4：啟動應用

```bash
# 下載依賴
go mod tidy

# 啟動應用（必須設置數據庫環境變數）
go run main.go
```

**重要**：此應用現在**強制使用數據庫模式**，必須設置 `USE_DB=true` 和 `USE_GORM=true`。

你應該會看到類似這樣的輸出：
```
[BOOT] Book service: GORM database mode (production ready)
[GORM] Connected to PostgreSQL successfully
[GORM] Starting database migration...
[GORM] Database migration completed successfully
[BOOT] listening on :8080
```

如果沒有設置必要的環境變數，應用會立即退出並顯示錯誤：
```
[BOOT] ERROR: This application requires database mode. Please set USE_DB=true and USE_GORM=true
```

## 🗄️ 數據庫結構說明

### 核心概念：GORM 如何工作

GORM 會根據你的 Go 結構體自動創建數據庫表。每個結構體對應一個表，每個欄位對應一個列。

### 主要表結構

#### 1. **books_gorm** - 書籍表
```go
type BookGORM struct {
    ID         uint           `gorm:"primaryKey" json:"id"`
    Title      string         `gorm:"size:255;not null;index" json:"title"`
    Author     string         `gorm:"size:255;not null;index" json:"author"`
    Price      float64        `gorm:"type:decimal(10,2);not null" json:"price"`
    ISBN       string         `gorm:"size:20;uniqueIndex" json:"isbn"`
    Category   string         `gorm:"size:100" json:"category"`
    CategoryID *uint          `gorm:"index" json:"category_id"`
    Published  *time.Time     `json:"published"`
    CreatedAt  time.Time      `json:"created_at"`      // 自動管理
    UpdatedAt  time.Time      `json:"updated_at"`      // 自動管理
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"` // 軟刪除
}
```

**對應的 SQL 表**：
```sql
CREATE TABLE books_gorm (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    isbn VARCHAR(20) UNIQUE,
    category VARCHAR(100),
    category_id INTEGER,
    published TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
```

#### 2. **users_gorm** - 用戶表
- 用戶信息和認證數據
- 與書籍的多對多關聯（通過 `user_books_gorm` 中間表）

#### 3. **categories_gorm** - 分類表
- 書籍分類信息
- 與書籍的一對多關聯

#### 4. **reviews_gorm** - 評論表
- 書籍評論和評分
- 與書籍和用戶的多對一關聯

#### 5. **user_books_gorm** - 用戶收藏中間表
- 處理用戶和書籍的多對多關係

## 🔌 API 端點完整指南

### 認證說明
所有 API 端點都需要 JWT 認證（除了健康檢查）。首先需要登入獲取 token：

```bash
# 登入獲取 token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "Matt", "password": "password"}'

# 回應範例
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-01T12:00:00Z"
  }
}
```

### 基本 CRUD 操作
所有原有的 `/api/v1/books` 端點都支持 GORM，提供相同的功能但更強大的後端實現：

```bash
# 獲取所有書籍
GET /api/v1/books

# 獲取特定書籍
GET /api/v1/books/{id}

# 創建新書籍
POST /api/v1/books
{
  "title": "1984",
  "author": "George Orwell",
  "price": 12.99
}

# 更新書籍
PUT /api/v1/books/{id}

# 部分更新書籍
PATCH /api/v1/books/{id}

# 刪除書籍（軟刪除）
DELETE /api/v1/books/{id}

# 硬刪除書籍（永久刪除，無法恢復）
DELETE /api/v1/gorm/books/{id}/permanent

# 級聯硬刪除（刪除書籍及所有相關記錄）
DELETE /api/v1/gorm/books/{id}/cascade
```

### 🚀 GORM 專用端點

#### 🔍 搜索和統計功能（增強版）
```bash
# 基礎搜索（支持標題、作者、ISBN、分類）
GET /api/v1/gorm/search?q=1984
# 回應：包含相關性評分的搜索結果

# 增強版搜索（支持多條件篩選）
GET /api/v1/gorm/search?q=orwell&min_price=5&max_price=20&category=Fiction&year=1949
# 支援參數：q, min_price, max_price, category, year

# 高級搜索（支持複雜查詢條件）
GET /api/v1/gorm/search-advanced?title=1984&author=orwell&min_price=10&order_by=price_asc&limit=10
# 支援參數：title, author, category, isbn, min_price, max_price, year, order_by, limit, offset

# 獲取書籍統計信息
GET /api/v1/gorm/statistics
# 回應：總書籍數、平均價格、最貴/最便宜書籍等

# 獲取作者統計
GET /api/v1/gorm/author-statistics
# 回應：每個作者的書籍數量、平均價格、總銷售額

# 數據庫健康檢查
GET /api/v1/gorm/database-health
# 回應：數據庫連接狀態、表統計、性能指標
```

#### 🏷️ 分類和篩選功能
```bash
# 根據分類獲取書籍
GET /api/v1/gorm/category/Fiction
GET /api/v1/gorm/category/Science

# 根據價格範圍獲取書籍
GET /api/v1/gorm/price-range?min_price=10&max_price=50

# 根據出版年份獲取書籍
GET /api/v1/gorm/published/1949
GET /api/v1/gorm/published/2023

# 獲取評分最高的書籍
GET /api/v1/gorm/top-rated?limit=10

# 獲取最近添加的書籍
GET /api/v1/gorm/recent?limit=5

# 獲取帶有評論的書籍
GET /api/v1/gorm/with-reviews
```

#### 📄 分頁和批量操作
```bash
# 分頁獲取書籍
GET /api/v1/gorm/paginated?page=1&page_size=10
# 回應：包含總數、頁數、當前頁數據

# 根據多個作者獲取書籍
GET /api/v1/gorm/by-authors?authors=George Orwell,Aldous Huxley,Ray Bradbury
```

#### ⚠️ 硬刪除功能（危險操作）
```bash
# 永久刪除書籍（從數據庫中完全移除，無法恢復）
DELETE /api/v1/gorm/books/{id}/permanent
# 回應：包含警告信息的確認消息

# 級聯硬刪除（刪除書籍及所有相關的評論記錄）
DELETE /api/v1/gorm/books/{id}/cascade
# 回應：包含警告信息的確認消息

# 注意：硬刪除操作無法撤銷，請謹慎使用！

## 🧪 完整測試範例

### 步驟 1：獲取認證 Token
```bash
# 登入獲取 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "Matt", "password": "password"}' \
  | jq -r '.data.token')

echo "Token: $TOKEN"
```

### 步驟 2：測試基本 CRUD 操作
```bash
# 1. 獲取所有書籍
echo "=== 獲取所有書籍 ==="
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/books | jq

# 2. 創建新書籍
echo "=== 創建新書籍 ==="
curl -X POST http://localhost:8080/api/v1/books \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "price": 15.99
  }' | jq

# 3. 獲取特定書籍
echo "=== 獲取特定書籍 ==="
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/books/1 | jq
```

### 步驟 3：測試 GORM 高級功能
```bash
# 1. 基礎搜索功能
echo "=== 基礎搜索書籍 ==="
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/gorm/search?q=1984" | jq

# 2. 增強版搜索（多條件篩選）
echo "=== 增強版搜索 ==="
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/gorm/search?q=orwell&min_price=5&max_price=20&category=Fiction" | jq

# 3. 高級搜索（複雜查詢）
echo "=== 高級搜索 ==="
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/gorm/search-advanced?title=1984&author=orwell&order_by=price_asc&limit=5" | jq

# 4. 統計信息
echo "=== 書籍統計 ==="
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/gorm/statistics | jq

# 5. 作者統計
echo "=== 作者統計 ==="
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/gorm/author-statistics | jq

# 6. 數據庫健康檢查
echo "=== 數據庫健康檢查 ==="
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/gorm/database-health | jq
```

### 步驟 4：測試分類和篩選功能
```bash
# 1. 根據分類獲取書籍
echo "=== 根據分類獲取書籍 ==="
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/gorm/category/Fiction | jq

# 2. 根據價格範圍獲取書籍
echo "=== 根據價格範圍獲取書籍 ==="
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/gorm/price-range?min_price=10&max_price=20" | jq

# 3. 分頁獲取書籍
echo "=== 分頁獲取書籍 ==="
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/gorm/paginated?page=1&page_size=3" | jq
```

### 步驟 5：測試錯誤處理
```bash
# 1. 無效的 token
echo "=== 測試無效 token ==="
curl -H "Authorization: Bearer invalid_token" \
  http://localhost:8080/api/v1/books

# 2. 不存在的書籍 ID
echo "=== 測試不存在的書籍 ==="
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/books/99999

# 3. 無效的搜索參數
echo "=== 測試無效搜索 ==="
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/gorm/search?q="
```

### 步驟 6：測試硬刪除功能（危險操作）
```bash
# ⚠️ 警告：以下操作會永久刪除數據，無法恢復！

# 1. 先創建一個測試書籍
echo "=== 創建測試書籍 ==="
TEST_BOOK_ID=$(curl -s -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "測試書籍",
    "author": "測試作者",
    "price": 9.99,
    "isbn": "TEST-123",
    "category": "測試分類"
  }' | jq -r '.data.id')

echo "測試書籍 ID: $TEST_BOOK_ID"

# 2. 軟刪除（可以恢復）
echo "=== 軟刪除測試 ==="
curl -X DELETE -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/books/$TEST_BOOK_ID | jq

# 3. 硬刪除（永久刪除）
echo "=== 硬刪除測試 ==="
curl -X DELETE -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/gorm/books/$TEST_BOOK_ID/permanent | jq

# 4. 級聯硬刪除（刪除書籍及相關記錄）
echo "=== 級聯硬刪除測試 ==="
curl -X DELETE -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/gorm/books/$TEST_BOOK_ID/cascade | jq

## ⚡ 性能優勢與最佳實踐

### 🚀 性能優勢

#### 1. **查詢優化**
- **自動 SQL 優化**：GORM 會自動優化生成的 SQL 查詢
- **預載入關聯**：使用 `Preload()` 避免 N+1 查詢問題
- **智能索引**：根據查詢模式自動建議索引創建
- **查詢緩存**：與 Redis 緩存完美整合

#### 2. **連接管理**
- **連接池**：自動管理數據庫連接池，提高並發性能
- **健康檢查**：定期檢查連接健康狀態
- **自動重連**：連接斷開時自動重新連接

#### 3. **內存優化**
- **批量操作**：支持批量插入、更新、刪除
- **分頁查詢**：內建分頁功能，避免一次性載入大量數據
- **選擇性載入**：使用 `Select()` 只載入需要的欄位

### 📚 最佳實踐

#### 1. **模型設計原則**
```go
// ✅ 好的設計
type BookGORM struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Title     string         `gorm:"size:255;not null;index" json:"title"`
    Author    string         `gorm:"size:255;not null;index" json:"author"`
    Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ❌ 避免的設計
type BadBook struct {
    ID    uint   `json:"id"`
    Title string `json:"title"` // 沒有索引，沒有長度限制
    Data  string `json:"data"`  // 沒有明確用途的欄位
}
```

#### 2. **查詢優化技巧**
```go
// ✅ 使用預載入避免 N+1 問題
var books []BookGORM
db.Preload("Reviews").Preload("CategoryRef").Find(&books)

// ✅ 只選擇需要的欄位
db.Select("id, title, author").Find(&books)

// ✅ 使用索引欄位進行查詢
db.Where("author = ?", "George Orwell").Find(&books)

// ❌ 避免全表掃描
db.Where("LOWER(title) LIKE ?", "%1984%").Find(&books)
```

#### 3. **錯誤處理模式**
```go
// ✅ 適當的錯誤處理
func (s *BookServiceGORM) CreateBook(book models.BookGORM) (*models.BookGORM, error) {
    if err := s.db.Create(&book).Error; err != nil {
        logging.Logger.Printf("[GORM] Failed to create book: %v", err)
        return nil, fmt.Errorf("failed to create book: %w", err)
    }
    return &book, nil
}
```

## 🔧 故障排除指南

### 常見問題與解決方案

#### 1. **連接問題**
```bash
# 問題：無法連接到數據庫
# 錯誤：dial tcp: connection refused

# 解決方案：
# 1. 檢查 PostgreSQL 是否運行
sudo systemctl status postgresql  # Linux
brew services list | grep postgres  # macOS

# 2. 檢查連接參數
echo $DB_HOST $DB_PORT $DB_USER $DB_NAME

# 3. 測試連接
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME
```

#### 2. **遷移問題**
```bash
# 問題：表結構衝突
# 錯誤：relation "books_gorm" already exists

# 解決方案：
# 1. 檢查現有表結構
psql -d books_db -c "\d books_gorm"

# 2. 手動刪除表（謹慎使用）
psql -d books_db -c "DROP TABLE IF EXISTS books_gorm CASCADE;"

# 3. 重新啟動應用
go run main.go
```

#### 3. **查詢性能問題**
```bash
# 問題：查詢速度慢
# 解決方案：

# 1. 啟用 GORM 日誌查看 SQL
export GORM_LOG_LEVEL=info

# 2. 檢查是否有適當的索引
psql -d books_db -c "\d books_gorm"

# 3. 分析查詢計劃
psql -d books_db -c "EXPLAIN ANALYZE SELECT * FROM books_gorm WHERE author = 'George Orwell';"
```

#### 4. **關聯查詢問題**
```bash
# 問題：關聯查詢失敗
# 錯誤：invalid field found for struct

# 解決方案：
# 1. 檢查外鍵定義
grep -n "foreignKey" models/gorm_models.go

# 2. 確保外鍵欄位存在
grep -n "CategoryID" models/gorm_models.go

# 3. 重新遷移
go run main.go
```

### 🔍 調試技巧

#### 1. **啟用詳細日誌**
```bash
# 在 .env 文件中添加
GORM_LOG_LEVEL=info
```

#### 2. **檢查數據庫狀態**
```bash
# 使用健康檢查端點
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/gorm/database-health
```

#### 3. **監控性能指標**
```bash
# 查看應用指標
curl http://localhost:8080/metrics
curl http://localhost:8080/api/metrics/detailed
```

## 🏗️ 生產環境架構

### 強制數據庫模式
此應用現在**專為生產環境設計**，強制使用 GORM + PostgreSQL：

- ✅ **數據持久化**：所有數據都存儲在 PostgreSQL 中
- ✅ **高級搜索**：支持複雜查詢、相關性評分、多條件篩選
- ✅ **關聯查詢**：自動處理表關係和預載入
- ✅ **軟刪除**：數據安全，支持恢復
- ✅ **自動遷移**：數據庫結構自動同步
- ✅ **緩存整合**：與 Redis 緩存完美配合
- ✅ **監控指標**：內建 Prometheus 指標收集
- ✅ **GORM Hooks**：自動化業務邏輯和數據一致性

### 與傳統方案的比較

| 功能特性 | 傳統 SQL | 內存模式 | GORM 生產模式 |
|---------|----------|----------|---------------|
| **數據持久化** | ✅ | ❌ | ✅ |
| **生產就緒** | ✅ | ❌ | ✅ |
| **代碼量** | 多 | 少 | 少 |
| **類型安全** | 低 | 中 | 高 |
| **自動遷移** | 無 | 無 | ✅ |
| **關聯查詢** | 手動 | 無 | 自動 |
| **軟刪除** | 手動 | 無 | 自動 |
| **時間戳** | 手動 | 無 | 自動 |
| **高級搜索** | 複雜 | 無 | ✅ |
| **緩存整合** | 手動 | 無 | 自動 |
| **監控指標** | 手動 | 無 | 自動 |

## 🎯 總結

GORM 整合讓你的 Go API 具備了現代 ORM 的所有優勢：

- **🚀 開發效率提升 70%**：更少的代碼，更多的功能
- **🛡️ 類型安全**：編譯時錯誤檢查，減少運行時問題
- **🔄 自動化**：遷移、關聯、時間戳全自動處理
- **📈 性能優化**：智能查詢優化和緩存整合
- **🔧 易於維護**：清晰的代碼結構和豐富的調試工具

無論你是 GORM 新手還是經驗豐富的開發者，這個整合都能讓你的 Go API 開發更加高效和愉快！

## 🔗 GORM Hooks - 生命週期管理

### 什麼是 GORM Hooks？

**GORM Hooks** 是 GORM 提供的生命週期回調函數，讓你在數據庫操作的特定時機自動執行自定義邏輯。就像"鉤子"一樣，在關鍵時刻"鉤住"並執行你的代碼。

### 🎯 Hooks 的完整生命週期

```
創建流程：
BeforeCreate → Validate → Create → AfterCreate

更新流程：
BeforeUpdate → Validate → Update → AfterUpdate

刪除流程：
BeforeDelete → Delete → AfterDelete
```

### 🚀 實際應用場景

#### 1. **BeforeCreate Hook - 創建前處理**
```go
func (b *BookGORM) BeforeCreate(tx *gorm.DB) error {
    // 數據清理和標準化
    b.Title = strings.TrimSpace(b.Title)
    b.Author = strings.TrimSpace(b.Author)
    
    // 數據驗證
    if b.Price < 0 {
        return errors.New("price must be positive")
    }
    
    // 自動設置默認值
    if b.Category == "" {
        b.Category = "General"
    }
    
    // 自動生成 ISBN
    if b.ISBN == "" {
        b.ISBN = generateISBN()
    }
    
    // 檢查重複數據
    var existingBook BookGORM
    err := tx.Where("title = ? AND author = ?", b.Title, b.Author).First(&existingBook).Error
    if err == nil {
        return errors.New("book with same title and author already exists")
    }
    
    return nil
}
```

#### 2. **AfterCreate Hook - 創建後處理**
```go
func (b *BookGORM) AfterCreate(tx *gorm.DB) error {
    // 更新統計數據
    updateBookCount(tx)
    updateAuthorStatistics(tx, b.Author)
    updateCategoryStatistics(tx, b.Category)
    
    // 自動創建分類記錄
    ensureCategoryExists(tx, b.Category)
    
    // 記錄審計日誌
    logAudit(tx, "BOOK_CREATED", b.ID, fmt.Sprintf("Book '%s' created", b.Title))
    
    // 發送業務通知
    sendNotification("BOOK_CREATED", map[string]interface{}{
        "book_id": b.ID,
        "title":   b.Title,
        "author":  b.Author,
    })
    
    // 更新緩存
    invalidateCache("books")
    
    return nil
}
```

#### 3. **BeforeUpdate Hook - 更新前處理**
```go
func (b *BookGORM) BeforeUpdate(tx *gorm.DB) error {
    // 獲取原始數據用於比較
    var originalBook BookGORM
    tx.First(&originalBook, b.ID)
    
    // 檢查價格變化
    if originalBook.Price != b.Price {
        // 如果價格大幅下降，觸發促銷通知
        if b.Price < originalBook.Price*0.8 {
            sendPromotionNotification(b.Title, b.Price, originalBook.Price)
        }
    }
    
    // 檢查 ISBN 重複
    if b.ISBN != "" {
        var existingBook BookGORM
        err := tx.Where("isbn = ? AND id != ?", b.ISBN, b.ID).First(&existingBook).Error
        if err == nil {
            return errors.New("ISBN already exists")
        }
    }
    
    return nil
}
```

#### 4. **BeforeDelete Hook - 刪除前檢查**
```go
func (b *BookGORM) BeforeDelete(tx *gorm.DB) error {
    // 檢查是否有相關評論
    var reviewCount int64
    tx.Model(&ReviewGORM{}).Where("book_id = ?", b.ID).Count(&reviewCount)
    if reviewCount > 0 {
        return errors.New("cannot delete book with reviews")
    }
    
    // 檢查是否有用戶收藏
    var favoriteCount int64
    tx.Model(&UserBookGORM{}).Where("book_id = ?", b.ID).Count(&favoriteCount)
    if favoriteCount > 0 {
        return errors.New("cannot delete book with user favorites")
    }
    
    return nil
}
```

### 💡 Hooks 的優勢

#### 1. **數據一致性**
- ✅ **自動驗證**：所有操作都經過相同的驗證邏輯
- ✅ **自動清理**：數據自動標準化和清理
- ✅ **重複檢查**：自動檢查重複數據

#### 2. **業務邏輯自動化**
- ✅ **統計更新**：自動更新相關統計數據
- ✅ **通知發送**：自動發送業務通知
- ✅ **審計日誌**：自動記錄操作日誌

#### 3. **代碼集中化**
- ✅ **邏輯集中**：業務邏輯集中在模型中
- ✅ **減少重複**：避免在每個 handler 中重複代碼
- ✅ **易於維護**：修改邏輯只需要改一個地方

#### 4. **事務安全**
- ✅ **原子操作**：在數據庫事務中執行
- ✅ **自動回滾**：如果 Hook 失敗，整個操作會回滾
- ✅ **數據完整性**：確保數據庫狀態一致

### 🧪 Hooks 測試範例

```bash
# 1. 測試 BeforeCreate Hook
echo "=== 測試 BeforeCreate Hook ==="
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "  GORM Hooks 測試  ",
    "author": "  Hook 作者  ",
    "price": 29.99
  }'

# 觀察日誌輸出：
# [GORM Hook] BeforeCreate: 開始創建書籍 'GORM Hooks 測試'
# [GORM Hook] BeforeCreate: 自動設置默認分類 'General'
# [GORM Hook] BeforeCreate: 自動生成 ISBN '123-456-789'
# [GORM Hook] BeforeCreate: 驗證通過，準備創建書籍 'GORM Hooks 測試'

# 2. 測試 AfterCreate Hook
# 觀察日誌輸出：
# [GORM Hook] AfterCreate: 書籍創建成功 'GORM Hooks 測試' (ID: 6)
# [GORM Hook Helper] updateBookCount: 書籍總數更新為 6
# [GORM Hook Helper] updateAuthorStatistics: 作者 'Hook 作者' 統計更新
# 📚 新書上架通知: GORM Hooks 測試

# 3. 測試 BeforeUpdate Hook
echo "=== 測試 BeforeUpdate Hook ==="
curl -X PUT http://localhost:8080/api/v1/books/6 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "GORM Hooks 測試",
    "author": "Hook 作者",
    "price": 19.99
  }'

# 觀察日誌輸出：
# [GORM Hook] BeforeUpdate: 價格變化 29.99 -> 19.99
# [GORM Hook] BeforeUpdate: 檢測到大幅降價，觸發促銷通知
# [GORM Hook Helper] sendPromotionNotification: 促銷通知 - 折扣: 33.3%

# 4. 測試 BeforeDelete Hook
echo "=== 測試 BeforeDelete Hook ==="
curl -X DELETE http://localhost:8080/api/v1/books/6 \
  -H "Authorization: Bearer $TOKEN"

# 觀察日誌輸出：
# [GORM Hook] BeforeDelete: 開始刪除書籍 'GORM Hooks 測試' (ID: 6)
# [GORM Hook] BeforeDelete: 刪除原因檢查完成
# [GORM Hook] AfterDelete: 書籍刪除成功 'GORM Hooks 測試' (ID: 6)
# 🗑️ 書籍下架通知: GORM Hooks 測試
```

### 🔧 自定義 Hooks

你可以根據需要添加更多自定義的 Hook 邏輯：

```go
// 自定義驗證函數
func (b *BookGORM) Validate() error {
    if b.Price < 0 {
        return errors.New("price must be positive")
    }
    if len(b.Title) < 1 {
        return errors.New("title cannot be empty")
    }
    return nil
}

// 自定義業務規則
func (b *BookGORM) CheckBusinessRules(tx *gorm.DB) error {
    // 檢查庫存
    if b.Stock < 0 {
        return errors.New("stock cannot be negative")
    }
    
    // 檢查價格範圍
    if b.Price > 1000 {
        return errors.New("price cannot exceed 1000")
    }
    
    return nil
}
```

### 📊 Hooks 性能影響

| Hook 類型 | 執行時機 | 性能影響 | 建議用途 |
|-----------|----------|----------|----------|
| **BeforeCreate** | 創建前 | 低 | 驗證、清理、設置默認值 |
| **AfterCreate** | 創建後 | 中 | 統計更新、通知發送 |
| **BeforeUpdate** | 更新前 | 低 | 變化檢查、重複驗證 |
| **AfterUpdate** | 更新後 | 中 | 統計更新、緩存清理 |
| **BeforeDelete** | 刪除前 | 低 | 依賴檢查、權限驗證 |
| **AfterDelete** | 刪除後 | 中 | 清理、統計更新 |

### ⚠️ 注意事項

1. **性能考慮**：Hooks 會在每次數據庫操作時執行，避免在 Hooks 中執行耗時操作
2. **錯誤處理**：Hook 返回錯誤會導致整個操作失敗，確保錯誤處理邏輯正確
3. **事務安全**：Hooks 在事務中執行，避免在 Hooks 中執行會影響事務的操作
4. **循環依賴**：避免在 Hooks 中觸發會再次調用相同 Hook 的操作
