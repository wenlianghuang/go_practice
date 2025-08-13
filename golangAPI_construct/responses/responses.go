package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
func Success(c *gin.Context, status int, data interface{}) {
	reqID := c.GetString("request_id")
	c.JSON(status, successEnvelope{
		Success:   true,
		Data:      data,
		RequestID: reqID,
	})
}

// Fail sends standardized error response.
func Fail(c *gin.Context, appErr *AppError) {
	if appErr == nil {
		appErr = NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
	reqID := c.GetString("request_id")
	c.JSON(appErr.Status, errorEnvelope{
		Success:   false,
		Error:     appErr,
		RequestID: reqID,
	})
}
