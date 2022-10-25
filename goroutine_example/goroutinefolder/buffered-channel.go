package goroutinefolder

import "fmt"

func Bufferedchannel() {
	ch := make(chan int, 1)
	ch <- 1
	fmt.Println("Ch: ", <-ch)
}
