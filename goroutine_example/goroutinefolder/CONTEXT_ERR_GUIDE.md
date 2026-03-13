# Context.Err() 詳解

## 🎯 什麼是 Context.Err()？

`Context.Err()` 返回 Context 結束的原因，幫助我們了解任務為什麼停止。

## 📋 兩種錯誤類型

Go 的 `context` 包定義了兩種標準錯誤：

| 錯誤類型 | 常數 | 何時發生 | 說明 |
|---------|------|---------|------|
| **DeadlineExceeded** | `context.DeadlineExceeded` | 超過設定的時限 | 任務執行超時 |
| **Canceled** | `context.Canceled` | 手動調用 `cancel()` | 主動取消任務 |

## 🔍 在下載器中的應用

### 1. 下載函數中的錯誤處理

```go
func simulateDownload(ctx context.Context, url string) (int, error) {
    downloadTime := time.Millisecond * 500
    
    select {
    case <-time.After(downloadTime):
        return 1024, nil  // 下載成功
        
    case <-ctx.Done():
        // Context 結束，返回原因
        return 0, ctx.Err()  // 🔑 返回具體錯誤
    }
}
```

### 2. Worker 中檢查錯誤類型

```go
bytes, err := simulateDownload(ctx, task.URL)
if err != nil {
    // 🔑 檢查具體的錯誤類型
    if err == context.DeadlineExceeded {
        fmt.Printf("⏰ 下載超時 (DeadlineExceeded)\n")
    } else if err == context.Canceled {
        fmt.Printf("🚫 下載被取消 (Canceled)\n")
    } else {
        fmt.Printf("❌ 其他錯誤: %v\n", err)
    }
    return
}
```

### 3. 在 select 中直接檢查

```go
select {
case task := <-tasks:
    // 處理任務...
    
case <-ctx.Done():
    // 🔑 Context 結束時檢查原因
    err := ctx.Err()
    if err == context.DeadlineExceeded {
        fmt.Printf("⏰ 超時停止\n")
    } else if err == context.Canceled {
        fmt.Printf("🚫 取消停止\n")
    }
    return
}
```

## 📊 實際輸出示例

### 短超時版本（1 秒）

```
📥 工人 1: 開始下載 [page5]
📥 工人 2: 開始下載 [page7]
✅ 工人 1: 完成 [page5] (4206 bytes)
⏰ 工人 2: 下載 [page7] 超時 (DeadlineExceeded)  ← 下載中超時
⏰ 工人 3: 收到超時信號 (DeadlineExceeded)，停止工作  ← select 收到超時
```

### 主動取消版本

```go
// 創建可取消的 Context
ctx, cancel := context.WithCancel(context.Background())

// 啟動工人...

// 某個條件下主動取消
if someCondition {
    cancel()  // 主動取消
}
```

輸出：
```
🚫 工人 1: 下載被取消 (Canceled)
🚫 工人 2: 收到取消信號 (Canceled)，停止工作
```

## 🎓 為什麼要區分錯誤類型？

### DeadlineExceeded（超時）

**原因：**
- 系統資源不足
- 任務執行太慢
- 網路延遲過高
- 超時設定太短

**應對策略：**
```go
if err == context.DeadlineExceeded {
    // 1. 記錄日誌，監控超時頻率
    log.Warn("Download timeout for", url)
    
    // 2. 可能需要調整超時時間
    // 3. 可能需要優化下載邏輯
    // 4. 可能需要重試機制
    
    metrics.IncrementTimeout()  // 增加超時計數
}
```

### Canceled（取消）

**原因：**
- 用戶主動取消
- 程序正常關閉
- 上層邏輯決定停止

**應對策略：**
```go
if err == context.Canceled {
    // 這是正常行為，不需要告警
    log.Info("Download canceled for", url)
    
    // 清理資源
    cleanup()
    
    // 不需要重試
    return
}
```

## 💡 最佳實踐

### 1. 總是檢查 ctx.Err()

```go
// ✅ 良好實踐
select {
case <-ctx.Done():
    switch ctx.Err() {
    case context.DeadlineExceeded:
        handleTimeout()
    case context.Canceled:
        handleCancellation()
    }
}

// ❌ 不好的實踐
select {
case <-ctx.Done():
    // 沒有檢查原因，無法區分超時和取消
    return
}
```

### 2. 使用 switch 語句

```go
// ✅ 清晰的錯誤處理
if err := ctx.Err(); err != nil {
    switch err {
    case context.DeadlineExceeded:
        log.Error("Timeout")
        metrics.RecordTimeout()
        return err
        
    case context.Canceled:
        log.Info("Canceled")
        return nil  // 取消不是錯誤
        
    default:
        log.Error("Unknown context error:", err)
        return err
    }
}
```

### 3. 在生產環境中添加監控

```go
if err == context.DeadlineExceeded {
    // 記錄超時指標
    prometheus.CounterVec.WithLabelValues("timeout").Inc()
    
    // 記錄慢查詢
    if duration > threshold {
        slowQueryLog.Record(url, duration)
    }
}
```

## 🔬 深入理解

### Context 的生命週期

```go
// 1. 創建 Context
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()  // 確保釋放資源

// 2. Context 正常運行
// ctx.Err() == nil

// 3a. 超時發生
// time.After(2*time.Second)
// ctx.Err() == context.DeadlineExceeded

// 3b. 或者手動取消
// cancel()
// ctx.Err() == context.Canceled
```

### 錯誤傳播

```go
func downloadFile(ctx context.Context, url string) error {
    // 內層函數返回 ctx.Err()
    data, err := fetchData(ctx, url)
    if err != nil {
        // err 可能是 context.DeadlineExceeded 或 context.Canceled
        return err
    }
    
    return saveData(data)
}

// 調用者可以檢查
if err := downloadFile(ctx, url); err != nil {
    if err == context.DeadlineExceeded {
        // 處理超時
    }
}
```

## 📈 實際場景示例

### 場景 1：批量下載（超時）

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

for _, url := range urls {
    err := download(ctx, url)
    if err == context.DeadlineExceeded {
        fmt.Printf("⏰ 總時限到了，停止下載\n")
        break  // 停止處理剩餘任務
    }
}
```

### 場景 2：用戶取消（優雅退出）

```go
ctx, cancel := context.WithCancel(context.Background())

// 監聽用戶取消信號
go func() {
    <-interruptSignal
    cancel()  // 用戶按 Ctrl+C
}()

// 工人檢查取消
if err == context.Canceled {
    fmt.Println("🚫 用戶取消了下載")
    saveProgress()  // 保存進度
    return nil      // 正常退出
}
```

### 場景 3：混合使用

```go
// 同時支持超時和取消
ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
defer cancel()

err := processTask(ctx)
switch err {
case context.DeadlineExceeded:
    log.Warn("Task timeout after 10s")
    return ErrTimeout
    
case context.Canceled:
    log.Info("Task canceled by user")
    return nil  // 取消不是錯誤
    
case nil:
    log.Info("Task completed successfully")
    return nil
    
default:
    log.Error("Task failed:", err)
    return err
}
```

## 🎯 總結

### Context.Err() 的重要性

1. **區分原因**：了解任務為何停止
2. **正確處理**：超時需要告警，取消是正常行為
3. **監控指標**：記錄超時率，優化系統
4. **用戶體驗**：顯示清楚的錯誤信息

### 記住

- ✅ **DeadlineExceeded** = 系統問題，需要監控和優化
- ✅ **Canceled** = 正常取消，不需要告警
- ✅ 總是檢查 `ctx.Err()` 來了解原因
- ✅ 根據不同錯誤類型採取不同行動

### 在下載器中的體現

```go
⏰ 工人 1: 下載 [page5] 超時 (DeadlineExceeded)  ← 明確顯示原因
🚫 工人 2: 收到取消信號 (Canceled)，停止工作   ← 區分取消和超時
```

這樣的錯誤處理讓程式更加健壯和易於維護！🚀
