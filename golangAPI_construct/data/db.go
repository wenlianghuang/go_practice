package data

import (
	"database/sql"
	"fmt"
	"golangAPI_construct/logging"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite" // 重要：載入 SQLite driver
)

func Open() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// 相對於「執行時的工作目錄」
		//dsn = "file:books.db?cache=shared&mode=rwc&_pragma=foreign_keys(ON)"
		// Use absolute path to avoid confusion
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		absPath, err := filepath.Abs(filepath.Join(cwd, "books.db"))
		if err != nil {
			return nil, err
		}
		// modernc sqlite DSN 建議使用 file: 前綴；Windows/Unix 通用用 / 分隔
		dsn = fmt.Sprintf("file:%s?cache=shared&mode=rwc&_pragma=foreign_keys(ON)", filepath.ToSlash(absPath))
	}
	//log.Printf("[DB Opening SQLite with DNS=%s]", dsn)
	logging.Logger.Printf("[DB] Opening SQLite with DSN=%s", dsn)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// SQLite 建議小連線池
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	log.Printf("[DB] connected dsn=%s", dsn)
	return db, nil
}
