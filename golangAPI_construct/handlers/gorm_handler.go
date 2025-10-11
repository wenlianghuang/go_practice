package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"golangAPI_construct/responses"
	"golangAPI_construct/services"
)

// GORMHandler GORM 專用處理器
// 提供 GORM 的高級功能，如搜索、統計、關聯查詢等
type GORMHandler struct {
	gormService *services.BookServiceGORM
}

// NewGORMHandler 創建 GORM 處理器實例
func NewGORMHandler(gormService *services.BookServiceGORM) *GORMHandler {
	return &GORMHandler{gormService: gormService}
}

// SearchBooks 搜索書籍
func (h *GORMHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "MISSING_QUERY", "Search query parameter 'q' is required"))
		return
	}

	results, err := h.gormService.SearchBooks(query)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "SEARCH_FAILED", "Failed to search books"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"query":   query,
		"results": results,
		"count":   len(results),
	})
}

// GetBookStatistics 獲取書籍統計信息
func (h *GORMHandler) GetBookStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.gormService.GetBookStatistics()
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "STATS_FAILED", "Failed to get book statistics"))
		return
	}

	responses.Success(w, r, http.StatusOK, stats)
}

// GetAuthorStatistics 獲取作者統計信息
func (h *GORMHandler) GetAuthorStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.gormService.GetAuthorStatistics()
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "AUTHOR_STATS_FAILED", "Failed to get author statistics"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"authors": stats,
		"count":   len(stats),
	})
}

// GetBooksByCategory 根據分類獲取書籍
func (h *GORMHandler) GetBooksByCategory(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")
	if category == "" {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "MISSING_CATEGORY", "Category parameter is required"))
		return
	}

	books, err := h.gormService.GetBooksByCategory(category)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "CATEGORY_FAILED", "Failed to get books by category"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"category": category,
		"books":    books,
		"count":    len(books),
	})
}

// GetBooksByPriceRange 根據價格範圍獲取書籍
func (h *GORMHandler) GetBooksByPriceRange(w http.ResponseWriter, r *http.Request) {
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")

	if minPriceStr == "" || maxPriceStr == "" {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "MISSING_PRICE_RANGE", "Both min_price and max_price parameters are required"))
		return
	}

	minPrice, err := strconv.ParseFloat(minPriceStr, 64)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_MIN_PRICE", "Invalid min_price parameter"))
		return
	}

	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_MAX_PRICE", "Invalid max_price parameter"))
		return
	}

	books, err := h.gormService.GetBooksByPriceRange(minPrice, maxPrice)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "PRICE_RANGE_FAILED", "Failed to get books by price range"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"min_price": minPrice,
		"max_price": maxPrice,
		"books":     books,
		"count":     len(books),
	})
}

// GetBooksByPublishedYear 根據出版年份獲取書籍
func (h *GORMHandler) GetBooksByPublishedYear(w http.ResponseWriter, r *http.Request) {
	yearStr := chi.URLParam(r, "year")
	if yearStr == "" {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "MISSING_YEAR", "Year parameter is required"))
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_YEAR", "Invalid year parameter"))
		return
	}

	books, err := h.gormService.GetBooksByPublishedDate(year)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "YEAR_FAILED", "Failed to get books by published year"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"year":  year,
		"books": books,
		"count": len(books),
	})
}

// GetTopRatedBooks 獲取評分最高的書籍
func (h *GORMHandler) GetTopRatedBooks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // 預設限制

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_LIMIT", "Invalid limit parameter"))
			return
		}
	}

	books, err := h.gormService.GetTopRatedBooks(limit)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "TOP_RATED_FAILED", "Failed to get top rated books"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"books": books,
		"count": len(books),
		"limit": limit,
	})
}

// GetRecentBooks 獲取最近添加的書籍
func (h *GORMHandler) GetRecentBooks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // 預設限制

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_LIMIT", "Invalid limit parameter"))
			return
		}
	}

	books, err := h.gormService.GetRecentBooks(limit)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "RECENT_FAILED", "Failed to get recent books"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"books": books,
		"count": len(books),
		"limit": limit,
	})
}

// GetBooksWithPagination 分頁獲取書籍
func (h *GORMHandler) GetBooksWithPagination(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	pageSize := 10

	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_PAGE", "Invalid page parameter"))
			return
		}
	}

	if pageSizeStr != "" {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 || pageSize > 100 {
			responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_PAGE_SIZE", "Invalid page_size parameter (must be 1-100)"))
			return
		}
	}

	books, total, err := h.gormService.GetBooksWithPagination(page, pageSize)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "PAGINATION_FAILED", "Failed to get paginated books"))
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"books":       books,
		"count":       len(books),
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
		"has_next":    int64(page) < totalPages,
		"has_prev":    page > 1,
	})
}

// GetBooksByMultipleAuthors 根據多個作者獲取書籍
func (h *GORMHandler) GetBooksByMultipleAuthors(w http.ResponseWriter, r *http.Request) {
	authorsParam := r.URL.Query().Get("authors")
	if authorsParam == "" {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "MISSING_AUTHORS", "Authors parameter is required (comma-separated)"))
		return
	}

	// 解析作者列表
	authors := []string{}
	for _, author := range splitCommaSeparated(authorsParam) {
		if author != "" {
			authors = append(authors, author)
		}
	}

	if len(authors) == 0 {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "EMPTY_AUTHORS", "No valid authors provided"))
		return
	}

	books, err := h.gormService.GetBooksByMultipleAuthors(authors)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "MULTIPLE_AUTHORS_FAILED", "Failed to get books by multiple authors"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"authors": authors,
		"books":   books,
		"count":   len(books),
	})
}

// GetBooksWithReviews 獲取帶有評論的書籍
func (h *GORMHandler) GetBooksWithReviews(w http.ResponseWriter, r *http.Request) {
	books, err := h.gormService.GetBooksWithReviews()
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "REVIEWS_FAILED", "Failed to get books with reviews"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"books": books,
		"count": len(books),
	})
}

// GetDatabaseHealth 獲取數據庫健康狀態
func (h *GORMHandler) GetDatabaseHealth(w http.ResponseWriter, r *http.Request) {
	err := h.gormService.GetDatabaseHealth()
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusServiceUnavailable, "DB_UNHEALTHY", "Database is not healthy"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"message": "Database connection is working properly",
	})
}

// splitCommaSeparated 分割逗號分隔的字串
func splitCommaSeparated(s string) []string {
	result := []string{}
	for _, part := range splitString(s, ",") {
		if trimmed := trimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// splitString 分割字串的輔助函數
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	result := []string{}
	start := 0

	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}

	result = append(result, s[start:])
	return result
}

// trimSpace 去除字串兩端空白的輔助函數
func trimSpace(s string) string {
	start := 0
	end := len(s)

	// 去除開頭空白
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	// 去除結尾空白
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
