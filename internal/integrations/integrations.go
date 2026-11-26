// Package integrations keeps third-party credentials (Cloudflare API token, ZeroSSL EAB keys)
// encrypted at rest and hands them to the features that need them. Secrets never leave the
// server through the API; the list shows only a short hint.
package integrations

import (
	"errors"
	"strings"
	"time"

	"orbit/internal/db"
	"orbit/internal/security"
)

const (
	Cloudflare = "cloudflare"
	ZeroSSL    = "zerossl"
)

// Catalog describes each supported integration for the dashboard.
var Catalog = map[string]Spec{
	Cloudflare: {Name: "Cloudflare", Secret: "API token", Used: "Origin certificates for domains behind Cloudflare, with automatic renewal"},
	ZeroSSL:    {Name: "ZeroSSL", Secret: "EAB HMAC key", Extra: "EAB key ID", Used: "ZeroSSL certificates for domains"},
}

type Spec struct {
	Name   string `json:"name"`
	Secret string `json:"secret_label"`
	Extra  string `json:"extra_label,omitempty"`
	Used   string `json:"used_for"`
}

// Status is what the dashboard shows. It never contains the secret.
type Status struct {
	Provider   string     `json:"provider"`
	Name       string     `json:"name"`
	Secret     string     `json:"secret_label"`
	Extra      string     `json:"extra_label,omitempty"`
	UsedFor    string     `json:"used_for"`
	Configured bool       `json:"configured"`
	Hint       string     `json:"hint,omitempty"`
	ExtraValue string     `json:"extra_value,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

type Manager struct {
	store *db.IntegrationStore
	vault *security.Vault
}

func New(store *db.IntegrationStore, vault *security.Vault) *Manager {
	return &Manager{store: store, vault: vault}
}

// List returns every supported integration, configured or not.
func (m *Manager) List() ([]Status, error) {
	saved, err := m.store.List()
	if err != nil {
		return nil, err
	}
	byName := map[string]*db.Integration{}
	for _, s := range saved {
		byName[s.Provider] = s
	}
	var out []Status
	for _, p := range []string{Cloudflare, ZeroSSL} {
		spec := Catalog[p]
		st := Status{Provider: p, Name: spec.Name, Secret: spec.Secret, Extra: spec.Extra, UsedFor: spec.Used}
		if s := byName[p]; s != nil {
			st.Configured = true
			st.Hint = s.Hint
			st.ExtraValue = s.Extra
			t := s.UpdatedAt
			st.UpdatedAt = &t
		}
		out = append(out, st)
	}
	return out, nil
}

// Set saves a credential, encrypted. extra is used for the ZeroSSL key ID.
func (m *Manager) Set(provider, secret, extra string) error {
	if _, ok := Catalog[provider]; !ok {
		return errors.New("unknown integration")
	}
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return errors.New("secret is required")
	}
	if provider == ZeroSSL && strings.TrimSpace(extra) == "" {
		return errors.New("the EAB key ID is required for ZeroSSL")
	}
	ct, nonce, err := m.vault.Encrypt([]byte(secret))
	if err != nil {
		return err
	}
	return m.store.Upsert(&db.Integration{
		Provider:  provider,
		SecretEnc: ct + ":" + nonce,
		Extra:     strings.TrimSpace(extra),
		Hint:      hint(secret),
	})
}

// Get returns the decrypted secret and extra value, or ok=false when nothing is saved.
func (m *Manager) Get(provider string) (secret, extra string, ok bool) {
	i, err := m.store.Get(provider)
	if err != nil || i == nil {
		return "", "", false
	}
	parts := strings.SplitN(i.SecretEnc, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	raw, err := m.vault.Decrypt(parts[0], parts[1])
	if err != nil {
		return "", "", false
	}
	return string(raw), i.Extra, true
}

func (m *Manager) Delete(provider string) error {
	if _, ok := Catalog[provider]; !ok {
		return errors.New("unknown integration")
	}
	return m.store.Delete(provider)
}

// hint keeps the last four characters so the operator can recognise the key.
func hint(secret string) string {
	if len(secret) <= 4 {
		return "••••"
	}
	return "••••" + secret[len(secret)-4:]
}
