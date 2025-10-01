package data

import (
	"database/sql"
	"fmt"
	"golangAPI_construct/logging"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func Open() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// PostgreSQL 預設連接字串
		/*
			host := getEnvOrDefault("DB_HOST", "localhost")
			port := getEnvOrDefault("DB_PORT", "5432")
			user := getEnvOrDefault("DB_USER", "postgres")
			password := getEnvOrDefault("DB_PASSWORD", "password")
			dbname := getEnvOrDefault("DB_NAME", "books_db")
			sslmode := getEnvOrDefault("DB_SSLMODE", "disable")
		*/
		// 直接從環境變數讀取，不設置預設值
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
			sslmode = "disable" // 只有這個可以有預設值
		}
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)
	}

	logging.Logger.Printf("[DB] Opening PostgreSQL with DSN=%s", maskPassword(dsn))
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// PostgreSQL 連線池設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	log.Printf("[DB] connected to PostgreSQL")
	return db, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 遮蔽密碼用於日誌記錄
func maskPassword(dsn string) string {
	// 簡單的密碼遮蔽，實際使用可以更複雜
	return "host=... user=... dbname=... (password masked)"
}
