package models

// Device represents a smart light device from any provider (LIFX, Hue, etc.)
type Device struct {
	ID           string                 `json:"id"`                 // Provider-specific device ID
	AccountID    string                 `json:"account_id"`         // Our account UUID
	Provider     string                 `json:"provider"`           // "lifx" or "hue"
	Label        string                 `json:"label"`              // User-friendly name
	Power        string                 `json:"power"`              // "on" or "off"
	Brightness   float64                `json:"brightness"`         // 0.0 - 1.0
	Color        *DeviceColor           `json:"color,omitempty"`    // Color information (if supported)
	Connected    bool                   `json:"connected"`          // Whether device is connected to network
	Reachable    bool                   `json:"reachable"`          // Whether device is reachable by cloud API
	Group        *DeviceGroup           `json:"group,omitempty"`    // Group/room information
	Location     *DeviceLocation        `json:"location,omitempty"` // Location/home information
	Capabilities []string               `json:"capabilities"`       // ["color", "temperature", "effects"]
	Metadata     map[string]interface{} `json:"metadata,omitempty"` // Provider-specific metadata
}

// DeviceColor represents the color state of a device
type DeviceColor struct {
	Hue        float64 `json:"hue"`        // 0-360 degrees
	Saturation float64 `json:"saturation"` // 0.0-1.0
	Kelvin     int     `json:"kelvin"`     // 1500-9000 (color temperature for white balance)
}

// DeviceGroup represents a group/room that devices belong to
type DeviceGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DeviceLocation represents the location/home that devices belong to
type DeviceLocation struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// IsOn returns true if the device is powered on
func (d *Device) IsOn() bool {
	return d.Power == "on"
}

// HasCapability checks if the device supports a specific capability
func (d *Device) HasCapability(capability string) bool {
	for _, cap := range d.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}

// SupportsColor returns true if the device supports color control
func (d *Device) SupportsColor() bool {
	return d.HasCapability("color")
}

// SupportsTemperature returns true if the device supports color temperature control
func (d *Device) SupportsTemperature() bool {
	return d.HasCapability("temperature")
}

// SupportsEffects returns true if the device supports effects (LIFX-specific)
func (d *Device) SupportsEffects() bool {
	return d.HasCapability("effects")
}
