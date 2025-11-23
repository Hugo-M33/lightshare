// Package lifx provides a client for interacting with the LIFX API
package lifx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	lifxAPIBaseURL = "https://api.lifx.com/v1"
	requestTimeout = 10 * time.Second
)

// AccountInfo contains information about a LIFX account
type AccountInfo struct {
	// Additional metadata
	Metadata map[string]interface{}
	// ProviderAccountID is the unique identifier from the provider
	ProviderAccountID string
	// Email associated with the account (if available)
	Email string
	// Label or name for the account
	Label string
}

// Client implements the Client interface for LIFX
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new LIFX client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

// LightsResponse represents the response from LIFX list lights endpoint
type LightsResponse []struct {
	Group struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"group"`
	Location struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"location"`
	ID    string `json:"id"`
	UUID  string `json:"uuid"`
	Label string `json:"label"`
	Power string `json:"power"`
	Color struct {
		Hue        float64 `json:"hue"`
		Saturation float64 `json:"saturation"`
		Kelvin     int     `json:"kelvin"`
	} `json:"color"`
	Brightness float64 `json:"brightness"`
	Connected  bool    `json:"connected"`
}

// ValidateToken validates the LIFX token by attempting to list lights
// This confirms the token is valid and has the necessary permissions
func (c *Client) ValidateToken(token string) (*AccountInfo, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("%s/lights/all", lifxAPIBaseURL), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call LIFX API: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error but don't override return error
			_ = closeErr
		}
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid token: unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var lights LightsResponse
	if err := json.NewDecoder(resp.Body).Decode(&lights); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// For LIFX, we use the first light's location ID as the account identifier
	// This helps distinguish between different LIFX accounts
	accountID := "lifx-account"
	accountLabel := "LIFX Account"

	if len(lights) > 0 && lights[0].Location.ID != "" {
		accountID = lights[0].Location.ID
		if lights[0].Location.Name != "" {
			accountLabel = lights[0].Location.Name
		}
	}

	return &AccountInfo{
		ProviderAccountID: accountID,
		Label:             accountLabel,
		Metadata: map[string]interface{}{
			"lights_count": len(lights),
		},
	}, nil
}

// GetAccountInfo retrieves account information for the LIFX account
// For LIFX, this is similar to ValidateToken since LIFX doesn't have a dedicated account info endpoint
func (c *Client) GetAccountInfo(token string) (*AccountInfo, error) {
	return c.ValidateToken(token)
}

// --- Phase 4: Device Control Implementation ---

// Device represents a LIFX light device
type Device struct {
	Color        *DeviceColor
	Group        *DeviceGroup
	Location     *DeviceLocation
	Metadata     map[string]interface{}
	ID           string
	Label        string
	Power        string
	Capabilities []string
	Brightness   float64
	Connected    bool
	Reachable    bool
}

// DeviceColor represents color information
type DeviceColor struct {
	Hue        float64 // 0-360
	Saturation float64 // 0.0-1.0
	Kelvin     int     // 1500-9000
}

// DeviceGroup represents a group/room
type DeviceGroup struct {
	ID   string
	Name string
}

// DeviceLocation represents a location/home
type DeviceLocation struct {
	ID   string
	Name string
}

// ListDevices returns all lights for the LIFX account
func (c *Client) ListDevices(token string) ([]*Device, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("%s/lights/all", lifxAPIBaseURL), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call LIFX API: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			_ = closeErr
		}
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid token: unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var lights LightsResponse
	if err := json.NewDecoder(resp.Body).Decode(&lights); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert LIFX response to unified Device format
	devices := make([]*Device, 0, len(lights))
	for _, light := range lights {
		capabilities := []string{}

		// All LIFX lights support brightness
		capabilities = append(capabilities, "brightness")

		// Check if device supports color
		// LIFX Color and LIFX+ support color, LIFX White doesn't
		if light.Color.Saturation > 0 || light.Label != "" {
			// Most LIFX devices support color
			capabilities = append(capabilities, "color")
		}

		// All LIFX lights support color temperature
		capabilities = append(capabilities, "temperature")

		// All LIFX lights support effects
		capabilities = append(capabilities, "effects")

		device := &Device{
			ID:         light.ID,
			Label:      light.Label,
			Power:      light.Power,
			Brightness: light.Brightness,
			Color: &DeviceColor{
				Hue:        light.Color.Hue,
				Saturation: light.Color.Saturation,
				Kelvin:     light.Color.Kelvin,
			},
			Connected:    light.Connected,
			Reachable:    light.Connected, // For LIFX, connected implies reachable
			Capabilities: capabilities,
		}

		if light.Group.ID != "" {
			device.Group = &DeviceGroup{
				ID:   light.Group.ID,
				Name: light.Group.Name,
			}
		}

		if light.Location.ID != "" {
			device.Location = &DeviceLocation{
				ID:   light.Location.ID,
				Name: light.Location.Name,
			}
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// GetDevice returns a specific light by ID
func (c *Client) GetDevice(token, deviceID string) (*Device, error) {
	selector := fmt.Sprintf("id:%s", deviceID)
	req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("%s/lights/%s", lifxAPIBaseURL, selector), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call LIFX API: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			_ = closeErr
		}
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid token: unauthorized")
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var lights LightsResponse
	if err := json.NewDecoder(resp.Body).Decode(&lights); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(lights) == 0 {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}

	light := lights[0]
	capabilities := []string{"brightness", "color", "temperature", "effects"}

	device := &Device{
		ID:           light.ID,
		Label:        light.Label,
		Power:        light.Power,
		Brightness:   light.Brightness,
		Color:        &DeviceColor{Hue: light.Color.Hue, Saturation: light.Color.Saturation, Kelvin: light.Color.Kelvin},
		Connected:    light.Connected,
		Reachable:    light.Connected,
		Capabilities: capabilities,
	}

	if light.Group.ID != "" {
		device.Group = &DeviceGroup{ID: light.Group.ID, Name: light.Group.Name}
	}
	if light.Location.ID != "" {
		device.Location = &DeviceLocation{ID: light.Location.ID, Name: light.Location.Name}
	}

	return device, nil
}

// SetPower turns lights on or off
func (c *Client) SetPower(token, selector string, state bool, duration float64) error {
	powerState := "off"
	if state {
		powerState = "on"
	}

	body := map[string]interface{}{
		"power":    powerState,
		"duration": duration,
	}

	return c.setState(token, selector, body)
}

// SetBrightness adjusts the brightness level
func (c *Client) SetBrightness(token, selector string, level float64, duration float64) error {
	body := map[string]interface{}{
		"brightness": level,
		"duration":   duration,
	}

	return c.setState(token, selector, body)
}

// SetColor sets the hue and saturation
func (c *Client) SetColor(token, selector string, color *DeviceColor, duration float64) error {
	// LIFX uses a string format: "hue:120 saturation:1.0"
	colorString := fmt.Sprintf("hue:%f saturation:%f", color.Hue, color.Saturation)

	body := map[string]interface{}{
		"color":    colorString,
		"duration": duration,
	}

	return c.setState(token, selector, body)
}

// SetColorTemperature sets the white balance
func (c *Client) SetColorTemperature(token, selector string, kelvin int, duration float64) error {
	colorString := fmt.Sprintf("kelvin:%d", kelvin)

	body := map[string]interface{}{
		"color":    colorString,
		"duration": duration,
	}

	return c.setState(token, selector, body)
}

// Pulse creates a pulsing effect
func (c *Client) Pulse(token, selector string, color *DeviceColor, cycles int, period float64) error {
	body := map[string]interface{}{
		"cycles": cycles,
		"period": period,
	}

	if color != nil {
		colorString := fmt.Sprintf("hue:%f saturation:%f", color.Hue, color.Saturation)
		body["color"] = colorString
	}

	return c.postEffect(token, selector, "pulse", body)
}

// Breathe creates a breathing effect
func (c *Client) Breathe(token, selector string, color *DeviceColor, cycles int, period float64) error {
	body := map[string]interface{}{
		"cycles": cycles,
		"period": period,
	}

	if color != nil {
		colorString := fmt.Sprintf("hue:%f saturation:%f", color.Hue, color.Saturation)
		body["color"] = colorString
	}

	return c.postEffect(token, selector, "breathe", body)
}

// setState is a helper method to set state on lights
func (c *Client) setState(token, selector string, body map[string]interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("%s/lights/%s/state", lifxAPIBaseURL, selector)
	req, err := http.NewRequestWithContext(context.Background(), "PUT", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call LIFX API: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			_ = closeErr
		}
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid token: unauthorized")
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("selector not found: %s", selector)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusMultiStatus {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// postEffect is a helper method to trigger effects
func (c *Client) postEffect(token, selector, effect string, body map[string]interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("%s/lights/%s/effects/%s", lifxAPIBaseURL, selector, effect)
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call LIFX API: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			_ = closeErr
		}
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid token: unauthorized")
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("selector not found: %s", selector)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusMultiStatus {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
