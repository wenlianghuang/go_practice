package routes

import (
	"context"
	"golangAPI_construct/data"
	"golangAPI_construct/handlers"
	"golangAPI_construct/logging"
	"golangAPI_construct/middleware"
	"golangAPI_construct/services"
	"os"

	"github.com/gin-gonic/gin"
)

// test git rebase.
func SetupRoutes() *gin.Engine {
	r := gin.New()

	// 設置全局中間件
	r.Use(
		middleware.RequestID(),
		middleware.ErrorHandler(),
		middleware.Logger(),
		middleware.CORS(),
		gin.Recovery(),
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
	r.GET("/api/health", bookHandler.HealthCheck)

	// API v1 路由組
	v1 := r.Group("/api/v1")
	{
		// 認證相關路由（不需要 JWT 驗證）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", middleware.JWTAuthMiddleware(), handlers.Logout)
			auth.POST("/refresh", middleware.JWTAuthMiddleware(), handlers.RefreshToken)
		}

		// 受保護的書籍路由（需要 JWT 驗證）
		books := v1.Group("/books", middleware.JWTAuthMiddleware())
		{
			books.GET("", bookHandler.GetBooks)
			books.POST("", bookHandler.CreateBook)
			books.GET("/:id", bookHandler.GetBookByID)
			books.PUT("/:id", bookHandler.UpdateBook)
			books.PATCH("/:id", bookHandler.PatchBook)
			books.DELETE("/:id", bookHandler.DeleteBook)
		}
	}

	return r
}
