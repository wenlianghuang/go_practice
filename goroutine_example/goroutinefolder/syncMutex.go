package goroutinefolder

import (
	"fmt"
	"sync"
)

type SyncNumber struct {
	v   int32
	mux sync.Mutex
}

func Syncmutex() {
	total := SyncNumber{v: 0}
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			total.mux.Lock()
			total.v++
			total.mux.Unlock()
			wg.Done()
		}()

	}

	wg.Wait() // wait until all goroutines finish(wg counter to 0 Add and Done)
	fmt.Printf("Final Value: %+v\n", total.v)
	total.mux.Lock()
	fmt.Printf("V: %+v\n", total.v)
	total.mux.Unlock()
}
