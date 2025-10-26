package middleware

import (
	"log"
	"net/http"
	"time"
)

// 日誌中間件
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 記錄請求開始
		log.Printf("🚀 Started %s %s", r.Method, r.URL.Path)

		// 調用下一個處理器
		next.ServeHTTP(w, r)

		// 記錄請求完成
		log.Printf("✅ Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// CORS 中間件
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 設置 CORS 頭
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 處理預檢請求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 請求驗證中間件
/**
 * ValidationMiddleware 驗證中間件
 *
 * 此中間件主要用來檢查當收到 POST 或 PUT 請求時，
 * Content-Type 是否設置為 "application/json"。
 * 若不是，將會以 400 Bad Request 回應錯誤訊息。
 *
 * @param next 下一個 http.Handler
 * @return http.Handler - 包裝過驗證邏輯的新 handler
 *
 * Example:
 *     router.Use(middleware.ValidationMiddleware)
 */
func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 檢查 Content-Type
		if r.Method == "POST" || r.Method == "PUT" {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// 錯誤處理中間件
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("❌ Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// 請求大小限制中間件
func RequestSizeMiddleware(maxSize int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxSize {
				http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
