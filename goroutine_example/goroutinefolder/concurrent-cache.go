package goroutinefolder

import (
	"fmt"
	"sync"
	"time"
)

// Cache 不安全的版本（會有 race condition）
type UnsafeCache struct {
	data map[string]string
}

func NewUnsafeCache() *UnsafeCache {
	return &UnsafeCache{
		data: make(map[string]string),
	}
}

func (c *UnsafeCache) Set(key, value string) {
	c.data[key] = value
}

func (c *UnsafeCache) Get(key string) string {
	return c.data[key]
}

// SafeCache 使用 Mutex 保護的安全版本
type SafeCache struct {
	data  map[string]string
	mutex sync.Mutex
}

func NewSafeCache() *SafeCache {
	return &SafeCache{
		data: make(map[string]string),
	}
}

func (c *SafeCache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}

func (c *SafeCache) Get(key string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.data[key]
}

// RWMutexCache 使用 RWMutex 的版本（讀寫鎖，性能更好）
type RWMutexCache struct {
	data  map[string]string
	mutex sync.RWMutex
}

func NewRWMutexCache() *RWMutexCache {
	return &RWMutexCache{
		data: make(map[string]string),
	}
}

func (c *RWMutexCache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}

func (c *RWMutexCache) Get(key string) string {
	c.mutex.RLock() // 讀鎖，允許多個 Goroutine 同時讀取
	defer c.mutex.RUnlock()
	return c.data[key]
}

// ConcurrentCacheUnsafe 演示不安全的版本（會有 race condition）
func ConcurrentCacheUnsafe() {
	fmt.Println("=== 不安全的 Cache 版本（會有 race condition）===")
	fmt.Println("⚠️  警告：請使用 'go run -race main.go ConcurrentCacheUnsafe' 來檢測 race condition")
	fmt.Println()

	cache := NewUnsafeCache()
	var wg sync.WaitGroup

	// 啟動 50 個 Goroutine 同時操作
	numGoroutines := 50
	operationsPerGoroutine := 100

	startTime := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 每個 Goroutine 執行多次 Set 和 Get 操作
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key_%d", id%10) // 使用 10 個不同的 key
				value := fmt.Sprintf("value_%d_%d", id, j)

				// 寫入操作
				cache.Set(key, value)

				// 讀取操作
				_ = cache.Get(key)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Printf("✓ 完成 %d 個 Goroutine，每個執行 %d 次操作\n", numGoroutines, operationsPerGoroutine)
	fmt.Printf("✓ 總操作次數: %d\n", numGoroutines*operationsPerGoroutine*2)
	fmt.Printf("✓ 執行時間: %v\n", duration)
	fmt.Println("\n⚠️  注意：雖然程序沒有崩潰，但這個版本有 race condition！")
	fmt.Println("   請使用 'go run -race main.go ConcurrentCacheUnsafe' 來檢測問題")
}

// ConcurrentCacheSafe 演示使用 Mutex 的安全版本
func ConcurrentCacheSafe() {
	fmt.Println("=== 使用 Mutex 的安全 Cache 版本 ===")
	fmt.Println()

	cache := NewSafeCache()
	var wg sync.WaitGroup

	// 啟動 50 個 Goroutine 同時操作
	numGoroutines := 50
	operationsPerGoroutine := 100

	startTime := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 每個 Goroutine 執行多次 Set 和 Get 操作
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key_%d", id%10) // 使用 10 個不同的 key
				value := fmt.Sprintf("value_%d_%d", id, j)

				// 寫入操作
				cache.Set(key, value)

				// 讀取操作
				_ = cache.Get(key)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Printf("✓ 完成 %d 個 Goroutine，每個執行 %d 次操作\n", numGoroutines, operationsPerGoroutine)
	fmt.Printf("✓ 總操作次數: %d (Set + Get)\n", numGoroutines*operationsPerGoroutine*2)
	fmt.Printf("✓ 執行時間: %v\n", duration)
	fmt.Println("✓ 這個版本使用 Mutex，是併發安全的！")
}

// ConcurrentCacheRWMutex 演示使用 RWMutex 的版本（性能更好）
func ConcurrentCacheRWMutex() {
	fmt.Println("=== 使用 RWMutex 的 Cache 版本（讀寫鎖）===")
	fmt.Println()

	cache := NewRWMutexCache()
	var wg sync.WaitGroup

	// 啟動 50 個 Goroutine 同時操作
	numGoroutines := 50
	operationsPerGoroutine := 100

	startTime := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 每個 Goroutine 執行多次 Set 和 Get 操作
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key_%d", id%10) // 使用 10 個不同的 key
				value := fmt.Sprintf("value_%d_%d", id, j)

				// 寫入操作
				cache.Set(key, value)

				// 讀取操作（多次讀取以展示 RWMutex 的優勢）
				_ = cache.Get(key)
				_ = cache.Get(key)
				_ = cache.Get(key)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Printf("✓ 完成 %d 個 Goroutine，每個執行 %d 次寫入和 %d 次讀取\n",
		numGoroutines, operationsPerGoroutine, operationsPerGoroutine*3)
	fmt.Printf("✓ 總操作次數: %d\n", numGoroutines*operationsPerGoroutine*4)
	fmt.Printf("✓ 執行時間: %v\n", duration)
	fmt.Println("✓ RWMutex 允許多個 Goroutine 同時讀取，性能更好！")
}

// ConcurrentCacheComparison 比較三種版本的性能
func ConcurrentCacheComparison() {
	fmt.Println("=== Cache 實作比較（安全版本） ===")
	fmt.Println()

	numGoroutines := 50
	operationsPerGoroutine := 100

	// 測試 Mutex 版本
	fmt.Println("1️⃣  測試 Mutex 版本...")
	cacheM := NewSafeCache()
	var wg sync.WaitGroup

	startM := time.Now()
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key_%d", id%10)
				value := fmt.Sprintf("value_%d_%d", id, j)
				cacheM.Set(key, value)
				_ = cacheM.Get(key)
				_ = cacheM.Get(key)
				_ = cacheM.Get(key)
			}
		}(i)
	}
	wg.Wait()
	durationM := time.Since(startM)

	// 測試 RWMutex 版本
	fmt.Println("2️⃣  測試 RWMutex 版本...")
	cacheRW := NewRWMutexCache()

	startRW := time.Now()
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("key_%d", id%10)
				value := fmt.Sprintf("value_%d_%d", id, j)
				cacheRW.Set(key, value)
				_ = cacheRW.Get(key)
				_ = cacheRW.Get(key)
				_ = cacheRW.Get(key)
			}
		}(i)
	}
	wg.Wait()
	durationRW := time.Since(startRW)

	// 結果比較
	fmt.Println("\n📊 性能比較結果：")
	fmt.Printf("   Mutex 版本:   %v\n", durationM)
	fmt.Printf("   RWMutex 版本: %v\n", durationRW)

	if durationRW < durationM {
		improvement := float64(durationM-durationRW) / float64(durationM) * 100
		fmt.Printf("   ✓ RWMutex 快了 %.2f%%！\n", improvement)
	} else {
		fmt.Println("   ⚠️  在這個測試中 Mutex 表現更好（可能因為操作太快）")
	}

	fmt.Println("\n💡 知識點：")
	fmt.Println("   - Mutex: 任何時候只有一個 Goroutine 可以訪問（讀或寫）")
	fmt.Println("   - RWMutex: 多個 Goroutine 可以同時讀取，但寫入時獨佔")
	fmt.Println("   - 當讀取操作遠多於寫入時，RWMutex 性能更好")
}
