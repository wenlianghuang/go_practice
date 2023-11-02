package goroutinefolder

import (
	"fmt"
	"sync"
	"time"
)

func MultiGoroutineoneval() {
	var lock sync.Mutex
	var wg sync.WaitGroup

	val := 0
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			lock.Lock()
			val++
			fmt.Printf("First goroutine val++ and val = %d\n", val)
			lock.Unlock()
			time.Sleep(3000)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			lock.Lock()
			val++
			lock.Unlock()
			fmt.Printf("Sec goroutine val++ and val= %d\n", val)
			time.Sleep(1000)
		}
	}()
	wg.Add(2)
	wg.Wait()
}
