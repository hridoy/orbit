package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"orbit/internal/db"
	"orbit/internal/security"
)

const (
	challengeTTL      = 5 * time.Minute
	maxChallengeTries = 5
)

// loginChallenge is a pending second factor. It is kept in memory and is not a bearer token, so
// it cannot be used to call the API.
type loginChallenge struct {
	username string
	expires  time.Time
	tries    int
}

func (h *AuthHandler) registerTwoFactorRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/login/2fa", h.handleLogin2FA)
	mux.HandleFunc("GET /api/v1/auth/2fa", h.handle2FAStatus)
	mux.HandleFunc("POST /api/v1/auth/2fa/setup", h.handle2FASetup)
	mux.HandleFunc("POST /api/v1/auth/2fa/enable", h.handle2FAEnable)
	mux.HandleFunc("POST /api/v1/auth/2fa/disable", h.handle2FADisable)
	mux.HandleFunc("DELETE /api/v1/users/{name}/2fa", h.handleAdminReset2FA)
}

// beginSecondFactor is called after a correct password for a member with two-factor on. It
// returns a challenge the client must answer with a code.
func (h *AuthHandler) beginSecondFactor(w http.ResponseWriter, username string) {
	raw := make([]byte, 32)
	_, _ = rand.Read(raw)
	id := hex.EncodeToString(raw)
	h.mu.Lock()
	h.challenges[id] = &loginChallenge{username: username, expires: time.Now().Add(challengeTTL)}
	h.mu.Unlock()
	writeAuthJSON(w, http.StatusOK, map[string]any{"two_factor_required": true, "challenge": id})
}

func (h *AuthHandler) handleLogin2FA(w http.ResponseWriter, r *http.Request) {
	ip := clientIP(r)
	if h.locked(ip) {
		writeAuthJSON(w, http.StatusTooManyRequests, map[string]string{"error": "too many failed logins; wait 15 minutes"})
		return
	}
	var body struct {
		Challenge string `json:"challenge"`
		Code      string `json:"code"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	h.mu.Lock()
	ch := h.challenges[body.Challenge]
	if ch == nil || time.Now().After(ch.expires) {
		delete(h.challenges, body.Challenge)
		h.mu.Unlock()
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "the sign-in step expired; log in again"})
		return
	}
	ch.tries++
	if ch.tries > maxChallengeTries {
		delete(h.challenges, body.Challenge)
		h.mu.Unlock()
		h.recordFailure(ip)
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "too many wrong codes; log in again"})
		return
	}
	username := ch.username
	h.mu.Unlock()

	rec, err := h.totp.Get(username)
	if err != nil || rec == nil || !rec.Enabled {
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "two-factor is not active; log in again"})
		return
	}
	secret, err := h.decryptSecret(rec)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read the two-factor secret"})
		return
	}
	step, ok := security.VerifyTOTP(secret, body.Code, time.Now())
	if !ok {
		h.recordFailure(ip)
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "wrong code"})
		return
	}
	if step <= rec.LastStep {
		// The same code was already used. Refuse replays within the window.
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "that code was already used; wait for the next one"})
		return
	}
	rec.LastStep = step
	if err := h.totp.Put(rec); err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.mu.Lock()
	delete(h.challenges, body.Challenge)
	h.mu.Unlock()
	h.clearFailures(ip)
	h.issue(w, username)
}

func (h *AuthHandler) handle2FAStatus(w http.ResponseWriter, r *http.Request) {
	name := memberName(security.GetClaims(r.Context()))
	rec, err := h.totp.Get(name)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]bool{
		"enabled": rec != nil && rec.Enabled,
		"pending": rec != nil && !rec.Enabled,
	})
}

// handle2FASetup creates a new secret that stays pending until the member proves a code works.
func (h *AuthHandler) handle2FASetup(w http.ResponseWriter, r *http.Request) {
	name := memberName(security.GetClaims(r.Context()))
	if name == "" {
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "sign in as a member first"})
		return
	}
	rec, err := h.totp.Get(name)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if rec != nil && rec.Enabled {
		writeAuthJSON(w, http.StatusConflict, map[string]string{"error": "two-factor is already on; turn it off first to set a new device"})
		return
	}
	secret, err := security.GenerateTOTPSecret()
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	ct, nonce, err := h.vault.Encrypt([]byte(secret))
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if err := h.totp.Put(&db.TOTPRecord{Username: name, SecretCiphertext: ct, SecretNonce: nonce}); err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]string{
		"secret": secret,
		"uri":    security.TOTPProvisioningURI(secret, name),
	})
}

// handle2FAEnable turns two-factor on after the member enters a code from the new secret.
func (h *AuthHandler) handle2FAEnable(w http.ResponseWriter, r *http.Request) {
	name := memberName(security.GetClaims(r.Context()))
	var body struct {
		Code string `json:"code"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	rec, err := h.totp.Get(name)
	if err != nil || rec == nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "start setup first"})
		return
	}
	if rec.Enabled {
		writeAuthJSON(w, http.StatusConflict, map[string]string{"error": "two-factor is already on"})
		return
	}
	secret, err := h.decryptSecret(rec)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read the two-factor secret"})
		return
	}
	step, ok := security.VerifyTOTP(secret, body.Code, time.Now())
	if !ok {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "wrong code; check the time on your phone and try again"})
		return
	}
	rec.Enabled = true
	rec.LastStep = step
	if err := h.totp.Put(rec); err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]bool{"enabled": true})
}

// handle2FADisable needs the password, so a stolen session token alone cannot turn it off.
func (h *AuthHandler) handle2FADisable(w http.ResponseWriter, r *http.Request) {
	name := memberName(security.GetClaims(r.Context()))
	var body struct {
		Password string `json:"password"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	u, err := h.users.Get(name)
	if err != nil || u == nil || !security.VerifyPassword(body.Password, u.PasswordHash) {
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "wrong password"})
		return
	}
	if err := h.totp.Delete(name); err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]bool{"enabled": false})
}

// handleAdminReset2FA lets an admin clear a member's two-factor when they lose their phone.
func (h *AuthHandler) handleAdminReset2FA(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if err := h.totp.Delete(name); err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]bool{"enabled": false})
}

func (h *AuthHandler) decryptSecret(rec *db.TOTPRecord) (string, error) {
	plain, err := h.vault.Decrypt(rec.SecretCiphertext, rec.SecretNonce)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
