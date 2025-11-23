// Package main is the entry point for the LightShare backend server.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/lightshare/backend/internal/config"
	"github.com/lightshare/backend/internal/handlers"
	"github.com/lightshare/backend/internal/middleware"
	"github.com/lightshare/backend/internal/repository"
	"github.com/lightshare/backend/internal/services"
	"github.com/lightshare/backend/pkg/crypto"
	"github.com/lightshare/backend/pkg/database"
	"github.com/lightshare/backend/pkg/email"
	"github.com/lightshare/backend/pkg/jwt"
	"github.com/lightshare/backend/pkg/logger"
	"github.com/lightshare/backend/pkg/redis"
)

var (
	version = "dev"
)

func main() {
	// Initialize logger
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.Init(logLevel)

	// Load configuration
	cfg := config.Load()

	// Initialize database
	logger.Info("Connecting to database...")
	db, err := database.New(database.Config{
		URL:             cfg.Database.URL,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	})
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("Failed to close database connection", "error", closeErr)
		}
	}()
	logger.Info("Database connected successfully")

	// Initialize Redis
	logger.Info("Connecting to Redis...")
	redisClient, err := redis.New(redis.Config{
		URL: cfg.Redis.URL,
	})
	if err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		// Clean up database connection before exiting
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("Failed to close database connection during cleanup", "error", closeErr)
		}
		//nolint:gocritic // exitAfterDefer is acceptable here as we manually clean up resources
		os.Exit(1)
	}
	defer func() {
		if redisClient != nil {
			if closeErr := redisClient.Close(); closeErr != nil {
				logger.Error("Failed to close Redis connection", "error", closeErr)
			}
		}
	}()
	logger.Info("Redis connected successfully")

	// Initialize services
	logger.Info("Initializing services...")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db.DB)
	accountRepo := repository.NewAccountRepository(db.DB)

	// Initialize JWT service
	jwtService := jwt.New(jwt.Config{
		Secret:            cfg.JWT.Secret,
		AccessExpiration:  cfg.JWT.AccessExpiration,
		RefreshExpiration: cfg.JWT.RefreshExpiration,
	})

	// Initialize email service
	emailService := email.New(&email.Config{
		SMTPHost:             cfg.Email.SMTPHost,
		SMTPPort:             cfg.Email.SMTPPort,
		SMTPUsername:         cfg.Email.SMTPUsername,
		SMTPPassword:         cfg.Email.SMTPPassword,
		FromEmail:            cfg.Email.FromEmail,
		FromName:             cfg.Email.FromName,
		BaseURL:              cfg.Email.BaseURL,
		MobileDeepLinkScheme: cfg.Email.MobileDeepLinkScheme,
	})

	// Load encryption key for provider tokens
	encryptionKey, err := crypto.LoadEncryptionKey()
	if err != nil {
		logger.Error("Failed to load encryption key", "error", err)
		logger.Info("To generate a new encryption key, run: openssl rand -hex 32")
		os.Exit(1)
	}

	// Initialize auth service
	authService := services.NewAuthService(userRepo, refreshTokenRepo, jwtService, emailService)

	// Initialize provider service
	providerService := services.NewProviderService(accountRepo, encryptionKey)

	// Initialize device service
	deviceService := services.NewDeviceService(
		accountRepo,
		redisClient,
		cfg.Devices.CacheTTL,
		cfg.Devices.RateLimitPerMin,
	)

	logger.Info("Services initialized successfully")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "LightShare API",
		ReadTimeout:           cfg.Server.ReadTimeout,
		WriteTimeout:          cfg.Server.WriteTimeout,
		DisableStartupMessage: false,
		ErrorHandler:          errorHandler,
	})

	// Setup middleware
	middleware.Setup(app)

	// Setup routes
	setupRoutes(app, authService, providerService, deviceService, jwtService)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		logger.Info("Starting server", "address", addr, "version", version)
		if err := app.Listen(addr); err != nil {
			logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}

	logger.Info("Server stopped")
}

func setupRoutes(app *fiber.App, authService *services.AuthService, providerService *services.ProviderService, deviceService *services.DeviceService, jwtService *jwt.Service) {
	// Health check endpoints
	app.Get("/health", handlers.Health(version))
	app.Get("/ready", handlers.Ready())

	// API v1 routes
	v1 := app.Group("/api/v1")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	providerHandler := handlers.NewProviderHandler(providerService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/signup", authHandler.Signup)
	auth.Post("/login", authHandler.Login)
	auth.Post("/verify-email", authHandler.VerifyEmail)
	auth.Post("/magic-link", authHandler.RequestMagicLink)
	auth.Post("/magic-link/verify", authHandler.LoginWithMagicLink)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", authHandler.Logout)

	// Protected auth routes
	authMiddleware := middleware.AuthMiddleware(jwtService)
	auth.Get("/me", authMiddleware, authHandler.Me)
	auth.Post("/logout-all", authMiddleware, authHandler.LogoutAll)

	// Provider routes (protected)
	providers := v1.Group("/providers", authMiddleware)
	providers.Post("/connect", providerHandler.ConnectProvider)

	// Account routes (protected)
	accounts := v1.Group("/accounts", authMiddleware)
	accounts.Get("", providerHandler.ListAccounts)
	accounts.Delete("/:id", providerHandler.DisconnectAccount)

	// Device routes (protected) - Phase 4
	// List all devices across all accounts
	v1.Get("/devices", authMiddleware, deviceHandler.ListDevices)

	// Account-specific device routes
	v1.Get("/accounts/:accountId/devices", authMiddleware, deviceHandler.ListAccountDevices)
	v1.Get("/accounts/:accountId/devices/:deviceId", authMiddleware, deviceHandler.GetDevice)
	v1.Post("/accounts/:accountId/devices/:selector/action", authMiddleware, deviceHandler.ExecuteAction)
	v1.Post("/accounts/:accountId/devices/refresh", authMiddleware, deviceHandler.RefreshDevices)
}

func errorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Log the error
	logger.Error("Request error",
		"error", err,
		"status", code,
		"path", c.Path(),
		"method", c.Method(),
	)

	// Return JSON error response
	return c.Status(code).JSON(fiber.Map{
		"error":  message,
		"status": code,
	})
}
