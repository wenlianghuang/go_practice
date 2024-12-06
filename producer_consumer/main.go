package main

import (
	"fmt"
	"time"
)

func producer(ch chan int) {
	for i := 0; i < 10; i++ {
		fmt.Printf("Producer: %d\n", i)
		ch <- i
		time.Sleep(100 * time.Millisecond)
	}
	close(ch)
}

func consumer(ch chan int) {
	for value := range ch {
		fmt.Printf("Consumer: %d\n", value)
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	ch := make(chan int)
	go producer(ch)
	consumer(ch)
}
