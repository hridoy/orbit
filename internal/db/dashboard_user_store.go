package db

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

var ErrUserExists = errors.New("a member with this username already exists")

type DashboardUser struct {
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	Disabled     bool      `json:"disabled"`
	CreatedAt    time.Time `json:"created_at"`
}

type DashboardUserStore struct {
	db *DB
}

func NewDashboardUserStore(db *DB) *DashboardUserStore {
	return &DashboardUserStore{db: db}
}

func (s *DashboardUserStore) Get(username string) (*DashboardUser, error) {
	var u DashboardUser
	var disabled int
	err := s.db.QueryRow(`SELECT username, password_hash, role, disabled, created_at FROM dashboard_users WHERE username = ?;`, username).
		Scan(&u.Username, &u.PasswordHash, &u.Role, &disabled, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.Disabled = disabled == 1
	return &u, nil
}

func (s *DashboardUserStore) List() ([]*DashboardUser, error) {
	rows, err := s.db.Query(`SELECT username, password_hash, role, disabled, created_at FROM dashboard_users ORDER BY created_at, username;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*DashboardUser
	for rows.Next() {
		var u DashboardUser
		var disabled int
		if err := rows.Scan(&u.Username, &u.PasswordHash, &u.Role, &disabled, &u.CreatedAt); err != nil {
			return nil, err
		}
		u.Disabled = disabled == 1
		out = append(out, &u)
	}
	return out, rows.Err()
}

func (s *DashboardUserStore) Count() (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM dashboard_users;`).Scan(&n)
	return n, err
}

func (s *DashboardUserStore) Create(u *DashboardUser) error {
	disabled := 0
	if u.Disabled {
		disabled = 1
	}
	_, err := s.db.Exec(`INSERT INTO dashboard_users (username, password_hash, role, disabled) VALUES (?, ?, ?, ?);`,
		u.Username, u.PasswordHash, u.Role, disabled)
	if err != nil && isUniqueViolation(err) {
		return ErrUserExists
	}
	return err
}

// Update changes the role and disabled flag. A nil field is left as it is.
func (s *DashboardUserStore) Update(username string, role *string, disabled *bool, passwordHash *string) error {
	u, err := s.Get(username)
	if err != nil {
		return err
	}
	if u == nil {
		return sql.ErrNoRows
	}
	if role != nil {
		u.Role = *role
	}
	if disabled != nil {
		u.Disabled = *disabled
	}
	if passwordHash != nil {
		u.PasswordHash = *passwordHash
	}
	d := 0
	if u.Disabled {
		d = 1
	}
	_, err = s.db.Exec(`UPDATE dashboard_users SET role = ?, disabled = ?, password_hash = ? WHERE username = ?;`,
		u.Role, d, u.PasswordHash, username)
	return err
}

func (s *DashboardUserStore) Delete(username string) error {
	_, err := s.db.Exec(`DELETE FROM dashboard_users WHERE username = ?;`, username)
	return err
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "constraint failed")
}
