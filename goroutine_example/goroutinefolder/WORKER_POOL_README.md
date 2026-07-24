# Worker Pool 範例說明

## 📚 學習目標

學習如何實作 Worker Pool 模式，這是 Go 併發程式設計中最常用的模式之一。

## 🎯 什麼是 Worker Pool？

Worker Pool 是一種並發模式，主要特點：
- **固定數量的 Worker**（工人）同時執行任務
- **任務隊列**（jobs channel）派發任務給閒置的工人
- **結果收集**（results channel）收集所有工人的處理結果
- **避免過度併發**：控制同時運行的 Goroutine 數量

## 📝 核心概念

```
任務來源 → jobs channel → Worker Pool → results channel → 主程式
           (任務隊列)    (固定工人數)    (結果收集)
```

### 關鍵組件

1. **jobs channel**：緩衝 channel，存放待處理的任務
2. **results channel**：緩衝 channel，存放處理結果
3. **WaitGroup**：等待所有 worker 完成
4. **監控者 Goroutine**：等待所有工人完成後關閉 results channel

## 🚀 提供的範例

### 1. `WorkerPoolExample` - 基本範例

最基本的 Worker Pool 實作，展示核心概念：
- 3 個工人處理 5 個任務
- 每個任務計算平方（j * j）
- 處理時間：500ms

```bash
go run main.go WorkerPoolExample
```

**輸出示例：**
```
🚀 啟動 3 個工人...
📦 派發 5 個任務...
👷 工人 1 正在處理任務 1
👷 工人 2 正在處理任務 2
👷 工人 3 正在處理任務 3
✅ 收到結果: 1
✅ 收到結果: 4
...
```

### 2. `WorkerPoolWithNames` - 字串任務版本

使用字串任務，更貼近真實應用：
- 3 個工人（小明、小華、小美）
- 5 個家事任務（洗碗、拖地、倒垃圾、洗衣服、煮飯）
- 處理時間：300ms

```bash
go run main.go WorkerPoolWithNames
```

### 3. `WorkerPoolConfigurableDefault` - 可配置版本

可自訂工人數、任務數和處理時間：
- 3 個工人，9 個任務
- 處理時間：200ms

```bash
go run main.go WorkerPoolConfigurableDefault
```

### 4. `WorkerPoolLarge` - 大型 Worker Pool

展示大規模併發處理：
- 10 個工人，50 個任務
- 處理時間：100ms
- 顯示加速比統計

```bash
go run main.go WorkerPoolLarge
```

**輸出統計：**
```
📊 統計資料：
   總耗時: 505ms
   完成任務數: 50
   理論最快時間: 5s (如果用單一工人)
   實際加速比: 9.89x
```

## 🔑 關鍵程式碼解析

### 完整流程

```go
// 1. 建立 channels
jobs := make(chan int, 10)
results := make(chan int, 10)
var wg sync.WaitGroup

// 2. 啟動 workers
for w := 1; w <= 3; w++ {
    wg.Add(1)
    go workerV2(w, jobs, results, &wg)
}

// 3. 派發任務並關閉 jobs
for j := 1; j <= 5; j++ {
    jobs <- j
}
close(jobs)  // 重要！告訴 workers 沒有更多任務

// 4. 監控者 Goroutine
go func() {
    wg.Wait()      // 等待所有 worker 完成
    close(results) // 關閉 results，通知主程式可以結束
}()

// 5. 收集結果
for res := range results {
    fmt.Println("收到結果:", res)
}
```

### Worker 函數

```go
func workerV2(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done() // 確保離開前通知 WaitGroup
    
    for j := range jobs {  // 持續從 jobs 接收任務，直到 channel 關閉
        // 處理任務
        time.Sleep(500 * time.Millisecond)
        results <- j * j  // 發送結果
    }
}
```

## 💡 重要觀念

### 為什麼需要監控者 Goroutine？

```go
go func() {
    wg.Wait()      // 等待所有 worker
    close(results) // 關閉 results channel
}()
```

- 不能在主程式中 `wg.Wait()`，因為主程式需要接收結果
- 必須有人在所有 worker 完成後關閉 `results` channel
- 關閉 `results` 後，主程式的 `for range results` 才會結束

### Channel 方向性

```go
jobs <-chan int      // 只讀 channel（worker 只能接收）
results chan<- int   // 只寫 channel（worker 只能發送）
```

這樣的設計：
- ✅ 類型安全，防止誤用
- ✅ 清楚表達意圖
- ✅ 編譯器會檢查錯誤

### Buffered vs Unbuffered Channel

```go
// Buffered（推薦）
jobs := make(chan int, 10)
results := make(chan int, 10)

// Unbuffered
jobs := make(chan int)
results := make(chan int)
```

**使用 Buffered 的優點：**
- 避免 Goroutine 阻塞
- 更好的性能
- 容量應 ≥ 任務數或結果數

## 📊 性能優勢

### 串行 vs 並行

**串行處理（1 個工人）：**
- 5 個任務 × 500ms = 2500ms

**並行處理（3 個工人）：**
- 約 1000ms（3 個任務同時進行）

**加速比 = 2.5x**

實際測試（50 任務，10 工人，100ms）：
- 理論串行時間：5000ms
- 實際並行時間：505ms
- **加速比：9.89x** 🚀

## 🎓 適用場景

Worker Pool 適合以下場景：

1. **批量處理**：處理大量獨立任務
2. **資源限制**：限制同時訪問資源的數量（如資料庫連接）
3. **網路請求**：並行發送多個 HTTP 請求
4. **圖片處理**：批量處理圖片轉換
5. **資料處理**：ETL（Extract, Transform, Load）流程

## 🔗 相關檔案

- `worker-pool-example.go` - 實作程式碼
- `main.go` - 主程式入口
- `workerpools.go` - 原有的簡化版本

## 📚 延伸學習

1. **Context 取消**：使用 `context.Context` 支援任務取消
2. **錯誤處理**：在 results 中傳遞錯誤資訊
3. **動態調整**：根據負載動態增減 worker 數量
4. **優先級隊列**：不同優先級的任務
5. **超時控制**：為每個任務設定超時

## ⚠️ 常見錯誤

### 1. 忘記關閉 jobs channel

```go
// ❌ 錯誤：workers 會永遠等待
for j := 1; j <= 5; j++ {
    jobs <- j
}
// 忘記 close(jobs)

// ✅ 正確
close(jobs)
```

### 2. 在主程式中 wg.Wait()

```go
// ❌ 錯誤：主程式會阻塞，無法接收結果
wg.Wait()
for res := range results {  // 永遠執行不到
    fmt.Println(res)
}

// ✅ 正確：使用監控者 Goroutine
go func() {
    wg.Wait()
    close(results)
}()
for res := range results {
    fmt.Println(res)
}
```

### 3. 忘記關閉 results channel

```go
// ❌ 錯誤：for range 永遠不會結束
go func() {
    wg.Wait()
    // 忘記 close(results)
}()

// ✅ 正確
go func() {
    wg.Wait()
    close(results)
}()
```

## 🎯 關鍵要點總結

1. **固定數量的 Worker**：避免過度併發
2. **使用 Buffered Channel**：提升性能
3. **正確關閉 Channel**：jobs 和 results 都要關閉
4. **監控者模式**：用獨立 Goroutine 監控完成狀態
5. **WaitGroup**：追蹤所有 worker 的完成狀態
6. **Channel 方向性**：明確讀寫權限

Happy Coding! 🎉
