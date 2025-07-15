package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title" binding:"required"`
	Author string  `json:"author" binding:"required"`
	Price  float64 `json:"price" binding:"required,min=0"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

var books = []Book{
	{ID: "1", Title: "1984", Author: "George Orwell", Price: 9.99},
	{ID: "2", Title: "Brave New World", Author: "Aldous Huxley", Price: 8.99},
	{ID: "3", Title: "To Kill a Mockingbird", Author: "Harper Lee", Price: 12.99},
}

var nextID = 4 // 用於生成新的 ID

func main() {
	router := gin.Default()

	// 中間件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// API 路由群組
	api := router.Group("/api/v1")
	{
		api.GET("/books", getBooks)          // GET - 獲取所有書籍
		api.POST("/books", addBook)          // POST - 新增書籍
		api.GET("/books/:id", getBookByID)   // GET - 根據ID獲取書籍
		api.PUT("/books/:id", updateBook)    // PUT - 更新書籍
		api.PATCH("/books/:id", patchBook)   // PATCH - 部分更新書籍
		api.DELETE("/books/:id", deleteBook) // DELETE - 刪除書籍
	}

	// 健康檢查端點
	router.GET("/health", healthCheck)

	router.Run("localhost:8080")
}

// 中間件：處理 CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// GET /api/v1/books - 獲取所有書籍
func getBooks(c *gin.Context) {
	// 支援查詢參數過濾
	author := c.Query("author")
	if author != "" {
		var filteredBooks []Book
		for _, book := range books {
			if book.Author == author {
				filteredBooks = append(filteredBooks, book)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"data":  filteredBooks,
			"count": len(filteredBooks),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  books,
		"count": len(books),
	})
}

// POST /api/v1/books - 新增書籍
func addBook(c *gin.Context) {
	var newBook Book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// 檢查是否已存在相同 ID
	for _, book := range books {
		if book.ID == newBook.ID {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "duplicate_id",
				Message: "Book with this ID already exists",
			})
			return
		}
	}

	// 如果沒有提供 ID，自動生成
	if newBook.ID == "" {
		newBook.ID = strconv.Itoa(nextID)
		nextID++
	}

	books = append(books, newBook)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"data":    newBook,
	})
}

// GET /api/v1/books/:id - 根據ID獲取書籍
func getBookByID(c *gin.Context) {
	id := c.Param("id")

	for _, book := range books {
		if book.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"data": book,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{
		Error:   "not_found",
		Message: "Book not found",
	})
}

// PUT /api/v1/books/:id - 更新書籍
func updateBook(c *gin.Context) {
	id := c.Param("id")
	var updatedBook Book

	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	for i, book := range books {
		if book.ID == id {
			// 保持原有的 ID
			updatedBook.ID = id
			books[i] = updatedBook
			c.JSON(http.StatusOK, gin.H{
				"message": "Book updated successfully",
				"data":    updatedBook,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{
		Error:   "not_found",
		Message: "Book not found",
	})
}

// PATCH - 部分更新書籍 (址更新提供的欄位)
func patchBook(c *gin.Context) {
	id := c.Param("id")
	var patchData map[string]interface{}

	if err := c.ShouldBindJSON(&patchData); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	for i, book := range books {
		if book.ID == id {
			if title, ok := patchData["title"].(string); ok {
				book.Title = title
			}
			if author, ok := patchData["author"].(string); ok {
				book.Author = author
			}
			if price, ok := patchData["price"].(float64); ok {
				book.Price = price
			}
			books[i] = book
			c.JSON(http.StatusOK, gin.H{
				"message": "Book patched successfully",
				"data":    book,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{
		Error:   "not_found",
		Message: "Book not found",
	})
}

// DELETE /api/v1/books/:id - 刪除書籍
func deleteBook(c *gin.Context) {
	id := c.Param("id")

	for i, book := range books {
		if book.ID == id {
			// 刪除書籍
			books = append(books[:i], books[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "Book deleted successfully",
				"data":    book,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, ErrorResponse{
		Error:   "not_found",
		Message: "Book not found",
	})
}

// GET /health - 健康檢查
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "bookstore-api",
		"version":   "1.0.0",
		"timestamp": gin.H{"books_count": len(books)},
	})
}
