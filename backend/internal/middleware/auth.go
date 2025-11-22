package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/lightshare/backend/pkg/jwt"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(jwtService *jwt.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "token expired",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		// Store user information in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)

		return c.Next()
	}
}

// GetUserID gets the user ID from the request context
func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}
	return userID, nil
}

// GetUserEmail gets the user email from the request context
func GetUserEmail(c *fiber.Ctx) (string, error) {
	email, ok := c.Locals("user_email").(string)
	if !ok {
		return "", fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}
	return email, nil
}

// GetUserRole gets the user role from the request context
func GetUserRole(c *fiber.Ctx) (string, error) {
	role, ok := c.Locals("user_role").(string)
	if !ok {
		return "", fiber.NewError(fiber.StatusUnauthorized, "user not authenticated")
	}
	return role, nil
}

// RequireRole creates a middleware that requires a specific role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, err := GetUserRole(c)
		if err != nil {
			return err
		}

		if role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "insufficient permissions",
			})
		}

		return c.Next()
	}
}
