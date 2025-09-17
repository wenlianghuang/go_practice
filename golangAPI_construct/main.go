// ...existing code...
package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"golangAPI_construct/logging"
	"golangAPI_construct/routes"
)

func main() {
	// 開發環境：自動載入專案根目錄 .env（不存在則忽略）
	_ = godotenv.Load()

	// 命令列旗標（旗標 > 環境變數 > 預設）
	useDB := flag.Bool("use-db", false, "use database backend (overrides USE_DB env)")
	dbDSN := flag.String("db-dsn", "", "database DSN (overrides DB_DSN env)")
	jwtSecret := flag.String("jwt-secret", "", "JWT secret (overrides JWT_SECRET env)")
	port := flag.String("port", "", "HTTP port (overrides PORT env)")
	flag.Parse()

	if *useDB {
		os.Setenv("USE_DB", "true")
	}
	if *dbDSN != "" {
		os.Setenv("DB_DSN", *dbDSN)
	}
	if *jwtSecret != "" {
		os.Setenv("JWT_SECRET", *jwtSecret)
	}
	if *port != "" {
		os.Setenv("PORT", *port)
	}

	logging.Init()

	router := routes.SetupRoutes()

	addr := ":" + getenvDefault("PORT", "8080")
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logging.Logger.Printf("[BOOT] listening on %s", addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Fatalf("[FATAL] listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.Logger.Fatalf("[FATAL] server shutdown: %v", err)
	}
	logging.Logger.Println("[BOOT] server stopped")
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
