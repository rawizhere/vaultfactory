// Package validator предоставляет валидацию входных данных.
package validator

import (
	"regexp"
	"strings"
)

// Validator предоставляет методы валидации.
type Validator struct{}

// NewValidator создает новый экземпляр Validator.
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateEmail проверяет корректность email адреса.
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{Field: "email", Message: "email is required"}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &ValidationError{Field: "email", Message: "invalid email format"}
	}

	return nil
}

// ValidatePassword проверяет корректность пароля.
func (v *Validator) ValidatePassword(password string) error {
	if password == "" {
		return &ValidationError{Field: "password", Message: "password is required"}
	}

	if len(password) < 8 {
		return &ValidationError{Field: "password", Message: "password must be at least 8 characters long"}
	}

	return nil
}

// ValidateDataName проверяет корректность имени данных.
func (v *Validator) ValidateDataName(name string) error {
	if strings.TrimSpace(name) == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}

	if len(name) > 255 {
		return &ValidationError{Field: "name", Message: "name must be less than 255 characters"}
	}

	return nil
}

// ValidationError представляет ошибку валидации.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
