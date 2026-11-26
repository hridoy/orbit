package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
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

// AuthHandler runs first-run setup and password login for the dashboard. Login returns a
// signed bearer token, so the rest of the API keeps using the same auth path.
type AuthHandler struct {
	admins   *db.AdminStore
	tokenMgr *security.TokenManager

	mu       sync.Mutex
	failures map[string][]time.Time // client IP -> recent failed logins
}

func NewAuthHandler(admins *db.AdminStore, tokenMgr *security.TokenManager) *AuthHandler {
	return &AuthHandler{admins: admins, tokenMgr: tokenMgr, failures: map[string][]time.Time{}}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/auth/status", h.handleStatus)
	mux.HandleFunc("POST /api/v1/auth/setup", h.handleSetup)
	mux.HandleFunc("POST /api/v1/auth/login", h.handleLogin)
}

func (h *AuthHandler) handleStatus(w http.ResponseWriter, r *http.Request) {
	acc, err := h.admins.Get()
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]any{"setup_required": acc == nil})
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

// handleSetup creates the operator account. It works once, and only from the local machine,
// so a stranger on the network cannot claim a fresh install.
func (h *AuthHandler) handleSetup(w http.ResponseWriter, r *http.Request) {
	if !isLoopback(r) {
		writeAuthJSON(w, http.StatusForbidden, map[string]string{"error": "first-run setup must be done from this machine (SSH tunnel)"})
		return
	}
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
	if err := h.admins.Create(c.Username, hash); err != nil {
		if errors.Is(err, db.ErrAdminExists) {
			writeAuthJSON(w, http.StatusConflict, map[string]string{"error": "setup is already complete; log in instead"})
			return
		}
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
	acc, err := h.admins.Get()
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if acc == nil {
		writeAuthJSON(w, http.StatusConflict, map[string]string{"error": "setup has not been completed yet"})
		return
	}
	// Check both fields every time so response timing does not reveal which one was wrong.
	userOK := c.Username == acc.Username
	passOK := security.VerifyPassword(c.Password, acc.PasswordHash)
	if !userOK || !passOK {
		h.recordFailure(ip)
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "wrong username or password"})
		return
	}
	h.clearFailures(ip)
	h.issue(w, acc.Username)
}

func (h *AuthHandler) issue(w http.ResponseWriter, username string) {
	id := make([]byte, 8)
	_, _ = rand.Read(id)
	token, err := h.tokenMgr.IssueToken("dashboard-"+hex.EncodeToString(id), "admin", dashboardSessionTTL)
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

func isLoopback(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
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
