package routes

import (
	"golangAPI_construct/handlers"
	"golangAPI_construct/middleware"
	"golangAPI_construct/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.New()
	// set up global middleware
	r.Use(
		middleware.RequestID(),
		middleware.ErrorHandler(),
		middleware.Logger(),
		middleware.CORS(), // Assuming CORS middleware is defined elsewhere
		gin.Recovery(),
	)

	bookService := services.NewBookService()
	bookHandler := handlers.NewBookHandler(bookService)

	/*
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

		v1 := r.Group("/api/v1")
		{
			v1.POST("/login", handlers.Login)
		}
	*/
	r.GET("/api/health", bookHandler.HealthCheck)
	// v1 group with JWT middleware
	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", handlers.Login)

		// Protected book routes
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
