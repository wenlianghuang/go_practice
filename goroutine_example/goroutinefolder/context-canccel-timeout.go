/*
Context 取消與超時（必備）
為 goroutine/worker 增加 context 支援，能在父層取消時快速收斂，避免 goroutine 洩漏。
展示 select { case <-ctx.Done(): ... } 的正確用法，以及搭配 time.After/Timeout 的模式。
*/
package goroutinefolder

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker_cancel_timeout(ctx context.Context, id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d: canceled\n", id)
			return
		case j, ok := <-jobs:
			if !ok {
				fmt.Printf("worker %d: jobs closed\n", id)
				return
			}
			// 模擬處理
			time.Sleep(150 * time.Millisecond)
			fmt.Printf("worker %d: processed %d\n", id, j)
		}
	}
}

func ContextCancelTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	jobs := make(chan int)
	var wg sync.WaitGroup

	// 開兩個 worker，都支援 ctx 取消
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go worker_cancel_timeout(ctx, i, jobs, &wg)
	}

	// 產生一些任務
	go func() {
		defer close(jobs)
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return
			case jobs <- i:
			}
		}
	}()

	wg.Wait()
	fmt.Println("done")
}
