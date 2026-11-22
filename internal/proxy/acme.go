package proxy

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/acme/autocert"
	"orbit/internal/db"
)

// SQLiteCertCache implements autocert.Cache using Orbit's embedded SQLite database.
type SQLiteCertCache struct {
	db *db.DB
}

func NewSQLiteCertCache(database *db.DB) *SQLiteCertCache {
	return &SQLiteCertCache{db: database}
}

func (c *SQLiteCertCache) Get(ctx context.Context, key string) ([]byte, error) {
	query := `SELECT cert_pem FROM tls_certificates WHERE domain = ?;`
	row := c.db.QueryRowContext(ctx, query, key)
	var data string
	err := row.Scan(&data)
	if err == sql.ErrNoRows {
		return nil, autocert.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}

func (c *SQLiteCertCache) Put(ctx context.Context, key string, data []byte) error {
	query := `
	INSERT INTO tls_certificates (domain, cert_pem, key_pem, expires_at, updated_at)
	VALUES (?, ?, '', datetime('now', '+90 days'), CURRENT_TIMESTAMP)
	ON CONFLICT(domain) DO UPDATE SET
		cert_pem = excluded.cert_pem,
		updated_at = CURRENT_TIMESTAMP;
	`
	_, err := c.db.ExecContext(ctx, query, key, string(data))
	return err
}

func (c *SQLiteCertCache) Delete(ctx context.Context, key string) error {
	query := `DELETE FROM tls_certificates WHERE domain = ?;`
	_, err := c.db.ExecContext(ctx, query, key)
	return err
}

// ACMEManager handles Let's Encrypt automatic certificate provisioning and renewal.
type ACMEManager struct {
	manager  *autocert.Manager
	appStore *db.AppStore
}

func NewACMEManager(database *db.DB, appStore *db.AppStore, email string) *ACMEManager {
	cache := NewSQLiteCertCache(database)

	m := &ACMEManager{
		appStore: appStore,
	}

	m.manager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      cache,
		HostPolicy: m.HostPolicy,
		Email:      email,
	}

	return m
}

// HostPolicy only issues certificates for domains actively registered in Orbit apps.
func (m *ACMEManager) HostPolicy(ctx context.Context, host string) error {
	cleanHost := strings.ToLower(strings.TrimSpace(host))
	apps, err := m.appStore.ListApps()
	if err != nil {
		return err
	}

	for _, app := range apps {
		if app.Domain != nil && strings.EqualFold(*app.Domain, cleanHost) {
			return nil // Domain is registered and authorized for TLS!
		}
	}

	return fmt.Errorf("host %s is not registered in Orbit", cleanHost)
}

func (m *ACMEManager) AutocertManager() *autocert.Manager {
	return m.manager
}

var _ autocert.Cache = (*SQLiteCertCache)(nil)
var _ = errors.New
