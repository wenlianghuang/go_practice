package goroutinefolder

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func tickerWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	t := time.NewTicker(200 * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("workder: ctx cancel")
			return
		case <-t.C:
			fmt.Println("worker: tick")
		}
	}
}

func SignalControlAndClosing() {
	// This is a placeholder function to demonstrate signal control and closing channels.
	// Actual implementation would depend on specific requirements.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go tickerWorker(ctx, &wg)

	fmt.Println("press Ctrl+C to stop...")
	wg.Wait()
	fmt.Println("graceful exist")
}
