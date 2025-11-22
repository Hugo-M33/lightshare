package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/lightshare/backend/internal/models"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked")
)

// RefreshTokenRepository handles refresh token database operations
type RefreshTokenRepository struct {
	db *sqlx.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *sqlx.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time, userAgent, ipAddress *string) (*models.RefreshToken, error) {
	token := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}

	query := `
		INSERT INTO refresh_tokens (
			id, user_id, token_hash, expires_at, created_at, user_agent, ip_address
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, user_id, token_hash, expires_at, created_at, revoked_at, user_agent, ip_address
	`

	err := r.db.GetContext(ctx, token, query,
		token.ID, token.UserID, token.TokenHash, token.ExpiresAt,
		token.CreatedAt, token.UserAgent, token.IPAddress,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return token, nil
}

// GetByTokenHash retrieves a refresh token by token hash
func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked_at, user_agent, ip_address
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	err := r.db.GetContext(ctx, &token, query, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// Check if token is revoked
	if token.RevokedAt != nil {
		return nil, ErrRefreshTokenRevoked
	}

	// Check if token is expired
	if token.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return &token, nil
}

// Revoke revokes a refresh token
func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	now := time.Now()
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE token_hash = $2
	`

	result, err := r.db.ExecContext(ctx, query, now, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE user_id = $2 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, now, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all refresh tokens: %w", err)
	}

	return nil
}

// DeleteExpired deletes all expired refresh tokens
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < $1 OR revoked_at < $1
	`

	// Delete tokens expired or revoked more than 7 days ago
	cutoff := time.Now().AddDate(0, 0, -7)
	_, err := r.db.ExecContext(ctx, query, cutoff)
	if err != nil {
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}

	return nil
}
