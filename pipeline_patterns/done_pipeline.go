package main

import (
	"fmt"
	"sync"
)

// ===== done 版本（獨立）=====

// 階段 1：生產者 - 將資料放入通道
func producer(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

// 階段 2：運算者 (Worker) - 處理資料
func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

// 階段 3：合併者 (Fan-in) - 匯整多個通道的結果
func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
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

func runWithDone() {
	done := make(chan struct{})
	defer close(done)

	in := producer(done, 1, 2, 3, 4, 5, 6, 7, 8)

	c1 := sq(done, in)
	c2 := sq(done, in)
	c3 := sq(done, in)

	out := merge(done, c1, c2, c3)

	for n := range out {
		fmt.Println("接收到 (done):", n)
	}
}

