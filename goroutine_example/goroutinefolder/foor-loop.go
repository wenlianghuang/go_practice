package goroutinefolder

import "fmt"

func Foorloop() {
	c := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			c <- i
		}
		close(c)
	}()

	for {
		v, ok := <-c
		if !ok {
			break
		}
		fmt.Println(v)
	}
}
