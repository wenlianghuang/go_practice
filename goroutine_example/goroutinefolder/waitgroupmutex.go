// mutex: 保護在共享資源下,只一個goroutine去找,避免error
// waitGroup: 用waitgroup讓多個goroutine並行做完所有的事,用Done是讓goroutine都做完以後才繼續往下走
package goroutinefolder

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu       sync.Mutex
	counters map[string]int
}

func (c *Counter) inc(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[name]++
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
