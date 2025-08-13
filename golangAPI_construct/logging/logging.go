package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var Logger *log.Logger

func Init() (*os.File, error) {
	logPath := os.Getenv("LOG_FILE")
	if logPath == "" {
		ts := time.Now().UTC().Format("2006_01_02_15_04_05")
		logPath = filepath.Join("logs", "Present_"+ts+".log")
	}
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	mw := io.MultiWriter(os.Stdout, f)
	Logger = log.New(mw, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
	Logger.Printf("[INIT] logging started file=%s", logPath)
	return f, nil
}

func DailySeparator() {
	if Logger != nil {
		Logger.Printf("======== %s ========", time.Now().UTC().Format(time.RFC3339))
	}
}
