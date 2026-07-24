package middleware

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golangAPI_construct/cache"
	"golangAPI_construct/logging"
)

// CacheConfig 緩存配置
// 定義緩存的行為和策略
type CacheConfig struct {
	// Duration 緩存過期時間
	Duration time.Duration

	// KeyGenerator 自定義鍵生成函數
	// 如果為 nil，使用預設的鍵生成策略
	KeyGenerator func(r *http.Request) string

	// SkipCache 決定是否跳過緩存的函數
	// 返回 true 時不會使用緩存
	SkipCache func(r *http.Request) bool

	// VaryHeaders 需要考慮的請求頭
	// 這些頭會影響緩存鍵的生成
	VaryHeaders []string
}

// CacheMiddleware 創建緩存中間件
// 這個中間件會自動緩存 GET 請求的響應，並在後續相同請求時直接返回緩存
func CacheMiddleware(cacheService cache.CacheService, config CacheConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 只緩存 GET 請求
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			// 檢查是否應該跳過緩存
			if config.SkipCache != nil && config.SkipCache(r) {
				next.ServeHTTP(w, r)
				return
			}

			// 生成緩存鍵
			cacheKey := generateCacheKey(r, config)

			// 嘗試從緩存獲取響應
			ctx := r.Context()
			cachedResponse, err := cacheService.Get(ctx, cacheKey)
			if err != nil {
				logging.Logger.Printf("[CACHE] Error getting cache for key %s: %v", cacheKey, err)
				next.ServeHTTP(w, r)
				return
			}

			// 如果找到緩存，直接返回
			if cachedResponse != "" {
				logging.Logger.Printf("[CACHE] Cache hit for key: %s", cacheKey)

				// 設置緩存相關的響應頭
				w.Header().Set("X-Cache", "HIT")
				w.Header().Set("X-Cache-Key", cacheKey)
				w.Header().Set("Content-Type", "application/json")

				// 寫入緩存的響應
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(cachedResponse))
				return
			}

			// 緩存未命中，執行原始處理器並緩存響應
			logging.Logger.Printf("[CACHE] Cache miss for key: %s", cacheKey)

			// 創建響應寫入器包裝器來捕獲響應
			wrapped := &cacheResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           make([]byte, 0),
			}

			// 執行下一個中間件/處理器
			next.ServeHTTP(wrapped, r)

			// 只緩存成功的響應（2xx 狀態碼）
			if wrapped.statusCode >= 200 && wrapped.statusCode < 300 {
				// 將響應存儲到緩存
				if err := cacheService.Set(ctx, cacheKey, string(wrapped.body), config.Duration); err != nil {
					logging.Logger.Printf("[CACHE] Error setting cache for key %s: %v", cacheKey, err)
				} else {
					logging.Logger.Printf("[CACHE] Cached response for key: %s", cacheKey)
				}

				// 設置緩存相關的響應頭
				wrapped.Header().Set("X-Cache", "MISS")
				wrapped.Header().Set("X-Cache-Key", cacheKey)
			}
		})
	}
}

// cacheResponseWriter 包裝 http.ResponseWriter 來捕獲響應內容
type cacheResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

// WriteHeader 捕獲狀態碼
func (crw *cacheResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// Write 捕獲響應內容
func (crw *cacheResponseWriter) Write(data []byte) (int, error) {
	crw.body = append(crw.body, data...)
	return crw.ResponseWriter.Write(data)
}

// generateCacheKey 生成緩存鍵
// 基於請求的 URL、方法和相關頭部生成唯一的緩存鍵
func generateCacheKey(r *http.Request, config CacheConfig) string {
	// 如果提供了自定義鍵生成函數，使用它
	if config.KeyGenerator != nil {
		return config.KeyGenerator(r)
	}

	// 預設的鍵生成策略
	var keyParts []string

	// 添加方法
	keyParts = append(keyParts, r.Method)

	// 添加路徑
	keyParts = append(keyParts, r.URL.Path)

	// 添加查詢參數（排序後）
	if r.URL.RawQuery != "" {
		keyParts = append(keyParts, "?"+r.URL.RawQuery)
	}

	// 添加相關的請求頭
	for _, headerName := range config.VaryHeaders {
		if headerValue := r.Header.Get(headerName); headerValue != "" {
			keyParts = append(keyParts, headerName+":"+headerValue)
		}
	}

	// 組合所有部分
	fullKey := strings.Join(keyParts, "|")

	// 生成 MD5 哈希以確保鍵長度一致
	hash := md5.Sum([]byte(fullKey))
	return fmt.Sprintf("cache:%x", hash)
}

// 預定義的緩存配置
var (
	// DefaultCacheConfig 預設緩存配置
	// 緩存 5 分鐘，考慮 Authorization 頭
	DefaultCacheConfig = CacheConfig{
		Duration:    5 * time.Minute,
		VaryHeaders: []string{"Authorization"},
		SkipCache: func(r *http.Request) bool {
			// 跳過包含 no-cache 頭的請求
			return r.Header.Get("Cache-Control") == "no-cache"
		},
	}

	// ShortCacheConfig 短期緩存配置
	// 緩存 1 分鐘，適合頻繁變化的數據
	ShortCacheConfig = CacheConfig{
		Duration:    1 * time.Minute,
		VaryHeaders: []string{"Authorization"},
		SkipCache: func(r *http.Request) bool {
			return r.Header.Get("Cache-Control") == "no-cache"
		},
	}

	// LongCacheConfig 長期緩存配置
	// 緩存 1 小時，適合相對穩定的數據
	LongCacheConfig = CacheConfig{
		Duration:    1 * time.Hour,
		VaryHeaders: []string{"Authorization"},
		SkipCache: func(r *http.Request) bool {
			return r.Header.Get("Cache-Control") == "no-cache"
		},
	}

	// PublicCacheConfig 公開緩存配置
	// 緩存 10 分鐘，不考慮 Authorization 頭（用於公開數據）
	PublicCacheConfig = CacheConfig{
		Duration:    10 * time.Minute,
		VaryHeaders: []string{}, // 不考慮任何頭
		SkipCache: func(r *http.Request) bool {
			return r.Header.Get("Cache-Control") == "no-cache"
		},
	}
)

// CacheInvalidationMiddleware 緩存失效中間件
// 當數據被修改時，自動清除相關的緩存
func CacheInvalidationMiddleware(cacheService cache.CacheService, patterns []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 只對修改操作進行緩存失效
			if r.Method == http.MethodGet || r.Method == http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}

			// 執行原始處理器
			wrapped := &cacheResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           make([]byte, 0),
			}

			next.ServeHTTP(wrapped, r)

			// 如果操作成功，清除相關緩存
			if wrapped.statusCode >= 200 && wrapped.statusCode < 300 {
				ctx := r.Context()
				for _, pattern := range patterns {
					// 這裡可以實現更複雜的緩存清除邏輯
					// 例如使用 Redis 的 SCAN 命令來查找匹配的鍵
					logging.Logger.Printf("[CACHE] Invalidating cache pattern: %s", pattern)

					// 簡單實現：清除所有以 pattern 開頭的緩存
					// 在實際應用中，你可能需要使用更精確的清除策略
					_ = cacheService.Delete(ctx, pattern)
				}
			}
		})
	}
}
