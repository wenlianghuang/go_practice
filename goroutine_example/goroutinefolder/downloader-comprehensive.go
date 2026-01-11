package goroutinefolder

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// DownloadTask 代表一個下載任務
type DownloadTask struct {
	URL string
	ID  int
}

// DownloadManager 管理下載狀態
type DownloadManager struct {
	downloaded   map[string]bool // 已下載的網址（需要 Mutex 保護）
	mutex        sync.Mutex      // 保護 map 的互斥鎖
	totalBytes   atomic.Uint64   // 總下載位元組數（使用 atomic）
	successCount atomic.Uint32   // 成功下載數
	skippedCount atomic.Uint32   // 跳過數（重複）
	timeoutCount atomic.Uint32   // 超時數
}

// NewDownloadManager 建立新的下載管理器
func NewDownloadManager() *DownloadManager {
	return &DownloadManager{
		downloaded: make(map[string]bool),
	}
}

// ⚠️ 警告：以下兩個方法分開使用會有 race condition！
// IsDownloaded 檢查網址是否已下載（線程安全，但不應與 MarkDownloaded 分開使用）
// 請使用 CheckAndMark 代替這兩個方法的組合
func (dm *DownloadManager) IsDownloaded(url string) bool {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	return dm.downloaded[url]
}

// MarkDownloaded 標記網址為已下載（線程安全，但不應與 IsDownloaded 分開使用）
// 請使用 CheckAndMark 代替這兩個方法的組合
func (dm *DownloadManager) MarkDownloaded(url string) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	dm.downloaded[url] = true
}

// ✅ 推薦使用：CheckAndMark 原子性地檢查並標記網址
// 這是正確的做法，避免了 "Check-Then-Act" 的 race condition
// 返回 true 表示可以下載（已佔位），false 表示已被其他工人下載
func (dm *DownloadManager) CheckAndMark(url string) bool {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// 🔑 關鍵：在同一個鎖內完成檢查和標記，保證原子性
	if dm.downloaded[url] {
		return false // 已下載，不可以再下載
	}

	dm.downloaded[url] = true // 立即標記（佔位），防止其他工人重複下載
	return true               // 可以下載
}

// simulateDownload 模擬下載網址（返回下載的位元組數）
func simulateDownload(ctx context.Context, url string) (int, error) {
	// 模擬下載時間 300-800ms
	downloadTime := time.Millisecond * time.Duration(300+rand.Intn(500))

	// 使用 select 檢查是否超時
	select {
	case <-time.After(downloadTime):
		// 下載完成，返回模擬的位元組數
		bytes := 1024 + rand.Intn(9216) // 1KB - 10KB
		return bytes, nil
	case <-ctx.Done():
		// Context 取消（超時）
		return 0, ctx.Err()
	}
}

// worker 工人函數，處理下載任務
func downloadWorker(
	id int,
	ctx context.Context,
	tasks <-chan DownloadTask,
	dm *DownloadManager,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				// Channel 已關閉，工人退出
				fmt.Printf("👷 工人 %d 完成所有任務\n", id)
				return
			}

			// ✅ 使用 CheckAndMark 原子性地檢查並標記（防止 race condition）
			// 這個方法在同一個鎖內完成「檢查是否已下載」和「標記為已下載」
			// 確保只有一個工人能夠獲得下載權限
			canDownload := dm.CheckAndMark(task.URL)
			if !canDownload {
				// 已被其他工人標記為下載中，跳過
				dm.skippedCount.Add(1)
				fmt.Printf("⏭️  工人 %d: 跳過 [%s]（已下載）\n", id, task.URL)
				continue
			}

			// 🎯 到這裡表示已經成功佔位，可以安心下載
			// 其他工人不可能同時下載相同的網址
			fmt.Printf("📥 工人 %d: 開始下載 [%s]\n", id, task.URL)
			bytes, err := simulateDownload(ctx, task.URL)

			if err != nil {
				// 超時或取消
				dm.timeoutCount.Add(1)
				fmt.Printf("⏰ 工人 %d: 下載 [%s] 超時\n", id, task.URL)
				return // 超時後工人退出
			}

			// 下載成功（網址已在 CheckAndMark 時標記）
			dm.totalBytes.Add(uint64(bytes))
			dm.successCount.Add(1)
			fmt.Printf("✅ 工人 %d: 完成 [%s] (%d bytes)\n", id, task.URL, bytes)

		case <-ctx.Done():
			// Context 超時，工人退出
			fmt.Printf("⏰ 工人 %d: 收到超時信號，停止工作\n", id)
			return
		}
	}
}

// DownloaderWithAllFeatures 綜合練習：下載器
func DownloaderWithAllFeatures() {
	fmt.Println("=== 綜合練習：網址下載器 ===")
	fmt.Println("📋 功能：Worker Pool + Mutex + Atomic + Context + WaitGroup")
	fmt.Println()

	// 模擬的網址列表（有重複）
	urls := []string{
		"https://example.com/page1",
		"https://example.com/page2",
		"https://example.com/page3",
		"https://example.com/page4",
		"https://example.com/page5",
		"https://example.com/page1", // 重複
		"https://example.com/page6",
		"https://example.com/page7",
		"https://example.com/page2", // 重複
		"https://example.com/page8",
		"https://example.com/page9",
		"https://example.com/page10",
		"https://example.com/page3", // 重複
		"https://example.com/page11",
		"https://example.com/page12",
		"https://example.com/page13",
		"https://example.com/page4", // 重複
		"https://example.com/page14",
		"https://example.com/page15",
		"https://example.com/page5", // 重複
	}

	// 建立下載管理器
	dm := NewDownloadManager()

	// 建立 Context，設定 2 秒超時
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 建立任務 Channel
	tasks := make(chan DownloadTask, len(urls))

	// 建立 WaitGroup
	var wg sync.WaitGroup

	// 啟動 5 個工人
	numWorkers := 5
	fmt.Printf("🚀 啟動 %d 個工人...\n", numWorkers)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go downloadWorker(w, ctx, tasks, dm, &wg)
	}

	// 派發任務
	fmt.Printf("📦 派發 %d 個下載任務...\n\n", len(urls))
	startTime := time.Now()

	go func() {
		for i, url := range urls {
			tasks <- DownloadTask{URL: url, ID: i + 1}
		}
		close(tasks) // 所有任務派發完畢，關閉 channel
	}()

	// 等待所有工人完成或超時
	wg.Wait()
	duration := time.Since(startTime)

	// 印出最終報表
	fmt.Println("\n" + repeatString("=", 60))
	fmt.Println("📊 最終報表")
	fmt.Println(repeatString("=", 60))
	fmt.Printf("⏱️  總耗時: %v\n", duration)
	fmt.Printf("👷 工人數量: %d\n", numWorkers)
	fmt.Printf("📋 總任務數: %d\n", len(urls))
	fmt.Println()

	success := dm.successCount.Load()
	skipped := dm.skippedCount.Load()
	timeout := dm.timeoutCount.Load()
	totalBytes := dm.totalBytes.Load()

	fmt.Printf("✅ 成功下載: %d 個網址\n", success)
	fmt.Printf("⏭️  跳過重複: %d 個網址\n", skipped)
	fmt.Printf("⏰ 超時任務: %d 個\n", timeout)
	fmt.Printf("📊 總流量: %s\n", formatBytes(totalBytes))

	if totalBytes > 0 {
		avgBytes := totalBytes / uint64(success)
		fmt.Printf("📈 平均大小: %s/個\n", formatBytes(avgBytes))
	}

	fmt.Println()
	fmt.Println("🔑 使用的技術：")
	fmt.Println("   ✓ Worker Pool: 5 個固定工人")
	fmt.Println("   ✓ Mutex: 保護 downloaded map（去重）")
	fmt.Println("   ✓ Atomic: 累加流量統計（無鎖併發）")
	fmt.Println("   ✓ Context: 2 秒超時控制")
	fmt.Println("   ✓ WaitGroup: 等待所有工人完成")
	fmt.Println("   ✓ Channel: 任務分發")

	if duration >= 2*time.Second {
		fmt.Println("\n⚠️  注意：達到 2 秒超時限制，部分任務未完成")
	}
}

// DownloaderWithShortTimeout 短超時版本（1 秒）
func DownloaderWithShortTimeout() {
	fmt.Println("=== 短超時版本（1 秒）===")
	fmt.Println()

	urls := []string{
		"https://example.com/page1",
		"https://example.com/page2",
		"https://example.com/page3",
		"https://example.com/page4",
		"https://example.com/page5",
		"https://example.com/page6",
		"https://example.com/page7",
		"https://example.com/page8",
		"https://example.com/page9",
		"https://example.com/page10",
	}

	dm := NewDownloadManager()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tasks := make(chan DownloadTask, len(urls))
	var wg sync.WaitGroup

	numWorkers := 5
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go downloadWorker(w, ctx, tasks, dm, &wg)
	}

	for i, url := range urls {
		tasks <- DownloadTask{URL: url, ID: i + 1}
	}
	close(tasks)

	startTime := time.Now()
	wg.Wait()
	duration := time.Since(startTime)

	fmt.Println("\n📊 報表：")
	fmt.Printf("⏱️  耗時: %v\n", duration)
	fmt.Printf("✅ 成功: %d\n", dm.successCount.Load())
	fmt.Printf("⏰ 超時: %d\n", dm.timeoutCount.Load())
	fmt.Printf("📊 流量: %s\n", formatBytes(dm.totalBytes.Load()))

	if duration >= 1*time.Second {
		fmt.Println("⚠️  達到 1 秒超時限制")
	}
}

// DownloaderWithLongTimeout 長超時版本（5 秒）
func DownloaderWithLongTimeout() {
	fmt.Println("=== 長超時版本（5 秒）===")
	fmt.Println()

	urls := []string{
		"https://example.com/page1",
		"https://example.com/page2",
		"https://example.com/page3",
		"https://example.com/page1", // 重複
		"https://example.com/page4",
		"https://example.com/page5",
		"https://example.com/page2", // 重複
		"https://example.com/page6",
		"https://example.com/page7",
		"https://example.com/page8",
	}

	dm := NewDownloadManager()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tasks := make(chan DownloadTask, len(urls))
	var wg sync.WaitGroup

	numWorkers := 3
	fmt.Printf("🚀 啟動 %d 個工人（較少工人，觀察去重效果）\n\n", numWorkers)

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go downloadWorker(w, ctx, tasks, dm, &wg)
	}

	for i, url := range urls {
		tasks <- DownloadTask{URL: url, ID: i + 1}
	}
	close(tasks)

	startTime := time.Now()
	wg.Wait()
	duration := time.Since(startTime)

	fmt.Println("\n📊 報表：")
	fmt.Printf("⏱️  耗時: %v\n", duration)
	fmt.Printf("✅ 成功: %d\n", dm.successCount.Load())
	fmt.Printf("⏭️  跳過: %d（去重機制生效）\n", dm.skippedCount.Load())
	fmt.Printf("📊 流量: %s\n", formatBytes(dm.totalBytes.Load()))
	fmt.Println("\n💡 在 5 秒內應該所有任務都能完成，重點觀察去重機制")
}

// formatBytes 格式化位元組數
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// repeatString 重複字串
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
