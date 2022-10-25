package goroutinefolder

import (
	"fmt"
	"time"
)

func Blockselect() {
	ch := make(chan string)

	go func() {
		fmt.Println("calculate goroutine starts calculating")
		time.Sleep(time.Second) // Heavy calculation
		fmt.Println("calculate goroutine ends calculating")

		ch <- "FINISH"
		time.Sleep(time.Second)
		fmt.Println("calculate goroutine finished")
	}()

	for {
		select {
		case <-ch: //Channel 中有資料執行此區域
			fmt.Println("main goroutine finished")
			return
		default:
			fmt.Println("WARNING...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
