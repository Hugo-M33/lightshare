package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/lightshare/backend/internal/models"
	"github.com/lightshare/backend/internal/repository"
	"github.com/lightshare/backend/pkg/crypto"
	"github.com/lightshare/backend/pkg/providers"
)

var (
	// ErrInvalidProvider is returned when an invalid provider type is specified
	ErrInvalidProvider = errors.New("invalid provider type")
	// ErrInvalidToken is returned when a provider token is invalid
	ErrInvalidToken = errors.New("invalid provider token")
	// ErrAccountNotOwned is returned when trying to access an account not owned by the user
	ErrAccountNotOwned = errors.New("account not owned by user")
)

// ProviderService handles provider connection operations
type ProviderService struct {
	accountRepo   repository.AccountRepositoryInterface
	encryptionKey []byte
}

// NewProviderService creates a new provider service
func NewProviderService(accountRepo repository.AccountRepositoryInterface, encryptionKey []byte) *ProviderService {
	return &ProviderService{
		accountRepo:   accountRepo,
		encryptionKey: encryptionKey,
	}
}

// ConnectProviderRequest represents a request to connect a provider
type ConnectProviderRequest struct {
	Provider string `json:"provider"`
	Token    string `json:"token"`
}

// ConnectProvider validates a provider token, encrypts it, and stores the account
func (s *ProviderService) ConnectProvider(ctx context.Context, userID uuid.UUID, req ConnectProviderRequest) (*models.Account, error) {
	// Validate provider type
	providerType := providers.Provider(req.Provider)
	if !providerType.IsValid() {
		return nil, ErrInvalidProvider
	}

	// Create provider client
	client, err := providers.NewClient(providerType)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider client: %w", err)
	}

	// Validate token by calling provider API
	accountInfo, err := client.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	// Encrypt the token
	encryptedToken, err := crypto.EncryptToken(req.Token, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt token: %w", err)
	}

	// Create account
	account, err := s.accountRepo.Create(ctx, &models.CreateAccountParams{
		OwnerUserID:       userID,
		Provider:          req.Provider,
		ProviderAccountID: accountInfo.ProviderAccountID,
		EncryptedToken:    encryptedToken,
		Metadata:          accountInfo.Metadata,
	})

	if err != nil {
		if errors.Is(err, repository.ErrAccountAlreadyExists) {
			return nil, errors.New("this provider account is already connected")
		}
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

// ListAccounts returns all accounts for a user
func (s *ProviderService) ListAccounts(ctx context.Context, userID uuid.UUID) ([]*models.Account, error) {
	accounts, err := s.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, nil
}

// DisconnectAccount disconnects a provider account
func (s *ProviderService) DisconnectAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	// Verify the account belongs to the user before deleting
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, repository.ErrAccountNotFound) {
			return repository.ErrAccountNotFound
		}
		return fmt.Errorf("failed to find account: %w", err)
	}

	if account.OwnerUserID != userID {
		return ErrAccountNotOwned
	}

	// Delete the account
	err = s.accountRepo.Delete(ctx, accountID, userID)
	if err != nil {
		return fmt.Errorf("failed to disconnect account: %w", err)
	}

	return nil
}
