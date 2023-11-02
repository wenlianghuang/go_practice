// ref:https://ppt.cc/fKxe0x
package waitgroupex

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func WaitGroupEx1() {
	wg := &sync.WaitGroup{}
	jobs := 100
	wg.Add(jobs) // work 100 job
	for i := 1; i <= 100; i++ {
		go doTask(wg)
	}
	wg.Wait() //Use to block the main function until all the goroutine is finished
	fmt.Println("jobs all done!")
}

func doTask(wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("one task done")
		wg.Done() //finish 1 job
	}()
	number := rand.Intn(10) + 1
	time.Sleep(time.Duration(number))
}
