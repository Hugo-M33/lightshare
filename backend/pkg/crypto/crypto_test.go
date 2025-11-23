package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptToken(t *testing.T) {
	// Generate a valid 32-byte key
	key := []byte("12345678901234567890123456789012") // 32 bytes

	testCases := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple token",
			plaintext: "my-secret-token-12345",
		},
		{
			name:      "long token",
			plaintext: "this-is-a-very-long-token-with-many-characters-0123456789",
		},
		{
			name:      "token with special chars",
			plaintext: "token!@#$%^&*()_+-=[]{}|;:,.<>?",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt the token
			encrypted, err := EncryptToken(tc.plaintext, key)
			if err != nil {
				t.Fatalf("EncryptToken failed: %v", err)
			}

			// Verify encrypted data is not empty
			if len(encrypted) == 0 {
				t.Fatal("Encrypted token is empty")
			}

			// Verify encrypted data is different from plaintext
			if bytes.Equal(encrypted, []byte(tc.plaintext)) {
				t.Fatal("Encrypted token is same as plaintext")
			}

			// Decrypt the token
			decrypted, err := DecryptToken(encrypted, key)
			if err != nil {
				t.Fatalf("DecryptToken failed: %v", err)
			}

			// Verify decrypted matches original
			if decrypted != tc.plaintext {
				t.Fatalf("Decrypted token doesn't match original. Got %q, want %q", decrypted, tc.plaintext)
			}
		})
	}
}

func TestEncryptTokenInvalidKey(t *testing.T) {
	testCases := []struct {
		name string
		key  []byte
	}{
		{
			name: "key too short",
			key:  []byte("short"),
		},
		{
			name: "key too long",
			key:  []byte("this-key-is-way-too-long-for-aes-256"),
		},
		{
			name: "empty key",
			key:  []byte(""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := EncryptToken("test-token", tc.key)
			if err == nil {
				t.Fatal("Expected error for invalid key, got nil")
			}
		})
	}
}

func TestDecryptTokenInvalidKey(t *testing.T) {
	validKey := []byte("12345678901234567890123456789012")
	encrypted, _ := EncryptToken("test-token", validKey)

	invalidKey := []byte("00000000000000000000000000000000")

	_, err := DecryptToken(encrypted, invalidKey)
	if err == nil {
		t.Fatal("Expected error when decrypting with wrong key, got nil")
	}
}

func TestDecryptTokenInvalidData(t *testing.T) {
	key := []byte("12345678901234567890123456789012")

	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "too short data",
			data: []byte("short"),
		},
		{
			name: "garbage data",
			data: []byte("this-is-not-encrypted-data"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecryptToken(tc.data, key)
			if err == nil {
				t.Fatal("Expected error for invalid data, got nil")
			}
		})
	}
}

func TestEncryptionUniqueness(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plaintext := "test-token"

	// Encrypt the same token twice
	encrypted1, err := EncryptToken(plaintext, key)
	if err != nil {
		t.Fatalf("First encryption failed: %v", err)
	}

	encrypted2, err := EncryptToken(plaintext, key)
	if err != nil {
		t.Fatalf("Second encryption failed: %v", err)
	}

	// The encrypted values should be different (due to random nonce)
	if bytes.Equal(encrypted1, encrypted2) {
		t.Fatal("Two encryptions of the same plaintext produced identical ciphertext (nonce not random)")
	}

	// But both should decrypt to the same plaintext
	decrypted1, _ := DecryptToken(encrypted1, key)
	decrypted2, _ := DecryptToken(encrypted2, key)

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Fatal("Decrypted values don't match original plaintext")
	}
}
