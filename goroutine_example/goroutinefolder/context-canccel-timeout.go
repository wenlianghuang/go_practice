/*
Context 取消與超時（必備）
為 goroutine/worker 增加 context 支援，能在父層取消時快速收斂，避免 goroutine 洩漏。
展示 select { case <-ctx.Done(): ... } 的正確用法，以及搭配 time.After/Timeout 的模式。

學習重點：
1. context.Context 是用來控制 goroutine 生命週期的標準工具。
2. context.WithTimeout 可以設定自動取消的時間點。
3. ctx.Done() 會回傳一個 channel，當 Context 被取消或超時時，該 channel 會被關閉。
*/
package goroutinefolder

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// worker_cancel_timeout 模擬一個會持續工作的 worker
// 它同時監聽兩個來源：ctx.Done() (取消訊號) 與 jobs channel (任務來源)
func worker_cancel_timeout(ctx context.Context, id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			// 收到 Context 的取消訊號（可能是超時或主動調用 cancel()）
			fmt.Printf("worker %d: 收到取消訊號，結束工作\n", id)
			return
		case j, ok := <-jobs:
			if !ok {
				// jobs channel 已被關閉
				fmt.Printf("worker %d: 任務頻道已關閉\n", id)
				return
			}
			// 模擬耗時的處理過程
			fmt.Printf("worker %d: 正在處理任務 %d...\n", id, j)
			time.Sleep(150 * time.Millisecond)
			fmt.Printf("worker %d: 完成處理任務 %d\n", id, j)
		}
	}
}

func ContextCancelTimeout() {
	// 建立一個 600 毫秒後會自動超時的 Context
	// ctx 用於傳遞取消訊號，cancel 函數則用於手動觸發取消
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)

	// 務必使用 defer 調用 cancel()
	// 即使 Context 因為超時而自動取消，調用 cancel 仍能確保底層資源被即時釋放
	defer cancel()

	jobs := make(chan int)
	var wg sync.WaitGroup

	// 啟動 2 個 worker，並傳入 ctx 讓它們具備「被控制」的能力
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go worker_cancel_timeout(ctx, i, jobs, &wg)
	}

	// 啟動一個 goroutine 來產生任務
	go func() {
		defer close(jobs)
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				// 如果 Context 已經取消或超時，發送者也應該停止工作，避免無謂的阻塞
				fmt.Println("發送者：Context 已超時，停止發送任務")
				return
			case jobs <- i:
				// 成功將任務放入頻道
			}
		}
	}()

	// 等待所有 worker 結束（由 Context 取消或任務處理完畢觸發）
	wg.Wait()
	fmt.Println("主程式執行結束 (Done)")
}