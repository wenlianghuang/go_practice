package goroutinefolder

import (
	"fmt"
	"sync"
	"time"
)

type SyncNumber struct {
	v   int32
	mux sync.Mutex
}

func Syncmutex() {
	total := SyncNumber{v: 0}

	for i := 0; i < 1000; i++ {
		go func() {
			total.mux.Lock()
			total.v++
			total.mux.Unlock()
		}()

	}

	time.Sleep(time.Second)
	total.mux.Lock()
	fmt.Printf("V: %+v\n", total.v)
	total.mux.Unlock()
}
