package goroutinefolder

import (
	"fmt"
	"sync"
	"time"
)

// jobs => read-only channel, receive task mission
// results => write-only channel, send task mission
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 5
	}
}

func Workerpools() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	close(jobs)

	for a := 1; a <= numJobs; a++ {
		// result(read-only) reciev from results channel
		// results(write-only) send to result
		result := <-results
		fmt.Println("Result:", result)
	}

}

// 更健壯的版本：不依賴 results 的 buffer,確保正確收尾

func WorkerpoolV2() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int)

	const workers = 3
	var wg sync.WaitGroup
	wg.Add(workers)

	// 啟動 workers
	for w := 1; w <= workers; w++ {
		go func(workerID int) {
			defer wg.Done()
			for j := range jobs {
				fmt.Println("worker", workerID, "started job", j)
				time.Sleep(time.Second)
				fmt.Println("worker", workerID, "finished job", j)
				results <- j * 5
			}
		}(w)
	}

	// 背景收尾：等全部worker結束後關閉results
	go func() {
		wg.Wait()
		close(results)
	}()

	// 投遞任務並關閉 jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// 收集結果
	for result := range results {
		fmt.Println("Result:", result)
	}
}
