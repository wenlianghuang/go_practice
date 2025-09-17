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
	// 決定使用記憶體會資料庫實作
	var bookService services.BookServiceInterface
	if os.Getenv("USE_DB") == "true" {
		db, err := data.Open()
		if err != nil {
			panic(err)
		}
		if err := data.Migrate(context.Background(), db); err != nil {
			panic(err)
		}
		// link to the book_db.go
		bookService = services.NewBookServiceDB(db)
	} else {
		// link to the book.go
		bookService = services.NewBookService()
		logging.Logger.Print("[BOOT] Book service: in-memory mode")
	}

	//bookService := services.NewBookService()
	bookHandler := handlers.NewBookHandler(bookService)

	// 啟動時印出目前資料筆數，方便確認來源是否為 DB
	//log.Printf("[BOOT] Books count at start: %d", bookService.GetBooksCount())
	logging.Logger.Printf("[BOOT] Books count at start: %d", len(bookService.GetAllBooks()))

	r.GET("/api/health", bookHandler.HealthCheck)
	// v1 group with JWT middleware => verify token for all /api/v1/books routes
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
