package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Account represents a connected smart lighting provider account
type Account struct {
	CreatedAt         time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time       `db:"updated_at" json:"updated_at"`
	Provider          string          `db:"provider" json:"provider"`
	ProviderAccountID string          `db:"provider_account_id" json:"provider_account_id"`
	EncryptedToken    []byte          `db:"encrypted_token" json:"-"`
	Metadata          json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	ID                uuid.UUID       `db:"id" json:"id"`
	OwnerUserID       uuid.UUID       `db:"owner_user_id" json:"owner_user_id"`
}

// AccountResponse represents the account data sent to clients
// This excludes sensitive fields like EncryptedToken
type AccountResponse struct {
	CreatedAt         time.Time              `json:"created_at"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	Provider          string                 `json:"provider"`
	ProviderAccountID string                 `json:"provider_account_id"`
	ID                uuid.UUID              `json:"id"`
}

// ToResponse converts an Account to an AccountResponse
func (a *Account) ToResponse() *AccountResponse {
	resp := &AccountResponse{
		ID:                a.ID,
		Provider:          a.Provider,
		ProviderAccountID: a.ProviderAccountID,
		CreatedAt:         a.CreatedAt,
	}

	// Parse metadata if present
	if len(a.Metadata) > 0 {
		var metadata map[string]interface{}
		if err := json.Unmarshal(a.Metadata, &metadata); err == nil {
			resp.Metadata = metadata
		}
	}

	return resp
}

// CreateAccountParams holds parameters for creating a new account
type CreateAccountParams struct {
	Metadata          map[string]interface{}
	Provider          string
	ProviderAccountID string
	EncryptedToken    []byte
	OwnerUserID       uuid.UUID
}
