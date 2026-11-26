package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddlewareLockdown(t *testing.T) {
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(i + 1)
	}
	tm, err := NewTokenManager(seed)
	if err != nil {
		t.Fatal(err)
	}
	token, _ := tm.IssueToken("t1", "admin", time.Hour)

	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	h := NewAuthMiddleware(tm).Wrap(ok)

	cases := []struct {
		name, path, auth string
		want             int
	}{
		{"telemetry needs token", "/api/v1/telemetry/metrics", "", http.StatusUnauthorized},
		{"telemetry with bearer", "/api/v1/telemetry/metrics", "Bearer " + token, http.StatusOK},
		{"query token rejected on normal API", "/api/v1/databases?token=" + token, "", http.StatusUnauthorized},
		{"query token allowed on SSE stream", "/api/v1/telemetry/stream?token=" + token, "", http.StatusOK},
		{"query token allowed on log stream", "/api/v1/apps/a1/logs/stream?token=" + token, "", http.StatusOK},
		{"healthz open", "/healthz", "", http.StatusOK},
	}
	for _, c := range cases {
		req := httptest.NewRequest("GET", c.path, nil)
		if c.auth != "" {
			req.Header.Set("Authorization", c.auth)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != c.want {
			t.Errorf("%s: want %d, got %d", c.name, c.want, rec.Code)
		}
	}
}
