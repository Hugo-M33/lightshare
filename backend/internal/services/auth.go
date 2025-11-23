package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/lightshare/backend/internal/models"
	"github.com/lightshare/backend/internal/repository"
	"github.com/lightshare/backend/pkg/crypto"
	"github.com/lightshare/backend/pkg/email"
	"github.com/lightshare/backend/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrWeakPassword       = errors.New("password too weak")
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	jwtService       *jwt.Service
	emailService     *email.Service
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	jwtService *jwt.Service,
	emailService *email.Service,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		emailService:     emailService,
	}
}

// SignupRequest represents a signup request
type SignupRequest struct {
	Email    string
	Password string
}

// SignupResponse represents a signup response
type SignupResponse struct {
	User    *models.User `json:"user"`
	Message string       `json:"message"`
}

// Signup creates a new user account
func (s *AuthService) Signup(ctx context.Context, req SignupRequest) (*SignupResponse, error) {
	// Validate email
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if !email.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email address")
	}

	// Validate password
	if len(req.Password) < 8 {
		return nil, ErrWeakPassword
	}

	// Hash password
	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate email verification token
	verificationToken, err := jwt.GenerateRandomToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Create user
	user, err := s.userRepo.Create(ctx, models.CreateUserParams{
		Email:                      req.Email,
		PasswordHash:               passwordHash,
		EmailVerificationToken:     verificationToken,
		EmailVerificationExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, errors.New("email already registered")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send verification email
	if err := s.emailService.SendVerificationEmail(user.Email, verificationToken); err != nil {
		// Log error but don't fail the signup
		// User can request a new verification email
		fmt.Printf("failed to send verification email: %v\n", err)
	}

	return &SignupResponse{
		User:    user,
		Message: "Account created successfully. Please check your email to verify your account.",
	}, nil
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string
	Password string
}

// LoginResponse represents a login response
type LoginResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
	TokenType    string       `json:"token_type"`
}

// Login authenticates a user with email and password
func (s *AuthService) Login(ctx context.Context, req LoginRequest, userAgent, ipAddress *string) (*LoginResponse, error) {
	// Normalize email
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Compare password
	if err := crypto.ComparePassword(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if email is verified
	if !user.EmailVerified {
		return nil, ErrEmailNotVerified
	}

	// Generate token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token in database
	refreshTokenHash := crypto.HashToken(tokenPair.RefreshToken)
	_, err = s.refreshTokenRepo.Create(ctx, user.ID, refreshTokenHash, tokenPair.ExpiresAt.Add(29*24*time.Hour), userAgent, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		TokenType:    tokenPair.TokenType,
	}, nil
}

// VerifyEmail verifies a user's email with the verification token and returns JWT tokens
func (s *AuthService) VerifyEmail(ctx context.Context, token string, userAgent, ipAddress *string) (*LoginResponse, error) {
	// Get user by verification token
	user, err := s.userRepo.GetByEmailVerificationToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrTokenExpired) {
			return nil, repository.ErrTokenExpired
		}
		return nil, fmt.Errorf("failed to get user by verification token: %w", err)
	}

	// Verify email (mark as verified and clear token)
	if err := s.userRepo.VerifyEmail(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}

	// Update user's email_verified status for the response
	user.EmailVerified = true

	// Generate token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token in database
	refreshTokenHash := crypto.HashToken(tokenPair.RefreshToken)
	_, err = s.refreshTokenRepo.Create(ctx, user.ID, refreshTokenHash, tokenPair.ExpiresAt.Add(29*24*time.Hour), userAgent, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		TokenType:    tokenPair.TokenType,
	}, nil
}

// RequestMagicLink sends a magic link to the user's email
func (s *AuthService) RequestMagicLink(ctx context.Context, emailAddr string) error {
	// Normalize email
	emailAddr = strings.TrimSpace(strings.ToLower(emailAddr))

	// Check if user exists
	user, err := s.userRepo.GetByEmail(ctx, emailAddr)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			// Don't reveal if email exists or not for security
			return nil
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Generate magic link token
	magicLinkToken, err := jwt.GenerateRandomToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate magic link token: %w", err)
	}

	// Set magic link token with 15 minute expiration
	expiresAt := time.Now().Add(15 * time.Minute)
	if err := s.userRepo.SetMagicLinkToken(ctx, user.Email, magicLinkToken, expiresAt); err != nil {
		return fmt.Errorf("failed to set magic link token: %w", err)
	}

	// Send magic link email
	if err := s.emailService.SendMagicLinkEmail(user.Email, magicLinkToken); err != nil {
		return fmt.Errorf("failed to send magic link email: %w", err)
	}

	return nil
}

// LoginWithMagicLink authenticates a user with a magic link token
func (s *AuthService) LoginWithMagicLink(ctx context.Context, token string, userAgent, ipAddress *string) (*LoginResponse, error) {
	// Get user by magic link token
	user, err := s.userRepo.GetByMagicLinkToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrTokenExpired) {
			return nil, errors.New("magic link expired")
		}
		return nil, fmt.Errorf("failed to get user by magic link: %w", err)
	}

	// Clear magic link token
	if err := s.userRepo.ClearMagicLinkToken(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to clear magic link token: %w", err)
	}

	// Generate token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token in database
	refreshTokenHash := crypto.HashToken(tokenPair.RefreshToken)
	_, err = s.refreshTokenRepo.Create(ctx, user.ID, refreshTokenHash, tokenPair.ExpiresAt.Add(29*24*time.Hour), userAgent, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		TokenType:    tokenPair.TokenType,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string, userAgent, ipAddress *string) (*LoginResponse, error) {
	// Validate refresh token
	_, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Check if refresh token exists and is not revoked
	refreshTokenHash := crypto.HashToken(refreshToken)
	storedToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return nil, errors.New("invalid refresh token")
		}
		if errors.Is(err, repository.ErrRefreshTokenRevoked) {
			return nil, errors.New("refresh token revoked")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate new token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Revoke old refresh token
	if err := s.refreshTokenRepo.Revoke(ctx, refreshTokenHash); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	// Store new refresh token
	newRefreshTokenHash := crypto.HashToken(tokenPair.RefreshToken)
	_, err = s.refreshTokenRepo.Create(ctx, user.ID, newRefreshTokenHash, tokenPair.ExpiresAt.Add(29*24*time.Hour), userAgent, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		TokenType:    tokenPair.TokenType,
	}, nil
}

// Logout logs out a user by revoking their refresh token
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	refreshTokenHash := crypto.HashToken(refreshToken)
	return s.refreshTokenRepo.Revoke(ctx, refreshTokenHash)
}

// LogoutAll logs out a user from all devices
func (s *AuthService) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return s.refreshTokenRepo.RevokeAllForUser(ctx, userID)
}
