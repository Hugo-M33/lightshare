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
	// ErrUserNotFound is returned when a user is not found in the database.
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists is returned when attempting to create a user with an email that already exists.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrTokenExpired is returned when a verification or magic link token has expired.
	ErrTokenExpired = errors.New("token expired")
	// ErrTokenNotFound is returned when a token is not found in the database.
	ErrTokenNotFound = errors.New("token not found")
)

// UserRepository handles user database operations
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, params models.CreateUserParams) (*models.User, error) {
	user := &models.User{
		ID:                         uuid.New(),
		Email:                      params.Email,
		PasswordHash:               params.PasswordHash,
		EmailVerified:              false,
		EmailVerificationToken:     &params.EmailVerificationToken,
		EmailVerificationExpiresAt: &params.EmailVerificationExpiresAt,
		Role:                       "user",
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}

	query := `
		INSERT INTO users (
			id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			role, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		RETURNING id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			magic_link_token, magic_link_expires_at,
			stripe_customer_id, role, created_at, updated_at
	`

	err := r.db.GetContext(ctx, user, query,
		user.ID, user.Email, user.PasswordHash, user.EmailVerified,
		user.EmailVerificationToken, user.EmailVerificationExpiresAt,
		user.Role, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return nil, ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			magic_link_token, magic_link_expires_at,
			stripe_customer_id, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			magic_link_token, magic_link_expires_at,
			stripe_customer_id, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetByEmailVerificationToken retrieves a user by email verification token
func (r *UserRepository) GetByEmailVerificationToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			magic_link_token, magic_link_expires_at,
			stripe_customer_id, role, created_at, updated_at
		FROM users
		WHERE email_verification_token = $1
			AND email_verification_expires_at > $2
	`

	err := r.db.GetContext(ctx, &user, query, token, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("failed to get user by verification token: %w", err)
	}

	return &user, nil
}

// VerifyEmail verifies a user's email using the verification token
func (r *UserRepository) VerifyEmail(ctx context.Context, token string) error {
	query := `
		UPDATE users
		SET email_verified = true,
			email_verification_token = NULL,
			email_verification_expires_at = NULL,
			updated_at = $1
		WHERE email_verification_token = $2
			AND email_verification_expires_at > $1
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), token)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTokenExpired
	}

	return nil
}

// SetMagicLinkToken sets a magic link token for password-less login
func (r *UserRepository) SetMagicLinkToken(ctx context.Context, email, token string, expiresAt time.Time) error {
	query := `
		UPDATE users
		SET magic_link_token = $1,
			magic_link_expires_at = $2,
			updated_at = $3
		WHERE email = $4
	`

	result, err := r.db.ExecContext(ctx, query, token, expiresAt, time.Now(), email)
	if err != nil {
		return fmt.Errorf("failed to set magic link token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// GetByMagicLinkToken retrieves a user by magic link token
func (r *UserRepository) GetByMagicLinkToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			magic_link_token, magic_link_expires_at,
			stripe_customer_id, role, created_at, updated_at
		FROM users
		WHERE magic_link_token = $1
			AND magic_link_expires_at > $2
	`

	err := r.db.GetContext(ctx, &user, query, token, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("failed to get user by magic link token: %w", err)
	}

	return &user, nil
}

// ClearMagicLinkToken clears the magic link token after use
func (r *UserRepository) ClearMagicLinkToken(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET magic_link_token = NULL,
			magic_link_expires_at = NULL,
			updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to clear magic link token: %w", err)
	}

	return nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET email = $1,
			password_hash = $2,
			email_verified = $3,
			email_verification_token = $4,
			email_verification_expires_at = $5,
			magic_link_token = $6,
			magic_link_expires_at = $7,
			stripe_customer_id = $8,
			role = $9,
			updated_at = $10
		WHERE id = $11
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email, user.PasswordHash, user.EmailVerified,
		user.EmailVerificationToken, user.EmailVerificationExpiresAt,
		user.MagicLinkToken, user.MagicLinkExpiresAt,
		user.StripeCustomerID, user.Role, user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
