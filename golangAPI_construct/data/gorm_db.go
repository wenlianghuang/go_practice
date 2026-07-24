package data

import (
	"fmt"
	"os"
	"time"

	"golangAPI_construct/logging"
	"golangAPI_construct/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenGORM 創建 GORM 數據庫連接
// 提供更強大的 ORM 功能，包括自動遷移、關聯查詢等
func OpenGORM() (*gorm.DB, error) {
	// 從環境變數讀取數據庫配置
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// 構建 DSN
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		sslmode := os.Getenv("DB_SSLMODE")

		// 檢查必要的環境變數
		if host == "" || port == "" || user == "" || password == "" || dbname == "" {
			return nil, fmt.Errorf("missing required database environment variables")
		}

		if sslmode == "" {
			sslmode = "disable"
		}
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)
	}

	// 配置 GORM 日誌
	var gormLogger logger.Interface
	if os.Getenv("GORM_LOG_LEVEL") == "silent" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.New(
			logging.Logger,
			logger.Config{
				SlowThreshold:             time.Second, // 慢查詢閾值
				LogLevel:                  logger.Info, // 日誌級別
				IgnoreRecordNotFoundError: true,        // 忽略記錄未找到錯誤
				Colorful:                  false,       // 禁用彩色輸出（適合日誌文件）
			},
		)
	}

	// 創建 GORM 配置
	config := &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// 連接數據庫
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 獲取底層 sql.DB 進行連接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	// 配置連接池
	sqlDB.SetMaxOpenConns(25)                  // 最大打開連接數
	sqlDB.SetMaxIdleConns(10)                  // 最大空閒連接數
	sqlDB.SetConnMaxLifetime(time.Hour)        // 連接最大生命週期
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 連接最大空閒時間

	// 測試連接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %v", err)
	}

	logging.Logger.Printf("[GORM] Connected to PostgreSQL successfully")
	return db, nil
}

// MigrateGORM 執行 GORM 自動遷移
// 自動創建和更新數據庫表結構
func MigrateGORM(db *gorm.DB) error {
	logging.Logger.Print("[GORM] Starting database migration...")

	// 自動遷移所有模型
	err := db.AutoMigrate(
		&models.BookGORM{},
		&models.UserGORM{},
		&models.CategoryGORM{},
		&models.ReviewGORM{},
		&models.UserBookGORM{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	logging.Logger.Print("[GORM] Database migration completed successfully")
	return nil
}

// SeedGORM 為 GORM 數據庫填充初始數據
func SeedGORM(db *gorm.DB) error {
	logging.Logger.Print("[GORM] Starting database seeding...")

	// 檢查是否已有數據
	var count int64
	db.Model(&models.BookGORM{}).Count(&count)
	if count > 0 {
		logging.Logger.Print("[GORM] Database already has data, skipping seeding")
		return nil
	}

	// 創建初始書籍數據
	books := []models.BookGORM{
		{
			Title:     "1984",
			Author:    "George Orwell",
			Price:     9.99,
			ISBN:      "978-0-452-28423-4",
			Category:  "Dystopian Fiction",
			Published: timePtr(1949, 6, 8),
		},
		{
			Title:     "Brave New World",
			Author:    "Aldous Huxley",
			Price:     8.99,
			ISBN:      "978-0-06-085052-4",
			Category:  "Science Fiction",
			Published: timePtr(1932, 1, 1),
		},
		{
			Title:     "To Kill a Mockingbird",
			Author:    "Harper Lee",
			Price:     12.99,
			ISBN:      "978-0-06-112008-4",
			Category:  "Fiction",
			Published: timePtr(1960, 7, 11),
		},
		{
			Title:     "The Great Gatsby",
			Author:    "F. Scott Fitzgerald",
			Price:     10.99,
			ISBN:      "978-0-7432-7356-5",
			Category:  "Fiction",
			Published: timePtr(1925, 4, 10),
		},
		{
			Title:     "Pride and Prejudice",
			Author:    "Jane Austen",
			Price:     11.99,
			ISBN:      "978-0-14-143951-8",
			Category:  "Romance",
			Published: timePtr(1813, 1, 28),
		},
	}

	// 批量創建書籍
	if err := db.CreateInBatches(books, 10).Error; err != nil {
		return fmt.Errorf("failed to seed books: %v", err)
	}

	// 創建初始分類
	categories := []models.CategoryGORM{
		{Name: "Fiction", Description: "Literary works of fiction"},
		{Name: "Science Fiction", Description: "Speculative fiction with scientific elements"},
		{Name: "Dystopian Fiction", Description: "Fiction depicting dystopian societies"},
		{Name: "Romance", Description: "Romantic literature"},
		{Name: "Mystery", Description: "Mystery and detective fiction"},
		{Name: "Fantasy", Description: "Fantasy literature"},
		{Name: "Non-Fiction", Description: "Factual and informative works"},
	}

	if err := db.CreateInBatches(categories, 10).Error; err != nil {
		return fmt.Errorf("failed to seed categories: %v", err)
	}

	// 創建初始用戶
	users := []models.UserGORM{
		{
			Username:     "admin",
			Email:        "admin@example.com",
			PasswordHash: "$2a$10$AQuMpFYbHBfGx2F2bS0.x.Nm.YTFzwjHaznp9uUCN9V5t3sweZ4w6", // password: "password"
			Roles:        "admin,editor",
			FirstName:    "Admin",
			LastName:     "User",
			IsActive:     true,
		},
		{
			Username:     "reader",
			Email:        "reader@example.com",
			PasswordHash: "$2a$10$AQuMpFYbHBfGx2F2bS0.x.Nm.YTFzwjHaznp9uUCN9V5t3sweZ4w6", // password: "password"
			Roles:        "user",
			FirstName:    "Book",
			LastName:     "Reader",
			IsActive:     true,
		},
	}

	if err := db.CreateInBatches(users, 10).Error; err != nil {
		return fmt.Errorf("failed to seed users: %v", err)
	}

	logging.Logger.Printf("[GORM] Database seeding completed successfully - %d books, %d categories, %d users",
		len(books), len(categories), len(users))
	return nil
}

// GetGORMStats 獲取 GORM 數據庫統計信息
func GetGORMStats(db *gorm.DB) map[string]interface{} {
	stats := make(map[string]interface{})

	// 書籍統計
	var bookCount int64
	db.Model(&models.BookGORM{}).Count(&bookCount)
	stats["books_count"] = bookCount

	// 用戶統計
	var userCount int64
	db.Model(&models.UserGORM{}).Count(&userCount)
	stats["users_count"] = userCount

	// 分類統計
	var categoryCount int64
	db.Model(&models.CategoryGORM{}).Count(&categoryCount)
	stats["categories_count"] = categoryCount

	// 評論統計
	var reviewCount int64
	db.Model(&models.ReviewGORM{}).Count(&reviewCount)
	stats["reviews_count"] = reviewCount

	// 平均價格
	var avgPrice float64
	db.Model(&models.BookGORM{}).Select("AVG(price)").Scan(&avgPrice)
	stats["average_price"] = avgPrice

	// 最高價格
	var maxPrice float64
	db.Model(&models.BookGORM{}).Select("MAX(price)").Scan(&maxPrice)
	stats["max_price"] = maxPrice

	// 最低價格
	var minPrice float64
	db.Model(&models.BookGORM{}).Select("MIN(price)").Scan(&minPrice)
	stats["min_price"] = minPrice

	return stats
}

// CloseGORM 關閉 GORM 數據庫連接
func CloseGORM(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// timePtr 創建時間指針的輔助函數
func timePtr(year int, month time.Month, day int) *time.Time {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return &t
}

// CreateGORMIndexes 創建 GORM 索引
func CreateGORMIndexes(db *gorm.DB) error {
	logging.Logger.Print("[GORM] Creating database indexes...")

	// 為書籍表創建索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_books_author ON books_gorm(author)").Error; err != nil {
		return fmt.Errorf("failed to create author index: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_books_category ON books_gorm(category)").Error; err != nil {
		return fmt.Errorf("failed to create category index: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_books_price ON books_gorm(price)").Error; err != nil {
		return fmt.Errorf("failed to create price index: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_books_published ON books_gorm(published)").Error; err != nil {
		return fmt.Errorf("failed to create published index: %v", err)
	}

	// 為用戶表創建索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users_gorm(email)").Error; err != nil {
		return fmt.Errorf("failed to create email index: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users_gorm(username)").Error; err != nil {
		return fmt.Errorf("failed to create username index: %v", err)
	}

	// 為評論表創建索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_reviews_book_id ON reviews_gorm(book_id)").Error; err != nil {
		return fmt.Errorf("failed to create book_id index: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews_gorm(user_id)").Error; err != nil {
		return fmt.Errorf("failed to create user_id index: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews_gorm(rating)").Error; err != nil {
		return fmt.Errorf("failed to create rating index: %v", err)
	}

	logging.Logger.Print("[GORM] Database indexes created successfully")
	return nil
}
