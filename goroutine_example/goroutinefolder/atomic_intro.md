# `sync/atomic` 原子操作介紹

在 Go 語言中，當多個 Goroutine 同時修改同一個變數時，會發生 **競態條件 (Race Condition)**，導致資料錯誤。雖然可以使用 `sync.Mutex` 鎖來保護，但對於簡單的計數器或標記，`sync/atomic` 提供了更輕量、更底層的硬體級別支援。

### 1. 什麼是原子操作 (Atomic Operation)？

「原子」的意思是 **不可分割**。
一般的程式碼如 `count++` 在電腦底層其實分為三個步驟：
1.  **讀取 (Load)**: 從記憶體讀取 `count` 目前的值。
2.  **修改 (Modify)**: CPU 將值加 1。
3.  **寫入 (Store)**: 將新值寫回記憶體。

如果在步驟 1 和 3 之間，另一個 Goroutine 也跑進來讀取或寫入，就會發生衝突。
**原子操作** 則保證這三個步驟是「瞬間完成」的，中間絕不會被其他 Goroutine 插隊。

### 2. 主要功能

`sync/atomic` 套件提供了幾種主要的操作類型（針對 `int32`, `int64`, `uint32` 等）：

| 操作類型 | 函式範例 | 說明 |
| :--- | :--- | :--- |
| **增加** | `AddInt64(&addr, n)` | 安全地將數值增加 n (n 可以是負數) |
| **讀取** | `LoadInt64(&addr)` | 安全地讀取變數目前的值 |
| **寫入** | `StoreInt64(&addr, val)` | 安全地設定變數的值 |
| **交換** | `SwapInt64(&addr, new)` | 設定新值，並回傳舊值 |
| **比較並交換** | `CompareAndSwapInt64(...)` | (CAS) 只有當舊值等於預期時才修改，常用於實作樂觀鎖 |

### 3. 何時使用 Atomic vs Mutex？

*   **使用 `sync/atomic`**:
    *   單純的計數器 (Counter)、累積量。
    *   簡單的狀態開關 (0 或 1)。
    *   追求極致效能且邏輯非常簡單時。

*   **使用 `sync.Mutex`**:
    *   需要保護一段複雜的程式碼邏輯（不只是一行變數修改）。
    *   涉及多個變數需要一起變更時。
    *   需要使用 `map`, `slice` 等複雜結構時。

### 4. 範例程式碼重點 (`goroutinefolder/atomic.go`)

在 `goroutinefolder/atomic.go` 範例中：
*   **普通變數 (`unsafeCounter++`)**: 會因為並發衝突導致最終結果小於預期。
*   **原子變數 (`atomic.AddInt64`)**: 保證每次加法都是安全的，最終結果必定正確。
