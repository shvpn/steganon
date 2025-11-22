package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// EncryptMessage encrypts a message using AES-256-GCM with SHA-256 hashed password
func EncryptMessage(message string, password string) ([]byte, error) {
	if password == "" {
		// No password, return plain message
		return []byte(message), nil
	}

	// Hash password with SHA-256 to get 32-byte key
	hash := sha256.Sum256([]byte(password))
	key := hash[:]

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt and append nonce
	ciphertext := gcm.Seal(nonce, nonce, []byte(message), nil)
	return ciphertext, nil
}

// DecryptMessage decrypts a message using AES-256-GCM with SHA-256 hashed password
func DecryptMessage(encrypted []byte, password string) (string, error) {
	if password == "" {
		// No password, return as plain text
		return string(encrypted), nil
	}

	// Hash password with SHA-256 to get 32-byte key
	hash := sha256.Sum256([]byte(password))
	key := hash[:]

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("invalid password or corrupted data")
	}

	return string(plaintext), nil
}
