package container

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/uptrace/bun"

	"github.com/tempizhere/vaultfactory/internal/server/auth"
	"github.com/tempizhere/vaultfactory/internal/server/handlers"
	"github.com/tempizhere/vaultfactory/internal/server/middleware"
	"github.com/tempizhere/vaultfactory/internal/server/repository"
	"github.com/tempizhere/vaultfactory/internal/server/service"
	"github.com/tempizhere/vaultfactory/internal/shared/config"
	"github.com/tempizhere/vaultfactory/internal/shared/crypto"
	"github.com/tempizhere/vaultfactory/internal/shared/interfaces"
	"github.com/tempizhere/vaultfactory/internal/shared/logger"
)

type Container struct {
	Config config.ConfigReader
	Logger logger.Logger
	DB     *bun.DB
	SQLDB  *sql.DB

	// Repositories
	UserRepo    interfaces.UserRepository
	SessionRepo interfaces.SessionRepository
	DataRepo    interfaces.DataRepository
	VersionRepo interfaces.VersionRepository

	// Services
	CryptoService *crypto.CryptoService
	JWTService    *auth.JWTService
	AuthService   interfaces.AuthService
	DataService   interfaces.DataService

	// Handlers
	AuthHandler *handlers.AuthHandler
	DataHandler *handlers.DataHandler

	// Middleware
	AuthMiddleware    *middleware.AuthMiddleware
	LoggingMiddleware *middleware.LoggingMiddleware

	// Router
	Router *mux.Router
}

func NewContainer(cfg config.ConfigReader, appLogger logger.Logger, db *bun.DB, sqldb *sql.DB) *Container {
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	dataRepo := repository.NewDataRepository(db)
	versionRepo := repository.NewVersionRepository(db)

	cryptoService := crypto.NewCryptoService()
	jwtService := auth.NewJWTService(cfg.GetJWTSecret(), cfg.GetJWTExpireDuration())
	authService := service.NewAuthService(userRepo, sessionRepo, cryptoService, jwtService, appLogger)
	dataService := service.NewDataService(dataRepo, versionRepo, cryptoService)

	authHandler := handlers.NewAuthHandler(authService)
	dataHandler := handlers.NewDataHandler(dataService)

	authMiddleware := middleware.NewAuthMiddleware(authService)
	loggingMiddleware := middleware.NewLoggingMiddleware(appLogger)

	router := setupRoutes(authHandler, dataHandler, authMiddleware, loggingMiddleware)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods("GET")

	return &Container{
		Config:            cfg,
		Logger:            appLogger,
		DB:                db,
		SQLDB:             sqldb,
		UserRepo:          userRepo,
		SessionRepo:       sessionRepo,
		DataRepo:          dataRepo,
		VersionRepo:       versionRepo,
		CryptoService:     cryptoService,
		JWTService:        jwtService,
		AuthService:       authService,
		DataService:       dataService,
		AuthHandler:       authHandler,
		DataHandler:       dataHandler,
		AuthMiddleware:    authMiddleware,
		LoggingMiddleware: loggingMiddleware,
		Router:            router,
	}
}

// setupRoutes устанавливает маршруты для API.
func setupRoutes(authHandler *handlers.AuthHandler, dataHandler *handlers.DataHandler, authMiddleware *middleware.AuthMiddleware, loggingMiddleware *middleware.LoggingMiddleware) *mux.Router {
	router := mux.NewRouter()

	router.Use(loggingMiddleware.Logging)

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

	return router
}
