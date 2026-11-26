package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"orbit/internal/db"
	"orbit/internal/security"
)

const (
	dashboardSessionTTL = 30 * 24 * time.Hour
	maxLoginFailures    = 5
	loginLockWindow     = 15 * time.Minute
	maxAuthBodyBytes    = 4 << 10
)

// AuthHandler runs first-run setup and password login for dashboard members. Login returns a
// signed bearer token that names the member, so role and status are checked on each request.
type AuthHandler struct {
	users    *db.DashboardUserStore
	tokenMgr *security.TokenManager

	mu       sync.Mutex
	failures map[string][]time.Time // client IP -> recent failed logins
}

func NewAuthHandler(users *db.DashboardUserStore, tokenMgr *security.TokenManager) *AuthHandler {
	return &AuthHandler{users: users, tokenMgr: tokenMgr, failures: map[string][]time.Time{}}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/auth/status", h.handleStatus)
	mux.HandleFunc("POST /api/v1/auth/setup", h.handleSetup)
	mux.HandleFunc("POST /api/v1/auth/login", h.handleLogin)
}

func (h *AuthHandler) handleStatus(w http.ResponseWriter, r *http.Request) {
	n, err := h.users.Count()
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]any{"setup_required": n == 0})
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func decodeCredentials(w http.ResponseWriter, r *http.Request) (credentials, bool) {
	var c credentials
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return c, false
	}
	c.Username = strings.TrimSpace(c.Username)
	return c, true
}

// handleSetup creates the first admin. It works only while no member exists, and the firewall
// limits who can reach the dashboard.
func (h *AuthHandler) handleSetup(w http.ResponseWriter, r *http.Request) {
	c, ok := decodeCredentials(w, r)
	if !ok {
		return
	}
	if c.Username == "" || len(c.Username) > 64 {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "username must be 1 to 64 characters"})
		return
	}
	hash, err := security.HashPassword(c.Password)
	if err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	n, err := h.users.Count()
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if n > 0 {
		writeAuthJSON(w, http.StatusConflict, map[string]string{"error": "setup is already complete; log in instead"})
		return
	}
	if err := h.users.Create(&db.DashboardUser{Username: c.Username, PasswordHash: hash, Role: "admin"}); err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.issue(w, c.Username)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	ip := clientIP(r)
	if h.locked(ip) {
		writeAuthJSON(w, http.StatusTooManyRequests, map[string]string{"error": "too many failed logins; wait 15 minutes"})
		return
	}
	c, ok := decodeCredentials(w, r)
	if !ok {
		return
	}
	u, err := h.users.Get(c.Username)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if u == nil {
		// Spend the same time as a real check so the response does not reveal which names exist.
		_ = security.VerifyPassword(c.Password, dummyHash)
		h.recordFailure(ip)
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "wrong username or password"})
		return
	}
	if !security.VerifyPassword(c.Password, u.PasswordHash) {
		h.recordFailure(ip)
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "wrong username or password"})
		return
	}
	if u.Disabled {
		writeAuthJSON(w, http.StatusForbidden, map[string]string{"error": "this account is disabled"})
		return
	}
	h.clearFailures(ip)
	h.issue(w, u.Username)
}

// dummyHash is checked when a username does not exist, so timing does not reveal it.
var dummyHash, _ = security.HashPassword("orbit-dummy-password")

func (h *AuthHandler) issue(w http.ResponseWriter, username string) {
	id := make([]byte, 8)
	_, _ = rand.Read(id)
	tokenID := security.MemberTokenPrefix + username + ":" + hex.EncodeToString(id)
	token, err := h.tokenMgr.IssueToken(tokenID, "viewer", dashboardSessionTTL)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]string{"token": token, "username": username})
}

func (h *AuthHandler) locked(ip string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.recent(ip)) >= maxLoginFailures
}

func (h *AuthHandler) recordFailure(ip string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.failures[ip] = append(h.recent(ip), time.Now())
}

func (h *AuthHandler) clearFailures(ip string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.failures, ip)
}

// recent drops failures older than the lock window. Caller must hold h.mu.
func (h *AuthHandler) recent(ip string) []time.Time {
	cutoff := time.Now().Add(-loginLockWindow)
	var kept []time.Time
	for _, t := range h.failures[ip] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	h.failures[ip] = kept
	return kept
}

// memberLookup answers the middleware: which role does this member have, and may they sign in.
type memberLookup struct {
	users *db.DashboardUserStore
}

func (m memberLookup) Lookup(username string) (string, bool, bool) {
	u, err := m.users.Get(username)
	if err != nil || u == nil {
		return "", false, false
	}
	return u.Role, !u.Disabled, true
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func writeAuthJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
