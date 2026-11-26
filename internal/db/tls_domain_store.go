package db

import (
	"database/sql"
	"time"
)

// TLSDomain is one public domain and how its certificate is obtained. Secret fields hold
// ciphertext; the tlsengine package encrypts and decrypts them.
type TLSDomain struct {
	Domain     string     `json:"domain"`
	Provider   string     `json:"provider"`
	Status     string     `json:"status"`
	LastError  string     `json:"last_error,omitempty"`
	CertPEM    string     `json:"cert_pem,omitempty"`
	KeyEnc     string     `json:"-"`
	EABKID     string     `json:"eab_kid,omitempty"`
	EABHMACEnc string     `json:"-"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type TLSDomainStore struct {
	db *DB
}

func NewTLSDomainStore(db *DB) *TLSDomainStore {
	return &TLSDomainStore{db: db}
}

func (s *TLSDomainStore) Upsert(d *TLSDomain) error {
	_, err := s.db.Exec(`
	INSERT INTO tls_domains (domain, provider, status, last_error, cert_pem, key_enc, eab_kid, eab_hmac_enc, expires_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(domain) DO UPDATE SET
		provider = excluded.provider,
		status = excluded.status,
		last_error = excluded.last_error,
		cert_pem = excluded.cert_pem,
		key_enc = excluded.key_enc,
		eab_kid = excluded.eab_kid,
		eab_hmac_enc = excluded.eab_hmac_enc,
		expires_at = excluded.expires_at,
		updated_at = CURRENT_TIMESTAMP;`,
		d.Domain, d.Provider, d.Status, d.LastError, d.CertPEM, d.KeyEnc, d.EABKID, d.EABHMACEnc, d.ExpiresAt)
	return err
}

func (s *TLSDomainStore) Get(domain string) (*TLSDomain, error) {
	rows, err := s.db.Query(`SELECT domain, provider, status, last_error, cert_pem, key_enc, eab_kid, eab_hmac_enc, expires_at, updated_at FROM tls_domains WHERE domain = ?;`, domain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list, err := scanTLSDomains(rows)
	if err != nil || len(list) == 0 {
		return nil, err
	}
	return list[0], nil
}

func (s *TLSDomainStore) List() ([]*TLSDomain, error) {
	rows, err := s.db.Query(`SELECT domain, provider, status, last_error, cert_pem, key_enc, eab_kid, eab_hmac_enc, expires_at, updated_at FROM tls_domains ORDER BY domain;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTLSDomains(rows)
}

func (s *TLSDomainStore) Delete(domain string) error {
	_, err := s.db.Exec(`DELETE FROM tls_domains WHERE domain = ?;`, domain)
	return err
}

func scanTLSDomains(rows *sql.Rows) ([]*TLSDomain, error) {
	var out []*TLSDomain
	for rows.Next() {
		var d TLSDomain
		var expires sql.NullTime
		if err := rows.Scan(&d.Domain, &d.Provider, &d.Status, &d.LastError, &d.CertPEM, &d.KeyEnc, &d.EABKID, &d.EABHMACEnc, &expires, &d.UpdatedAt); err != nil {
			return nil, err
		}
		if expires.Valid {
			t := expires.Time
			d.ExpiresAt = &t
		}
		out = append(out, &d)
	}
	return out, rows.Err()
}
