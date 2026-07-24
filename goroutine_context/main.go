package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	// 父 context：給 worker2、worker3；cancelAll 會一併結束他們
	ctx, cancelAll := context.WithCancel(context.Background())
	defer cancelAll()

	// 子 context：只給 worker1。只呼叫 stopW1 時，僅 worker1 會收到 Done
	w1Ctx, stopW1 := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		worker(w1Ctx, "worker1")
	}()
	go func() {
		defer wg.Done()
		worker(ctx, "worker2")
	}()
	go func() {
		defer wg.Done()
		worker(ctx, "worker3")
	}()

	fmt.Println("▶ 先跑 2 秒…")
	time.Sleep(2 * time.Second)

	fmt.Println("⏹ 只取消 worker1（stopW1），worker2/3 繼續用父 ctx")
	stopW1()

	fmt.Println("▶ worker2、worker3 再跑 5 秒…")
	time.Sleep(5 * time.Second)

	fmt.Println("⏹ cancelAll：結束 worker2、worker3")
	cancelAll()

	wg.Wait()
	fmt.Println("main: 所有 worker 已退出")
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("  %s: stopped (%v)\n", name, ctx.Err())
			return
		case <-time.After(1 * time.Second):
			fmt.Printf("  %s: working\n", name)
		}
	}
}
