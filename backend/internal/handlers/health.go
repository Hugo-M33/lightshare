package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// Health returns the health check handler
func Health(version string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Version:   version,
		})
	}
}

// ReadyResponse represents the readiness check response
type ReadyResponse struct {
	Checks map[string]string `json:"checks"`
	Status string            `json:"status"`
	Ready  bool              `json:"ready"`
}

// Ready returns the readiness check handler
// This will be extended to check database and Redis connections
func Ready() fiber.Handler {
	return func(c *fiber.Ctx) error {
		checks := map[string]string{
			"database": "ok",
			"redis":    "ok",
		}

		// TODO: Add actual health checks for database and Redis

		allHealthy := true
		for _, status := range checks {
			if status != "ok" {
				allHealthy = false
				break
			}
		}

		response := ReadyResponse{
			Status: "ready",
			Checks: checks,
			Ready:  allHealthy,
		}

		if !allHealthy {
			response.Status = "not_ready"
			return c.Status(fiber.StatusServiceUnavailable).JSON(response)
		}

		return c.JSON(response)
	}
}
