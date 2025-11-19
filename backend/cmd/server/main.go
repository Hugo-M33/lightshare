package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/lightshare/backend/internal/config"
	"github.com/lightshare/backend/internal/handlers"
	"github.com/lightshare/backend/internal/middleware"
	"github.com/lightshare/backend/pkg/logger"
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
	setupRoutes(app)

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
	if err := app.Shutdown(); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}
	logger.Info("Server stopped")
}

func setupRoutes(app *fiber.App) {
	// Health check endpoints
	app.Get("/health", handlers.Health(version))
	app.Get("/ready", handlers.Ready())

	// API v1 routes
	v1 := app.Group("/v1")

	// Auth routes (to be implemented in Phase 2)
	_ = v1.Group("/auth")

	// Account routes (to be implemented in Phase 3)
	_ = v1.Group("/accounts")

	// Provider routes (to be implemented in Phase 3)
	_ = v1.Group("/providers")
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
		"error":   message,
		"status":  code,
	})
}
