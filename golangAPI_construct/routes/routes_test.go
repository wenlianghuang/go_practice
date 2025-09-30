package routes

import (
	"golangAPI_construct/handlers"
	"golangAPI_construct/middleware"

	"github.com/gin-gonic/gin"
)

// SetupVersionedRoutes 設置版本化路由
func SetupVersionedRoutes(r *gin.Engine, bookHandler *handlers.BookHandler, authHandler *handlers.AuthHandler) {
	// API v1 路由群組
	v1 := r.Group("/api/v1")
	{
		setupV1Routes(v1, bookHandler, authHandler)
	}

	// API v2 路由群組 (未來版本)
	v2 := r.Group("/api/v2")
	{
		setupV2Routes(v2, bookHandler, authHandler)
	}

	// 默認路由重定向到最新版本
	r.GET("/api/books", func(c *gin.Context) {
		c.Redirect(302, "/api/v1/books")
	})
}

// setupV1Routes 設置 v1 版本路由
func setupV1Routes(v1 *gin.RouterGroup, bookHandler *handlers.BookHandler, authHandler *handlers.AuthHandler) {
	// 認證路由
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	// 書籍路由
	books := v1.Group("/books")
	books.Use(middleware.AuthMiddleware()) // 需要認證
	{
		books.GET("", bookHandler.GetBooks)          // GET /api/v1/books
		books.POST("", bookHandler.CreateBook)       // POST /api/v1/books
		books.GET("/:id", bookHandler.GetBook)       // GET /api/v1/books/:id
		books.PUT("/:id", bookHandler.UpdateBook)    // PUT /api/v1/books/:id
		books.DELETE("/:id", bookHandler.DeleteBook) // DELETE /api/v1/books/:id
	}
}

// setupV2Routes 設置 v2 版本路由 (示例：新功能)
func setupV2Routes(v2 *gin.RouterGroup, bookHandler *handlers.BookHandler, authHandler *handlers.AuthHandler) {
	// 認證路由 (v2 可能有不同的認證機制)
	auth := v2.Group("/auth")
	{
		auth.POST("/login", authHandler.LoginV2)
		auth.POST("/register", authHandler.RegisterV2)
		auth.POST("/refresh", authHandler.RefreshToken) // v2 新增的刷新令牌功能
	}

	// 書籍路由 (v2 可能有額外功能)
	books := v2.Group("/books")
	books.Use(middleware.AuthMiddleware())
	{
		books.GET("", bookHandler.GetBooksV2)    // 可能支持更多查詢選項
		books.POST("", bookHandler.CreateBookV2) // 可能有不同的數據結構
		books.GET("/:id", bookHandler.GetBookV2)
		books.PUT("/:id", bookHandler.UpdateBookV2)
		books.DELETE("/:id", bookHandler.DeleteBookV2)
		books.GET("/:id/reviews", bookHandler.GetBookReviews)    // v2 新功能：書評
		books.POST("/:id/reviews", bookHandler.CreateBookReview) // v2 新功能
	}

	// v2 新增的統計功能
	stats := v2.Group("/stats")
	stats.Use(middleware.AuthMiddleware())
	{
		stats.GET("/books", bookHandler.GetBookStats)
		stats.GET("/popular", bookHandler.GetPopularBooks)
	}
}
