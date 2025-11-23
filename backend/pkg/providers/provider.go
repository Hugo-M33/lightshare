// Package providers defines the common interface for smart lighting providers
// and implements the factory pattern for creating provider clients.
package providers

import (
	"fmt"

	"github.com/lightshare/backend/pkg/providers/lifx"
)

// Provider represents the type of smart lighting provider
type Provider string

// Supported provider types
const (
	// ProviderLIFX represents the LIFX smart lighting provider
	ProviderLIFX Provider = "lifx"
	// ProviderHue represents the Philips Hue smart lighting provider
	ProviderHue Provider = "hue"
)

// IsValid checks if the provider type is valid
func (p Provider) IsValid() bool {
	return p == ProviderLIFX || p == ProviderHue
}

// String returns the string representation of the provider
func (p Provider) String() string {
	return string(p)
}

// AccountInfo contains information about a provider account
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

// Device represents a smart light device (unified across providers)
type Device struct {
	ID           string
	Label        string
	Power        string  // "on" or "off"
	Brightness   float64 // 0.0-1.0
	Color        *DeviceColor
	Connected    bool
	Reachable    bool
	Group        *DeviceGroup
	Location     *DeviceLocation
	Capabilities []string
	Metadata     map[string]interface{}
}

// DeviceColor represents color information for a device
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

// Client defines the interface that all provider clients must implement
type Client interface {
	// ValidateToken validates the token by making a test API call
	// Returns AccountInfo if valid, error otherwise
	ValidateToken(token string) (*AccountInfo, error)

	// GetAccountInfo retrieves account information using the token
	GetAccountInfo(token string) (*AccountInfo, error)

	// --- Phase 4: Device Control Methods ---

	// ListDevices returns all lights/devices for the account
	ListDevices(token string) ([]*Device, error)

	// GetDevice returns a specific device by ID
	GetDevice(token, deviceID string) (*Device, error)

	// SetPower turns device(s) on or off
	// selector: "all", "id:d073d5", "group_id:xxx", "location_id:xxx"
	// state: true for on, false for off
	// duration: transition time in seconds
	SetPower(token, selector string, state bool, duration float64) error

	// SetBrightness adjusts device brightness
	// level: 0.0-1.0
	// duration: transition time in seconds
	SetBrightness(token, selector string, level float64, duration float64) error

	// SetColor sets device color (hue/saturation)
	// duration: transition time in seconds
	SetColor(token, selector string, color *DeviceColor, duration float64) error

	// SetColorTemperature sets white balance
	// kelvin: 1500-9000
	// duration: transition time in seconds
	SetColorTemperature(token, selector string, kelvin int, duration float64) error

	// --- Effects (LIFX-specific, will return error for Hue) ---

	// Pulse creates a pulsing effect
	// cycles: number of times to pulse
	// period: time for one cycle in seconds
	Pulse(token, selector string, color *DeviceColor, cycles int, period float64) error

	// Breathe creates a breathing effect
	// cycles: number of times to breathe
	// period: time for one cycle in seconds
	Breathe(token, selector string, color *DeviceColor, cycles int, period float64) error
}

// lifxClientAdapter adapts the LIFX client to the Client interface
type lifxClientAdapter struct {
	client *lifx.Client
}

func (a *lifxClientAdapter) ValidateToken(token string) (*AccountInfo, error) {
	info, err := a.client.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return &AccountInfo{
		ProviderAccountID: info.ProviderAccountID,
		Email:             info.Email,
		Label:             info.Label,
		Metadata:          info.Metadata,
	}, nil
}

func (a *lifxClientAdapter) GetAccountInfo(token string) (*AccountInfo, error) {
	info, err := a.client.GetAccountInfo(token)
	if err != nil {
		return nil, err
	}
	return &AccountInfo{
		ProviderAccountID: info.ProviderAccountID,
		Email:             info.Email,
		Label:             info.Label,
		Metadata:          info.Metadata,
	}, nil
}

// ListDevices returns all devices for the account
func (a *lifxClientAdapter) ListDevices(token string) ([]*Device, error) {
	lifxDevices, err := a.client.ListDevices(token)
	if err != nil {
		return nil, err
	}

	devices := make([]*Device, len(lifxDevices))
	for i, d := range lifxDevices {
		devices[i] = convertLIFXDevice(d)
	}
	return devices, nil
}

// GetDevice returns a specific device by ID
func (a *lifxClientAdapter) GetDevice(token, deviceID string) (*Device, error) {
	lifxDevice, err := a.client.GetDevice(token, deviceID)
	if err != nil {
		return nil, err
	}
	return convertLIFXDevice(lifxDevice), nil
}

// SetPower turns device(s) on or off
func (a *lifxClientAdapter) SetPower(token, selector string, state bool, duration float64) error {
	return a.client.SetPower(token, selector, state, duration)
}

// SetBrightness adjusts device brightness
func (a *lifxClientAdapter) SetBrightness(token, selector string, level float64, duration float64) error {
	return a.client.SetBrightness(token, selector, level, duration)
}

// SetColor sets device color
func (a *lifxClientAdapter) SetColor(token, selector string, color *DeviceColor, duration float64) error {
	lifxColor := &lifx.DeviceColor{
		Hue:        color.Hue,
		Saturation: color.Saturation,
		Kelvin:     color.Kelvin,
	}
	return a.client.SetColor(token, selector, lifxColor, duration)
}

// SetColorTemperature sets white balance
func (a *lifxClientAdapter) SetColorTemperature(token, selector string, kelvin int, duration float64) error {
	return a.client.SetColorTemperature(token, selector, kelvin, duration)
}

// Pulse creates a pulsing effect
func (a *lifxClientAdapter) Pulse(token, selector string, color *DeviceColor, cycles int, period float64) error {
	var lifxColor *lifx.DeviceColor
	if color != nil {
		lifxColor = &lifx.DeviceColor{
			Hue:        color.Hue,
			Saturation: color.Saturation,
			Kelvin:     color.Kelvin,
		}
	}
	return a.client.Pulse(token, selector, lifxColor, cycles, period)
}

// Breathe creates a breathing effect
func (a *lifxClientAdapter) Breathe(token, selector string, color *DeviceColor, cycles int, period float64) error {
	var lifxColor *lifx.DeviceColor
	if color != nil {
		lifxColor = &lifx.DeviceColor{
			Hue:        color.Hue,
			Saturation: color.Saturation,
			Kelvin:     color.Kelvin,
		}
	}
	return a.client.Breathe(token, selector, lifxColor, cycles, period)
}

// convertLIFXDevice converts a LIFX device to the generic Device type
func convertLIFXDevice(d *lifx.Device) *Device {
	device := &Device{
		ID:           d.ID,
		Label:        d.Label,
		Power:        d.Power,
		Brightness:   d.Brightness,
		Connected:    d.Connected,
		Reachable:    d.Reachable,
		Capabilities: d.Capabilities,
		Metadata:     d.Metadata,
	}

	if d.Color != nil {
		device.Color = &DeviceColor{
			Hue:        d.Color.Hue,
			Saturation: d.Color.Saturation,
			Kelvin:     d.Color.Kelvin,
		}
	}

	if d.Group != nil {
		device.Group = &DeviceGroup{
			ID:   d.Group.ID,
			Name: d.Group.Name,
		}
	}

	if d.Location != nil {
		device.Location = &DeviceLocation{
			ID:   d.Location.ID,
			Name: d.Location.Name,
		}
	}

	return device
}

// NewClient creates a new provider client based on the provider type
func NewClient(provider Provider) (Client, error) {
	switch provider {
	case ProviderLIFX:
		return &lifxClientAdapter{client: lifx.NewClient()}, nil
	case ProviderHue:
		return nil, fmt.Errorf("hue provider not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
