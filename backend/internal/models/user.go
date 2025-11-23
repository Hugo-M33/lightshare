// Package models defines data structures representing domain entities.
package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID                         uuid.UUID  `db:"id" json:"id"`
	CreatedAt                  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt                  time.Time  `db:"updated_at" json:"updated_at"`
	MagicLinkExpiresAt         *time.Time `db:"magic_link_expires_at" json:"-"`
	EmailVerificationExpiresAt *time.Time `db:"email_verification_expires_at" json:"-"`
	EmailVerificationToken     *string    `db:"email_verification_token" json:"-"`
	MagicLinkToken             *string    `db:"magic_link_token" json:"-"`
	StripeCustomerID           *string    `db:"stripe_customer_id" json:"stripe_customer_id,omitempty"`
	Role                       string     `db:"role" json:"role"`
	Email                      string     `db:"email" json:"email"`
	PasswordHash               string     `db:"password_hash" json:"-"`
	EmailVerified              bool       `db:"email_verified" json:"email_verified"`
}

// CreateUserParams holds parameters for creating a new user
type CreateUserParams struct {
	EmailVerificationExpiresAt time.Time
	Email                      string
	PasswordHash               string
	EmailVerificationToken     string
}

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	UserID    uuid.UUID  `db:"user_id" json:"user_id"`
	ExpiresAt time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	RevokedAt *time.Time `db:"revoked_at" json:"revoked_at,omitempty"`
	UserAgent *string    `db:"user_agent" json:"user_agent,omitempty"`
	IPAddress *string    `db:"ip_address" json:"ip_address,omitempty"`
	TokenHash string     `db:"token_hash" json:"-"`
}
