package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"golangAPI_construct/logging"

	"github.com/redis/go-redis/v9"
)

// CacheService 定義緩存服務的接口
// 這個接口提供了統一的緩存操作，可以輕鬆切換不同的緩存實現
type CacheService interface {
	// 基本操作
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// 批量操作
	GetMultiple(ctx context.Context, keys []string) (map[string]string, error)
	SetMultiple(ctx context.Context, data map[string]interface{}, expiration time.Duration) error
	DeleteMultiple(ctx context.Context, keys []string) error

	// 高級操作
	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// 健康檢查
	Ping(ctx context.Context) error
	Close() error
}

// RedisCacheService Redis 緩存服務實現
// 這個結構體封裝了 Redis 客戶端，提供了完整的緩存功能
type RedisCacheService struct {
	client *redis.Client
	prefix string // 鍵前綴，用於避免鍵衝突
}

// NewRedisCacheService 創建新的 Redis 緩存服務實例
// 從環境變數讀取 Redis 配置，如果連接失敗會返回錯誤
func NewRedisCacheService() (*RedisCacheService, error) {
	// 從環境變數讀取 Redis 配置
	host := getEnvOrDefault("REDIS_HOST", "localhost")
	port := getEnvOrDefault("REDIS_PORT", "6379")
	password := os.Getenv("REDIS_PASSWORD") // 密碼是可選的
	db := getEnvOrDefault("REDIS_DB", "0")
	prefix := getEnvOrDefault("REDIS_PREFIX", "golangapi:")

	// 解析資料庫編號
	dbNum, err := strconv.Atoi(db)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %v", err)
	}

	// 創建 Redis 客戶端
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       dbNum,

		// 連接池配置
		PoolSize:     10, // 連接池大小
		MinIdleConns: 5,  // 最小空閒連接數
		MaxIdleConns: 10, // 最大空閒連接數
		MaxRetries:   3,  // 最大重試次數

		// 超時配置
		DialTimeout:  5 * time.Second, // 連接超時
		ReadTimeout:  3 * time.Second, // 讀取超時
		WriteTimeout: 3 * time.Second, // 寫入超時
		PoolTimeout:  4 * time.Second, // 連接池超時
	})

	// 測試連接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	logging.Logger.Printf("[CACHE] Connected to Redis at %s:%s (DB: %d)", host, port, dbNum)

	return &RedisCacheService{
		client: rdb,
		prefix: prefix,
	}, nil
}

// buildKey 構建完整的緩存鍵
// 添加前綴以避免與其他應用的鍵衝突
func (r *RedisCacheService) buildKey(key string) string {
	return r.prefix + key
}

// Get 獲取緩存值
func (r *RedisCacheService) Get(ctx context.Context, key string) (string, error) {
	fullKey := r.buildKey(key)
	result := r.client.Get(ctx, fullKey)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return "", nil // 鍵不存在，返回空字串
		}
		return "", result.Err()
	}
	return result.Val(), nil
}

// Set 設置緩存值
// value 可以是任何可以 JSON 序列化的類型
func (r *RedisCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	fullKey := r.buildKey(key)

	// 將值轉換為 JSON 字串
	var jsonValue string
	switch v := value.(type) {
	case string:
		jsonValue = v
	default:
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %v", err)
		}
		jsonValue = string(jsonBytes)
	}

	return r.client.Set(ctx, fullKey, jsonValue, expiration).Err()
}

// Delete 刪除緩存鍵
func (r *RedisCacheService) Delete(ctx context.Context, key string) error {
	fullKey := r.buildKey(key)
	return r.client.Del(ctx, fullKey).Err()
}

// Exists 檢查鍵是否存在
func (r *RedisCacheService) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := r.buildKey(key)
	result := r.client.Exists(ctx, fullKey)
	return result.Val() > 0, result.Err()
}

// GetMultiple 批量獲取緩存值
func (r *RedisCacheService) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	if len(keys) == 0 {
		return make(map[string]string), nil
	}

	// 構建完整的鍵列表
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = r.buildKey(key)
	}

	// 批量獲取
	result := r.client.MGet(ctx, fullKeys...)
	if result.Err() != nil {
		return nil, result.Err()
	}

	// 構建結果映射
	values := result.Val()
	resultMap := make(map[string]string)
	for i, key := range keys {
		if i < len(values) && values[i] != nil {
			resultMap[key] = values[i].(string)
		}
	}

	return resultMap, nil
}

// SetMultiple 批量設置緩存值
func (r *RedisCacheService) SetMultiple(ctx context.Context, data map[string]interface{}, expiration time.Duration) error {
	if len(data) == 0 {
		return nil
	}

	// 使用 Pipeline 提高性能
	pipe := r.client.Pipeline()

	for key, value := range data {
		fullKey := r.buildKey(key)

		// 將值轉換為 JSON 字串
		var jsonValue string
		switch v := value.(type) {
		case string:
			jsonValue = v
		default:
			jsonBytes, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal value for key %s: %v", key, err)
			}
			jsonValue = string(jsonBytes)
		}

		pipe.Set(ctx, fullKey, jsonValue, expiration)
	}

	// 執行 Pipeline
	_, err := pipe.Exec(ctx)
	return err
}

// DeleteMultiple 批量刪除緩存鍵
func (r *RedisCacheService) DeleteMultiple(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	// 構建完整的鍵列表
	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = r.buildKey(key)
	}

	return r.client.Del(ctx, fullKeys...).Err()
}

// Increment 增加計數器
func (r *RedisCacheService) Increment(ctx context.Context, key string) (int64, error) {
	fullKey := r.buildKey(key)
	result := r.client.Incr(ctx, fullKey)
	return result.Val(), result.Err()
}

// Decrement 減少計數器
func (r *RedisCacheService) Decrement(ctx context.Context, key string) (int64, error) {
	fullKey := r.buildKey(key)
	result := r.client.Decr(ctx, fullKey)
	return result.Val(), result.Err()
}

// Expire 設置鍵的過期時間
func (r *RedisCacheService) Expire(ctx context.Context, key string, expiration time.Duration) error {
	fullKey := r.buildKey(key)
	return r.client.Expire(ctx, fullKey, expiration).Err()
}

// Ping 測試 Redis 連接
func (r *RedisCacheService) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close 關閉 Redis 連接
func (r *RedisCacheService) Close() error {
	return r.client.Close()
}

// MemoryCacheService 記憶體緩存服務實現（用於開發和測試）
// 當 Redis 不可用時，可以使用這個實現作為備選方案
type MemoryCacheService struct {
	data   map[string]cacheItem
	prefix string
}

type cacheItem struct {
	value      string
	expiration time.Time
}

// NewMemoryCacheService 創建記憶體緩存服務
func NewMemoryCacheService() *MemoryCacheService {
	return &MemoryCacheService{
		data:   make(map[string]cacheItem),
		prefix: "memory:",
	}
}

func (m *MemoryCacheService) buildKey(key string) string {
	return m.prefix + key
}

func (m *MemoryCacheService) Get(ctx context.Context, key string) (string, error) {
	fullKey := m.buildKey(key)
	item, exists := m.data[fullKey]
	if !exists {
		return "", nil
	}

	// 檢查是否過期
	if time.Now().After(item.expiration) {
		delete(m.data, fullKey)
		return "", nil
	}

	return item.value, nil
}

func (m *MemoryCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	fullKey := m.buildKey(key)

	var jsonValue string
	switch v := value.(type) {
	case string:
		jsonValue = v
	default:
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %v", err)
		}
		jsonValue = string(jsonBytes)
	}

	m.data[fullKey] = cacheItem{
		value:      jsonValue,
		expiration: time.Now().Add(expiration),
	}

	return nil
}

func (m *MemoryCacheService) Delete(ctx context.Context, key string) error {
	fullKey := m.buildKey(key)
	delete(m.data, fullKey)
	return nil
}

func (m *MemoryCacheService) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := m.buildKey(key)
	item, exists := m.data[fullKey]
	if !exists {
		return false, nil
	}

	// 檢查是否過期
	if time.Now().After(item.expiration) {
		delete(m.data, fullKey)
		return false, nil
	}

	return true, nil
}

// 其他方法的簡化實現...
func (m *MemoryCacheService) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, key := range keys {
		if value, err := m.Get(ctx, key); err == nil && value != "" {
			result[key] = value
		}
	}
	return result, nil
}

func (m *MemoryCacheService) SetMultiple(ctx context.Context, data map[string]interface{}, expiration time.Duration) error {
	for key, value := range data {
		if err := m.Set(ctx, key, value, expiration); err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryCacheService) DeleteMultiple(ctx context.Context, keys []string) error {
	for _, key := range keys {
		m.Delete(ctx, key)
	}
	return nil
}

func (m *MemoryCacheService) Increment(ctx context.Context, key string) (int64, error) {
	value, _ := m.Get(ctx, key)
	count, _ := strconv.ParseInt(value, 10, 64)
	count++
	m.Set(ctx, key, strconv.FormatInt(count, 10), time.Hour)
	return count, nil
}

func (m *MemoryCacheService) Decrement(ctx context.Context, key string) (int64, error) {
	value, _ := m.Get(ctx, key)
	count, _ := strconv.ParseInt(value, 10, 64)
	count--
	m.Set(ctx, key, strconv.FormatInt(count, 10), time.Hour)
	return count, nil
}

func (m *MemoryCacheService) Expire(ctx context.Context, key string, expiration time.Duration) error {
	value, err := m.Get(ctx, key)
	if err != nil {
		return err
	}
	return m.Set(ctx, key, value, expiration)
}

func (m *MemoryCacheService) Ping(ctx context.Context) error {
	return nil // 記憶體緩存總是可用的
}

func (m *MemoryCacheService) Close() error {
	return nil
}

// getEnvOrDefault 獲取環境變數或返回預設值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewCacheService 創建緩存服務實例
// 根據環境變數決定使用 Redis 還是記憶體緩存
func NewCacheService() (CacheService, error) {
	// 檢查是否啟用緩存
	if os.Getenv("CACHE_ENABLED") != "true" {
		logging.Logger.Print("[CACHE] Cache disabled, using memory cache")
		return NewMemoryCacheService(), nil
	}

	// 嘗試創建 Redis 緩存服務
	redisService, err := NewRedisCacheService()
	if err != nil {
		logging.Logger.Printf("[CACHE] Failed to connect to Redis: %v, falling back to memory cache", err)
		return NewMemoryCacheService(), nil
	}

	return redisService, nil
}
