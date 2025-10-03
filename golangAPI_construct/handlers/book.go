package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"golangAPI_construct/models"
	"golangAPI_construct/responses"
	"golangAPI_construct/services"
)

var startTime = time.Now()

type BookHandler struct {
	//service *services.BookService
	service services.BookServiceInterface
}

func NewBookHandler(service services.BookServiceInterface) *BookHandler {
	return &BookHandler{service: service}
}

// GetBooks supports optional ?author= filtering.
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	var result []models.Book
	if author != "" {
		result = h.service.GetBooksByAuthor(author)
	} else {
		result = h.service.GetAllBooks()
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"items": result,
		"count": len(result),
	})
}

// CreateBook creates a new book with basic validation.
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid request body"))
		return
	}
	if newBook.Title == "" || newBook.Author == "" || newBook.Price < 0 {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_FIELDS", "title, author required; price must be >= 0"))
		return
	}

	book, err := h.service.CreateBook(newBook)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "CREATE_FAILED", err.Error()))
		return
	}

	responses.Success(w, r, http.StatusCreated, book)
}

// GetBookByID returns a single book.
func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	book, err := h.service.GetBookByID(id)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusNotFound, "NOT_FOUND", err.Error()))
		return
	}
	responses.Success(w, r, http.StatusOK, book)
}

// UpdateBook replaces an existing book.
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var updatedBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid request body"))
		return
	}
	if updatedBook.Title == "" || updatedBook.Author == "" || updatedBook.Price < 0 {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_FIELDS", "title, author required; price must be >= 0"))
		return
	}

	book, err := h.service.UpdateBook(id, updatedBook)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusNotFound, "NOT_FOUND", err.Error()))
		return
	}

	responses.Success(w, r, http.StatusOK, book)
}

// PatchBook partially updates fields of a book.
func (h *BookHandler) PatchBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var patchData models.BookPatch
	if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_JSON", "invalid request body"))
		return
	}
	if patchData.Price != nil && *patchData.Price < 0 {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_FIELDS", "price must be >= 0"))
		return
	}

	book, err := h.service.PatchBook(id, patchData)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusNotFound, "NOT_FOUND", err.Error()))
		return
	}

	responses.Success(w, r, http.StatusOK, book)
}

// DeleteBook deletes a book.
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	book, err := h.service.DeleteBook(id)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusNotFound, "NOT_FOUND", err.Error()))
		return
	}
	responses.Success(w, r, http.StatusOK, map[string]interface{}{"deleted": book})
}

// HealthCheck returns service health info.
func (h *BookHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime).Seconds()
	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"status":       "healthy",
		"service":      "bookstore-api",
		"version":      "1.0.0",
		"books_count":  h.service.GetBooksCount(),
		"uptime_sec":   uptime,
		"current_time": time.Now().Format(time.RFC3339),
	})
}
