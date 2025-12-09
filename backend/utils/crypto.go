package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

func EncryptMessage(message string, password string) ([]byte, error) {
	if password == "" {
		return []byte(message), nil
	}
	hash := sha256.Sum256([]byte(password))
	key := hash[:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize()) //use gcm.NonceSize() which is 12 bytes for performances.
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(message), nil)
	return ciphertext, nil
}

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
