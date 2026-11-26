package db

import (
	"database/sql"
	"errors"
)

// ErrAdminExists is returned when the first-run setup has already created the operator account.
var ErrAdminExists = errors.New("admin account already exists")

type AdminUser struct {
	Username     string
	PasswordHash string
}

type AdminStore struct {
	db *DB
}

func NewAdminStore(db *DB) *AdminStore {
	return &AdminStore{db: db}
}

// Get returns the operator account, or nil when first-run setup has not happened yet.
func (s *AdminStore) Get() (*AdminUser, error) {
	var u AdminUser
	err := s.db.QueryRow(`SELECT username, password_hash FROM admin_users WHERE id = 1;`).
		Scan(&u.Username, &u.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Create stores the operator account once. A second call fails with ErrAdminExists.
func (s *AdminStore) Create(username, passwordHash string) error {
	res, err := s.db.Exec(`
	INSERT INTO admin_users (id, username, password_hash) VALUES (1, ?, ?)
	ON CONFLICT(id) DO NOTHING;`, username, passwordHash)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrAdminExists
	}
	return nil
}
