package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golangAPI_construct/routes"
)

func main() {
	// 建立路由
	router := routes.SetupRoutes()

	// 讀取埠號（預設 8080）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 建立 http.Server 以支援優雅關閉與逾時控制
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second, // 限制讀取 request header + body 的時間
		WriteTimeout: 10 * time.Second, // 限制寫回應的時間
		IdleTimeout:  60 * time.Second, // keep-alive 連線閒置上限
	}

	log.Printf("[BOOT] server listening on %s (PID=%d)", srv.Addr, os.Getpid())

	// 啟動伺服器（非阻塞）
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] listen error: %v", err)
		}
	}()

	// 捕捉系統訊號（Ctrl+C -> SIGINT；容器停止 -> SIGTERM）
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞直到收到關閉訊號
	sig := <-stop
	log.Printf("[SHUTDOWN] signal received: %s", sig.String())

	// 給在途請求最多 8 秒完成
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[WARN] graceful shutdown timeout or error: %v; forcing close", err)
		_ = srv.Close()
	} else {
		log.Println("[SHUTDOWN] graceful shutdown complete")
	}
}
