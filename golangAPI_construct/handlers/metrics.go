package handlers

import (
	"net/http"

	"golangAPI_construct/middleware"
	"golangAPI_construct/responses"
)

// MetricsHandler 處理 metrics 相關的 HTTP 請求
// 這個處理器提供兩種格式的指標數據：
// 1. Prometheus 格式（用於監控系統）
// 2. JSON 格式（用於健康檢查和調試）
type MetricsHandler struct{}

// NewMetricsHandler 創建新的 MetricsHandler 實例
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// PrometheusMetrics 返回 Prometheus 格式的指標
// 這個端點通常被 Prometheus 服務器調用來收集指標數據
// 格式符合 Prometheus 的 text-based exposition format
func (h *MetricsHandler) PrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	// 設置正確的 Content-Type
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	// 獲取 Prometheus 格式的指標數據
	metricsData := middleware.FormatPrometheusMetrics()

	// 直接寫入響應
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metricsData))
}

// HealthMetrics 返回 JSON 格式的健康指標
// 這個端點提供應用程式健康狀態的快速概覽
// 適合用於負載均衡器的健康檢查
func (h *MetricsHandler) HealthMetrics(w http.ResponseWriter, r *http.Request) {
	// 獲取健康指標數據
	healthData := middleware.GetHealthMetrics()

	// 使用統一的成功響應格式
	responses.Success(w, r, http.StatusOK, healthData)
}

// DetailedMetrics 返回詳細的 JSON 格式指標
// 這個端點提供完整的指標數據，包括每個端點的詳細統計
// 適合用於調試和詳細分析
func (h *MetricsHandler) DetailedMetrics(w http.ResponseWriter, r *http.Request) {
	// 獲取完整的指標數據
	metrics := middleware.GetMetrics()

	// 構建詳細的響應數據
	detailedData := map[string]interface{}{
		"summary": map[string]interface{}{
			"uptime_seconds":     metrics.UptimeSeconds,
			"total_requests":     metrics.TotalRequests,
			"total_errors":       metrics.TotalErrors,
			"active_connections": metrics.ActiveConnections,
			"endpoints_tracked":  len(metrics.RequestCount),
		},
		"endpoints": map[string]interface{}{
			"request_counts":    metrics.RequestCount,
			"request_durations": metrics.RequestDuration,
			"error_counts":      metrics.ErrorCount,
		},
	}

	// 使用統一的成功響應格式
	responses.Success(w, r, http.StatusOK, detailedData)
}
