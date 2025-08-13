package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"golangAPI_construct/models"
	"golangAPI_construct/responses"
	"golangAPI_construct/services"
)

var startTime = time.Now()

type BookHandler struct {
	service *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
	return &BookHandler{service: service}
}

// GetBooks supports optional ?author= filtering.
func (h *BookHandler) GetBooks(c *gin.Context) {
	author := c.Query("author")

	var result []models.Book
	if author != "" {
		result = h.service.GetBooksByAuthor(author)
	} else {
		result = h.service.GetAllBooks()
	}

	responses.Success(c, http.StatusOK, gin.H{
		"items": result,
		"count": len(result),
	})
}

// CreateBook creates a new book with basic validation.
func (h *BookHandler) CreateBook(c *gin.Context) {
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.Error(responses.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid request body"))
		return
	}
	if newBook.Title == "" || newBook.Author == "" || newBook.Price < 0 {
		c.Error(responses.NewAppError(http.StatusBadRequest, "INVALID_FIELDS", "title, author required; price must be >= 0"))
		return
	}

	book, err := h.service.CreateBook(newBook)
	if err != nil {
		c.Error(responses.NewAppError(http.StatusInternalServerError, "CREATE_FAILED", err.Error()))
		return
	}

	responses.Success(c, http.StatusCreated, book)
}

// GetBookByID returns a single book.
func (h *BookHandler) GetBookByID(c *gin.Context) {
	id := c.Param("id")
	book, err := h.service.GetBookByID(id)
	if err != nil {
		c.Error(responses.NewAppError(http.StatusNotFound, "NOT_FOUND", err.Error()))
		return
	}
	responses.Success(c, http.StatusOK, book)
}

// UpdateBook replaces an existing book.
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.Error(responses.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid request body"))
		return
	}
	if updatedBook.Title == "" || updatedBook.Author == "" || updatedBook.Price < 0 {
		c.Error(responses.NewAppError(http.StatusBadRequest, "INVALID_FIELDS", "title, author required; price must be >= 0"))
		return
	}

	book, err := h.service.UpdateBook(id, updatedBook)
	if err != nil {
		c.Error(responses.NewAppError(http.StatusNotFound, "NOT_FOUND", err.Error()))
		return
	}

	responses.Success(c, http.StatusOK, book)
}

// PatchBook partially updates fields of a book.
func (h *BookHandler) PatchBook(c *gin.Context) {
	id := c.Param("id")
	var patchData models.BookPatch
	if err := c.ShouldBindJSON(&patchData); err != nil {
		c.Error(responses.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid request body"))
		return
	}
	if patchData.Price != nil && *patchData.Price < 0 {
		c.Error(responses.NewAppError(http.StatusBadRequest, "INVALID_FIELDS", "price must be >= 0"))
		return
	}

	book, err := h.service.PatchBook(id, patchData)
	if err != nil {
		c.Error(responses.NewAppError(http.StatusNotFound, "NOT_FOUND", err.Error()))
		return
	}

	responses.Success(c, http.StatusOK, book)
}

// DeleteBook deletes a book.
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	book, err := h.service.DeleteBook(id)
	if err != nil {
		c.Error(responses.NewAppError(http.StatusNotFound, "NOT_FOUND", err.Error()))
		return
	}
	responses.Success(c, http.StatusOK, gin.H{"deleted": book})
}

// HealthCheck returns service health info.
func (h *BookHandler) HealthCheck(c *gin.Context) {
	uptime := time.Since(startTime).Seconds()
	responses.Success(c, http.StatusOK, gin.H{
		"status":       "healthy",
		"service":      "bookstore-api",
		"version":      "1.0.0",
		"books_count":  h.service.GetBooksCount(),
		"uptime_sec":   uptime,
		"current_time": time.Now().Format(time.RFC3339),
	})
}
