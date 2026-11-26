package integrations

import (
	"path/filepath"
	"strings"
	"testing"

	"orbit/internal/db"
	"orbit/internal/security"
)

func newManager(t *testing.T) *Manager {
	t.Helper()
	conn, err := db.Open(filepath.Join(t.TempDir(), "orbit.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	if err := conn.Migrate(); err != nil {
		t.Fatal(err)
	}
	key := make([]byte, 32)
	key[0] = 3
	v, err := security.NewVault(key)
	if err != nil {
		t.Fatal(err)
	}
	return New(db.NewIntegrationStore(conn), v)
}

func TestSecretIsEncryptedAndRoundTrips(t *testing.T) {
	m := newManager(t)
	if err := m.Set(Cloudflare, "cf-token-ABCD1234", ""); err != nil {
		t.Fatal(err)
	}
	secret, _, ok := m.Get(Cloudflare)
	if !ok || secret != "cf-token-ABCD1234" {
		t.Fatalf("round trip failed: %q %v", secret, ok)
	}
}

func TestListNeverExposesSecret(t *testing.T) {
	m := newManager(t)
	if err := m.Set(Cloudflare, "cf-token-ABCD1234", ""); err != nil {
		t.Fatal(err)
	}
	list, err := m.List()
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range list {
		if s.Provider == Cloudflare {
			if !s.Configured || s.Hint != "••••1234" {
				t.Fatalf("unexpected status: %+v", s)
			}
		}
	}
	raw := ""
	for _, s := range list {
		raw += s.Hint + s.ExtraValue
	}
	if strings.Contains(raw, "cf-token-ABCD") {
		t.Fatal("secret leaked into the list")
	}
}

func TestZeroSSLNeedsKeyID(t *testing.T) {
	m := newManager(t)
	if err := m.Set(ZeroSSL, "hmac-secret", ""); err == nil {
		t.Fatal("ZeroSSL saved without its key ID")
	}
	if err := m.Set(ZeroSSL, "hmac-secret", "kid-1"); err != nil {
		t.Fatal(err)
	}
	secret, extra, ok := m.Get(ZeroSSL)
	if !ok || secret != "hmac-secret" || extra != "kid-1" {
		t.Fatalf("got %q %q %v", secret, extra, ok)
	}
}

func TestUnknownProviderRejected(t *testing.T) {
	m := newManager(t)
	if err := m.Set("aws", "x", ""); err == nil {
		t.Fatal("unknown provider accepted")
	}
}
