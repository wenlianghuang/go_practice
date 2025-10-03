package responses

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

type successEnvelope struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id"`
}

type errorEnvelope struct {
	Success   bool      `json:"success"`
	Error     *AppError `json:"error"`
	RequestID string    `json:"request_id"`
}

// Success sends standardized success response.
func Success(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	reqID := r.Context().Value("request_id")
	var reqIDStr string
	if reqID != nil {
		reqIDStr = reqID.(string)
	} else {
		reqIDStr = r.Header.Get("X-Request-ID")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := successEnvelope{
		Success:   true,
		Data:      data,
		RequestID: reqIDStr,
	}

	json.NewEncoder(w).Encode(response)
}

// Fail sends standardized error response.
func Fail(w http.ResponseWriter, r *http.Request, appErr *AppError) {
	if appErr == nil {
		appErr = NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}

	reqID := r.Context().Value("request_id")
	var reqIDStr string
	if reqID != nil {
		reqIDStr = reqID.(string)
	} else {
		reqIDStr = r.Header.Get("X-Request-ID")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)

	response := errorEnvelope{
		Success:   false,
		Error:     appErr,
		RequestID: reqIDStr,
	}

	json.NewEncoder(w).Encode(response)
}
