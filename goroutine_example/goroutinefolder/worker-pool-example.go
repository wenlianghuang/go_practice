package goroutinefolder

import (
	"fmt"
	"sync"
	"time"
)

// workerV2 是執行任務的工人函式（使用 WaitGroup）
func workerV2(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done() // 工人離開前報備
	for j := range jobs {
		fmt.Printf("👷 工人 %d 正在處理任務 %d\n", id, j)
		time.Sleep(time.Millisecond * 500)
		results <- j * j
	}
}

// WorkerPoolExample 演示基本的 Worker Pool 模式
func WorkerPoolExample() {
	fmt.Println("=== Worker Pool 基本範例 ===")
	fmt.Println()

	jobs := make(chan int, 10)
	results := make(chan int, 10)
	var wg sync.WaitGroup

	// 1. 啟動 3 個工人
	numWorkers := 3
	fmt.Printf("🚀 啟動 %d 個工人...\n", numWorkers)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go workerV2(w, jobs, results, &wg)
	}

	// 2. 派發 5 個任務並關閉 jobs
	numJobs := 5
	fmt.Printf("📦 派發 %d 個任務...\n\n", numJobs)
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// 3. 啟動監控者 Goroutine 🕵️‍♂️
	go func() {
		wg.Wait()      // 等待 3 個工人都執行完畢
		close(results) // 關閉結果通道，通知 main 結束接收
	}()

	// 4. 主程式優雅地接收所有結果
	for res := range results {
		fmt.Printf("✅ 收到結果: %d\n", res)
	}

	fmt.Println("\n🎉 所有任務與結果處理完畢！")
}

// workerWithName 是帶有名稱的工人函式（進階版）
func workerWithName(name string, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("👷 工人 [%s] 正在處理任務: %s\n", name, job)
		time.Sleep(time.Millisecond * 300)
		results <- fmt.Sprintf("%s 已完成 '%s'", name, job)
	}
}

// WorkerPoolWithNames 演示使用字串任務的 Worker Pool
func WorkerPoolWithNames() {
	fmt.Println("=== Worker Pool 字串任務範例 ===")
	fmt.Println()

	jobs := make(chan string, 10)
	results := make(chan string, 10)
	var wg sync.WaitGroup

	// 工人名稱列表
	workerNames := []string{"小明", "小華", "小美"}

	// 1. 啟動工人
	fmt.Printf("🚀 啟動 %d 個工人...\n", len(workerNames))
	for _, name := range workerNames {
		wg.Add(1)
		go workerWithName(name, jobs, results, &wg)
	}

	// 2. 派發任務
	tasks := []string{"洗碗", "拖地", "倒垃圾", "洗衣服", "煮飯"}
	fmt.Printf("📦 派發 %d 個任務...\n\n", len(tasks))
	for _, task := range tasks {
		jobs <- task
	}
	close(jobs)

	// 3. 監控者
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. 接收結果
	for res := range results {
		fmt.Printf("✅ %s\n", res)
	}

	fmt.Println("\n🎉 所有家事完成！")
}

// WorkerPoolConfigurable 可配置的 Worker Pool
func WorkerPoolConfigurable(numWorkers, numJobs int, processingTime time.Duration) {
	fmt.Println("=== 可配置的 Worker Pool ===")
	fmt.Printf("工人數量: %d, 任務數量: %d, 處理時間: %v\n\n", numWorkers, numJobs, processingTime)

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	startTime := time.Now()

	// 1. 啟動工人
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range jobs {
				fmt.Printf("👷 工人 %d 處理任務 %d\n", id, j)
				time.Sleep(processingTime)
				results <- j * 2
			}
		}(w)
	}

	// 2. 派發任務
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// 3. 監控者
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. 接收結果
	count := 0
	for res := range results {
		count++
		fmt.Printf("✅ 結果 %d: %d\n", count, res)
	}

	duration := time.Since(startTime)

	fmt.Printf("\n📊 統計資料：\n")
	fmt.Printf("   總耗時: %v\n", duration)
	fmt.Printf("   完成任務數: %d\n", count)
	fmt.Printf("   理論最快時間: %v (如果用單一工人)\n", processingTime*time.Duration(numJobs))
	fmt.Printf("   實際加速比: %.2fx\n", float64(processingTime*time.Duration(numJobs))/float64(duration))
}

// WorkerPoolConfigurableDefault 使用預設參數的可配置 Worker Pool
func WorkerPoolConfigurableDefault() {
	WorkerPoolConfigurable(3, 9, 200*time.Millisecond)
}

// WorkerPoolLarge 大型 Worker Pool 演示
func WorkerPoolLarge() {
	fmt.Println("=== 大型 Worker Pool (10 工人, 50 任務) ===")
	fmt.Println()
	WorkerPoolConfigurable(10, 50, 100*time.Millisecond)
}
