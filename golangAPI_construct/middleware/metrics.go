package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"golangAPI_construct/logging"
)

// Metrics 結構體用於存儲應用程式的各種指標
// 這些指標可以幫助監控應用程式的健康狀態和性能
type Metrics struct {
	// HTTP 請求相關指標
	RequestCount    map[string]int64   // 每個端點的請求總數
	RequestDuration map[string]float64 // 每個端點的平均響應時間（毫秒）
	ErrorCount      map[string]int64   // 每個端點的錯誤總數

	// 系統相關指標
	ActiveConnections int64   // 當前活躍連接數
	TotalRequests     int64   // 總請求數
	TotalErrors       int64   // 總錯誤數
	UptimeSeconds     float64 // 應用程式運行時間（秒）

	// 互斥鎖，用於保護並發訪問
	mu sync.RWMutex

	// 啟動時間，用於計算運行時間
	startTime time.Time
}

// 全局 metrics 實例
var globalMetrics *Metrics

// 初始化全局 metrics
func init() {
	globalMetrics = &Metrics{
		RequestCount:    make(map[string]int64),
		RequestDuration: make(map[string]float64),
		ErrorCount:      make(map[string]int64),
		startTime:       time.Now(),
	}
}

// MetricsMiddleware 創建一個 metrics 中間件
// 這個中間件會收集每個 HTTP 請求的詳細信息，包括：
// - 請求計數
// - 響應時間
// - 錯誤計數
// - 狀態碼分布
func MetricsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 增加活躍連接數
			globalMetrics.mu.Lock()
			globalMetrics.ActiveConnections++
			globalMetrics.TotalRequests++
			globalMetrics.mu.Unlock()

			// 創建響應寫入器包裝器來捕獲狀態碼
			wrapped := &metricsResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// 執行下一個中間件/處理器
			next.ServeHTTP(wrapped, r)

			// 計算請求處理時間
			duration := time.Since(start)
			durationMs := float64(duration.Nanoseconds()) / 1e6 // 轉換為毫秒

			// 構建端點標識符（方法 + 路徑）
			endpoint := r.Method + " " + r.URL.Path

			// 更新指標
			globalMetrics.mu.Lock()

			// 更新請求計數
			globalMetrics.RequestCount[endpoint]++

			// 更新平均響應時間（使用移動平均）
			if currentAvg, exists := globalMetrics.RequestDuration[endpoint]; exists {
				// 移動平均：新平均值 = (舊平均值 * (n-1) + 新值) / n
				count := globalMetrics.RequestCount[endpoint]
				globalMetrics.RequestDuration[endpoint] = (currentAvg*float64(count-1) + durationMs) / float64(count)
			} else {
				globalMetrics.RequestDuration[endpoint] = durationMs
			}

			// 如果是錯誤響應（4xx 或 5xx），增加錯誤計數
			if wrapped.statusCode >= 400 {
				globalMetrics.ErrorCount[endpoint]++
				globalMetrics.TotalErrors++
			}

			// 減少活躍連接數
			globalMetrics.ActiveConnections--

			// 更新運行時間
			globalMetrics.UptimeSeconds = time.Since(globalMetrics.startTime).Seconds()

			globalMetrics.mu.Unlock()

			// 記錄詳細的請求日誌（用於調試）
			logging.Logger.Printf("[METRICS] %s %s - Status: %d, Duration: %.2fms, Active: %d",
				r.Method, r.URL.Path, wrapped.statusCode, durationMs, globalMetrics.ActiveConnections)
		})
	}
}

// metricsResponseWriter 包裝 http.ResponseWriter 來捕獲狀態碼
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 捕獲狀態碼
func (mrw *metricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}

// GetMetrics 返回當前的指標數據
// 這個函數可以被 /metrics 端點調用來獲取 Prometheus 格式的指標
func GetMetrics() *Metrics {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	// 創建指標的副本以避免數據競爭
	metrics := &Metrics{
		RequestCount:      make(map[string]int64),
		RequestDuration:   make(map[string]float64),
		ErrorCount:        make(map[string]int64),
		ActiveConnections: globalMetrics.ActiveConnections,
		TotalRequests:     globalMetrics.TotalRequests,
		TotalErrors:       globalMetrics.TotalErrors,
		UptimeSeconds:     globalMetrics.UptimeSeconds,
		startTime:         globalMetrics.startTime,
	}

	// 複製 map 數據
	for k, v := range globalMetrics.RequestCount {
		metrics.RequestCount[k] = v
	}
	for k, v := range globalMetrics.RequestDuration {
		metrics.RequestDuration[k] = v
	}
	for k, v := range globalMetrics.ErrorCount {
		metrics.ErrorCount[k] = v
	}

	return metrics
}

// FormatPrometheusMetrics 將指標格式化為 Prometheus 格式
// 這個函數生成符合 Prometheus 規範的指標輸出
func FormatPrometheusMetrics() string {
	metrics := GetMetrics()

	var output string

	// 添加註釋和幫助信息
	output += "# HELP http_requests_total Total number of HTTP requests\n"
	output += "# TYPE http_requests_total counter\n"

	// 輸出每個端點的請求計數
	for endpoint, count := range metrics.RequestCount {
		output += `http_requests_total{endpoint="` + endpoint + `"} ` + strconv.FormatInt(count, 10) + "\n"
	}

	output += "\n# HELP http_request_duration_milliseconds Average HTTP request duration in milliseconds\n"
	output += "# TYPE http_request_duration_milliseconds gauge\n"

	// 輸出每個端點的平均響應時間
	for endpoint, duration := range metrics.RequestDuration {
		output += `http_request_duration_milliseconds{endpoint="` + endpoint + `"} ` + strconv.FormatFloat(duration, 'f', 2, 64) + "\n"
	}

	output += "\n# HELP http_errors_total Total number of HTTP errors\n"
	output += "# TYPE http_errors_total counter\n"

	// 輸出每個端點的錯誤計數
	for endpoint, count := range metrics.ErrorCount {
		output += `http_errors_total{endpoint="` + endpoint + `"} ` + strconv.FormatInt(count, 10) + "\n"
	}

	output += "\n# HELP http_active_connections Current number of active HTTP connections\n"
	output += "# TYPE http_active_connections gauge\n"
	output += "http_active_connections " + strconv.FormatInt(metrics.ActiveConnections, 10) + "\n"

	output += "\n# HELP http_total_requests Total number of HTTP requests since startup\n"
	output += "# TYPE http_total_requests counter\n"
	output += "http_total_requests " + strconv.FormatInt(metrics.TotalRequests, 10) + "\n"

	output += "\n# HELP http_total_errors Total number of HTTP errors since startup\n"
	output += "# TYPE http_total_errors counter\n"
	output += "http_total_errors " + strconv.FormatInt(metrics.TotalErrors, 10) + "\n"

	output += "\n# HELP application_uptime_seconds Application uptime in seconds\n"
	output += "# TYPE application_uptime_seconds gauge\n"
	output += "application_uptime_seconds " + strconv.FormatFloat(metrics.UptimeSeconds, 'f', 2, 64) + "\n"

	return output
}

// GetHealthMetrics 返回用於健康檢查的簡化指標
// 這個函數提供應用程式健康狀態的快速概覽
func GetHealthMetrics() map[string]interface{} {
	metrics := GetMetrics()

	// 計算錯誤率
	errorRate := 0.0
	if metrics.TotalRequests > 0 {
		errorRate = float64(metrics.TotalErrors) / float64(metrics.TotalRequests) * 100
	}

	return map[string]interface{}{
		"status":             "healthy",
		"uptime_seconds":     metrics.UptimeSeconds,
		"total_requests":     metrics.TotalRequests,
		"total_errors":       metrics.TotalErrors,
		"error_rate_percent": errorRate,
		"active_connections": metrics.ActiveConnections,
		"endpoints_tracked":  len(metrics.RequestCount),
	}
}
