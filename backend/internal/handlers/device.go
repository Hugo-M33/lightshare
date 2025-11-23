package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lightshare/backend/internal/models"
	"github.com/lightshare/backend/internal/services"
)

// DeviceHandler handles device-related HTTP requests
type DeviceHandler struct {
	deviceService *services.DeviceService
}

// NewDeviceHandler creates a new device handler
func NewDeviceHandler(deviceService *services.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// ListDevices lists all devices for the authenticated user
// GET /api/v1/devices
func (h *DeviceHandler) ListDevices(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}

	devices, err := h.deviceService.ListDevices(c.Context(), userID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to list devices")
	}

	return c.JSON(fiber.Map{
		"devices": devices,
	})
}

// ListAccountDevices lists devices for a specific account
// GET /api/v1/accounts/:accountId/devices
func (h *DeviceHandler) ListAccountDevices(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}

	accountID := c.Params("accountId")
	if accountID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "account ID is required")
	}

	devices, err := h.deviceService.ListAccountDevices(c.Context(), userID.String(), accountID)
	if err != nil {
		if err.Error() == "account not found: account not found" {
			return fiber.NewError(fiber.StatusNotFound, "account not found")
		}
		if err.Error() == "unauthorized: user does not own this account" {
			return fiber.NewError(fiber.StatusForbidden, "unauthorized")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to list devices")
	}

	return c.JSON(fiber.Map{
		"devices": devices,
	})
}

// GetDevice returns a specific device
// GET /api/v1/accounts/:accountId/devices/:deviceId
func (h *DeviceHandler) GetDevice(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}

	accountID := c.Params("accountId")
	deviceID := c.Params("deviceId")

	if accountID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "account ID is required")
	}
	if deviceID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "device ID is required")
	}

	device, err := h.deviceService.GetDevice(c.Context(), userID.String(), accountID, deviceID)
	if err != nil {
		if err.Error() == "account not found: account not found" {
			return fiber.NewError(fiber.StatusNotFound, "account not found")
		}
		if err.Error() == "unauthorized: user does not own this account" {
			return fiber.NewError(fiber.StatusForbidden, "unauthorized")
		}
		if err.Error() == "rate limit exceeded: max 30 requests per minute" {
			return fiber.NewError(fiber.StatusTooManyRequests, "rate limit exceeded")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get device")
	}

	return c.JSON(device)
}

// ExecuteAction executes a control action on device(s)
// POST /api/v1/accounts/:accountId/devices/:selector/action
func (h *DeviceHandler) ExecuteAction(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}

	accountID := c.Params("accountId")
	selector := c.Params("selector")

	if accountID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "account ID is required")
	}
	if selector == "" {
		return fiber.NewError(fiber.StatusBadRequest, "selector is required")
	}

	var action models.ActionRequest
	if err := c.BodyParser(&action); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Validate action
	if err := action.ValidateParameters(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err := h.deviceService.ExecuteAction(c.Context(), userID.String(), accountID, selector, &action)
	if err != nil {
		if err.Error() == "account not found: account not found" {
			return fiber.NewError(fiber.StatusNotFound, "account not found")
		}
		if err.Error() == "unauthorized: user does not own this account" {
			return fiber.NewError(fiber.StatusForbidden, "unauthorized")
		}
		if err.Error() == "rate limit exceeded: max 30 requests per minute" {
			return fiber.NewError(fiber.StatusTooManyRequests, "rate limit exceeded")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to execute action")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "action executed successfully",
	})
}

// RefreshDevices forces a cache refresh for an account
// POST /api/v1/accounts/:accountId/devices/refresh
func (h *DeviceHandler) RefreshDevices(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}

	accountID := c.Params("accountId")
	if accountID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "account ID is required")
	}

	devices, err := h.deviceService.RefreshDevices(c.Context(), userID.String(), accountID)
	if err != nil {
		if err.Error() == "account not found: account not found" {
			return fiber.NewError(fiber.StatusNotFound, "account not found")
		}
		if err.Error() == "unauthorized: user does not own this account" {
			return fiber.NewError(fiber.StatusForbidden, "unauthorized")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to refresh devices")
	}

	return c.JSON(fiber.Map{
		"devices": devices,
	})
}
