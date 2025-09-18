package goroutinefolder

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func gen(ctx context.Context, n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return
			case out <- i:
			}
		}
	}()
	return out
}
func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case <-ctx.Done():
				return
			case out <- n * n:
			}
		}
	}()
	return out
}
func fanIn(ctx context.Context, chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	transfer := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case <-ctx.Done():
				return
			case out <- n:
			}
		}
	}
	wg.Add(len(chans))
	for _, c := range chans {
		go transfer(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
func PipelineFaninFanout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	src := gen(ctx, 100)
	// fan-out: 3 個平方工人
	out1 := square(ctx, src)
	out2 := square(ctx, src)
	out3 := square(ctx, src)
	// fan-in: 將三個平方工人的結果合併
	merged := fanIn(ctx, out1, out2, out3)

	for v := range merged {
		fmt.Println(v)
		time.Sleep(60 * time.Millisecond)
	}
	fmt.Println("pipeline done")

}
