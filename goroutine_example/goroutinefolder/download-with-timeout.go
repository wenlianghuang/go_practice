package goroutinefolder

import (
	"fmt"
	"time"
)

// 模擬下載檔案的函式
func downloadFile(filename string, duration time.Duration, ch chan string) {
	fmt.Printf("開始下載 %s (需要 %v)\n", filename, duration)
	time.Sleep(duration) // 模擬下載時間
	ch <- filename       // 下載完成，發送檔名到 channel
}

// DownloadWithTimeout 演示使用 select 和 time.After 處理超時
func DownloadWithTimeout() {
	fmt.Println("=== 檔案下載超時控制範例 ===")

	// 建立一個 channel 用來接收下載進度
	downloadCh := make(chan string)

	// 啟動三個 Goroutine，模擬下載不同檔案
	go downloadFile("檔案 A", 800*time.Millisecond, downloadCh)
	go downloadFile("檔案 B", 1200*time.Millisecond, downloadCh)
	go downloadFile("檔案 C", 2000*time.Millisecond, downloadCh)

	// 設定總時限
	timeout := 1500 * time.Millisecond
	fmt.Printf("設定超時時限: %v\n\n", timeout)

	// 記錄成功下載的檔案數量
	successCount := 0
	totalFiles := 3

	// 使用 select 監聽下載完成或超時
	for successCount < totalFiles {
		select {
		case filename := <-downloadCh:
			// 成功接收到下載完成的檔名
			successCount++
			fmt.Printf("✓ 已成功下載 [%s] (%d/%d)\n", filename, successCount, totalFiles)

		case <-time.After(timeout):
			// 超時發生
			fmt.Printf("\n⚠ 下載超時，取消剩餘任務\n")
			fmt.Printf("成功: %d/%d，失敗: %d/%d\n", successCount, totalFiles, totalFiles-successCount, totalFiles)
			return
		}
	}

	fmt.Println("\n所有檔案下載完成！")
}

// DownloadWithTimeoutV2 演示使用 context 來取消剩餘任務（進階版）
func DownloadWithTimeoutV2() {
	fmt.Println("=== 檔案下載超時控制範例 (一次性超時) ===")

	// 建立一個 channel 用來接收下載進度
	downloadCh := make(chan string, 3) // 使用 buffered channel 避免 goroutine 阻塞

	// 啟動三個 Goroutine，模擬下載不同檔案
	go downloadFile("檔案 A", 800*time.Millisecond, downloadCh)
	go downloadFile("檔案 B", 1200*time.Millisecond, downloadCh)
	go downloadFile("檔案 C", 2000*time.Millisecond, downloadCh)

	// 設定總時限（從現在開始計時）
	deadline := time.After(1500 * time.Millisecond)
	fmt.Printf("設定超時時限: 1.5 秒\n\n")

	// 記錄成功下載的檔案數量
	successCount := 0
	totalFiles := 3

	// 使用 for-select 持續監聽
	for successCount < totalFiles {
		select {
		case filename := <-downloadCh:
			// 成功接收到下載完成的檔名
			successCount++
			fmt.Printf("✓ 已成功下載 [%s] (%d/%d)\n", filename, successCount, totalFiles)

		case <-deadline:
			// 超時發生（只會觸發一次）
			fmt.Printf("\n⚠ 下載超時，取消剩餘任務\n")
			fmt.Printf("成功: %d/%d，失敗: %d/%d\n", successCount, totalFiles, totalFiles-successCount, totalFiles)

			// 等待一小段時間，讓還在執行的 goroutine 有機會完成並發送（可選）
			time.Sleep(100 * time.Millisecond)

			// 檢查是否有額外的下載在超時後完成
			close(downloadCh)
			for filename := range downloadCh {
				fmt.Printf("  (超時後完成: %s)\n", filename)
			}
			return
		}
	}

	fmt.Println("\n所有檔案下載完成！")
}
