package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Task 代表單一 User 請求
type Task struct {
	ID       int
	Input    string
	ResultCh chan string
}

// BatchInferenceService 支援將多個 Task 打包成 Batch 送給 AI 模型
type BatchInferenceService struct {
	taskQueue    chan Task
	maxBatch     int           // 最大 Batch Size (例如一次最多打包 8 個)
	batchTimeout time.Duration // 最長等待時間 (例如最多等 50ms 就必須送出)
	workerNum    int
	wg           sync.WaitGroup
}

// 建立一個 BatchInferenceService 實例
// workerNum: 工作者數量
// queueSize: 佇列大小
// maxBatch: 最大 Batch 大小
// batchTimeout: 最大等待時間
func NewBatchInferenceService(workerNum, queueSize, maxBatch int, batchTimeout time.Duration) *BatchInferenceService {
	return &BatchInferenceService{
		taskQueue:    make(chan Task, queueSize),
		maxBatch:     maxBatch,
		batchTimeout: batchTimeout,
		workerNum:    workerNum,
	}
}

// Start 啟動 Batch Workers
func (s *BatchInferenceService) Start(ctx context.Context) {
	for i := 1; i <= s.workerNum; i++ {
		s.wg.Add(1)
		go s.batchWorker(ctx, i)
	}
}

// batchWorker 負責收集單一 Task 並組合成 Batch 進行 AI 推論
func (s *BatchInferenceService) batchWorker(ctx context.Context, workerID int) {
	defer s.wg.Done()

	for {
		// 檢查 Context 是否取消，且 Queue 是否已清空
		select {
		case <-ctx.Done():
			if len(s.taskQueue) == 0 {
				fmt.Printf("[Worker %d] 收到關機指令且 Queue 已空，安全退出。\n", workerID)
				return
			}
		default:
		}

		// 1. 收集 Batch
		batch := s.collectBatch(ctx)
		if len(batch) == 0 {
			// 代表沒有任務且 Context 已經被 Cancel 了
			return
		}

		// 2. 模擬 GPU Batch 推論（無論 Batch 大小是 1 還是 8，GPU 耗時幾乎都是 0.5s）
		fmt.Printf("⚡ [Worker %d] 正在處理 Batch 推論 (Batch 大小: %d)...\n", workerID, len(batch))
		time.Sleep(500 * time.Millisecond)

		// 3. 分發結果給各自的 User
		for _, task := range batch {
			task.ResultCh <- fmt.Sprintf("Task %d 完成 (Batch Size: %d, Worker: %d)", task.ID, len(batch), workerID)
		}
	}
}

// collectBatch 是最精髓的地方：同時滿足「湊滿數量」或「超時觸發」
func (s *BatchInferenceService) collectBatch(ctx context.Context) []Task {
	var batch []Task

	// 阻塞等待第一個 Task 進來
	select {
	case task, ok := <-s.taskQueue:
		if !ok {
			return nil
		}
		batch = append(batch, task)
	case <-ctx.Done():
		// 關機流程中，如果 Queue 裡面還有剩餘 Task，繼續拿出來處理
		select {
		case task, ok := <-s.taskQueue:
			if !ok {
				return nil
			}
			batch = append(batch, task)
		default:
			return nil
		}
	}

	// 啟動定時器，避免為了湊滿 Batch 讓第一個 Task 等太久
	timer := time.NewTimer(s.batchTimeout)
	defer timer.Stop()

	for len(batch) < s.maxBatch {
		select {
		case task, ok := <-s.taskQueue:
			if !ok {
				return batch
			}
			batch = append(batch, task)
		case <-timer.C:
			// 時間到！就算沒湊滿 maxBatch 也直接送出
			return batch
		}
	}

	return batch
}

// Submit 供 External/HTTP 提交任務
func (s *BatchInferenceService) Submit(ctx context.Context, task Task) (string, error) {
	select {
	case s.taskQueue <- task:
	case <-ctx.Done():
		return "", fmt.Errorf("服務正在關閉或請求超時，拒絕接收新任務")
	}

	select {
	case res := <-task.ResultCh:
		return res, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// Shutdown 優雅關機
func (s *BatchInferenceService) Shutdown() {
	close(s.taskQueue) // 關閉 Channel，不再接受新任務
	s.wg.Wait()        // 等待所有 Worker 將剩餘 Batch 處理完畢
	fmt.Println("✅ 服務已完全優雅關閉（Graceful Shutdown Complete）")
}

func main() {
	// 建立一個 Worker 數量為 40、最大 Batch 為 10、Batch 等待時間 100ms 的服務
	svc := NewBatchInferenceService(40, 1000, 10, 100*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	svc.Start(ctx)

	// 模擬併發發送 1000 個請求
	var wg sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			reqCtx, reqCancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer reqCancel()

			task := Task{ID: id, Input: "audio_data", ResultCh: make(chan string, 1)}
			res, err := svc.Submit(reqCtx, task)
			if err != nil {
				fmt.Printf("❌ Task %d 失敗: %v\n", id, err)
				return
			}
			_ = res
		}(i)
	}

	// 監聽系統訊號（例如 Ctrl+C 或 k8s SIGTERM）
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\n⚠️ 收到關機訊號，啟動優雅關機程序...")
		cancel()       // 通知 Worker ctx 結束
		svc.Shutdown() // 等待 Worker 清空佇列
		os.Exit(0)
	}()

	wg.Wait()
	fmt.Println("🎉 所有請求處理 completed！嘗試優雅關機...")
	cancel()
	svc.Shutdown()
}
