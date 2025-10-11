# Go API Construct - Redis 緩存系統

## 功能概述

這個 Go API 項目現在包含了完整的 Redis 緩存系統，可以顯著提升應用程式的性能和可擴展性。

## 緩存功能

### 1. 自動緩存
- **GET 請求緩存**：自動緩存 GET 請求的響應
- **智能鍵生成**：基於 URL、查詢參數和相關頭部生成唯一緩存鍵
- **過期策略**：支持可配置的緩存過期時間

### 2. 緩存失效
- **自動失效**：當數據被修改時自動清除相關緩存
- **批量清除**：支持批量清除多個緩存鍵
- **模式匹配**：支持基於模式的緩存清除

### 3. 多種緩存後端
- **Redis**：生產環境推薦，高性能分散式緩存
- **記憶體緩存**：開發環境和 Redis 不可用時的備選方案

## 環境變數配置

創建 `.env` 文件並添加以下配置：

```bash
# 啟用緩存
CACHE_ENABLED=true

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_PREFIX=golangapi:
```

## 緩存策略

### 預定義配置

1. **DefaultCacheConfig**：5分鐘緩存，考慮 Authorization 頭
2. **ShortCacheConfig**：1分鐘緩存，適合頻繁變化的數據
3. **LongCacheConfig**：1小時緩存，適合相對穩定的數據
4. **PublicCacheConfig**：10分鐘緩存，不考慮 Authorization 頭

### 使用範例

```go
// 在路由中使用緩存
r.With(middleware.CacheMiddleware(cacheService, middleware.DefaultCacheConfig)).Get("/books", handler)

// 緩存失效
r.With(middleware.CacheInvalidationMiddleware(cacheService, []string{"books"})).Post("/books", handler)
```

## 性能優勢

1. **響應時間**：緩存命中時響應時間可減少 90% 以上
2. **資料庫負載**：減少資料庫查詢次數，降低資料庫負載
3. **並發能力**：提高應用程式的並發處理能力
4. **用戶體驗**：更快的頁面加載和 API 響應

## 監控和調試

### 緩存狀態頭
- `X-Cache: HIT` - 緩存命中
- `X-Cache: MISS` - 緩存未命中
- `X-Cache-Key` - 使用的緩存鍵

### 日誌記錄
```
[CACHE] Cache hit for key: cache:abc123
[CACHE] Cache miss for key: cache:def456
[CACHE] Cached response for key: cache:ghi789
```

## 部署建議

### 開發環境
```bash
CACHE_ENABLED=false  # 使用記憶體緩存
```

### 生產環境
```bash
CACHE_ENABLED=true
REDIS_HOST=your-redis-host
REDIS_PASSWORD=your-secure-password
REDIS_DB=0
```

## 故障處理

1. **Redis 不可用**：自動降級到記憶體緩存
2. **緩存錯誤**：記錄錯誤但不影響主要功能
3. **連接超時**：使用合理的超時設置避免阻塞

## 最佳實踐

1. **緩存鍵設計**：使用有意義的前綴和命名
2. **過期時間**：根據數據更新頻率設置合適的過期時間
3. **緩存失效**：及時清除過期或無效的緩存
4. **監控**：監控緩存命中率和性能指標
5. **測試**：在測試環境中驗證緩存行為

## 擴展功能

未來可以添加的功能：
- 分散式緩存鎖
- 緩存預熱
- 緩存統計和監控
- 多級緩存策略
- 緩存壓縮
