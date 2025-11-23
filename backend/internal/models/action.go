package models

import (
	"fmt"
)

// ActionRequest represents a control action request from the client
type ActionRequest struct {
	Parameters map[string]interface{} `json:"parameters" validate:"required"`
	Action     string                 `json:"action" validate:"required"`
}

// Supported action types
const (
	ActionPower       = "power"       // Turn on/off
	ActionBrightness  = "brightness"  // Adjust brightness
	ActionColor       = "color"       // Set color (hue/saturation)
	ActionTemperature = "temperature" // Set color temperature (kelvin)
	ActionEffect      = "effect"      // Trigger effect (pulse, breathe, etc.)
)

// Supported effect names
const (
	EffectPulse   = "pulse"
	EffectBreathe = "breathe"
)

// IsValidAction checks if the action type is supported
func (a *ActionRequest) IsValidAction() bool {
	switch a.Action {
	case ActionPower, ActionBrightness, ActionColor, ActionTemperature, ActionEffect:
		return true
	default:
		return false
	}
}

// ValidateParameters validates that required parameters are present and valid
func (a *ActionRequest) ValidateParameters() error {
	if !a.IsValidAction() {
		return fmt.Errorf("invalid action type: %s", a.Action)
	}

	switch a.Action {
	case ActionPower:
		return a.validatePowerParameters()
	case ActionBrightness:
		return a.validateBrightnessParameters()
	case ActionColor:
		return a.validateColorParameters()
	case ActionTemperature:
		return a.validateTemperatureParameters()
	case ActionEffect:
		return a.validateEffectParameters()
	default:
		return fmt.Errorf("unknown action: %s", a.Action)
	}
}

func (a *ActionRequest) validatePowerParameters() error {
	state, ok := a.Parameters["state"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'state' parameter (must be string)")
	}
	if state != PowerStateOn && state != PowerStateOff {
		return fmt.Errorf("invalid state value: %s (must be 'on' or 'off')", state)
	}
	return nil
}

func (a *ActionRequest) validateBrightnessParameters() error {
	level, ok := a.Parameters["level"].(float64)
	if !ok {
		return fmt.Errorf("missing or invalid 'level' parameter (must be number)")
	}
	if level < 0.0 || level > 1.0 {
		return fmt.Errorf("invalid brightness level: %f (must be 0.0-1.0)", level)
	}
	return nil
}

func (a *ActionRequest) validateColorParameters() error {
	hue, hueOk := a.Parameters["hue"].(float64)
	saturation, satOk := a.Parameters["saturation"].(float64)

	if !hueOk {
		return fmt.Errorf("missing or invalid 'hue' parameter (must be number)")
	}
	if !satOk {
		return fmt.Errorf("missing or invalid 'saturation' parameter (must be number)")
	}

	if hue < 0.0 || hue > 360.0 {
		return fmt.Errorf("invalid hue value: %f (must be 0-360)", hue)
	}
	if saturation < 0.0 || saturation > 1.0 {
		return fmt.Errorf("invalid saturation value: %f (must be 0.0-1.0)", saturation)
	}

	return nil
}

func (a *ActionRequest) validateTemperatureParameters() error {
	kelvin, ok := a.Parameters["kelvin"].(float64)
	if !ok {
		return fmt.Errorf("missing or invalid 'kelvin' parameter (must be number)")
	}
	if kelvin < 1500 || kelvin > 9000 {
		return fmt.Errorf("invalid kelvin value: %f (must be 1500-9000)", kelvin)
	}
	return nil
}

func (a *ActionRequest) validateEffectParameters() error {
	name, ok := a.Parameters["name"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'name' parameter (must be string)")
	}

	if name != EffectPulse && name != EffectBreathe {
		return fmt.Errorf("invalid effect name: %s (must be 'pulse' or 'breathe')", name)
	}

	// Color is optional for effects, but if provided should be valid
	if colorData, hasColor := a.Parameters["color"].(map[string]interface{}); hasColor {
		if hue, hueOk := colorData["hue"].(float64); hueOk {
			if hue < 0.0 || hue > 360.0 {
				return fmt.Errorf("invalid effect color hue: %f (must be 0-360)", hue)
			}
		}
		if sat, satOk := colorData["saturation"].(float64); satOk {
			if sat < 0.0 || sat > 1.0 {
				return fmt.Errorf("invalid effect color saturation: %f (must be 0.0-1.0)", sat)
			}
		}
	}

	return nil
}

// GetPowerState returns the desired power state for power actions
func (a *ActionRequest) GetPowerState() (bool, error) {
	if a.Action != ActionPower {
		return false, fmt.Errorf("not a power action")
	}
	state, ok := a.Parameters["state"].(string)
	if !ok {
		return false, fmt.Errorf("invalid state parameter")
	}
	return state == PowerStateOn, nil
}

// GetBrightnessLevel returns the brightness level for brightness actions
func (a *ActionRequest) GetBrightnessLevel() (float64, error) {
	if a.Action != ActionBrightness {
		return 0, fmt.Errorf("not a brightness action")
	}
	level, ok := a.Parameters["level"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid level parameter")
	}
	return level, nil
}

// GetDuration returns the duration parameter (optional, defaults to 0.5 seconds)
func (a *ActionRequest) GetDuration() float64 {
	if duration, ok := a.Parameters["duration"].(float64); ok {
		return duration
	}
	return 0.5 // Default transition duration
}
