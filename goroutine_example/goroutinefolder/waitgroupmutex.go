// mutex: 保護在共享資源下,只一個goroutine去找,避免error
// waitGroup: 用waitgroup讓多個goroutine並行做完所有的事,用Done是讓goroutine都做完以後才繼續往下走
package goroutinefolder

import (
	"fmt"
	"sort"
	"sync"
)

type Counter struct {
	mu       sync.RWMutex
	counters map[string]int
}

func (c *Counter) inc(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[name]++
}

// 進階 API: 批次加總,減少鎖競爭
func (c *Counter) Add(name string, n int) {
	c.mu.Lock()
	c.counters[name] += n
	c.mu.Unlock()
}

// 進階 API: 只讀存取可用 RLock 並行
func (c *Counter) Get(name string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v := c.counters[name]
	return v
}

// 進階 API: 回傳snapshot,避免呼叫端長時間持有鎖
func (c *Counter) Snapshot() map[string]int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	snap := make(map[string]int)
	for k, v := range c.counters {
		snap[k] = v
	}
	return snap
}
func WaitgroupMutex() {
	var wg sync.WaitGroup
	c := Counter{counters: make(map[string]int)}

	doIncrement := func(name string, n int) {
		for i := 0; i < n; i++ {
			c.inc(name)
		}
		wg.Done()
	}

	wg.Add(3)
	go doIncrement("a", 100)
	go doIncrement("a", 100)
	go doIncrement("c", 100)
	wg.Wait()
	fmt.Println(c.counters)
}

// Advanced Example: worker pool + channel,並用 Snapshot 穩定輸出
func WaitgroupMutexAdvanced() {
	c := Counter{counters: make(map[string]int)}

	type task struct {
		name string
		n    int
	}

	jobs := make(chan task, 8)

	const workers = 4
	var wgWorkers sync.WaitGroup
	wgWorkers.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {

			defer wgWorkers.Done()
			for t := range jobs {
				// 批次加總,一次拿鎖完成
				c.Add(t.name, t.n)
			}
		}()
	}

	// 產生任務
	for _, t := range []task{
		{"a", 100}, {"a", 100}, {"c", 100},
		{"b", 250}, {"c", 700}, {"a", 200},
	} {
		jobs <- t
	}

	close(jobs)      // 關閉 channel,讓 worker 知道沒有任務了
	wgWorkers.Wait() // 等待所有 worker 做完

	// 安全讀取結果
	snap := c.Snapshot()
	keys := make([]string, 0, len(snap))
	for k := range snap {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, snap[k])

	}

}
