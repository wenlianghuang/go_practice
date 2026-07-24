# Goroutine 教學範例專案 — 完整說明

本專案為 **Go 併發程式設計（Goroutine / Channel / sync）** 的教學用程式庫，內含大量可獨立執行的範例。每個範例都對應一個函式，透過 **`main.go`** 以命令列參數選擇要執行的範例。

---

## 目錄

1. [main.go 的用途](#maingo-的用途)
2. [專案結構總覽](#專案結構總覽)
3. [依主題分類的檔案說明](#依主題分類的檔案說明)
4. [如何執行範例](#如何執行範例)
5. [延伸閱讀（現有 README）](#延伸閱讀現有-readme)

---

## main.go 的用途

**`main.go`** 是此教學專案的**單一入口**，用途如下：

- **不執行任何具體業務邏輯**，只負責：
  - 根據命令列第一個參數（函式名稱）從 `goroutinefolder` 套件中選出對應的示範函式並執行。
- **沒有帶參數時**：印出所有可用的函式名稱（已排序），並提示用法。
- **帶一個參數時**：例如 `go run main.go Bufferedchannel`，會執行 `goroutinefolder.Bufferedchannel()`。

因此，**所有教學範例都寫在 `goroutinefolder/` 底下**，`main.go` 純粹是「**範例執行器 / Demo Runner**」，方便你快速切換、比較不同主題的程式行為。

### 使用方式

```bash
# 列出所有可用範例
go run main.go

# 執行指定範例（函式名稱區分大小寫）
go run main.go <FunctionName>
```

例如：

```bash
go run main.go AtomicOperation
go run main.go Bufferedchannel
go run main.go DownloaderWithAllFeatures
```

---

## 專案結構總覽

```
goroutine_example/
├── main.go                 # 入口：依參數執行 goroutinefolder 中的示範函式
├── go.mod
├── go.sum
├── README.md               # 本說明文件
└── goroutinefolder/
    ├── atomic.go                    # 原子操作 (sync/atomic)
    ├── atomic_intro.md              # 原子操作概念說明
    ├── blockselect.go               # select 非阻塞 + default
    ├── buffered-channel.go          # 緩衝 channel 基本用法
    ├── channel-bloch-push.go        # 送資料進 channel 時的阻塞
    ├── channel-block-pull.go       # 從 channel 取資料時的阻塞
    ├── channel-block-pull-wait.go   # channel + WaitGroup 等待
    ├── closingchannels.go           # 關閉 channel 與 done 訊號
    ├── concurrent-cache.go          # 併發快取 (Mutex / RWMutex)
    ├── CONCURRENT_CACHE_README.md   # 併發快取說明
    ├── context-canccel-timeout.go   # Context 取消與超時
    ├── CONTEXT_ERR_GUIDE.md         # Context.Err() 說明
    ├── directions.go                # channel 方向 (只讀/只寫)
    ├── download-with-timeout.go     # 下載超時 (select / context)
    ├── downloader-comprehensive.go  # 綜合下載器 (Worker+Mutex+Atomic+Context)
    ├── DOWNLOADER_README.md         # 下載器與 Race 修正說明
    ├── foor-loop.go                 # for-range 接收 channel
    ├── multigoroutinwithoneval.go   # 多 goroutine 共用一變數 + Mutex
    ├── pipline-fanin-fanout.go      # Pipeline、Fan-out、Fan-in
    ├── RACE_CONDITION_FIX.md        # Check-Then-Act 競態修正
    ├── signal-control-and-closing.go # 訊號控制與優雅結束
    ├── sum_with_channel.go          # 用 channel 收集並行計算結果
    ├── syncMutex.go                 # sync.Mutex 基本
    ├── syncMutex2.go                # Mutex 保護 map
    ├── waitgroupmutex.go            # WaitGroup + RWMutex / Counter
    ├── worker-pool-example.go       # Worker Pool 多種範例
    ├── workerpools.go               # Worker Pool 基礎與 V2
    └── WORKER_POOL_README.md        # Worker Pool 說明
```

---

## 依主題分類的檔案說明

### 一、Channel 基礎

| 檔案 | 對應函式 | 說明 |
|------|----------|------|
| **buffered-channel.go** | `Bufferedchannel` | 緩衝 channel：`make(chan int, 1)`，送一個值不必立刻有人接收，示範無阻塞的送/收。 |
| **channel-bloch-push.go** | `Channelblockpush` | **送資料時阻塞**：goroutine 把 `"FINISH"` 送進無緩衝 channel 後，若 main 還沒接收，該 goroutine 會卡在 `ch <- "FINISH"`，直到 main 執行 `<-ch`。 |
| **channel-block-pull.go** | `Channelblockpull` | **收資料時阻塞**：main 先執行 `<-ch`，此時 goroutine 還沒送資料，main 會阻塞，直到 goroutine 送入 `"FINISH"`。 |
| **channel-block-pull-wait.go** | `Channelblockpullwait` | 在「拉取阻塞」的基礎上加上 **WaitGroup** 與 **done channel**：main 收到結果後關閉 `done`，並 `wg.Wait()` 確保計算用 goroutine 完全結束再繼續。 |
| **closingchannels.go** | `Closingchannels` | 關閉 channel 的慣例：發送端 `close(jobs)`，接收端用 `j, more := <-jobs` 判斷 `more == false` 表示收完，並用 `done` channel 通知 main 工作結束。 |
| **directions.go** | `Directions` | **Channel 方向**：`chan<- T`（只寫）、`<-chan T`（只讀）。示範 `ping` 只寫、`pong` 從 ping 讀再寫到 pong，main 從 pong 讀。 |
| **foor-loop.go** | `Foorloop` | 用 **for + `v, ok := <-c`** 接收 channel，並在 `!ok` 時 break；搭配發送端 `close(c)` 表示不再送資料。 |
| **sum_with_channel.go** | `SumWithChannel` | 將切片切成兩半，用兩個 goroutine 分別計算加總，結果透過 **同一個 channel** 送回，main 收兩次得到兩段的和再相加。 |

### 二、select 與非阻塞

| 檔案 | 對應函式 | 說明 |
|------|----------|------|
| **blockselect.go** | `Blockselect` | **select + default**：在 `for` 裡用 `select` 監聽 `<-ch`；有資料就 return，沒有就走 `default` 印出 "WARNING..." 並 sleep，避免 main 一直卡在 `<-ch`。 |

### 三、sync：Mutex、WaitGroup、RWMutex

| 檔案 | 對應函式 | 說明 |
|------|----------|------|
| **syncMutex.go** | `Syncmutex` | **sync.Mutex** 基本：自訂結構體帶 `sync.Mutex`，多個 goroutine 對同一變數做 `Lock() -> 修改 -> Unlock()`，並用 **WaitGroup** 等待全部完成，最後印出正確累加結果。 |
| **syncMutex2.go** | `SyncMutex2` | **Mutex 保護 map**：`Container` 內有 `map[string]int`，多個 goroutine 呼叫 `inc(name)`，用同一個 Mutex 保護對 map 的讀寫。 |
| **waitgroupmutex.go** | `WaitgroupMutex`, `WaitgroupMutexAdvanced` | **WaitGroup + RWMutex**：`Counter` 用 `RWMutex`、`inc`/`Add`/`Get`/`Snapshot` 等 API；進階版結合 **Worker Pool + channel**，用 `Snapshot()` 安全取得當前計數。 |
| **multigoroutinwithoneval.go** | `MultiGoroutineoneval` | **多個 goroutine 共用一個變數**：兩個 goroutine 同時對 `val` 做 `val++`，用 **Mutex** 保護，並用 **WaitGroup** 等待兩者結束。 |

### 四、原子操作 (sync/atomic)

| 檔案 | 對應函式 | 說明 |
|------|----------|------|
| **atomic.go** | `AtomicOperation` | **sync/atomic**：同一批 goroutine 同時做「普通變數 `unsafeCounter++`」與「`atomic.AddInt64(&safeCounter, 1)`」；對比結果說明競態條件與原子操作的正確性。 |
| **atomic_intro.md** | — | 說明何謂原子操作、何時用 atomic vs Mutex、以及本專案中 atomic 範例的重點。 |

### 五、併發快取 (Concurrent Cache)

| 檔案 | 對應函式 | 說明 |
|------|----------|------|
| **concurrent-cache.go** | `ConcurrentCacheUnsafe`, `ConcurrentCacheSafe`, `ConcurrentCacheRWMutex`, `ConcurrentCacheComparison` | **三種實作**：無鎖（示範 race）、**Mutex** 保護、**RWMutex** 讀多寫少；最後一個函式比較 Mutex 與 RWMutex 在讀多寫少下的效能。 |
| **CONCURRENT_CACHE_README.md** | — | 說明各版本差異、如何用 `-race` 檢測、Mutex vs RWMutex 的適用場景。 |

### 六、Context 取消與超時

| 檔案 | 對應函式 | 說明 |
|------|----------|------|
| **context-canccel-timeout.go** | `ContextCancelTimeout` | **context.WithTimeout**：建立會自動取消的 Context，worker 用 `select` 監聽 `ctx.Done()` 與 `jobs`；發送端也會在 `ctx.Done()` 時停止送任務，示範生命週期收斂。 |
| **CONTEXT_ERR_GUIDE.md** | — | 說明 **Context.Err()**、`DeadlineExceeded` 與 `Canceled` 的區分與在錯誤處理、監控上的用法。 |

### 七、下載與超時

| 檔案 | 對應函式 | 說明 |
|------|----------|------|
| **download-with-timeout.go** | `DownloadWithTimeout`, `DownloadWithTimeoutV2` | **select + time.After** 做下載超時；V2 用一次性 deadline、buffered channel，並在超時後關閉 channel 做收尾。 |
| **downloader-comprehensive.go** | `DownloaderWithAllFeatures`, `DownloaderWithShortTimeout`, `DownloaderWithLongTimeout` | **綜合下載器**：Worker Pool + **Mutex**（去重 map）+ **Atomic**（流量與計數）+ **Context**（全局限時與單次任務超時）+ **WaitGroup**；內含 **CheckAndMark** 避免 Check-Then-Act 競態。 |
| **DOWNLOADER_README.md** | — | 需求說明、各版本差異、CheckAndMark 與 Context.Err() 的用法、常見錯誤與最佳實踐。 |
| **RACE_CONDITION_FIX.md** | — | 說明下載器裡 **Check-Then-Act** 的競態問題，以及用 **CheckAndMark** 原子化檢查與標記的修正方式。 |

### 八、Pipeline、Fan-out、Fan-in

| 檔案 | 對應函式 | 說明 |
|------|----------|------|
| **pipline-fanin-fanout.go** | `PipelineFaninFanout` | **Pipeline**：`gen` 產生整數 → **Fan-out** 多個 `square` 從同一 channel 讀並平方 → **Fan-in** 將多個 channel 合併成一個；全程帶 **Context** 支援取消與超時。 |

### 九、訊號與優雅結束

| 檔案 | 對應函式 | 說明 |
|------|----------|------|
| **signal-control-and-closing.go** | `SignalControlAndClosing` | **signal.NotifyContext**：監聽 /SIGTERM（如 Ctrl+C），Context 被取消後 worker 的 `select` 收到 `ctx.Done()` 並結束，最後 `wg.Wait()` 優雅退出。 |

### 十、Worker Pool

| 檔案 | 對應函式 | 說明 |
|------|----------|------|
| **workerpools.go** | `Workerpools`, `WorkerpoolV2` | 基礎 Worker Pool：固定數個 worker 從 `jobs` 讀任務、把結果寫入 `results`；V2 用 **WaitGroup** 與背景 goroutine 在全部 worker 結束後 **close(results)**，主程式用 `for range results` 收結果。 |
| **worker-pool-example.go** | `WorkerPoolExample`, `WorkerPoolWithNames`, `WorkerPoolConfigurableDefault`, `WorkerPoolLarge` | 多種 Worker Pool 示範：基本數字任務、字串任務（名字/家事）、可配置工人數/任務數/處理時間、以及較大規模的效能示範。 |
| **WORKER_POOL_README.md** | — | Worker Pool 概念、流程、為何需要「監控者 goroutine」關閉 results、常見錯誤與正確寫法。 |

---

## 如何執行範例

在 `goroutine_example` 目錄下：

```bash
# 1. 列出所有可用的函式名稱
go run main.go

# 2. 執行單一範例（函式名稱需與 main.go 中 map 的 key 完全一致）
go run main.go Bufferedchannel
go run main.go AtomicOperation
go run main.go ContextCancelTimeout
go run main.go DownloaderWithAllFeatures
go run main.go WorkerPoolExample
```

使用 **race detector** 觀察競態（例如不安全快取或下載器）：

```bash
go run -race main.go ConcurrentCacheUnsafe
go run -race main.go DownloaderWithAllFeatures
```

---

## 延伸閱讀（現有 README）

- **atomic_intro.md** — 原子操作概念與何時用 atomic / Mutex  
- **CONCURRENT_CACHE_README.md** — 併發快取三種實作與 race 檢測  
- **CONTEXT_ERR_GUIDE.md** — Context.Err()、DeadlineExceeded、Canceled 的區分與實務用法  
- **DOWNLOADER_README.md** — 綜合下載器需求、設計與 CheckAndMark、Context 錯誤處理  
- **RACE_CONDITION_FIX.md** — Check-Then-Act 競態與 CheckAndMark 修正  
- **WORKER_POOL_README.md** — Worker Pool 流程、監控者模式與常見錯誤  

---

## 總結

- **main.go**：只做「依命令列參數選擇並執行 `goroutinefolder` 裡的示範函式」，是**範例執行器**，不包含教學邏輯本身。  
- **goroutinefolder/**：包含所有教學用 **.go** 與 **.md**；每個 .go 檔案都對應一到多個可從 main 呼叫的函式，涵蓋 channel、select、Mutex、RWMutex、WaitGroup、atomic、Context、超時、Worker Pool、Pipeline、訊號處理等主題。  
- 本 **README.md** 提供專案總覽、每個檔案的職責說明，以及執行方式與延伸閱讀索引，方便依主題查找與教學使用。
