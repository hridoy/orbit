package server

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync/atomic"

	"orbit/internal/db"
	"orbit/internal/tlsengine"
)

const (
	panelDomainKey = "panel_domain"
	lockKey        = "lock_to_panel_domain"
)

// SettingsHandler manages server-wide settings. The panel domain is the public name of the
// dashboard. Setting it issues a certificate for it and routes it to the dashboard.
type SettingsHandler struct {
	store  *db.ServerSettingsStore
	tls    *tlsengine.Engine
	domain atomic.Value // current panel domain, "" when unset
	locked atomic.Bool  // when true, the dashboard answers only on the panel domain
}

func NewSettingsHandler(store *db.ServerSettingsStore, tls *tlsengine.Engine) (*SettingsHandler, error) {
	h := &SettingsHandler{store: store, tls: tls}
	d, err := store.Get(panelDomainKey)
	if err != nil {
		return nil, err
	}
	h.domain.Store(d)
	lock, err := store.Get(lockKey)
	if err != nil {
		return nil, err
	}
	h.locked.Store(lock == "1" && d != "")
	return h, nil
}

// Locked reports whether remote access to the dashboard is limited to the panel domain.
func (h *SettingsHandler) Locked() bool {
	return h.locked.Load()
}

// Blocks reports whether a request must be refused because the dashboard is locked to its
// panel domain. Local requests (the SSH tunnel and the public proxy) are never blocked, so the
// operator can always get back in.
func (h *SettingsHandler) Blocks(r *http.Request) bool {
	if !h.Locked() {
		return false
	}
	if isLoopbackAddr(r.RemoteAddr) {
		return false
	}
	return hostOnly(r.Host) != h.PanelDomain()
}

func isLoopbackAddr(remote string) bool {
	host, _, err := net.SplitHostPort(remote)
	if err != nil {
		host = remote
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// PanelDomain returns the configured panel domain, or "".
func (h *SettingsHandler) PanelDomain() string {
	d, _ := h.domain.Load().(string)
	return d
}

func (h *SettingsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/server/settings", h.handleGet)
	mux.HandleFunc("PUT /api/v1/server/settings", h.handlePut)
}

func (h *SettingsHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	writeAuthJSON(w, http.StatusOK, map[string]any{"panel_domain": h.PanelDomain(), "lock_to_panel_domain": h.Locked()})
}

type settingsBody struct {
	PanelDomain string `json:"panel_domain"`
	LockToPanel *bool  `json:"lock_to_panel_domain,omitempty"`
	Provider    string `json:"provider"`
	EABKID      string `json:"eab_kid,omitempty"`
	EABHMAC     string `json:"eab_hmac,omitempty"`
	CertPEM     string `json:"cert_pem,omitempty"`
	KeyPEM      string `json:"key_pem,omitempty"`
}

func (h *SettingsHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	var body settingsBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// An empty panel domain turns the setting off, and the lock with it. The dashboard stays on the server IP.
	if strings.TrimSpace(body.PanelDomain) == "" {
		if err := h.store.Set(panelDomainKey, ""); err != nil {
			writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		_ = h.store.Set(lockKey, "")
		h.domain.Store("")
		h.locked.Store(false)
		writeAuthJSON(w, http.StatusOK, map[string]any{"panel_domain": "", "lock_to_panel_domain": false})
		return
	}

	domain, err := tlsengine.NormalizeDomain(body.PanelDomain)
	if err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	rec, err := h.tls.Add(tlsengine.AddRequest{
		Domain: domain, Provider: body.Provider, EABKID: body.EABKID, EABHMAC: body.EABHMAC,
		CertPEM: body.CertPEM, KeyPEM: body.KeyPEM,
	})
	if err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := h.store.Set(panelDomainKey, domain); err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.domain.Store(domain)

	// Locking is allowed only once the panel's certificate is live. Otherwise the operator
	// could lock themselves out of the IP address before HTTPS works.
	if body.LockToPanel != nil {
		if *body.LockToPanel && !h.tls.HasCertificate(domain) {
			writeAuthJSON(w, http.StatusConflict, map[string]string{"error": "wait until the certificate for " + domain + " is active, then turn the lock on"})
			return
		}
		val := "0"
		if *body.LockToPanel {
			val = "1"
		}
		if err := h.store.Set(lockKey, val); err != nil {
			writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		h.locked.Store(*body.LockToPanel)
	}
	code := http.StatusAccepted
	if rec.Status == tlsengine.StatusActive {
		code = http.StatusOK
	}
	writeAuthJSON(w, code, map[string]any{"panel_domain": domain, "certificate": rec})
}

// hostOnly strips a port from a Host header.
func hostOnly(hostport string) string {
	if h, _, err := net.SplitHostPort(hostport); err == nil {
		return strings.ToLower(h)
	}
	return strings.ToLower(hostport)
}
