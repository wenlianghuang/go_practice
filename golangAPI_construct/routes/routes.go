package routes

import (
	"golangAPI_construct/handlers"
	"golangAPI_construct/middleware"
	"golangAPI_construct/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS())

	bookService := services.NewBookService()
	bookHandler := handlers.NewBookHandler(bookService)

	api := router.Group("/api")
	{
		api.GET("/health", bookHandler.HealthCheck)

		api.GET("/books", bookHandler.GetBooks)
		api.POST("/books", bookHandler.CreateBook)
		api.GET("/books/:id", bookHandler.GetBookByID)
		api.PUT("/books/:id", bookHandler.UpdateBook)
		api.PATCH("/books/:id", bookHandler.PatchBook)
		api.DELETE("/books/:id", bookHandler.DeleteBook)
	}

	return router
}
