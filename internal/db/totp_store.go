package db

import (
	"database/sql"
	"errors"
)

// TOTPRecord is a member's two-factor state. The secret is stored encrypted.
type TOTPRecord struct {
	Username         string
	SecretCiphertext string
	SecretNonce      string
	Enabled          bool
	LastStep         int64
}

type TOTPStore struct {
	db *DB
}

func NewTOTPStore(db *DB) *TOTPStore {
	return &TOTPStore{db: db}
}

// Get returns the record, or nil when the member has never set up two-factor.
func (s *TOTPStore) Get(username string) (*TOTPRecord, error) {
	var r TOTPRecord
	var enabled int
	err := s.db.QueryRow(`SELECT username, secret_ciphertext, secret_nonce, enabled, last_step FROM dashboard_totp WHERE username = ?;`, username).
		Scan(&r.Username, &r.SecretCiphertext, &r.SecretNonce, &enabled, &r.LastStep)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	r.Enabled = enabled == 1
	return &r, nil
}

// Put stores a new pending secret (enabled=false) or replaces the existing record.
func (s *TOTPStore) Put(r *TOTPRecord) error {
	enabled := 0
	if r.Enabled {
		enabled = 1
	}
	_, err := s.db.Exec(`INSERT INTO dashboard_totp (username, secret_ciphertext, secret_nonce, enabled, last_step, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(username) DO UPDATE SET secret_ciphertext = excluded.secret_ciphertext,
			secret_nonce = excluded.secret_nonce, enabled = excluded.enabled,
			last_step = excluded.last_step, updated_at = CURRENT_TIMESTAMP;`,
		r.Username, r.SecretCiphertext, r.SecretNonce, enabled, r.LastStep)
	return err
}

func (s *TOTPStore) Delete(username string) error {
	_, err := s.db.Exec(`DELETE FROM dashboard_totp WHERE username = ?;`, username)
	return err
}
