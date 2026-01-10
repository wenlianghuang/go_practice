package goroutinefolder

import (
	"fmt"
)

// sum 函式接收一個切片與一個 channel
func sum(s []int, c chan int) {
	total := 0
	for _, v := range s {
		total += v
	}
	// 將計算結果送入 channel
	c <- total
}

func SumWithChannel() {
	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c
	fmt.Println(x, y, x+y)
}
