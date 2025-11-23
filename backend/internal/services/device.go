package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lightshare/backend/internal/models"
	"github.com/lightshare/backend/internal/repository"
	"github.com/lightshare/backend/pkg/providers"
	"github.com/redis/go-redis/v9"
)

// DeviceService handles device-related business logic
type DeviceService struct {
	accountRepo     *repository.AccountRepository
	cache           *redis.Client
	cacheTTL        time.Duration
	rateLimitPerMin int
}

// NewDeviceService creates a new device service
func NewDeviceService(
	accountRepo *repository.AccountRepository,
	cache *redis.Client,
	cacheTTL time.Duration,
	rateLimitPerMin int,
) *DeviceService {
	return &DeviceService{
		accountRepo:     accountRepo,
		cache:           cache,
		cacheTTL:        cacheTTL,
		rateLimitPerMin: rateLimitPerMin,
	}
}

// ListDevices returns all devices for a user's accounts
func (s *DeviceService) ListDevices(ctx context.Context, userID string) ([]*models.Device, error) {
	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get all accounts for user
	accounts, err := s.accountRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	allDevices := make([]*models.Device, 0)

	// Fetch devices for each account
	for _, account := range accounts {
		// Check cache first
		devices, err := s.getCachedDevices(ctx, account.ID.String())
		if err == nil {
			// Cache hit
			allDevices = append(allDevices, devices...)
			continue
		}

		// Cache miss - fetch from provider
		devices, err = s.fetchDevicesFromProvider(ctx, account)
		if err != nil {
			// Log error but continue with other accounts
			continue
		}

		// Cache the devices
		if err := s.setCachedDevices(ctx, account.ID.String(), devices); err != nil {
			// Log error but continue
			_ = err
		}

		allDevices = append(allDevices, devices...)
	}

	return allDevices, nil
}

// ListAccountDevices returns devices for a specific account
func (s *DeviceService) ListAccountDevices(ctx context.Context, userID, accountID string) ([]*models.Device, error) {
	// Get account and verify ownership
	account, err := s.accountRepo.FindByIDString(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	if account.OwnerUserID.String() != userID {
		return nil, fmt.Errorf("unauthorized: user does not own this account")
	}

	// Check cache first
	devices, err := s.getCachedDevices(ctx, accountID)
	if err == nil {
		return devices, nil
	}

	// Cache miss - fetch from provider
	devices, err = s.fetchDevicesFromProvider(ctx, account)
	if err != nil {
		return nil, err
	}

	// Cache the devices
	if err := s.setCachedDevices(ctx, accountID, devices); err != nil {
		// Log error but continue
		_ = err
	}

	return devices, nil
}

// GetDevice returns a specific device by ID
func (s *DeviceService) GetDevice(ctx context.Context, userID, accountID, deviceID string) (*models.Device, error) {
	// Get account and verify ownership
	account, err := s.accountRepo.FindByIDString(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	if account.OwnerUserID.String() != userID {
		return nil, fmt.Errorf("unauthorized: user does not own this account")
	}

	// Check rate limit
	if rateLimitErr := s.checkRateLimit(ctx, accountID); rateLimitErr != nil {
		return nil, rateLimitErr
	}

	// Get decrypted token
	token, err := s.accountRepo.GetDecryptedToken(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Create provider client
	client, err := providers.NewClient(providers.Provider(account.Provider))
	if err != nil {
		return nil, fmt.Errorf("failed to create provider client: %w", err)
	}

	// Get device from provider
	providerDevice, err := client.GetDevice(token, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device from provider: %w", err)
	}

	// Convert to our device model
	device := s.convertProviderDevice(providerDevice, accountID, account.Provider)

	return device, nil
}

// ExecuteAction executes a control action on device(s)
func (s *DeviceService) ExecuteAction(ctx context.Context, userID, accountID, selector string, action *models.ActionRequest) error {
	// Validate action
	if err := action.ValidateParameters(); err != nil {
		return fmt.Errorf("invalid action parameters: %w", err)
	}

	// Get account and verify ownership
	account, err := s.accountRepo.FindByIDString(ctx, accountID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	if account.OwnerUserID.String() != userID {
		return fmt.Errorf("unauthorized: user does not own this account")
	}

	// Check rate limit
	if rateLimitErr := s.checkRateLimit(ctx, accountID); rateLimitErr != nil {
		return rateLimitErr
	}

	// Get decrypted token
	token, err := s.accountRepo.GetDecryptedToken(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	// Create provider client
	client, err := providers.NewClient(providers.Provider(account.Provider))
	if err != nil {
		return fmt.Errorf("failed to create provider client: %w", err)
	}

	// Execute action based on type
	if err := s.executeProviderAction(client, token, selector, action); err != nil {
		return err
	}

	// Invalidate cache for this account
	if err := s.invalidateCache(ctx, accountID); err != nil {
		// Log error but don't fail the request
		_ = err
	}

	return nil
}

// RefreshDevices forces a cache refresh for an account
func (s *DeviceService) RefreshDevices(ctx context.Context, userID, accountID string) ([]*models.Device, error) {
	// Get account and verify ownership
	account, err := s.accountRepo.FindByIDString(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	if account.OwnerUserID.String() != userID {
		return nil, fmt.Errorf("unauthorized: user does not own this account")
	}

	// Invalidate cache
	if invalidateErr := s.invalidateCache(ctx, accountID); invalidateErr != nil {
		// Log error but continue
		_ = invalidateErr
	}

	// Fetch fresh data from provider
	devices, err := s.fetchDevicesFromProvider(ctx, account)
	if err != nil {
		return nil, err
	}

	// Cache the devices
	if err := s.setCachedDevices(ctx, accountID, devices); err != nil {
		// Log error but continue
		_ = err
	}

	return devices, nil
}

// --- Private helper methods ---

// fetchDevicesFromProvider fetches devices from the provider API
func (s *DeviceService) fetchDevicesFromProvider(ctx context.Context, account *models.Account) ([]*models.Device, error) {
	// Check rate limit
	if err := s.checkRateLimit(ctx, account.ID.String()); err != nil {
		return nil, err
	}

	// Get decrypted token
	token, err := s.accountRepo.GetDecryptedToken(ctx, account.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Create provider client
	client, err := providers.NewClient(providers.Provider(account.Provider))
	if err != nil {
		return nil, fmt.Errorf("failed to create provider client: %w", err)
	}

	// Get devices from provider
	providerDevices, err := client.ListDevices(token)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices from provider: %w", err)
	}

	// Convert to our device model
	devices := make([]*models.Device, len(providerDevices))
	for i, pd := range providerDevices {
		devices[i] = s.convertProviderDevice(pd, account.ID.String(), account.Provider)
	}

	return devices, nil
}

// convertProviderDevice converts a provider device to our device model
func (s *DeviceService) convertProviderDevice(pd *providers.Device, accountID, provider string) *models.Device {
	device := &models.Device{
		ID:           pd.ID,
		AccountID:    accountID,
		Provider:     provider,
		Label:        pd.Label,
		Power:        pd.Power,
		Brightness:   pd.Brightness,
		Connected:    pd.Connected,
		Reachable:    pd.Reachable,
		Capabilities: pd.Capabilities,
		Metadata:     pd.Metadata,
	}

	if pd.Color != nil {
		device.Color = &models.DeviceColor{
			Hue:        pd.Color.Hue,
			Saturation: pd.Color.Saturation,
			Kelvin:     pd.Color.Kelvin,
		}
	}

	if pd.Group != nil {
		device.Group = &models.DeviceGroup{
			ID:   pd.Group.ID,
			Name: pd.Group.Name,
		}
	}

	if pd.Location != nil {
		device.Location = &models.DeviceLocation{
			ID:   pd.Location.ID,
			Name: pd.Location.Name,
		}
	}

	return device
}

// executeProviderAction executes an action via the provider client
func (s *DeviceService) executeProviderAction(client providers.Client, token, selector string, action *models.ActionRequest) error {
	duration := action.GetDuration()

	switch action.Action {
	case models.ActionPower:
		state, err := action.GetPowerState()
		if err != nil {
			return err
		}
		return client.SetPower(token, selector, state, duration)

	case models.ActionBrightness:
		level, err := action.GetBrightnessLevel()
		if err != nil {
			return err
		}
		return client.SetBrightness(token, selector, level, duration)

	case models.ActionColor:
		hue, _ := action.Parameters["hue"].(float64)
		saturation, _ := action.Parameters["saturation"].(float64)
		kelvin := 3500 // Default kelvin value
		if k, ok := action.Parameters["kelvin"].(float64); ok {
			kelvin = int(k)
		}
		color := &providers.DeviceColor{
			Hue:        hue,
			Saturation: saturation,
			Kelvin:     kelvin,
		}
		return client.SetColor(token, selector, color, duration)

	case models.ActionTemperature:
		kelvin, _ := action.Parameters["kelvin"].(float64)
		return client.SetColorTemperature(token, selector, int(kelvin), duration)

	case models.ActionEffect:
		name, _ := action.Parameters["name"].(string)
		cycles := 3 // Default cycles
		if c, ok := action.Parameters["cycles"].(float64); ok {
			cycles = int(c)
		}
		period := 1.0 // Default period (seconds)
		if p, ok := action.Parameters["period"].(float64); ok {
			period = p
		}

		var color *providers.DeviceColor
		if colorData, ok := action.Parameters["color"].(map[string]interface{}); ok {
			hue, _ := colorData["hue"].(float64)
			saturation, _ := colorData["saturation"].(float64)
			kelvin := 3500
			if k, ok := colorData["kelvin"].(float64); ok {
				kelvin = int(k)
			}
			color = &providers.DeviceColor{
				Hue:        hue,
				Saturation: saturation,
				Kelvin:     kelvin,
			}
		}

		switch name {
		case models.EffectPulse:
			return client.Pulse(token, selector, color, cycles, period)
		case models.EffectBreathe:
			return client.Breathe(token, selector, color, cycles, period)
		default:
			return fmt.Errorf("unknown effect: %s", name)
		}

	default:
		return fmt.Errorf("unknown action: %s", action.Action)
	}
}

// getCachedDevices retrieves devices from cache
func (s *DeviceService) getCachedDevices(ctx context.Context, accountID string) ([]*models.Device, error) {
	key := fmt.Sprintf("devices:account:%s", accountID)
	data, err := s.cache.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var devices []*models.Device
	if err := json.Unmarshal(data, &devices); err != nil {
		return nil, err
	}

	return devices, nil
}

// setCachedDevices stores devices in cache
func (s *DeviceService) setCachedDevices(ctx context.Context, accountID string, devices []*models.Device) error {
	key := fmt.Sprintf("devices:account:%s", accountID)
	data, err := json.Marshal(devices)
	if err != nil {
		return err
	}

	return s.cache.Set(ctx, key, data, s.cacheTTL).Err()
}

// invalidateCache removes devices from cache
func (s *DeviceService) invalidateCache(ctx context.Context, accountID string) error {
	key := fmt.Sprintf("devices:account:%s", accountID)
	return s.cache.Del(ctx, key).Err()
}

// checkRateLimit checks if the account has exceeded the rate limit
func (s *DeviceService) checkRateLimit(ctx context.Context, accountID string) error {
	key := fmt.Sprintf("ratelimit:account:%s", accountID)

	// Increment counter
	count, err := s.cache.Incr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	// Set expiry on first request
	if count == 1 {
		s.cache.Expire(ctx, key, 60*time.Second)
	}

	// Check limit
	if count > int64(s.rateLimitPerMin) {
		return fmt.Errorf("rate limit exceeded: max %d requests per minute", s.rateLimitPerMin)
	}

	return nil
}
