package db

import (
	"database/sql"
	"errors"
)

// ServerSettingsStore keeps small server-wide values such as the panel domain.
type ServerSettingsStore struct {
	db *DB
}

func NewServerSettingsStore(db *DB) *ServerSettingsStore {
	return &ServerSettingsStore{db: db}
}

// Get returns the value for key, or "" when it is not set.
func (s *ServerSettingsStore) Get(key string) (string, error) {
	var v string
	err := s.db.QueryRow(`SELECT value FROM server_settings WHERE key = ?;`, key).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return v, err
}

// Set stores a value. An empty value clears the setting.
func (s *ServerSettingsStore) Set(key, value string) error {
	if value == "" {
		_, err := s.db.Exec(`DELETE FROM server_settings WHERE key = ?;`, key)
		return err
	}
	_, err := s.db.Exec(`
	INSERT INTO server_settings (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP;`, key, value)
	return err
}
