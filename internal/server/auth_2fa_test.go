package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"orbit/internal/security"
)

// memberRequest runs a handler as the signed-in member, the way the auth middleware does.
func memberRequest(h http.HandlerFunc, method, path string, body any, member string) *httptest.ResponseRecorder {
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewReader(raw))
	req.RemoteAddr = "127.0.0.1:50000"
	claims := &security.TokenClaims{TokenID: security.MemberTokenPrefix + member + ":abc", Role: "admin"}
	req = req.WithContext(context.WithValue(req.Context(), security.ClaimsContextKey, claims))
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func decodeMap(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &m); err != nil {
		t.Fatalf("bad json %q: %v", rec.Body.String(), err)
	}
	return m
}

func TestTwoFactorLoginFlow(t *testing.T) {
	h, _ := newAuthTestHandler(t)
	local := "127.0.0.1:50000"
	if rec := post(h.handleSetup, "/api/v1/auth/setup", map[string]string{"username": "ops", "password": "correct-horse"}, local); rec.Code != http.StatusOK {
		t.Fatalf("setup: %d", rec.Code)
	}

	// Set up 2FA as the signed-in member.
	setup := decodeMap(t, memberRequest(h.handle2FASetup, http.MethodPost, "/api/v1/auth/2fa/setup", nil, "ops"))
	secret, _ := setup["secret"].(string)
	if secret == "" {
		t.Fatalf("no secret returned: %v", setup)
	}

	// A code cannot enable it before setup is confirmed with a valid code.
	if rec := memberRequest(h.handle2FAEnable, http.MethodPost, "/api/v1/auth/2fa/enable", map[string]string{"code": "000000"}, "ops"); rec.Code == http.StatusOK {
		t.Fatal("enable accepted a wrong code")
	}
	now := time.Now()
	code, _ := security.TOTPCode(secret, now)
	if rec := memberRequest(h.handle2FAEnable, http.MethodPost, "/api/v1/auth/2fa/enable", map[string]string{"code": code}, "ops"); rec.Code != http.StatusOK {
		t.Fatalf("enable: %d %s", rec.Code, rec.Body)
	}

	// Password alone no longer gets a session; it returns a challenge.
	rec := post(h.handleLogin, "/api/v1/auth/login", map[string]string{"username": "ops", "password": "correct-horse"}, "127.0.0.1:50001")
	login := decodeMap(t, rec)
	challenge, _ := login["challenge"].(string)
	if rec.Code != http.StatusOK || challenge == "" || login["token"] != nil {
		t.Fatalf("login should need 2FA, got %d %v", rec.Code, login)
	}

	// Wrong code is rejected.
	if rec := post(h.handleLogin2FA, "/api/v1/auth/login/2fa", map[string]string{"challenge": challenge, "code": "000000"}, "127.0.0.1:50001"); rec.Code != http.StatusUnauthorized {
		t.Fatalf("wrong code: got %d", rec.Code)
	}

	// Correct code issues a session.
	// Enabling used the current step, so sign-in uses the next one (replay protection).
	code, _ = security.TOTPCode(secret, time.Now().Add(30*time.Second))
	rec = post(h.handleLogin2FA, "/api/v1/auth/login/2fa", map[string]string{"challenge": challenge, "code": code}, "127.0.0.1:50001")
	if rec.Code != http.StatusOK || decodeMap(t, rec)["token"] == nil {
		t.Fatalf("correct code: got %d %s", rec.Code, rec.Body)
	}

	// The same challenge cannot be used again.
	if rec := post(h.handleLogin2FA, "/api/v1/auth/login/2fa", map[string]string{"challenge": challenge, "code": code}, "127.0.0.1:50001"); rec.Code == http.StatusOK {
		t.Fatal("challenge reused")
	}

	// The same code cannot be used on a second sign-in.
	rec = post(h.handleLogin, "/api/v1/auth/login", map[string]string{"username": "ops", "password": "correct-horse"}, "127.0.0.1:50002")
	second, _ := decodeMap(t, rec)["challenge"].(string)
	if rec := post(h.handleLogin2FA, "/api/v1/auth/login/2fa", map[string]string{"challenge": second, "code": code}, "127.0.0.1:50002"); rec.Code == http.StatusOK {
		t.Fatal("replayed code accepted")
	}

	// Disable needs the password.
	if rec := memberRequest(h.handle2FADisable, http.MethodPost, "/api/v1/auth/2fa/disable", map[string]string{"password": "nope"}, "ops"); rec.Code == http.StatusOK {
		t.Fatal("disable accepted a wrong password")
	}
	if rec := memberRequest(h.handle2FADisable, http.MethodPost, "/api/v1/auth/2fa/disable", map[string]string{"password": "correct-horse"}, "ops"); rec.Code != http.StatusOK {
		t.Fatalf("disable: %d", rec.Code)
	}
	rec = post(h.handleLogin, "/api/v1/auth/login", map[string]string{"username": "ops", "password": "correct-horse"}, "127.0.0.1:50003")
	if decodeMap(t, rec)["token"] == nil {
		t.Fatal("after disable, password login should issue a session")
	}
}

func TestTwoFactorChallengeExpires(t *testing.T) {
	h, _ := newAuthTestHandler(t)
	local := "127.0.0.1:50000"
	_ = post(h.handleSetup, "/api/v1/auth/setup", map[string]string{"username": "ops", "password": "correct-horse"}, local)
	h.challenges["stale"] = &loginChallenge{username: "ops", expires: time.Now().Add(-time.Second)}
	if rec := post(h.handleLogin2FA, "/api/v1/auth/login/2fa", map[string]string{"challenge": "stale", "code": "123456"}, local); rec.Code != http.StatusUnauthorized {
		t.Fatalf("expired challenge: got %d", rec.Code)
	}
}
