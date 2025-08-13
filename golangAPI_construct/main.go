package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golangAPI_construct/logging"
	"golangAPI_construct/routes"
)

func main() {
	// 初始化日誌
	logFile, err := logging.Init()
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer logFile.Close()

	router := routes.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logging.Logger.Printf("[BOOT] server listening on %s (PID=%d)", srv.Addr, os.Getpid())

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Fatalf("[FATAL] listen error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	logging.Logger.Printf("[SHUTDOWN] signal received: %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logging.Logger.Printf("[WARN] graceful shutdown issue: %v; forcing close", err)
		_ = srv.Close()
	} else {
		logging.Logger.Println("[SHUTDOWN] graceful shutdown complete")
	}
}
