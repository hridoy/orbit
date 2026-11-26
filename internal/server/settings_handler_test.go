package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"orbit/internal/db"
	"orbit/internal/security"
	"orbit/internal/tlsengine"
)

func newSettingsTest(t *testing.T) *SettingsHandler {
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
	key[0] = 9
	v, err := security.NewVault(key)
	if err != nil {
		t.Fatal(err)
	}
	eng := tlsengine.New(db.NewTLSDomainStore(conn), v)
	h, err := NewSettingsHandler(db.NewServerSettingsStore(conn), eng)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func putSettings(h *SettingsHandler, body any) *httptest.ResponseRecorder {
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/server/settings", bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	h.handlePut(rec, req)
	return rec
}

func TestPanelDomainRejectsInvalidName(t *testing.T) {
	h := newSettingsTest(t)
	if rec := putSettings(h, map[string]string{"panel_domain": "not a domain"}); rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid name: got %d, want 400", rec.Code)
	}
	if h.PanelDomain() != "" {
		t.Fatal("invalid name was stored")
	}
}

func TestPanelDomainCanBeCleared(t *testing.T) {
	h := newSettingsTest(t)
	h.domain.Store("panel.example.com")
	if rec := putSettings(h, map[string]string{"panel_domain": ""}); rec.Code != http.StatusOK {
		t.Fatalf("clear: got %d", rec.Code)
	}
	if h.PanelDomain() != "" {
		t.Fatal("panel domain not cleared")
	}
}

func TestHostOnlyStripsPort(t *testing.T) {
	if got := hostOnly("Panel.Example.com:443"); got != "panel.example.com" {
		t.Fatalf("got %q", got)
	}
	if got := hostOnly("panel.example.com"); got != "panel.example.com" {
		t.Fatalf("got %q", got)
	}
}
