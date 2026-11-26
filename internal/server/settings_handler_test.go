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

func TestLockBlocksRemoteIPOnly(t *testing.T) {
	h := newSettingsTest(t)
	h.domain.Store("panel.example.com")
	h.locked.Store(true)

	remoteIP := httptest.NewRequest("GET", "/", nil)
	remoteIP.Host = "163.227.238.133:8787"
	remoteIP.RemoteAddr = "203.0.113.7:50000"
	if !h.Blocks(remoteIP) {
		t.Fatal("remote request to the IP was not blocked")
	}

	panel := httptest.NewRequest("GET", "/", nil)
	panel.Host = "panel.example.com"
	panel.RemoteAddr = "203.0.113.7:50000"
	if h.Blocks(panel) {
		t.Fatal("panel domain request was blocked")
	}

	tunnel := httptest.NewRequest("GET", "/", nil)
	tunnel.Host = "127.0.0.1:8787"
	tunnel.RemoteAddr = "127.0.0.1:50001"
	if h.Blocks(tunnel) {
		t.Fatal("local SSH tunnel request was blocked")
	}

	h.locked.Store(false)
	if h.Blocks(remoteIP) {
		t.Fatal("unlocked dashboard blocked a request")
	}
}

func TestLockNeedsActiveCertificate(t *testing.T) {
	h := newSettingsTest(t)
	lock := true
	raw, _ := json.Marshal(map[string]any{"panel_domain": "panel.example.com", "lock_to_panel_domain": lock})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/server/settings", bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	h.handlePut(rec, req)
	if rec.Code == http.StatusOK || rec.Code == http.StatusCreated {
		t.Fatalf("lock accepted before the certificate was active: %d", rec.Code)
	}
	if h.Locked() {
		t.Fatal("dashboard locked without an active certificate")
	}
}
