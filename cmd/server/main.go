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

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"

	"github.com/tempizhere/vaultfactory/internal/server/auth"
	"github.com/tempizhere/vaultfactory/internal/server/handlers"
	"github.com/tempizhere/vaultfactory/internal/server/middleware"
	"github.com/tempizhere/vaultfactory/internal/server/repository"
	"github.com/tempizhere/vaultfactory/internal/server/service"
	"github.com/tempizhere/vaultfactory/internal/shared/config"
	"github.com/tempizhere/vaultfactory/internal/shared/crypto"
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

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Логируем информацию о сборке
	logger.Info("Starting VaultFactory Server",
		zap.String("version", getBuildInfo(buildVersion)),
		zap.String("build_date", getBuildInfo(buildDate)),
		zap.String("build_commit", getBuildInfo(buildCommit)))

	sqldb := openDB(cfg.GetDSN())
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := runMigrations(db); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	cryptoService := crypto.NewCryptoService()
	jwtService := auth.NewJWTService(cfg.GetJWTSecret(), cfg.Security.JWTExpireDuration)

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	dataRepo := repository.NewDataRepository(db)
	versionRepo := repository.NewVersionRepository(db)

	authService := service.NewAuthService(userRepo, sessionRepo, cryptoService, jwtService)
	dataService := service.NewDataService(dataRepo, versionRepo, cryptoService)

	authHandler := handlers.NewAuthHandler(authService)
	dataHandler := handlers.NewDataHandler(dataService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	router := setupRoutes(authHandler, dataHandler, authMiddleware)

	server := &http.Server{
		Addr:         cfg.GetServerAddr(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Starting server", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func setupRoutes(authHandler *handlers.AuthHandler, dataHandler *handlers.DataHandler, authMiddleware *middleware.AuthMiddleware) *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")
	auth.HandleFunc("/refresh", authHandler.Refresh).Methods("POST")
	auth.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	data := api.PathPrefix("/data").Subrouter()
	data.Use(authMiddleware.RequireAuth)
	data.HandleFunc("", dataHandler.CreateData).Methods("POST")
	data.HandleFunc("", dataHandler.GetUserData).Methods("GET")
	data.HandleFunc("/sync", dataHandler.SyncData).Methods("GET")
	data.HandleFunc("/{id}", dataHandler.GetData).Methods("GET")
	data.HandleFunc("/{id}", dataHandler.UpdateData).Methods("PUT")
	data.HandleFunc("/{id}", dataHandler.DeleteData).Methods("DELETE")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
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
