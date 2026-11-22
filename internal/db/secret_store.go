package db

import (
	"database/sql"
	"fmt"
	"time"

	"orbit/internal/security"
)

type SecretRecord struct {
	SecretID        string    `json:"secret_id"`
	ScopeType       string    `json:"scope_type"` // 'global', 'app', 'database'
	ScopeID         string    `json:"scope_id"`
	SecretKey       string    `json:"secret_key"`
	CiphertextBase64 string   `json:"-"`
	NonceBase64     string    `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type SecretStore struct {
	db    *DB
	vault *security.Vault
}

func NewSecretStore(db *DB, vault *security.Vault) *SecretStore {
	return &SecretStore{db: db, vault: vault}
}

// SetSecret encrypts and persists a secret key/value pair.
func (s *SecretStore) SetSecret(secretID, scopeType, scopeID, key, value string) error {
	cipherB64, nonceB64, err := s.vault.Encrypt([]byte(value))
	if err != nil {
		return fmt.Errorf("encrypt secret %s: %w", key, err)
	}

	query := `
	INSERT INTO secrets (secret_id, scope_type, scope_id, secret_key, ciphertext_base64, nonce_base64, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(scope_type, scope_id, secret_key) DO UPDATE SET
		ciphertext_base64 = excluded.ciphertext_base64,
		nonce_base64 = excluded.nonce_base64,
		updated_at = CURRENT_TIMESTAMP;
	`
	_, err = s.db.Exec(query, secretID, scopeType, scopeID, key, cipherB64, nonceB64)
	return err
}

// GetSecret decrypts and retrieves a single secret value.
func (s *SecretStore) GetSecret(scopeType, scopeID, key string) (string, error) {
	query := `
	SELECT ciphertext_base64, nonce_base64
	FROM secrets WHERE scope_type = ? AND scope_id = ? AND secret_key = ?;
	`
	row := s.db.QueryRow(query, scopeType, scopeID, key)
	var cipherB64, nonceB64 string
	err := row.Scan(&cipherB64, &nonceB64)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	plaintextBytes, err := s.vault.Decrypt(cipherB64, nonceB64)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret %s: %w", key, err)
	}
	return string(plaintextBytes), nil
}

// GetScopeEnvMap loads all secrets for a given scope decrypted into an in-memory map.
// Crucial: this is injected directly into process memory (cgroups runner) without touching disk.
func (s *SecretStore) GetScopeEnvMap(scopeType, scopeID string) (map[string]string, error) {
	query := `
	SELECT secret_key, ciphertext_base64, nonce_base64
	FROM secrets WHERE scope_type = ? AND scope_id = ?;
	`
	rows, err := s.db.Query(query, scopeType, scopeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	envMap := make(map[string]string)
	for rows.Next() {
		var k, cipherB64, nonceB64 string
		if err := rows.Scan(&k, &cipherB64, &nonceB64); err != nil {
			return nil, err
		}
		decrypted, err := s.vault.Decrypt(cipherB64, nonceB64)
		if err != nil {
			return nil, fmt.Errorf("decrypt secret %s: %w", k, err)
		}
		envMap[k] = string(decrypted)
	}
	return envMap, nil
}

// DeleteSecret removes a secret.
func (s *SecretStore) DeleteSecret(scopeType, scopeID, key string) error {
	query := `DELETE FROM secrets WHERE scope_type = ? AND scope_id = ? AND secret_key = ?;`
	_, err := s.db.Exec(query, scopeType, scopeID, key)
	return err
}
