package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1 定義Server
	srv := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(5 * time.Second)
			fmt.Fprintf(w, "處理完成！")
		}),
	}

	// 2. 在背景 (Goroutine) 啟動 Server
	// 這樣 main 才能繼續往下執行「等待訊號」的邏輯
	go func() {
		fmt.Println("Server started on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Listen 錯誤: %s\n", err)
		}
	}()

	// 3. 設定訊號監聽 (跟剛才學的一樣)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 4. 等待訊號 (阻塞)
	<-quit

	// 5. 收到訊號後，優雅關閉 Server
	fmt.Println("收到訊號，準備關閉 Server...")

	// 6. 設定超時時間，並在超時後強制關閉 Server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server 關閉錯誤: %s\n", err)
	}
	fmt.Println("Server 已關閉")
}
