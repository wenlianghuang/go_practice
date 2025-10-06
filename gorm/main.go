package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// --- 1. 定義模型 (Models) ---

// User 模型 - gorm.Model 內建了 ID, CreatedAt, UpdatedAt, DeletedAt
type User struct {
	gorm.Model
	Name       string     `gorm:"size:255;not null"`
	Email      string     `gorm:"unique;not null"`
	CreditCard CreditCard // 一對一關係：一個 User 擁有一張 CreditCard
}

// CreditCard 模型
type CreditCard struct {
	gorm.Model
	Number string `gorm:"unique;not null"`
	UserID uint   // 外鍵，GORM 會自動關聯到 User 的 ID
}

func main() {
	// --- 2. 連接到 PostgreSQL 資料庫 ---

	// DSN (Data Source Name) - 請替換成您自己的資料庫連線資訊
	// 格式: "host=主機 user=使用者 password=密碼 dbname=資料庫名稱 port=埠號 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := "host=localhost user=matthuang password=wenliang75 dbname=test_db port=5432 sslmode=disable TimeZone=Asia/Taipei"

	// 建立新的 Logger，以便我們可以清楚看到 GORM 執行的 SQL
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 閾值
			LogLevel:                  logger.Info, // Log 等級
			IgnoreRecordNotFoundError: true,        // 忽略 ErrRecordNotFound 錯誤
			Colorful:                  true,        // 啟用彩色日誌
		},
	)

	// 連接資料庫
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // 使用我們自訂的 Logger
	})
	if err != nil {
		panic("無法連接到資料庫")
	}

	fmt.Println("✅ 資料庫連接成功！")

	// --- 3. 自動遷移 (Auto Migration) ---
	// GORM 會自動檢查並創建 User 和 CreditCard 資料表
	err = db.AutoMigrate(&User{}, &CreditCard{})
	if err != nil {
		panic("資料庫遷移失敗")
	}
	fmt.Println("✅ 資料庫遷移成功！")

	// --- 4. 創建記錄 (Create) ---
	fmt.Println("\n--- 創建使用者與信用卡 ---")
	user := User{
		Name:  "Jinzhu",
		Email: "jinzhu@example.com",
		CreditCard: CreditCard{
			Number: "1234-5678-9012-3456",
		},
	}
	// GORM 會在一個事務中同時創建 User 和 CreditCard
	result := db.Create(&user)
	if result.Error != nil {
		fmt.Println("創建失敗:", result.Error)
	} else {
		fmt.Printf("🎉 成功創建使用者: %s (ID: %d)\n", user.Name, user.ID)
	}

	// --- 5. 查詢記錄 (Read) ---
	fmt.Println("\n--- 查詢使用者 ---")
	var retrievedUser User
	// First 會根據主鍵查詢第一筆記錄
	// Preload("CreditCard") 會使用預加載 (Eager Loading) 一併查詢關聯的信用卡資料
	// 這樣可以避免 N+1 查詢問題
	db.Preload("CreditCard").First(&retrievedUser, "email = ?", "jinzhu@example.com")

	if retrievedUser.ID > 0 {
		fmt.Printf("🔍 查詢到使用者: %s, Email: %s\n", retrievedUser.Name, retrievedUser.Email)
		fmt.Printf("💳 他的信用卡號碼是: %s\n", retrievedUser.CreditCard.Number)
	} else {
		fmt.Println("找不到該使用者")
	}

	// --- 6. 更新記錄 (Update) ---
	fmt.Println("\n--- 更新使用者名稱 ---")
	// 使用 Model 指定要更新的對象，然後用 Update 更新單一欄位
	db.Model(&retrievedUser).Update("Name", "Jinzhu (Admin)")
	fmt.Printf("🔧 使用者名稱已更新為: %s\n", retrievedUser.Name)

	// --- 7. 刪除記錄 (Delete) ---
	fmt.Println("\n--- 刪除使用者 ---")
	// GORM 的刪除是「軟刪除 (Soft Delete)」
	// 記錄不會真的從資料庫消失，而是將 DeletedAt 欄位標記上時間
	db.Delete(&retrievedUser)
	fmt.Printf("🗑️ 已軟刪除使用者: %s\n", retrievedUser.Name)

	// --- 驗證軟刪除 ---
	var checkUser User
	// 正常的 First 查詢會自動過濾掉被軟刪除的記錄
	db.First(&checkUser, "email = ?", "jinzhu@example.com")
	if checkUser.ID == 0 {
		fmt.Println("✅ 驗證成功：正常查詢無法找到被軟刪除的使用者。")
	}

	// 如果想查詢包含軟刪除的記錄，可以使用 Unscoped
	db.Unscoped().First(&checkUser, "email = ?", "jinzhu@example.com")
	if checkUser.ID != 0 {
		fmt.Println("✅ 驗證成功：使用 Unscoped() 可以找到被軟刪除的使用者。")
	}
}
