package routes

import (
	"golangAPI_construct/cache"
	"golangAPI_construct/data"
	"golangAPI_construct/handlers"
	"golangAPI_construct/logging"
	"golangAPI_construct/middleware"
	"golangAPI_construct/services"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// test git rebase.
func SetupRoutes() http.Handler {
	r := chi.NewRouter()

	// 設置全局中間件
	r.Use(
		middleware.RequestID(),
		middleware.ErrorHandler(), // 包含 panic recovery 和統一錯誤處理
		middleware.Logger(),
		middleware.MetricsMiddleware(), // 添加 metrics 收集中間件
		middleware.CORS(),
		// middleware.Recoverer(), // 移除：與 ErrorHandler 重複
	)

	// 決定使用記憶體或 GORM 實作（已移除傳統 SQL 選項）
	var bookService services.BookServiceInterface
	if os.Getenv("USE_DB") == "true" && os.Getenv("USE_GORM") == "true" {
		// 使用 GORM（推薦的數據庫模式）
		gormDB, err := data.OpenGORM()
		if err != nil {
			panic(err)
		}
		if err := data.MigrateGORM(gormDB); err != nil {
			panic(err)
		}
		if err := data.SeedGORM(gormDB); err != nil {
			logging.Logger.Printf("[GORM] Warning: failed to seed database: %v", err)
		}
		if err := data.CreateGORMIndexes(gormDB); err != nil {
			logging.Logger.Printf("[GORM] Warning: failed to create indexes: %v", err)
		}
		bookService = services.NewBookServiceGORM(gormDB)
		logging.Logger.Print("[BOOT] Book service: GORM database mode (enhanced search capabilities)")
	} else {
		// 使用內存模式（開發/測試用）
		bookService = services.NewBookService()
		logging.Logger.Print("[BOOT] Book service: in-memory mode")
		logging.Logger.Print("[BOOT] Note: Set USE_DB=true and USE_GORM=true to enable database features")
	}

	// 初始化緩存服務
	cacheService, err := cache.NewCacheService()
	if err != nil {
		logging.Logger.Printf("[CACHE] Failed to initialize cache service: %v", err)
		// 繼續運行，但不使用緩存
		cacheService = nil
	} else {
		logging.Logger.Print("[CACHE] Cache service initialized successfully")
	}

	bookHandler := handlers.NewBookHandler(bookService)
	metricsHandler := handlers.NewMetricsHandler()

	// 如果使用 GORM，創建 GORM 專用處理器
	var gormHandler *handlers.GORMHandler
	if os.Getenv("USE_GORM") == "true" {
		if gormService, ok := bookService.(*services.BookServiceGORM); ok {
			gormHandler = handlers.NewGORMHandler(gormService)
		}
	}

	// 啟動時印出目前資料筆數
	logging.Logger.Printf("[BOOT] Books count at start: %d", len(bookService.GetAllBooks()))

	// 健康檢查端點（不需要認證）
	r.Get("/api/health", bookHandler.HealthCheck)

	// Metrics 端點（不需要認證，用於監控）
	// 這些端點提供應用程式的性能指標和健康狀態
	r.Get("/metrics", metricsHandler.PrometheusMetrics)            // Prometheus 格式指標
	r.Get("/api/metrics/health", metricsHandler.HealthMetrics)     // JSON 格式健康指標
	r.Get("/api/metrics/detailed", metricsHandler.DetailedMetrics) // JSON 格式詳細指標

	// API v1 路由組
	r.Route("/api/v1", func(r chi.Router) {
		// 認證相關路由（不需要 JWT 驗證）
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", handlers.Login)
			r.With(middleware.JWTAuthMiddleware()).Post("/logout", handlers.Logout)
			r.With(middleware.JWTAuthMiddleware()).Post("/refresh", handlers.RefreshToken)
		})

		// 受保護的書籍路由（需要 JWT 驗證）
		r.Route("/books", func(r chi.Router) {
			// 所有書籍相關路由都需要 JWT 認證
			r.Use(middleware.JWTAuthMiddleware())

			// 如果緩存服務可用，為 GET 路由添加緩存
			if cacheService != nil {
				// GET 路由：使用緩存提高性能
				r.With(middleware.CacheMiddleware(cacheService, middleware.DefaultCacheConfig)).Get("/", bookHandler.GetBooks)
				r.With(middleware.CacheMiddleware(cacheService, middleware.DefaultCacheConfig)).Get("/{id}", bookHandler.GetBookByID)

				// 修改操作：添加緩存失效中間件
				r.With(
					middleware.CacheInvalidationMiddleware(cacheService, []string{"books"}),
					middleware.RequestValidator(middleware.BookValidationRules),
				).Post("/", bookHandler.CreateBook)

				r.With(
					middleware.CacheInvalidationMiddleware(cacheService, []string{"books"}),
					middleware.RequestValidator(middleware.BookValidationRules),
				).Put("/{id}", bookHandler.UpdateBook)

				r.With(
					middleware.CacheInvalidationMiddleware(cacheService, []string{"books"}),
					middleware.RequestValidator(middleware.BookValidationRules),
				).Patch("/{id}", bookHandler.PatchBook)

				r.With(middleware.CacheInvalidationMiddleware(cacheService, []string{"books"})).Delete("/{id}", bookHandler.DeleteBook)
			} else {
				// 沒有緩存服務時的原始路由
				r.Get("/", bookHandler.GetBooks)
				r.Get("/{id}", bookHandler.GetBookByID)
				r.With(middleware.RequestValidator(middleware.BookValidationRules)).Post("/", bookHandler.CreateBook)
				r.With(middleware.RequestValidator(middleware.BookValidationRules)).Put("/{id}", bookHandler.UpdateBook)
				r.With(middleware.RequestValidator(middleware.BookValidationRules)).Patch("/{id}", bookHandler.PatchBook)
				r.Delete("/{id}", bookHandler.DeleteBook)
			}
		})

		// GORM 專用路由（需要認證）
		if gormHandler != nil {
			r.Route("/gorm", func(r chi.Router) {
				r.Use(middleware.JWTAuthMiddleware())

				// 搜索和統計功能（增強版）
				r.Get("/search", gormHandler.SearchBooks)                  // 增強版搜索，支持多條件篩選
				r.Get("/search-advanced", gormHandler.SearchBooksAdvanced) // 高級搜索，支持複雜查詢
				r.Get("/statistics", gormHandler.GetBookStatistics)
				r.Get("/author-statistics", gormHandler.GetAuthorStatistics)
				r.Get("/database-health", gormHandler.GetDatabaseHealth)

				// 分類和篩選功能
				r.Get("/category/{category}", gormHandler.GetBooksByCategory)
				r.Get("/price-range", gormHandler.GetBooksByPriceRange)
				r.Get("/published/{year}", gormHandler.GetBooksByPublishedYear)
				r.Get("/top-rated", gormHandler.GetTopRatedBooks)
				r.Get("/recent", gormHandler.GetRecentBooks)
				r.Get("/with-reviews", gormHandler.GetBooksWithReviews)

				// 分頁和批量操作
				r.Get("/paginated", gormHandler.GetBooksWithPagination)
				r.Get("/by-authors", gormHandler.GetBooksByMultipleAuthors)
			})
		}
	})

	return r
}
