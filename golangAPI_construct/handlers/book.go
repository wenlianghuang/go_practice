package handlers

import (
	"golangAPI_construct/models"
	"golangAPI_construct/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		"data":  result,
		"count": len(result),
	})
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var newBook models.Book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	book, err := h.service.CreateBook(newBook)
	if err != nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error:   "duplicate_id",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Book created successfully",
		"data":    book,
	})
}

func (h *BookHandler) GetBookByID(c *gin.Context) {
	id := c.Param("id")

	book, err := h.service.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var updatedBook models.Book

	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	book, err := h.service.UpdateBook(id, updatedBook)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book updated successfully (PUT)",
		"data":    book,
	})
}

func (h *BookHandler) PatchBook(c *gin.Context) {
	id := c.Param("id")
	var patchData models.BookPatch

	if err := c.ShouldBindJSON(&patchData); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	book, err := h.service.PatchBook(id, patchData)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book updated successfully (PATCH)",
		"data":    book,
	})
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	book, err := h.service.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
		"data":    book,
	})
}

func (h *BookHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "bookstore-api",
		"version": "1.0.0",
		"data":    gin.H{"books_count": h.service.GetBooksCount()},
	})
}
