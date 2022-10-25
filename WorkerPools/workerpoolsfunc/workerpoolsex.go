package workerpoolsfunc

import (
	"fmt"
	"time"
)

func Worker(id int, jobs <-chan int, result chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "start job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finish job", j)
		result <- j * 2
	}
}
