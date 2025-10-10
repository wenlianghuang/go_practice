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
			r.Use(middleware.JWTAuthMiddleware())
			r.Get("/", bookHandler.GetBooks)
			r.With(middleware.RequestValidator(middleware.BookValidationRules)).Post("/", bookHandler.CreateBook)
			r.Get("/{id}", bookHandler.GetBookByID)
			r.With(middleware.RequestValidator(middleware.BookValidationRules)).Put("/{id}", bookHandler.UpdateBook)
			r.With(middleware.RequestValidator(middleware.BookValidationRules)).Patch("/{id}", bookHandler.PatchBook)
			r.Delete("/{id}", bookHandler.DeleteBook)
		})
	})

	return r
}
