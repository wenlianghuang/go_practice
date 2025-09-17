package goroutinefolder

import (
	"fmt"
)

func Closingchannels() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				// 工作 Goroutine 向 done Channel 發送一個 true 值，表示它已經完成了所有工作。
				done <- true
				return
			}
		}
	}()

	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
	}
	close(jobs)
	fmt.Println("sent all jobs")

	// 這是一個阻塞操作。主 Goroutine 會在這裡暫停，一直等到它能從 done Channel 接收到一個值為止。
	<-done
}
