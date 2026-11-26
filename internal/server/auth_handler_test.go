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
)

func newAuthTestHandler(t *testing.T) (*AuthHandler, *security.TokenManager) {
	t.Helper()
	conn, err := db.Open(filepath.Join(t.TempDir(), "orbit.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	if err := conn.Migrate(); err != nil {
		t.Fatal(err)
	}
	seed := make([]byte, 32)
	seed[0] = 1
	tm, err := security.NewTokenManager(seed)
	if err != nil {
		t.Fatal(err)
	}
	return NewAuthHandler(db.NewAdminStore(conn), tm), tm
}

func post(h http.HandlerFunc, path string, body any, remote string) *httptest.ResponseRecorder {
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	req.RemoteAddr = remote
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func TestSetupThenLogin(t *testing.T) {
	h, tm := newAuthTestHandler(t)
	local := "127.0.0.1:50000"

	if rec := post(h.handleSetup, "/api/v1/auth/setup", map[string]string{"username": "ops", "password": "short"}, local); rec.Code != http.StatusBadRequest {
		t.Fatalf("short password: got %d, want 400", rec.Code)
	}
	rec := post(h.handleSetup, "/api/v1/auth/setup", map[string]string{"username": "ops", "password": "correct-horse"}, local)
	if rec.Code != http.StatusOK {
		t.Fatalf("setup: got %d body %s", rec.Code, rec.Body)
	}
	var out map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	claims, err := tm.VerifyToken(out["token"])
	if err != nil || claims.Role != "admin" {
		t.Fatalf("setup token invalid: %v", err)
	}
	if rec := post(h.handleSetup, "/api/v1/auth/setup", map[string]string{"username": "evil", "password": "another-pass"}, local); rec.Code != http.StatusConflict {
		t.Fatalf("second setup: got %d, want 409", rec.Code)
	}

	if rec := post(h.handleLogin, "/api/v1/auth/login", map[string]string{"username": "ops", "password": "wrong-pass1"}, "198.51.100.7:1"); rec.Code != http.StatusUnauthorized {
		t.Fatalf("wrong password: got %d, want 401", rec.Code)
	}
	if rec := post(h.handleLogin, "/api/v1/auth/login", map[string]string{"username": "ops", "password": "correct-horse"}, "198.51.100.7:1"); rec.Code != http.StatusOK {
		t.Fatalf("login: got %d", rec.Code)
	}
}

func TestLoginLockoutAfterFailures(t *testing.T) {
	h, _ := newAuthTestHandler(t)
	post(h.handleSetup, "/api/v1/auth/setup", map[string]string{"username": "ops", "password": "correct-horse"}, "127.0.0.1:1")
	remote := "198.51.100.8:1"
	for i := 0; i < maxLoginFailures; i++ {
		post(h.handleLogin, "/api/v1/auth/login", map[string]string{"username": "ops", "password": "nope-nope"}, remote)
	}
	if rec := post(h.handleLogin, "/api/v1/auth/login", map[string]string{"username": "ops", "password": "correct-horse"}, remote); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("locked login: got %d, want 429", rec.Code)
	}
}
