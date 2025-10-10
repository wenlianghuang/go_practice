package routes

import (
	"context"
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
		middleware.CORS(),
		// middleware.Recoverer(), // 移除：與 ErrorHandler 重複
	)

	// 決定使用記憶體或資料庫實作
	var bookService services.BookServiceInterface
	if os.Getenv("USE_DB") == "true" {
		db, err := data.Open()
		if err != nil {
			panic(err)
		}
		if err := data.Migrate(context.Background(), db); err != nil {
			panic(err)
		}
		bookService = services.NewBookServiceDB(db)
		logging.Logger.Print("[BOOT] Book service: database mode")
	} else {
		bookService = services.NewBookService()
		logging.Logger.Print("[BOOT] Book service: in-memory mode")
	}

	bookHandler := handlers.NewBookHandler(bookService)

	// 啟動時印出目前資料筆數
	logging.Logger.Printf("[BOOT] Books count at start: %d", len(bookService.GetAllBooks()))

	// 健康檢查端點（不需要認證）
	r.Get("/api/health", bookHandler.HealthCheck)

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

			// GET 路由不需要驗證請求體，所以不使用 RequestValidator
			r.Get("/", bookHandler.GetBooks)
			r.Get("/{id}", bookHandler.GetBookByID)

			// POST 路由：創建新書籍，需要驗證請求體數據
			// 使用 BookValidationRules 確保 title, author, price 欄位符合要求
			r.With(middleware.RequestValidator(middleware.BookValidationRules)).Post("/", bookHandler.CreateBook)

			// PUT 路由：完整更新書籍，需要驗證請求體數據
			r.With(middleware.RequestValidator(middleware.BookValidationRules)).Put("/{id}", bookHandler.UpdateBook)

			// PATCH 路由：部分更新書籍，需要驗證請求體數據
			r.With(middleware.RequestValidator(middleware.BookValidationRules)).Patch("/{id}", bookHandler.PatchBook)

			// DELETE 路由不需要驗證請求體
			r.Delete("/{id}", bookHandler.DeleteBook)
		})
	})

	return r
}
