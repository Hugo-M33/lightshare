package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/lightshare/backend/internal/middleware"
	"github.com/lightshare/backend/internal/repository"
	"github.com/lightshare/backend/internal/services"
	"github.com/lightshare/backend/pkg/logger"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// parseRequestBody parses the request body and sends an error response if parsing fails.
// Returns true if an error occurred (and error response was sent), false otherwise.
func parseRequestBody(c *fiber.Ctx, req interface{}) bool {
	if err := c.BodyParser(req); err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
		return true
	}
	return false
}

// SignupRequest represents the signup request body
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup handles user signup
func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var req SignupRequest
	if parseRequestBody(c, &req) {
		return nil
	}

	// Call auth service
	resp, err := h.authService.Signup(c.Context(), services.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, services.ErrWeakPassword) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "password must be at least 8 characters",
			})
		}
		if err.Error() == "email already registered" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "email already registered",
			})
		}
		if err.Error() == "invalid email address" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid email address",
			})
		}
		logger.Error("Failed to signup user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create account",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if parseRequestBody(c, &req) {
		return nil
	}

	// Get user agent and IP address
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()

	// Call auth service
	resp, err := h.authService.Login(c.Context(), services.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}, &userAgent, &ipAddress)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid email or password",
			})
		}
		if errors.Is(err, services.ErrEmailNotVerified) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "email not verified",
			})
		}
		logger.Error("Failed to login user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to login",
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// VerifyEmailRequest represents the verify email request body
type VerifyEmailRequest struct {
	Token string `json:"token"`
}

// VerifyEmail handles email verification
func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	var req VerifyEmailRequest
	if parseRequestBody(c, &req) {
		return nil
	}

	// Get user agent and IP address
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()

	// Call auth service
	resp, err := h.authService.VerifyEmail(c.Context(), req.Token, &userAgent, &ipAddress)
	if err != nil {
		if errors.Is(err, repository.ErrTokenExpired) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "verification token expired",
			})
		}
		logger.Error("Failed to verify email", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to verify email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// MagicLinkRequest represents the magic link request body
type MagicLinkRequest struct {
	Email string `json:"email"`
}

// RequestMagicLink handles magic link request
func (h *AuthHandler) RequestMagicLink(c *fiber.Ctx) error {
	var req MagicLinkRequest
	if parseRequestBody(c, &req) {
		return nil
	}

	// Call auth service
	err := h.authService.RequestMagicLink(c.Context(), req.Email)
	if err != nil {
		logger.Error("Failed to send magic link", "error", err)
		// Don't reveal if email exists or not
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "if the email exists, a magic link has been sent",
	})
}

// LoginWithMagicLinkRequest represents the magic link login request body
type LoginWithMagicLinkRequest struct {
	Token string `json:"token"`
}

// LoginWithMagicLink handles login with magic link
func (h *AuthHandler) LoginWithMagicLink(c *fiber.Ctx) error {
	var req LoginWithMagicLinkRequest
	if parseRequestBody(c, &req) {
		return nil
	}

	// Get user agent and IP address
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()

	// Call auth service
	resp, err := h.authService.LoginWithMagicLink(c.Context(), req.Token, &userAgent, &ipAddress)
	if err != nil {
		if err.Error() == "magic link expired" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "magic link expired",
			})
		}
		logger.Error("Failed to login with magic link", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid magic link",
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// RefreshTokenRequest represents the refresh token request body
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if parseRequestBody(c, &req) {
		return nil
	}

	// Get user agent and IP address
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()

	// Call auth service
	resp, err := h.authService.RefreshToken(c.Context(), req.RefreshToken, &userAgent, &ipAddress)
	if err != nil {
		if err.Error() == "invalid refresh token" || err.Error() == "refresh token revoked" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		logger.Error("Failed to refresh token", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to refresh token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// LogoutRequest represents the logout request body
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req LogoutRequest
	if parseRequestBody(c, &req) {
		return nil
	}

	// Call auth service
	err := h.authService.Logout(c.Context(), req.RefreshToken)
	if err != nil {
		logger.Error("Failed to logout user", "error", err)
		// Don't fail on logout errors
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logged out successfully",
	})
}

// LogoutAll handles logout from all devices
func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	// Get user ID from context
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	// Call auth service
	err = h.authService.LogoutAll(c.Context(), userID)
	if err != nil {
		logger.Error("Failed to logout all", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to logout from all devices",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logged out from all devices successfully",
	})
}

// Me returns the current user's information
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	email, err := middleware.GetUserEmail(c)
	if err != nil {
		return err
	}

	role, err := middleware.GetUserRole(c)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":    userID,
		"email": email,
		"role":  role,
	})
}
