package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var Logger *log.Logger

// Init sets up the global logger. Returns the opened *os.File so call can defer Close().

func Init() (*os.File, error) {
	logPath := os.Getenv("LOG_FILE")
	if logPath == "" {
		ts := time.Now().UTC().Format("20060102T150405Z")
		logPath = filepath.Join("logs", "app_"+ts+".log")
	}

	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	mw := io.MultiWriter(os.Stdout, f)

	// 自訂旗標: 日期時間 (UTC)
	Logger = log.New(mw, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
	Logger.Printf("[INIT] logging started file=%s", logPath)
	return f, nil
}

// Optional helper to add simple time-based separator (call daily via cron/goroutine if needed)
func DailySeparator() {
	if Logger != nil {
		Logger.Printf("======== %s ========", time.Now().UTC().Format(time.RFC3339))
	}
}
