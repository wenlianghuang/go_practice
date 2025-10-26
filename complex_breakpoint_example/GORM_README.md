# GORM + PostgreSQL 集成指南

## 已完成的工作

### 1. 添加 GORM 和 PostgreSQL 依賴
- ✅ 已更新 `go.mod` 加入 GORM 和 PostgreSQL 驅動
- ✅ 已安裝所有依賴

### 2. 更新配置系統
- ✅ 已更新 `config/config.go` 支持數據庫連接字符串
- ✅ 支持通過環境變量配置數據庫參數

### 3. 創建 GORM 數據庫實現
- ✅ 創建 `database/database_gorm.go` - 完整的 GORM 實現
- ✅ 實現自動遷移（AutoMigrate）
- ✅ 實現事務支持
- ✅ 創建 `services/services_gorm.go` - GORM 版本的服務層

## 需要完成的工作

### 1. 統一 ID 類型
**問題**: 當前 `models/models.go` 使用 `uint` 類型作為 ID，但 `database/database.go`（內存數據庫）仍使用 `int` 類型，導致類型不匹配。

**解決方案選擇**:
- **選項 A**: 將所有數據庫（內存和 GORM）統一使用 `uint` 類型
  - 需要修改 `database/database.go` 中所有 ID 相關的類型
  - 修改 map 的 key 類型從 `int` 改為 `uint`
  - 這是推薦的方案，因為 GORM 默認使用 `uint` 作為主鍵

- **選項 B**: 保持模型使用 `int`，在 GORM 層做轉換
  - 需要為 GORM 創建單獨的模型結構
  - 在 GORM 數據庫層做類型轉換

### 2. 更新 Handlers 和 Routes
- 需要創建 GORM 版本的 handlers（或修改現有的 handlers 以支持兩種服務類型）
- 修改 `routes/routes.go` 以動態選擇使用哪種數據庫實現

### 3. 更新 main.go
目前的 main.go 已支持基本結構，但需要：
- 實現 Handler 類型的動態創建
- 或創建適配器模式以橋接兩種服務類型

## 如何使用 GORM 數據庫

### 選項 1: 使用內存數據庫（推薦，無需額外安裝）

目前項目默認使用內存數據庫，無需安裝任何額外軟件：

```bash
# 直接運行，使用內存數據庫
go run main.go
```

**優點**:
- ✅ 無需安裝任何軟件
- ✅ 開箱即用
- ✅ 適合開發和測試
- ✅ 完全功能

### 選項 2: 安裝本地 PostgreSQL

#### macOS
```bash
# 使用 Homebrew 安裝
brew install postgresql@15
brew services start postgresql@15

# 創建數據庫
createdb breakpoint_db

# 設置密碼（如果需要的話）
psql postgres
# 在 psql 中執行:
# ALTER USER $USER PASSWORD 'postgres';
```

#### Linux (Ubuntu/Debian)
```bash
# 安裝 PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# 啟動服務
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 創建數據庫
sudo -u postgres createdb breakpoint_db
```

#### Windows
1. 從 [PostgreSQL 官網](https://www.postgresql.org/download/windows/) 下載安裝程序
2. 運行安裝程序，按提示完成安裝
3. 創建數據庫：
```bash
# 在命令提示符中
createdb -U postgres breakpoint_db
```

### 選項 3: 使用 Docker（需先安裝 Docker）

如果您之後想安裝 Docker，可以從以下方式安裝：

#### macOS
```bash
# 安裝 Docker Desktop
brew install --cask docker
# 或者從 https://www.docker.com/products/docker-desktop 下載
```

#### Linux
```bash
# Ubuntu/Debian
sudo apt install docker.io
sudo systemctl start docker

# 或用戶添加到 docker 組
sudo usermod -aG docker $USER
```

#### Windows
從 [Docker Desktop](https://www.docker.com/products/docker-desktop) 下載安裝

### 配置和運行（如果使用 PostgreSQL）

```bash
# 設置環境變量
export DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=breakpoint_db sslmode=disable"

# 如果使用 Docker
docker run --name postgres-db \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=breakpoint_db \
  -p 5432:5432 -d postgres

# 運行應用（需要修改 main.go 支持 GORM）
# USE_GORM=true go run main.go
```

## 架構說明

### 數據庫層
1. **內存數據庫** (`database/database.go`): 用於測試和開發
2. **GORM 數據庫** (`database/database_gorm.go`): 用於生產環境，連接 PostgreSQL

### 服務層
1. **標準服務** (`services/services.go`): 為內存數據庫提供服務
2. **GORM 服務** (`services/services_gorm.go`): 為 GORM 數據庫提供服務

兩者實現相同的接口，可以通過工廠模式或依賴注入動態選擇。

## 自動遷移

GORM 會自動創建以下表：
- `users` - 用戶表
- `transactions` - 交易表
- `bank_accounts` - 銀行帳戶表
- `loan_applications` - 貸款申請表

## 事務支持

所有財務操作（存款、提款、轉帳）都使用數據庫事務，確保數據一致性。

## 下一步建議

1. **完成 ID 類型統一**（選擇選項 A，推薦）
2. **創建 GORM Handlers** 或修改現有 Handlers 支持接口
3. **完整測試 GORM 實現**
4. **添加數據庫連接池配置**
5. **添加查詢性能優化**

## 文件結構

```
complex_breakpoint_example/
├── config/
│   └── config.go           # 支持數據庫配置
├── database/
│   ├── database.go          # 內存數據庫實現
│   └── database_gorm.go     # GORM 數據庫實現
├── models/
│   ├── models.go            # 數據模型（目前使用 uint）
│   └── gorm_models.go       # (已刪除，不再需要)
├── services/
│   ├── services.go          # 標準服務
│   └── services_gorm.go     # GORM 服務
├── handlers/
│   └── handlers.go          # 需要更新以支持兩種服務
├── routes/
│   └── routes.go            # 需要更新以支持兩種服務
├── main.go                  # 需要完善 GORM 支持
└── GORM_README.md           # 本文件
```

## 測試

```bash
# 測試內存數據庫
go test ./...

# 測試 GORM（需要 PostgreSQL 運行）
USE_GORM=true go test ./database/...
```

## 注意事項

1. GORM 使用軟刪除（軟刪除標記已添加到模型中）
2. 所有時間戳字段由 GORM 自動管理
3. 需要手動管理數據庫連接的關閉
4. 建議在生產環境中使用連接池

