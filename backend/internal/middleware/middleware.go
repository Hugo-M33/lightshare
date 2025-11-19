package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/lightshare/backend/pkg/logger"
)

// Setup sets up all middleware for the Fiber app
func Setup(app *fiber.App) {
	// Recover from panics
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Request ID
	app.Use(requestid.New())

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		ExposeHeaders:    "X-Request-ID,X-RateLimit-Limit,X-RateLimit-Remaining,X-RateLimit-Reset",
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	// Request logging
	app.Use(RequestLogger())
}

// RequestLogger returns a middleware that logs HTTP requests
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request ID
		requestID := c.GetRespHeader("X-Request-ID")

		// Log the request
		logger.Info("HTTP request",
			"request_id", requestID,
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"latency_ms", latency.Milliseconds(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
		)

		return err
	}
}
