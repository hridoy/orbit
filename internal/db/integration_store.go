package db

import (
	"database/sql"
	"errors"
	"time"
)

// Integration is one saved third-party credential. SecretEnc holds ciphertext.
type Integration struct {
	Provider  string    `json:"provider"`
	SecretEnc string    `json:"-"`
	Extra     string    `json:"extra"`
	Hint      string    `json:"hint"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IntegrationStore struct {
	db *DB
}

func NewIntegrationStore(db *DB) *IntegrationStore {
	return &IntegrationStore{db: db}
}

func (s *IntegrationStore) Get(provider string) (*Integration, error) {
	var i Integration
	err := s.db.QueryRow(`SELECT provider, secret_enc, extra, hint, updated_at FROM integrations WHERE provider = ?;`, provider).
		Scan(&i.Provider, &i.SecretEnc, &i.Extra, &i.Hint, &i.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (s *IntegrationStore) List() ([]*Integration, error) {
	rows, err := s.db.Query(`SELECT provider, secret_enc, extra, hint, updated_at FROM integrations ORDER BY provider;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Integration
	for rows.Next() {
		var i Integration
		if err := rows.Scan(&i.Provider, &i.SecretEnc, &i.Extra, &i.Hint, &i.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &i)
	}
	return out, rows.Err()
}

func (s *IntegrationStore) Upsert(i *Integration) error {
	_, err := s.db.Exec(`
	INSERT INTO integrations (provider, secret_enc, extra, hint, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(provider) DO UPDATE SET secret_enc = excluded.secret_enc, extra = excluded.extra,
		hint = excluded.hint, updated_at = CURRENT_TIMESTAMP;`, i.Provider, i.SecretEnc, i.Extra, i.Hint)
	return err
}

func (s *IntegrationStore) Delete(provider string) error {
	_, err := s.db.Exec(`DELETE FROM integrations WHERE provider = ?;`, provider)
	return err
}
