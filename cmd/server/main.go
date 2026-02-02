// Package main содержит точку входа для серверного приложения VaultFactory.
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"

	"github.com/tempizhere/vaultfactory/internal/server/container"
	"github.com/tempizhere/vaultfactory/internal/shared/config"
	"github.com/tempizhere/vaultfactory/internal/shared/constants"
	"github.com/tempizhere/vaultfactory/internal/shared/logger"
	"github.com/tempizhere/vaultfactory/internal/shared/models"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	cfg, err := config.LoadConfig("configs/server.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаем логгер на основе конфигурации
	appLogger, err := logger.NewLogger(cfg.GetLoggingLevel(), cfg.GetLoggingFormat())
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer func() { _ = appLogger.Sync() }()

	// Логируем информацию о сборке
	appLogger.Info("Starting VaultFactory Server",
		zap.String("version", getBuildInfo(buildVersion)),
		zap.String("build_date", getBuildInfo(buildDate)),
		zap.String("build_commit", getBuildInfo(buildCommit)))

	sqldb := openDB(cfg.GetDSN())
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	if err := db.Ping(); err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := runMigrations(db); err != nil {
		appLogger.Fatal("Failed to run migrations", zap.Error(err))
	}

	container := container.NewContainer(cfg, appLogger, db, sqldb)

	server := &http.Server{
		Addr:         cfg.GetServerAddr(),
		Handler:      container.Router,
		ReadTimeout:  time.Duration(constants.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(constants.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(constants.IdleTimeoutSeconds) * time.Second,
	}

	go func() {
		appLogger.Info("Starting server", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("Server exited")
}

func openDB(dsn string) *sql.DB {
	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	return sqldb
}

func runMigrations(db *bun.DB) error {
	ctx := context.Background()

	// Создаём таблицы
	_, err := db.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*models.UserSession)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*models.DataItem)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*models.DataVersion)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func getBuildInfo(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}
