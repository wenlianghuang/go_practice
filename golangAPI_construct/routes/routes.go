package routes

import (
	"golangAPI_construct/handlers"
	"golangAPI_construct/middleware"
	"golangAPI_construct/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// 中間件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// 初始化服務和處理器
	bookService := services.NewBookService()
	bookHandler := handlers.NewBookHandler(bookService)

	// API 路由群組
	api := router.Group("/api/v1")
	{
		api.GET("/books", bookHandler.GetBooks)
		api.POST("/books", bookHandler.CreateBook)
		api.GET("/books/:id", bookHandler.GetBookByID)
		api.PUT("/books/:id", bookHandler.UpdateBook)
		api.PATCH("/books/:id", bookHandler.PatchBook)
		api.DELETE("/books/:id", bookHandler.DeleteBook)
	}

	// 健康檢查端點
	router.GET("/health", bookHandler.HealthCheck)

	return router
}
