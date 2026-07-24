package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ===== context 版本（獨立）=====

func producerCtx(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func sqCtx(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func sqCtxWithDelay(ctx context.Context, in <-chan int, delay time.Duration) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if delay > 0 {
				timer := time.NewTimer(delay)
				select {
				case <-timer.C:
				case <-ctx.Done():
					timer.Stop()
					return
				}
			}

			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func mergeCtx(ctx context.Context, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func watchCtxDone(name string, ctx context.Context) {
	// 教學用：觀察不同節點的 context 何時被取消、取消原因是什麼
	<-ctx.Done()
	fmt.Printf("[CancellationTree] %s done: %v\n", name, ctx.Err())
}

func demoCancellationTree() {
	/*
		Cancellation Tree（取消樹）的核心概念：

		- context 是一棵「父 -> 子」的樹。
		- 父節點被 cancel/timeout 時，所有子孫節點都會被取消（向下傳播）。
		- 子節點被 cancel 時，不會影響父節點，也不會影響兄弟節點（不會向上或水平傳播）。

		下面用 root -> childA, childB -> grandchildB 示範：
	*/
	root, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	childA, cancelA := context.WithCancel(root)
	childB, cancelB := context.WithCancel(root)
	defer cancelA()
	defer cancelB()

	grandchildB, cancelGrandchildB := context.WithTimeout(childB, 120*time.Millisecond)
	defer cancelGrandchildB()

	go watchCtxDone("root", root)
	go watchCtxDone("childA", childA)
	go watchCtxDone("childB", childB)
	go watchCtxDone("grandchildB(timeout)", grandchildB)

	// 先取消 childA：只會影響 childA，不會影響 root / childB / grandchildB。
	cancelA()
	time.Sleep(30 * time.Millisecond)

	// 再取消 root：會影響整棵樹向下傳播，childB 以及 grandchildB 都會被取消。
	rootCancel()

	// 留一點時間讓 goroutine 印出訊息（示範用途，避免程式太快結束看不到輸出）
	time.Sleep(50 * time.Millisecond)
}

// demoContextValues 示範 Context Value 的幾個重點用法：
//
// 1. value 會從父 context 傳到子 context（同一個 key）。
// 2. 子 context 可以覆寫同一個 key 的 value（只影響自己與其子孫）。
// 3. 如果 key 不存在，Value 會回傳 nil。
func demoContextValues() {
	// 一般會用自訂型別當作 key，避免跟別的 package 衝突
	type ctxKey string

	const (
		keyRequestID ctxKey = "requestID"
		keyUserID    ctxKey = "userID"
	)

	// root 上放一個 requestID
	root := context.WithValue(context.Background(), keyRequestID, "req-123")

	// childA 繼承 root 的 value，也可以再加自己的 userID
	childA := context.WithValue(root, keyUserID, "user-A")

	// childB 從 root 繼承，但把 requestID 覆寫成另一個值
	childB := context.WithValue(root, keyRequestID, "req-456")

	// grandchildB 從 childB 繼承（會看到覆寫後的 requestID）
	grandchildB := context.WithValue(childB, keyUserID, "user-B")

	fmt.Println("=== Context Values Demo ===")
	fmt.Printf("root:        requestID=%v, userID=%v\n",
		root.Value(keyRequestID), root.Value(keyUserID)) // userID 沒設定，會是 <nil>
	fmt.Printf("childA:      requestID=%v, userID=%v\n",
		childA.Value(keyRequestID), childA.Value(keyUserID))
	fmt.Printf("childB:      requestID=%v, userID=%v\n",
		childB.Value(keyRequestID), childB.Value(keyUserID)) // userID 沒設定
	fmt.Printf("grandchildB: requestID=%v, userID=%v\n",
		grandchildB.Value(keyRequestID), grandchildB.Value(keyUserID))

	fmt.Println("=== End Context Values Demo ===")
}

func runWithContext() {
	// Context Values 教學：先示範 value 的傳遞與覆寫
	demoContextValues()

	// Cancellation Tree 教學：再跑一次取消樹 demo
	demoCancellationTree()

	// 方便測試：把 timeout 或 delay 調大/調小，就能觀察超時取消是否生效
	timeout := 600 * time.Millisecond
	perItemDelay := 250 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	in := producerCtx(ctx, 1, 2, 3, 4, 5, 6, 7, 8)

	c1 := sqCtxWithDelay(ctx, in, perItemDelay)
	c2 := sqCtxWithDelay(ctx, in, perItemDelay)
	c3 := sqCtxWithDelay(ctx, in, perItemDelay)

	out := mergeCtx(ctx, c1, c2, c3)

	count := 0
	for n := range out {
		fmt.Println("接收到 (context):", n)
		count++
	}

	if err := ctx.Err(); err != nil {
		fmt.Printf("context 結束原因: %v (收到 %d 筆)\n", err, count)
	} else {
		fmt.Printf("context 正常結束 (收到 %d 筆)\n", count)
	}
}
