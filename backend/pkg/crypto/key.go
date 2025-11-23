package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

// LoadEncryptionKey loads the encryption key from environment variable
// For local development, expects ENCRYPTION_KEY to be a 64-character hex string (32 bytes)
// In production, this would integrate with AWS KMS or similar
func LoadEncryptionKey() ([]byte, error) {
	keyHex := os.Getenv("ENCRYPTION_KEY")
	if keyHex == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY environment variable not set")
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes (64 hex chars), got %d bytes", len(key))
	}

	return key, nil
}

// GenerateEncryptionKey generates a new random 32-byte encryption key
// This is a utility function for setting up new environments
func GenerateEncryptionKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate encryption key: %w", err)
	}
	return hex.EncodeToString(key), nil
}
