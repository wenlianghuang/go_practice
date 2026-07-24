package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golangAPI_construct/models"
	"golangAPI_construct/responses"
	"golangAPI_construct/security"
	"golangAPI_construct/services"

	"gorm.io/gorm"
)

const defaultTokenTTL = 2 * time.Hour

type AuthHandler struct {
	users      *services.UserService
	tokenStore security.TokenStore
}

func NewAuthHandler(users *services.UserService, tokenStore security.TokenStore) *AuthHandler {
	return &AuthHandler{
		users:      users,
		tokenStore: tokenStore,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string   `json:"token"`
	ExpiresAt int64    `json:"expires_at"`
	User      string   `json:"user"`
	Roles     []string `json:"roles"`
}

type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RegisterResponse struct {
	ID        uint     `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	CreatedAt int64    `json:"created_at"`
}

type UserDetailsResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Roles     []string   `json:"roles"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	IsActive  bool       `json:"is_active"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt int64      `json:"created_at"`
	UpdatedAt int64      `json:"updated_at"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body"))
		return
	}

	user, err := h.users.FindByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid username or password"))
			return
		}
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "USER_LOOKUP_FAILED", "Failed to lookup user"))
		return
	}

	if !security.CheckPassword(user.PasswordHash, req.Password) {
		responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid username or password"))
		return
	}

	roles := services.ParseRoles(user.Roles)
	token, err := security.GenerateTokenWithClaims(user.Username, user.ID, roles, defaultTokenTTL)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate token"))
		return
	}

	responses.Success(w, r, http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(defaultTokenTTL).Unix(),
		User:      user.Username,
		Roles:     roles,
	})
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body"))
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "username, email and password are required"))
		return
	}

	if existing, err := h.users.FindByUsername(req.Username); err == nil && existing != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusConflict, "USERNAME_TAKEN", "Username already exists"))
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "USER_LOOKUP_FAILED", "Failed to lookup user"))
		return
	}

	if existing, err := h.users.FindByEmail(req.Email); err == nil && existing != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusConflict, "EMAIL_TAKEN", "Email already exists"))
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "USER_LOOKUP_FAILED", "Failed to lookup user"))
		return
	}

	hashed, err := security.HashPassword(req.Password)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "HASH_FAILED", "Failed to hash password"))
		return
	}

	user := models.UserGORM{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashed,
		Roles:        "user",
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
	}
	created, err := h.users.CreateUser(user)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "USER_CREATE_FAILED", "Failed to create user"))
		return
	}

	responses.Success(w, r, http.StatusCreated, RegisterResponse{
		ID:        created.ID,
		Username:  created.Username,
		Email:     created.Email,
		Roles:     services.ParseRoles(created.Roles),
		CreatedAt: created.CreatedAt.Unix(),
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*security.Claims)
	if !ok || claims == nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "NOT_AUTHENTICATED", "User not authenticated"))
		return
	}

	if h.tokenStore != nil && claims.ID != "" && claims.ExpiresAt != nil {
		_ = h.tokenStore.Revoke(claims.ID, claims.ExpiresAt.Time)
	}

	responses.Success(w, r, http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*security.Claims)
	if !ok || claims == nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusUnauthorized, "NOT_AUTHENTICATED", "User not authenticated"))
		return
	}

	var userID uint64
	if claims.Subject != "" {
		if id, err := strconv.ParseUint(claims.Subject, 10, 64); err == nil {
			userID = id
		}
	}

	token, err := security.GenerateTokenWithClaims(claims.Username, uint(userID), claims.Roles, defaultTokenTTL)
	if err != nil {
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate token"))
		return
	}

	responses.Success(w, r, http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(defaultTokenTTL).Unix(),
		User:      claims.Username,
		Roles:     claims.Roles,
	})
}

func (h *AuthHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.URL.Query().Get("username"))
	if username == "" {
		responses.Fail(w, r, responses.NewAppError(http.StatusBadRequest, "MISSING_USERNAME", "username query parameter is required"))
		return
	}

	user, err := h.users.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responses.Fail(w, r, responses.NewAppError(http.StatusNotFound, "USER_NOT_FOUND", "User not found"))
			return
		}
		responses.Fail(w, r, responses.NewAppError(http.StatusInternalServerError, "USER_LOOKUP_FAILED", "Failed to lookup user"))
		return
	}

	responses.Success(w, r, http.StatusOK, UserDetailsResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Roles:     services.ParseRoles(user.Roles),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	})
}
