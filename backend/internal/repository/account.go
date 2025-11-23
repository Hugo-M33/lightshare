package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/lightshare/backend/internal/models"
	"github.com/lightshare/backend/pkg/crypto"
)

var (
	// ErrAccountNotFound is returned when an account is not found in the database
	ErrAccountNotFound = errors.New("account not found")
	// ErrAccountAlreadyExists is returned when attempting to create a duplicate account
	ErrAccountAlreadyExists = errors.New("account already exists for this provider")
)

// AccountRepositoryInterface defines the interface for account repository operations
type AccountRepositoryInterface interface {
	Create(ctx context.Context, params *models.CreateAccountParams) (*models.Account, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Account, error)
	FindByID(ctx context.Context, accountID uuid.UUID) (*models.Account, error)
	Delete(ctx context.Context, accountID, userID uuid.UUID) error
}

// AccountRepository handles account database operations
type AccountRepository struct {
	db *sqlx.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create creates a new account
func (r *AccountRepository) Create(ctx context.Context, params *models.CreateAccountParams) (*models.Account, error) {
	account := &models.Account{
		ID:                uuid.New(),
		OwnerUserID:       params.OwnerUserID,
		Provider:          params.Provider,
		ProviderAccountID: params.ProviderAccountID,
		EncryptedToken:    params.EncryptedToken,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Serialize metadata to JSON if present
	if params.Metadata != nil {
		metadataJSON, err := json.Marshal(params.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		account.Metadata = metadataJSON
	}

	query := `
		INSERT INTO accounts (
			id, owner_user_id, provider, provider_account_id,
			encrypted_token, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		RETURNING id, owner_user_id, provider, provider_account_id,
			encrypted_token, metadata, created_at, updated_at
	`

	err := r.db.GetContext(ctx, account, query,
		account.ID, account.OwnerUserID, account.Provider, account.ProviderAccountID,
		account.EncryptedToken, account.Metadata, account.CreatedAt, account.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "pq: duplicate key value violates unique constraint \"accounts_owner_user_id_provider_provider_account_id_key\"" {
			return nil, ErrAccountAlreadyExists
		}
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

// FindByUserID retrieves all accounts for a user
func (r *AccountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Account, error) {
	var accounts []*models.Account
	query := `
		SELECT id, owner_user_id, provider, provider_account_id,
			encrypted_token, metadata, created_at, updated_at
		FROM accounts
		WHERE owner_user_id = $1
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &accounts, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find accounts by user id: %w", err)
	}

	return accounts, nil
}

// FindByID retrieves a specific account by ID
func (r *AccountRepository) FindByID(ctx context.Context, accountID uuid.UUID) (*models.Account, error) {
	var account models.Account
	query := `
		SELECT id, owner_user_id, provider, provider_account_id,
			encrypted_token, metadata, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &account, query, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, fmt.Errorf("failed to find account by id: %w", err)
	}

	return &account, nil
}

// Delete deletes an account
func (r *AccountRepository) Delete(ctx context.Context, accountID, userID uuid.UUID) error {
	query := `
		DELETE FROM accounts
		WHERE id = $1 AND owner_user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, accountID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAccountNotFound
	}

	return nil
}

// FindByIDString retrieves an account by string ID (convenience method for Phase 4)
func (r *AccountRepository) FindByIDString(ctx context.Context, accountID string) (*models.Account, error) {
	id, err := uuid.Parse(accountID)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}
	return r.FindByID(ctx, id)
}

// GetDecryptedToken retrieves and decrypts the token for an account
func (r *AccountRepository) GetDecryptedToken(ctx context.Context, accountID string) (string, error) {
	account, err := r.FindByIDString(ctx, accountID)
	if err != nil {
		return "", err
	}

	// Decrypt the token
	token, err := crypto.DecryptToken(account.EncryptedToken)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt token: %w", err)
	}

	return token, nil
}
