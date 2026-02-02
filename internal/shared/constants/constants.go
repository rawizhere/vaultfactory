// Package constants содержит константы приложения.
package constants

const (
	// HTTP Status Codes
	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusInternalServerError = 500

	// Timeouts
	DefaultTimeoutSeconds = 30
	ReadTimeoutSeconds    = 15
	WriteTimeoutSeconds   = 15
	IdleTimeoutSeconds    = 60

	// Password requirements
	MinPasswordLength = 8
	MaxNameLength     = 255

	// JWT
	DefaultJWTExpireHours         = 24
	DefaultRefreshTokenExpireDays = 30

	// Database
	DefaultPageSize = 50
	MaxPageSize     = 100
)
