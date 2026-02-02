package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/tempizhere/vaultfactory/internal/shared/errors"
)

// ErrorHandler обрабатывает ошибки HTTP запросов.
type ErrorHandler struct{}

// NewErrorHandler создает новый экземпляр ErrorHandler.
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleError обрабатывает ошибку и отправляет JSON ответ.
func (h *ErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *errors.AppError
	var ok bool

	if appErr, ok = err.(*errors.AppError); !ok {
		appErr = errors.NewInternalServerError("Internal server error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
	}

	_ = json.NewEncoder(w).Encode(response)
}
