// Package errors предоставляет типы ошибок для приложения.
package errors

import (
	"fmt"
	"net/http"
)

// AppError представляет ошибку приложения с HTTP кодом.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewBadRequest создает ошибку 400 Bad Request.
func NewBadRequest(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Err:     err,
	}
}

// NewUnauthorized создает ошибку 401 Unauthorized.
func NewUnauthorized(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
		Err:     err,
	}
}

// NewForbidden создает ошибку 403 Forbidden.
func NewForbidden(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
		Err:     err,
	}
}

// NewNotFound создает ошибку 404 Not Found.
func NewNotFound(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
		Err:     err,
	}
}

// NewConflict создает ошибку 409 Conflict.
func NewConflict(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
		Err:     err,
	}
}

// NewInternalServerError создает ошибку 500 Internal Server Error.
func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}
