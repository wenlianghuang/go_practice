package routes

import (
	"golangAPI_construct/handlers"
	"golangAPI_construct/middleware"
	"golangAPI_construct/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.ErrorHandler(),
		middleware.Logger(),
		middleware.CORS(), // Assuming CORS middleware is defined elsewhere
		gin.Recovery(),
	)

	bookService := services.NewBookService()
	bookHandler := handlers.NewBookHandler(bookService)

	api := r.Group("/api")
	{
		api.GET("/health", bookHandler.HealthCheck)

		api.GET("/books", bookHandler.GetBooks)
		api.POST("/books", bookHandler.CreateBook)
		api.GET("/books/:id", bookHandler.GetBookByID)
		api.PUT("/books/:id", bookHandler.UpdateBook)
		api.PATCH("/books/:id", bookHandler.PatchBook)
		api.DELETE("/books/:id", bookHandler.DeleteBook)
	}

	return r
}
