package s3

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"orbit/internal/security"
)

const S3CredentialsFile = "s3_credentials.enc"

// S3Credentials is the S3 access key pair. It is random (not derived from the master key),
// so leaking an S3 secret does not expose the master key.
type S3Credentials struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type sealedCredentials struct {
	Nonce      string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
}

const s3CredentialsLabel = "orbit/s3/credentials/v1"

func credentialAEAD(masterKey []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(security.DeriveKey(masterKey, s3CredentialsLabel))
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func randomHex(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// SealS3Credentials encrypts creds under masterKey and writes them to path (0600).
// It does not rename or replace anything else, so callers control atomicity.
func SealS3Credentials(path string, masterKey []byte, creds *S3Credentials) error {
	aead, err := credentialAEAD(masterKey)
	if err != nil {
		return err
	}
	plain, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	sealed := sealedCredentials{
		Nonce:      hex.EncodeToString(nonce),
		Ciphertext: hex.EncodeToString(aead.Seal(nil, nonce, plain, nil)),
	}
	data, err := json.Marshal(sealed)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// OpenS3Credentials decrypts the credentials file at path with masterKey.
func OpenS3Credentials(path string, masterKey []byte) (*S3Credentials, error) {
	aead, err := credentialAEAD(masterKey)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var sealed sealedCredentials
	if err := json.Unmarshal(raw, &sealed); err != nil {
		return nil, fmt.Errorf("parse s3 credentials: %w", err)
	}
	nonce, err1 := hex.DecodeString(sealed.Nonce)
	ct, err2 := hex.DecodeString(sealed.Ciphertext)
	if err1 != nil || err2 != nil || len(nonce) != aead.NonceSize() {
		return nil, fmt.Errorf("malformed s3 credentials file")
	}
	plain, err := aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt s3 credentials: %w", err)
	}
	var creds S3Credentials
	if err := json.Unmarshal(plain, &creds); err != nil {
		return nil, fmt.Errorf("parse decrypted s3 credentials: %w", err)
	}
	return &creds, nil
}

// NewS3Credentials generates a random access key and secret.
func NewS3Credentials() (*S3Credentials, error) {
	suffix, err := randomHex(10)
	if err != nil {
		return nil, err
	}
	secret, err := randomHex(32)
	if err != nil {
		return nil, err
	}
	return &S3Credentials{AccessKey: "ORB" + strings.ToUpper(suffix), SecretKey: secret}, nil
}

// loadOrCreateS3Credentials returns the S3 key pair in dataDir, generating it on first use.
func loadOrCreateS3Credentials(dataDir string, masterKey []byte) (*S3Credentials, error) {
	if len(masterKey) != 32 {
		return nil, fmt.Errorf("master key must be 32 bytes, got %d", len(masterKey))
	}
	path := filepath.Join(dataDir, S3CredentialsFile)
	if _, err := os.Stat(path); err == nil {
		return OpenS3Credentials(path, masterKey)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	creds, err := NewS3Credentials()
	if err != nil {
		return nil, err
	}
	if err := SealS3Credentials(path, masterKey, creds); err != nil {
		return nil, fmt.Errorf("write s3 credentials: %w", err)
	}
	return creds, nil
}
