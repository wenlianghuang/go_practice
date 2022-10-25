package goroutine

import (
	"fmt"
	"time"
)

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

func Goroutine() {
	f("direct")

	go f("goroutine")

	go func(msg string, msg2 string) {
		fmt.Println(msg)
		fmt.Println(msg2)
	}("I am going know", "What should I do?")

	time.Sleep(time.Second)
	fmt.Println("done")
}
