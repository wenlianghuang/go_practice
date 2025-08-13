package handlers

import (
	"net/http"
	"time"

	"golangAPI_construct/models"
	"golangAPI_construct/services"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

type BookHandler struct {
	service *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	author := c.Query("author")

	var result []models.Book
	if author != "" {
		result = h.service.GetBooksByAuthor(author)
	} else {
		result = h.service.GetAllBooks()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"count":   len(result),
	})
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	if newBook.Title == "" || newBook.Author == "" || newBook.Price < 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_fields", Message: "title, author required; price must be >= 0"})
		return
	}

	book, err := h.service.CreateBook(newBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "create_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    book,
	})
}

func (h *BookHandler) GetBookByID(c *gin.Context) {
	id := c.Param("id")
	book, err := h.service.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    book,
	})
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	if updatedBook.Title == "" || updatedBook.Author == "" || updatedBook.Price < 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_fields", Message: "title, author required; price must be >= 0"})
		return
	}

	book, err := h.service.UpdateBook(id, updatedBook)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    book,
	})
}

func (h *BookHandler) PatchBook(c *gin.Context) {
	id := c.Param("id")
	var patchData models.BookPatch
	if err := c.ShouldBindJSON(&patchData); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}
	if patchData.Price != nil && *patchData.Price < 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid_fields", Message: "price must be >= 0"})
		return
	}

	book, err := h.service.PatchBook(id, patchData)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    book,
	})
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	book, err := h.service.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "not_found", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    book,
	})
}

func (h *BookHandler) HealthCheck(c *gin.Context) {
	uptime := time.Since(startTime).Seconds()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"status":       "healthy",
			"service":      "bookstore-api",
			"version":      "1.0.0",
			"books_count":  h.service.GetBooksCount(),
			"uptime_sec":   uptime,
			"current_time": time.Now().Format(time.RFC3339),
		},
	})
}
