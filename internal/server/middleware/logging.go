// Package middleware содержит HTTP middleware для сервера.
package middleware

import (
	"net/http"
	"time"

	"github.com/tempizhere/vaultfactory/internal/shared/logger"
	"go.uber.org/zap"
)

// LoggingMiddleware предоставляет middleware для логирования HTTP запросов.
type LoggingMiddleware struct {
	logger logger.Logger
}

// NewLoggingMiddleware создает новый экземпляр LoggingMiddleware.
func NewLoggingMiddleware(logger logger.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Logging логирует HTTP запросы и ответы.
func (m *LoggingMiddleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		m.logger.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
			zap.Int("status_code", wrapped.statusCode),
			zap.Duration("duration", duration),
			zap.Int("response_size", wrapped.size),
		)

		if wrapped.statusCode >= 500 {
			m.logger.Error("HTTP server error",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status_code", wrapped.statusCode),
				zap.Duration("duration", duration),
			)
		}
	})
}

// responseWriter обертка для http.ResponseWriter для перехвата статус кода.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// WriteHeader перехватывает вызов WriteHeader для получения статус кода.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write перехватывает вызов Write для подсчета размера ответа.
func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}
