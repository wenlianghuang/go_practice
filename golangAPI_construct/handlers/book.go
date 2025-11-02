package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"golangAPI_construct/middleware"
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

func hasPrivilegedRole(roles []string) bool {
	for _, role := range roles {
		switch strings.TrimSpace(strings.ToLower(role)) {
		case "admin", "editor":
			return true
		}
	}
	return false
}

// GetBooks supports optional ?author= filtering.
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	var userID uint
	if userIDStr, _ := r.Context().Value("user_id").(string); userIDStr != "" {
		if parsed, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			userID = uint(parsed)
		}
	}
	roles, _ := r.Context().Value("roles").([]string)

	result, err := h.service.GetBooksForUser(userID, roles)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "BOOKS_FETCH_FAILED", "Failed to fetch books"))
		return
	}

	if author != "" {
		filtered := make([]models.Book, 0, len(result))
		authorLower := strings.ToLower(author)
		for _, book := range result {
			if strings.Contains(strings.ToLower(book.Author), authorLower) {
				filtered = append(filtered, book)
			}
		}
		result = filtered
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"items": result,
		"count": len(result),
	})
}

// CreateBook creates a new book with validation middleware.
// 這個函數現在依賴於 RequestValidator 中間件來驗證請求數據
// 驗證後的數據會存儲在 request context 中，我們可以直接使用
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var userID uint
	if userIDStr, _ := r.Context().Value("user_id").(string); userIDStr != "" {
		if parsed, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			userID = uint(parsed)
		}
	}
	roles, _ := r.Context().Value("roles").([]string)

	// 從 context 中獲取已驗證的數據
	// RequestValidator 中間件會將驗證後的 JSON 數據存儲在這裡
	// 🔴 斷點 4：檢查從中間件獲取的驗證數據
	validatedData, ok := r.Context().Value(middleware.ValidatedDataKey).(map[string]interface{})
	if !ok {
		// 如果沒有找到驗證後的數據，說明中間件沒有正確執行
		// 這通常不應該發生，但為了安全起見我們還是要檢查
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "VALIDATION_ERROR", "Failed to get validated data"))
		return
	}

	// 從驗證後的數據創建書籍對象
	// 由於數據已經通過驗證，我們可以安全地進行類型斷言
	// 驗證規則確保了這些欄位存在且類型正確
	newBook := models.Book{
		Title:  validatedData["title"].(string),  // 已驗證：非空字串，長度 1-200
		Author: validatedData["author"].(string), // 已驗證：非空字串，長度 1-100
		Price:  validatedData["price"].(float64), // 已驗證：數字，範圍 0-10000
	}

	// 處理可選欄位
	if isbn, ok := validatedData["isbn"].(string); ok && isbn != "" {
		newBook.ISBN = isbn
	}
	if category, ok := validatedData["category"].(string); ok && category != "" {
		newBook.Category = category
	}
	if publishedStr, ok := validatedData["published"].(string); ok && publishedStr != "" {
		if published, err := time.Parse(time.RFC3339, publishedStr); err == nil {
			newBook.Published = &published
		}
	}

	// 調用服務層創建書籍
	createdBook, err := h.service.CreateBook(newBook)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "CREATE_FAILED", err.Error()))
		return
	}

	if createdBook != nil && !hasPrivilegedRole(roles) && userID > 0 {
		if bookID, err := strconv.ParseUint(createdBook.ID, 10, 64); err == nil {
			if svc, ok := h.service.(*services.BookServiceGORM); ok {
				if err := svc.AddBookToUserFavorites(userID, uint(bookID)); err != nil {
					responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "FAVORITE_LINK_FAILED", "Failed to link book to user"))
					return
				}
			}
		}
	}

	// 返回成功響應
	responses.Success(w, r, http.StatusCreated, createdBook)
}

// GetBookByID returns a single book.
func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var userID uint
	if userIDStr, _ := r.Context().Value("user_id").(string); userIDStr != "" {
		if parsed, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			userID = uint(parsed)
		}
	}
	roles, _ := r.Context().Value("roles").([]string)

	book, err := h.service.GetBookForUser(id, userID, roles)
	if err != nil {
		if err == services.ErrBookNotFound {
			responses.Fail(w, r, responses.NewAppError(http.StatusNotFound, "NOT_FOUND", err.Error()))
			return
		}
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "BOOK_FETCH_FAILED", "Failed to fetch book"))
		return
	}
	responses.Success(w, r, http.StatusOK, book)
}

// UpdateBook replaces an existing book.
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var userID uint
	if userIDStr, _ := r.Context().Value("user_id").(string); userIDStr != "" {
		if parsed, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			userID = uint(parsed)
		}
	}
	roles, _ := r.Context().Value("roles").([]string)

	if _, err := h.service.GetBookForUser(id, userID, roles); err != nil {
		if err == services.ErrBookNotFound {
			responses.Fail(w, r, responses.NewAppError(http.StatusNotFound, "BOOK_NOT_FOUND", "Book not found"))
			return
		}
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "BOOK_FETCH_FAILED", "Failed to verify book ownership"))
		return
	}

	validatedData, ok := r.Context().Value(middleware.ValidatedDataKey).(map[string]interface{})
	if !ok {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "VALIDATION_ERROR", "Failed to get validated data"))
		return
	}

	// 創建完整的書籍對象
	book := models.Book{
		Title:  validatedData["title"].(string),
		Author: validatedData["author"].(string),
		Price:  validatedData["price"].(float64),
	}

	// 處理可選欄位
	if isbn, ok := validatedData["isbn"].(string); ok && isbn != "" {
		book.ISBN = isbn
	}
	if category, ok := validatedData["category"].(string); ok && category != "" {
		book.Category = category
	}
	if publishedStr, ok := validatedData["published"].(string); ok && publishedStr != "" {
		if published, err := time.Parse(time.RFC3339, publishedStr); err == nil {
			book.Published = &published
		}
	}

	updatedBook, err := h.service.UpdateBook(id, book)
	if err != nil {
		if err == services.ErrBookNotFound {
			responses.Fail(w, r, responses.NewAppError(http.StatusNotFound, "BOOK_NOT_FOUND", "Book not found"))
			return
		}
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "UPDATE_FAILED", err.Error()))
		return
	}

	responses.Success(w, r, http.StatusOK, updatedBook)
}

// PatchBook partially updates fields of a book.
func (h *BookHandler) PatchBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var userID uint
	if userIDStr, _ := r.Context().Value("user_id").(string); userIDStr != "" {
		if parsed, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			userID = uint(parsed)
		}
	}
	roles, _ := r.Context().Value("roles").([]string)

	if _, err := h.service.GetBookForUser(id, userID, roles); err != nil {
		if err == services.ErrBookNotFound {
			responses.Fail(w, r, responses.NewAppError(http.StatusNotFound, "BOOK_NOT_FOUND", "Book not found"))
			return
		}
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "BOOK_FETCH_FAILED", "Failed to verify book ownership"))
		return
	}

	validatedData, ok := r.Context().Value(middleware.ValidatedDataKey).(map[string]interface{})
	if !ok {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "VALIDATION_ERROR", "Failed to get validated data"))
		return
	}

	// 創建部分更新對象
	patch := models.BookPatch{}

	if title, ok := validatedData["title"].(string); ok {
		patch.Title = &title
	}
	if author, ok := validatedData["author"].(string); ok {
		patch.Author = &author
	}
	if price, ok := validatedData["price"].(float64); ok {
		patch.Price = &price
	}
	if isbn, ok := validatedData["isbn"].(string); ok {
		patch.ISBN = &isbn
	}
	if category, ok := validatedData["category"].(string); ok {
		patch.Category = &category
	}
	if publishedStr, ok := validatedData["published"].(string); ok && publishedStr != "" {
		if published, err := time.Parse(time.RFC3339, publishedStr); err == nil {
			patch.Published = &published
		}
	}

	updatedBook, err := h.service.PatchBook(id, patch)
	if err != nil {
		if err == services.ErrBookNotFound {
			responses.Fail(w, r, responses.NewAppError(http.StatusNotFound, "BOOK_NOT_FOUND", "Book not found"))
			return
		}
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "PATCH_FAILED", err.Error()))
		return
	}

	responses.Success(w, r, http.StatusOK, updatedBook)
}

// DeleteBook deletes a book.
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var userID uint
	if userIDStr, _ := r.Context().Value("user_id").(string); userIDStr != "" {
		if parsed, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			userID = uint(parsed)
		}
	}
	roles, _ := r.Context().Value("roles").([]string)

	if _, err := h.service.GetBookForUser(id, userID, roles); err != nil {
		if err == services.ErrBookNotFound {
			responses.Fail(w, r, responses.NewAppError(http.StatusNotFound, "NOT_FOUND", "Book not found"))
			return
		}
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "BOOK_FETCH_FAILED", "Failed to verify book ownership"))
		return
	}

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
