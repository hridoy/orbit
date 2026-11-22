package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidMasterKey = errors.New("master key must be exactly 32 bytes for AES-256")
	ErrDecryptionFailed = errors.New("failed to decrypt ciphertext")
)

// Vault handles AES-256-GCM envelope encryption for sensitive runtime secrets.
type Vault struct {
	aead cipher.AEAD
}

// NewVault initializes an AES-256-GCM cipher using the provided 32-byte master key.
func NewVault(masterKey []byte) (*Vault, error) {
	if len(masterKey) != 32 {
		return nil, ErrInvalidMasterKey
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM AEAD: %w", err)
	}

	return &Vault{aead: gcm}, nil
}

// Encrypt encrypts plain bytes and returns base64-encoded ciphertext and nonce.
func (v *Vault) Encrypt(plaintext []byte) (ciphertextB64, nonceB64 string, err error) {
	nonce := make([]byte, v.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", fmt.Errorf("failed to generate random nonce: %w", err)
	}

	ciphertext := v.aead.Seal(nil, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), base64.StdEncoding.EncodeToString(nonce), nil
}

// Decrypt decodes base64 ciphertext and nonce, then decrypts back to plaintext.
func (v *Vault) Decrypt(ciphertextB64, nonceB64 string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, fmt.Errorf("invalid ciphertext base64: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(nonceB64)
	if err != nil {
		return nil, fmt.Errorf("invalid nonce base64: %w", err)
	}

	if len(nonce) != v.aead.NonceSize() {
		return nil, ErrDecryptionFailed
	}

	plaintext, err := v.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}
