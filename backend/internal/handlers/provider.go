package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/lightshare/backend/internal/repository"
	"github.com/lightshare/backend/internal/services"
	"github.com/lightshare/backend/pkg/logger"
)

// ProviderHandler handles provider connection endpoints
type ProviderHandler struct {
	providerService *services.ProviderService
}

// NewProviderHandler creates a new provider handler
func NewProviderHandler(providerService *services.ProviderService) *ProviderHandler {
	return &ProviderHandler{
		providerService: providerService,
	}
}

// ConnectProviderRequest represents the connect provider request body
type ConnectProviderRequest struct {
	Provider string `json:"provider"`
	Token    string `json:"token"`
}

// ConnectProvider handles provider connection
func (h *ProviderHandler) ConnectProvider(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req ConnectProviderRequest
	if parseRequestBody(c, &req) {
		return nil
	}

	// Validate request
	if req.Provider == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "provider is required",
		})
	}
	if req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "token is required",
		})
	}

	// Call provider service
	account, err := h.providerService.ConnectProvider(c.Context(), userID, services.ConnectProviderRequest{
		Provider: req.Provider,
		Token:    req.Token,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidProvider) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid provider type",
			})
		}
		if errors.Is(err, services.ErrInvalidToken) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid provider token",
			})
		}
		if err.Error() == "this provider account is already connected" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "this provider account is already connected",
			})
		}
		logger.Error("Failed to connect provider", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to connect provider",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(account.ToResponse())
}

// ListAccounts handles listing all connected accounts
func (h *ProviderHandler) ListAccounts(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Call provider service
	accounts, err := h.providerService.ListAccounts(c.Context(), userID)
	if err != nil {
		logger.Error("Failed to list accounts", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to list accounts",
		})
	}

	// Convert to response format
	accountResponses := make([]interface{}, 0, len(accounts))
	for _, account := range accounts {
		accountResponses = append(accountResponses, account.ToResponse())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accounts": accountResponses,
	})
}

// DisconnectAccount handles disconnecting a provider account
func (h *ProviderHandler) DisconnectAccount(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get account ID from URL param
	accountIDStr := c.Params("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid account id",
		})
	}

	// Call provider service
	err = h.providerService.DisconnectAccount(c.Context(), userID, accountID)
	if err != nil {
		if errors.Is(err, repository.ErrAccountNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "account not found",
			})
		}
		if errors.Is(err, services.ErrAccountNotOwned) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "account not owned by user",
			})
		}
		logger.Error("Failed to disconnect account", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to disconnect account",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "account disconnected successfully",
	})
}
