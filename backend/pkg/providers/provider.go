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

// Client defines the interface that all provider clients must implement
type Client interface {
	// ValidateToken validates the token by making a test API call
	// Returns AccountInfo if valid, error otherwise
	ValidateToken(token string) (*AccountInfo, error)

	// GetAccountInfo retrieves account information using the token
	GetAccountInfo(token string) (*AccountInfo, error)
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
