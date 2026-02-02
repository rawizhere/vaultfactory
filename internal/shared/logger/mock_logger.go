// Package logger предоставляет централизованное логирование на основе Zap.
package logger

import (
	"go.uber.org/zap"
)

// MockLogger реализует интерфейс Logger для тестирования.
type MockLogger struct{}

// NewMockLogger создает новый mock логгер.
func NewMockLogger() Logger {
	return &MockLogger{}
}

// Debug логирует сообщение уровня Debug.
func (m *MockLogger) Debug(msg string, fields ...zap.Field) {}

// Info логирует сообщение уровня Info.
func (m *MockLogger) Info(msg string, fields ...zap.Field) {}

// Warn логирует сообщение уровня Warn.
func (m *MockLogger) Warn(msg string, fields ...zap.Field) {}

// Error логирует сообщение уровня Error.
func (m *MockLogger) Error(msg string, fields ...zap.Field) {}

// Fatal логирует сообщение уровня Fatal и завершает программу.
func (m *MockLogger) Fatal(msg string, fields ...zap.Field) {}

// With создает новый логгер с дополнительными полями.
func (m *MockLogger) With(fields ...zap.Field) Logger {
	return m
}

// Sync синхронизирует буферы логгера.
func (m *MockLogger) Sync() error {
	return nil
}
