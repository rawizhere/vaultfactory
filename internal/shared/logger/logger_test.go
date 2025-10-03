package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	t.Run("create logger with info level", func(t *testing.T) {
		logger, err := NewLogger("info", "json")
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("create logger with debug level", func(t *testing.T) {
		logger, err := NewLogger("debug", "console")
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("create logger with invalid level defaults to info", func(t *testing.T) {
		logger, err := NewLogger("invalid", "json")
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})
}

func TestNewProductionLogger(t *testing.T) {
	logger, err := NewProductionLogger()
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestNewDevelopmentLogger(t *testing.T) {
	logger, err := NewDevelopmentLogger()
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestZapLogger_Methods(t *testing.T) {
	logger, err := NewLogger("info", "json")
	assert.NoError(t, err)

	// Тестируем, что методы не паникуют
	assert.NotPanics(t, func() {
		logger.Debug("debug message", zap.String("key", "value"))
		logger.Info("info message", zap.String("key", "value"))
		logger.Warn("warn message", zap.String("key", "value"))
		logger.Error("error message", zap.String("key", "value"))
	})

	// Тестируем With
	withLogger := logger.With(zap.String("context", "test"))
	assert.NotNil(t, withLogger)

	// Тестируем Sync
	_ = logger.Sync()
}

func TestMockLogger(t *testing.T) {
	mockLogger := NewMockLogger()
	assert.NotNil(t, mockLogger)

	// Тестируем, что методы не паникуют
	assert.NotPanics(t, func() {
		mockLogger.Debug("debug message", zap.String("key", "value"))
		mockLogger.Info("info message", zap.String("key", "value"))
		mockLogger.Warn("warn message", zap.String("key", "value"))
		mockLogger.Error("error message", zap.String("key", "value"))
		mockLogger.Fatal("fatal message", zap.String("key", "value"))
	})

	// Тестируем With
	withLogger := mockLogger.With(zap.String("context", "test"))
	assert.NotNil(t, withLogger)

	// Тестируем Sync
	err := mockLogger.Sync()
	assert.NoError(t, err)
}
