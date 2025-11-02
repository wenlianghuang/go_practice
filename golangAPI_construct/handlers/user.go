package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golangAPI_construct/responses"
	"golangAPI_construct/services"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	users    *services.UserService
	bookGorm *services.BookServiceGORM
}

func NewUserHandler(users *services.UserService, bookGorm *services.BookServiceGORM) *UserHandler {
	return &UserHandler{
		users:    users,
		bookGorm: bookGorm,
	}
}

type favoriteRequest struct {
	BookID uint `json:"book_id"`
}

func (h *UserHandler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID"))
		return
	}

	books, err := h.bookGorm.GetBooksByUserFavorites(uint(userID))
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "FAVORITES_FETCH_FAILED", "Failed to fetch user favorites"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"count":   len(books),
		"books":   books,
	})
}

func (h *UserHandler) AddUserFavorite(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID"))
		return
	}

	var req favoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body"))
		return
	}
	if req.BookID == 0 {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "MISSING_BOOK_ID", "book_id is required"))
		return
	}

	if err := h.bookGorm.AddBookToUserFavorites(uint(userID), req.BookID); err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "FAVORITE_ADD_FAILED", "Failed to add favorite"))
		return
	}

	responses.Success(w, r, http.StatusCreated, map[string]interface{}{
		"user_id": userID,
		"book_id": req.BookID,
		"status":  "favorite added",
	})
}

func (h *UserHandler) RemoveUserFavorite(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	bookIDStr := chi.URLParam(r, "bookID")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID"))
		return
	}
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil || bookID == 0 {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_BOOK_ID", "Invalid book ID"))
		return
	}

	if err := h.bookGorm.RemoveBookFromUserFavorites(uint(userID), uint(bookID)); err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "FAVORITE_REMOVE_FAILED", "Failed to remove favorite"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"book_id": bookID,
		"status":  "favorite removed",
	})
}

func (h *UserHandler) GetUsersByBook(w http.ResponseWriter, r *http.Request) {
	bookIDStr := chi.URLParam(r, "bookID")
	if bookIDStr == "" {
		bookIDStr = chi.URLParam(r, "id")
	}
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil || bookID == 0 {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_BOOK_ID", "Invalid book ID"))
		return
	}

	users, err := h.bookGorm.GetUsersByBook(uint(bookID))
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "BOOK_USERS_FETCH_FAILED", "Failed to fetch users for book"))
		return
	}

	responses.Success(w, r, http.StatusOK, map[string]interface{}{
		"book_id": bookID,
		"count":   len(users),
		"users":   users,
	})
}
