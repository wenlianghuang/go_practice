package goroutinefolder

import (
	"fmt"
	"sync"
	"time"
)

func Channelblockpullwait() {
	//var mu sync.Mutex
	var wg sync.WaitGroup
	ch := make(chan string)
	done := make(chan string)
	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println("calculate goroutine starts calculating")
		time.Sleep(time.Second) // Heavy calculation
		fmt.Println("calculate goroutine ends calculating")

		ch <- "FINISH"
		<-done // waiting for main goroutine
		fmt.Println("calculate goroutine finished")
	}()

	fmt.Println("main goroutine is waiting for channel to receive value")
	fmt.Println(<-ch) // goroutine 執行會在此被迫等待
	close(done)
	wg.Wait() // wait until calculate goroutine finish       // notify calcuate goroutine has recieved the message
	fmt.Println("main goroutine finished")

}
