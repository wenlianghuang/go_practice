# Concurrent Cache 練習說明

## 📚 練習目標

學習如何使用 Mutex 和 RWMutex 來保護共享資源（map），避免 race condition。

## 🎯 包含的版本

### 1. `ConcurrentCacheUnsafe` - 不安全版本
- ❌ **沒有使用任何鎖機制**
- ⚠️ 會產生 `concurrent map writes` 錯誤
- 用於演示為什麼需要同步機制

### 2. `ConcurrentCacheSafe` - Mutex 版本
- ✅ 使用 `sync.Mutex` 保護
- 任何時候只有一個 Goroutine 可以訪問（讀或寫）
- 簡單、安全、適合大多數場景

### 3. `ConcurrentCacheRWMutex` - 讀寫鎖版本
- ✅ 使用 `sync.RWMutex` 保護
- 允許多個 Goroutine 同時讀取
- 寫入時獨佔鎖
- 讀多寫少的場景性能更好

### 4. `ConcurrentCacheComparison` - 性能比較
- 比較 Mutex 和 RWMutex 的性能差異
- 展示在讀多寫少場景下 RWMutex 的優勢

## 🚀 如何運行

```bash
# 測試不安全版本（會崩潰）
go run main.go ConcurrentCacheUnsafe

# 測試 Mutex 版本（安全）
go run main.go ConcurrentCacheSafe

# 測試 RWMutex 版本（安全且在讀多寫少時更快）
go run main.go ConcurrentCacheRWMutex

# 性能比較
go run main.go ConcurrentCacheComparison
```

## 🔍 使用 Race Detector

Go 提供了內建的 race detector 來檢測 race condition：

```bash
# 使用 race detector 運行（會檢測到不安全版本的問題）
go run -race main.go ConcurrentCacheUnsafe

# 驗證安全版本沒有 race condition
go run -race main.go ConcurrentCacheSafe
```

## 💡 關鍵知識點

### Mutex vs RWMutex

| 特性 | Mutex | RWMutex |
|------|-------|---------|
| 同時讀取 | ❌ 不允許 | ✅ 允許多個 |
| 同時寫入 | ❌ 不允許 | ❌ 不允許 |
| 讀取時寫入 | ❌ 不允許 | ❌ 不允許 |
| 適用場景 | 讀寫比例相近 | 讀多寫少 |
| 性能 | 一般 | 讀多時更好 |

### 為什麼需要鎖？

Go 的 map **不是**併發安全的：
- 多個 Goroutine 同時讀取：✅ 可以
- 多個 Goroutine 同時寫入：❌ 會崩潰
- 一個讀取一個寫入：❌ 會崩潰

### 常見錯誤

```
fatal error: concurrent map writes
```

這個錯誤表示多個 Goroutine 同時寫入 map，必須使用鎖來保護。

## 📊 測試結果示例

```
=== Cache 實作比較（安全版本） ===

1️⃣  測試 Mutex 版本...
2️⃣  測試 RWMutex 版本...

📊 性能比較結果：
   Mutex 版本:   2.845333ms
   RWMutex 版本: 1.900791ms
   ✓ RWMutex 快了 33.20%！

💡 知識點：
   - Mutex: 任何時候只有一個 Goroutine 可以訪問（讀或寫）
   - RWMutex: 多個 Goroutine 可以同時讀取，但寫入時獨佔
   - 當讀取操作遠多於寫入時，RWMutex 性能更好
```

## 🎓 學習重點

1. **Race Condition** - 多個 Goroutine 訪問共享資源時的問題
2. **Mutex** - 互斥鎖，確保同一時間只有一個 Goroutine 訪問
3. **RWMutex** - 讀寫鎖，允許多個讀取者，但寫入時獨佔
4. **defer unlock** - 確保鎖總是會被釋放
5. **WaitGroup** - 等待所有 Goroutine 完成

## 🔗 相關檔案

- `concurrent-cache.go` - 實作程式碼
- `main.go` - 主程式入口
