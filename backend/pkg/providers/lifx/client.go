// Package lifx provides a client for interacting with the LIFX API
package lifx

import (
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
	ID         string  `json:"id"`
	UUID       string  `json:"uuid"`
	Label      string  `json:"label"`
	Power      string  `json:"power"`
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
