/*
sync/atomic 原子操作範例
展示如何使用原子操作來避免競態條件（Race Condition），
以及與普通變數操作的對比。
*/
package goroutinefolder

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// AtomicOperation 比較非安全計數器與原子計數器
func AtomicOperation() {
	var unsafeCounter int64 // 普通變數 (不安全)
	var safeCounter int64   // 將被用於原子操作的變數 (安全)

	var wg sync.WaitGroup
	totalGoroutines := 1000

	wg.Add(totalGoroutines)

	fmt.Printf("啟動 %d 個 Goroutine 進行並發累加...\n", totalGoroutines)

	for i := 0; i < totalGoroutines; i++ {
		go func() {
			defer wg.Done()

			// --- 1. 非原子操作 (Unsafe) ---
			// 這行代碼在機器碼層面其實是三個步驟：讀取 -> 修改 -> 寫入
			// 當多個 Goroutine 同時執行時，會發生「競態條件」(Race Condition)
			// 導致數值被覆蓋，結果通常會小於預期。
			unsafeCounter++

			// --- 2. 原子操作 (Safe) ---
			// 使用 atomic 套件直接在底層 CPU 指令級別鎖定記憶體位置
			// 保證這個「加 1」的動作是不可分割的 (Atomic)。
			atomic.AddInt64(&safeCounter, 1)

			// 稍微暫停一下，增加並發衝突的機會
			time.Sleep(time.Millisecond)
		}()
	}

	wg.Wait()

	fmt.Println("\n=== 執行結果比較 ===")
	fmt.Printf("預期數值: %d\n", totalGoroutines)
	
	// 這裡讀取 unsafeCounter 其實理論上也是不安全的，但在範例結束時讀取通常沒問題
	fmt.Printf("普通變數 (Unsafe): %d  <-- 發生 Race Condition，數值錯誤\n", unsafeCounter)
	
	// 讀取原子變數時，建議也使用 LoadInt64 以確保記憶體一致性
	safeVal := atomic.LoadInt64(&safeCounter)
	fmt.Printf("原子變數 (Safe):   %d  <-- 數值正確\n", safeVal)
}
