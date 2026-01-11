# 綜合練習：網址下載器

## 🎯 練習目標

這是一個綜合性的 Goroutine 練習，整合了所有重要的併發概念：
- **Worker Pool**：固定數量的工人處理任務
- **Mutex**：保護共享資源（去重 map）
- **Atomic**：無鎖的流量統計
- **Context**：超時控制
- **WaitGroup**：等待所有工人完成
- **Channel**：任務分發和結果收集

## 📋 需求說明

### 1. Worker Pool（工作池）
- ✅ 固定啟動 5 個工人
- ✅ 並發處理多個下載任務

### 2. 任務分配
- ✅ 處理 20 個網址
- ✅ 有些網址會重複（模擬重複抓取）

### 3. 去重機制（Mutex）
- ✅ 使用 `map[string]bool` 記錄已下載網址
- ✅ 使用 `sync.Mutex` 保護 map（併發安全）
- ✅ 重複的網址會被跳過

### 4. 流量統計（Atomic）
- ✅ 使用 `atomic.Uint64` 累加總下載位元組
- ✅ 使用 `atomic.Uint32` 記錄成功/跳過/超時數
- ✅ 無鎖操作，性能更好

### 5. 超時控制（Context）
- ✅ 設定 2 秒總限時
- ✅ 超時則所有工人停止
- ✅ 使用 `context.WithTimeout`

### 6. 生命週期（WaitGroup）
- ✅ 確保所有工人完成後才印出最終報表
- ✅ 使用 `sync.WaitGroup` 追蹤工人狀態

## 🚀 提供的版本

### 1. `DownloaderWithAllFeatures` - 完整功能版本

符合所有需求的完整實作：
- 5 個工人
- 20 個任務（包含重複網址）
- 2 秒超時
- 完整統計報表

```bash
go run main.go DownloaderWithAllFeatures
```

### 2. `DownloaderWithShortTimeout` - 短超時版本

測試超時機制：
- 5 個工人
- 10 個任務
- **1 秒超時**（更容易觸發超時）

```bash
go run main.go DownloaderWithShortTimeout
```

### 3. `DownloaderWithLongTimeout` - 長超時版本

觀察去重機制：
- 3 個工人（較少）
- 10 個任務（包含重複）
- 5 秒超時（足夠完成所有任務）

```bash
go run main.go DownloaderWithLongTimeout
```

## 📊 輸出示例

### 完整功能版本輸出

```
=== 綜合練習：網址下載器 ===
📋 功能：Worker Pool + Mutex + Atomic + Context + WaitGroup

🚀 啟動 5 個工人...
📦 派發 20 個下載任務...

📥 工人 2: 開始下載 [https://example.com/page1]
📥 工人 3: 開始下載 [https://example.com/page2]
✅ 工人 3: 完成 [https://example.com/page2] (8104 bytes)
⏭️  工人 5: 跳過 [https://example.com/page2]（已下載）
⏰ 工人 1: 下載 [https://example.com/page14] 超時
👷 工人 2 完成所有任務

============================================================
📊 最終報表
============================================================
⏱️  總耗時: 2.001s
👷 工人數量: 5
📋 總任務數: 20

✅ 成功下載: 15 個網址
⏭️  跳過重複: 4 個網址
⏰ 超時任務: 1 個
📊 總流量: 92.24 KB
📈 平均大小: 6.15 KB/個

🔑 使用的技術：
   ✓ Worker Pool: 5 個固定工人
   ✓ Mutex: 保護 downloaded map（去重）
   ✓ Atomic: 累加流量統計（無鎖併發）
   ✓ Context: 2 秒超時控制
   ✓ WaitGroup: 等待所有工人完成
   ✓ Channel: 任務分發

⚠️  注意：達到 2 秒超時限制，部分任務未完成
```

## 🔑 關鍵程式碼解析

### ⚠️ 重要：防止 Race Condition 的正確做法

在併發環境下，**最容易犯的錯誤**是 "Check-Then-Act" 模式：

#### ❌ 錯誤做法（有 Race Condition）

```go
// 步驟 1：檢查
if dm.IsDownloaded(task.URL) {
    return  // 已下載
}
// ⚠️ 危險空隙！其他工人可能同時通過檢查

// 步驟 2：下載
download(task.URL)

// 步驟 3：標記
dm.MarkDownloaded(task.URL)
```

**問題分析：**

| 時間 | 工人 A | 工人 B | 結果 |
|-----|--------|--------|------|
| T1 | `IsDownloaded("page1")` → false ✓ | | |
| T2 | | `IsDownloaded("page1")` → false ✓ | **兩個都通過！** |
| T3 | 開始下載 "page1" | 開始下載 "page1" | **重複下載** 💥 |
| T4 | `MarkDownloaded("page1")` | `MarkDownloaded("page1")` | 浪費資源 |

#### ✅ 正確做法（使用 CheckAndMark）

```go
// 原子性地檢查並標記（在同一個鎖內完成）
func (dm *DownloadManager) CheckAndMark(url string) bool {
    dm.mutex.Lock()
    defer dm.mutex.Unlock()
    
    // 🔑 關鍵：檢查和標記在同一個鎖內，不可分割
    if dm.downloaded[url] {
        return false  // 已下載
    }
    
    dm.downloaded[url] = true  // 立即佔位
    return true  // 可以下載
}

// 使用方式
canDownload := dm.CheckAndMark(task.URL)
if !canDownload {
    return  // 已被其他工人佔位
}
// ✅ 安全：已經佔位，保證只有這個工人會下載
download(task.URL)
```

**改進效果：**

| 時間 | 工人 A | 工人 B | 結果 |
|-----|--------|--------|------|
| T1 | `CheckAndMark("page1")` → true ✓<br>（同時完成檢查+標記） | | |
| T2 | | `CheckAndMark("page1")` → false ✗<br>（已被 A 標記） | **B 被阻止** ✅ |
| T3 | 開始下載 "page1" | 跳過（已被佔位） | **只有 A 下載** ✅ |

### 🎓 併發安全原則：原子性（Atomicity）

> **相關的檢查和修改操作必須在同一個鎖的保護下完成，確保原子性。**

這是併發程式設計中最重要的原則之一！

---

### 1. DownloadManager 結構

```go
type DownloadManager struct {
    downloaded   map[string]bool  // 已下載的網址
    mutex        sync.Mutex       // 保護 map
    totalBytes   atomic.Uint64    // 總位元組數
    successCount atomic.Uint32    // 成功計數
    skippedCount atomic.Uint32    // 跳過計數
    timeoutCount atomic.Uint32    // 超時計數
}
```

**為什麼這樣設計？**
- `map` 不是併發安全的，必須用 `Mutex` 保護
- 數字統計用 `atomic`，避免鎖競爭，性能更好

### 2. 去重機制（Mutex）

```go
func (dm *DownloadManager) IsDownloaded(url string) bool {
    dm.mutex.Lock()
    defer dm.mutex.Unlock()
    return dm.downloaded[url]
}

func (dm *DownloadManager) MarkDownloaded(url string) {
    dm.mutex.Lock()
    defer dm.mutex.Unlock()
    dm.downloaded[url] = true
}
```

**關鍵點：**
- 讀取和寫入 map 都需要鎖保護
- 使用 `defer` 確保鎖總是被釋放

### 3. 流量統計（Atomic）

```go
dm.totalBytes.Add(uint64(bytes))    // 累加位元組
dm.successCount.Add(1)               // 成功計數 +1
dm.skippedCount.Add(1)               // 跳過計數 +1

// 讀取
total := dm.totalBytes.Load()
success := dm.successCount.Load()
```

**為什麼用 Atomic？**
- 簡單的數字累加不需要 Mutex（性能更好）
- `Add()` 和 `Load()` 都是原子操作（thread-safe）

### 4. 超時控制（Context）

```go
// 建立 2 秒超時的 Context
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

// 在下載時檢查超時
select {
case <-time.After(downloadTime):
    return bytes, nil  // 下載完成
case <-ctx.Done():
    return 0, ctx.Err()  // 超時
}
```

**Context 的作用：**
- 統一控制所有 Goroutine 的生命週期
- 超時時自動通知所有工人停止

### 5. Worker 生命週期

```go
func downloadWorker(id int, ctx context.Context, tasks <-chan DownloadTask, 
                    dm *DownloadManager, wg *sync.WaitGroup) {
    defer wg.Done()  // 工人退出前通知
    
    for {
        select {
        case task, ok := <-tasks:
            if !ok {
                return  // Channel 關閉，退出
            }
            // 處理任務...
            
        case <-ctx.Done():
            return  // Context 超時，退出
        }
    }
}
```

**工人退出的兩種情況：**
1. 任務 channel 關閉（所有任務完成）
2. Context 超時（時間到了）

## 💡 技術要點

### 1. Mutex vs Atomic

| 特性 | Mutex | Atomic |
|------|-------|--------|
| 用途 | 保護複雜資料結構（map, slice） | 簡單數字操作 |
| 性能 | 較慢（需要鎖） | 較快（無鎖） |
| 適用 | 讀寫 map | 計數器、累加器 |

```go
// ✅ 正確：map 用 Mutex
dm.mutex.Lock()
dm.downloaded[url] = true
dm.mutex.Unlock()

// ✅ 正確：計數器用 Atomic
dm.successCount.Add(1)

// ❌ 錯誤：不能對 map 用 atomic
// atomic 只支援數字類型
```

### 2. Context 傳遞

```go
// ✅ 正確：Context 傳遞給所有工人
for w := 1; w <= 5; w++ {
    go downloadWorker(w, ctx, tasks, dm, &wg)
}

// ✅ 正確：在任務中檢查 Context
select {
case <-ctx.Done():
    return  // 超時退出
}
```

### 3. WaitGroup 的正確使用

```go
// 1. 在啟動工人前 Add
wg.Add(1)
go downloadWorker(...)

// 2. 在工人函數中 defer Done
func downloadWorker(...) {
    defer wg.Done()  // 確保一定會執行
    ...
}

// 3. 等待所有工人完成
wg.Wait()
```

## 🎓 學習重點

### 併發安全

1. **Map 不是併發安全的**
   - 多個 Goroutine 同時讀寫會 panic
   - 必須使用 Mutex 保護

2. **Atomic 操作是併發安全的**
   - 用於簡單的數字操作
   - 性能比 Mutex 好

### 超時控制

1. **Context 的優勢**
   - 統一控制所有 Goroutine
   - 可以傳遞取消信號
   - 設定超時非常簡單

2. **Select 的威力**
   - 同時監聽多個 channel
   - 實現超時、取消等複雜邏輯

### Worker Pool 模式

1. **固定數量的工人**
   - 避免創建過多 Goroutine
   - 控制資源使用

2. **任務隊列**
   - 使用 buffered channel
   - 工人自動從隊列取任務

## 📈 性能分析

### 下載時間

**串行處理（1 個工人）：**
- 20 個任務 × 平均 500ms = 10 秒

**並行處理（5 個工人）：**
- 約 2 秒（受限於超時設定）
- 完成 15 個任務

**加速比：**
- 在相同時間內完成更多任務
- 充分利用等待時間（I/O 操作）

### Atomic vs Mutex

測試場景：1000000 次計數操作

| 方法 | 時間 | 相對性能 |
|------|------|----------|
| Atomic | ~10ms | 1.0x |
| Mutex | ~50ms | 0.2x |

**結論：Atomic 比 Mutex 快 5 倍**（僅限簡單數字操作）

## ⚠️ 常見錯誤

### 1. Check-Then-Act Race Condition（最常見！）

這是本練習修正的**最重要問題**！

```go
// ❌ 錯誤：分開檢查和標記
if dm.IsDownloaded(url) {
    return
}
// 危險空隙：其他工人可能同時通過檢查
dm.MarkDownloaded(url)

// ✅ 正確：原子性地檢查並標記
if !dm.CheckAndMark(url) {
    return  // 已被佔位
}
// 安全：已經佔位，不會重複下載
```

**為什麼會這樣？**

即使 `IsDownloaded` 和 `MarkDownloaded` 各自都是線程安全的（都有鎖保護），但是：

1. `IsDownloaded()` 執行完後**釋放鎖**
2. 在這個空隙中，其他工人可以執行 `IsDownloaded()`
3. 多個工人都得到 `false`，都認為可以下載
4. 結果：**重複下載**！

這就是經典的 **"時間窗口"（Time Window）** 問題。

**解決方案：**
- 使用 `CheckAndMark` 在**同一個鎖內**完成檢查和標記
- 確保操作的**原子性（Atomicity）**

### 2. 忘記保護 Map

```go
// ❌ 錯誤：直接訪問 map（併發時會 panic）
if dm.downloaded[url] {
    return
}

// ✅ 正確：使用 CheckAndMark（同時解決保護和原子性問題）
if !dm.CheckAndMark(url) {
    return
}
```

### 3. 濫用 Mutex

```go
// ❌ 錯誤：簡單計數也用 Mutex
dm.mutex.Lock()
dm.count++
dm.mutex.Unlock()

// ✅ 正確：簡單計數用 Atomic
dm.count.Add(1)
```

### 3. 忘記檢查 Context

```go
// ❌ 錯誤：沒有檢查超時
time.Sleep(500 * time.Millisecond)
return bytes, nil

// ✅ 正確：使用 select 檢查
select {
case <-time.After(500 * time.Millisecond):
    return bytes, nil
case <-ctx.Done():
    return 0, ctx.Err()
}
```

### 4. WaitGroup 計數錯誤

```go
// ❌ 錯誤：在 Goroutine 內 Add
go func() {
    wg.Add(1)  // 可能在 Wait 之後才執行
    defer wg.Done()
}()

// ✅ 正確：在啟動前 Add
wg.Add(1)
go func() {
    defer wg.Done()
}()
```

## 🎯 適用場景

這個模式適合：

1. **批量網路請求**
   - 爬蟲、API 批量調用
   - 圖片批量下載

2. **資料處理**
   - 批量讀取檔案
   - 資料庫批量查詢

3. **任務調度**
   - 定時任務管理
   - 工作流引擎

4. **微服務**
   - 服務間並行調用
   - 降級和熔斷

## 📚 延伸學習

1. **錯誤處理**：如何收集所有工人的錯誤？
2. **優先級隊列**：不同優先級的任務如何處理？
3. **動態調整**：根據負載動態增減工人數？
4. **進度報告**：如何實時報告下載進度？
5. **斷點續傳**：下載失敗如何重試？

## 🎉 總結

這個練習整合了 Go 併發程式設計的核心概念：

| 概念 | 用途 | 關鍵工具 |
|------|------|----------|
| Worker Pool | 控制併發數 | Goroutine + Channel |
| 去重機制 | 保護共享 map | Mutex |
| 流量統計 | 無鎖計數 | Atomic |
| 超時控制 | 生命週期管理 | Context |
| 同步等待 | 等待完成 | WaitGroup |

掌握這個模式，你就掌握了 Go 併發程式設計的精髓！🚀
