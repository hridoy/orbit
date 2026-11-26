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

const panelDomainKey = "panel_domain"

// SettingsHandler manages server-wide settings. The panel domain is the public name of the
// dashboard. Setting it issues a certificate for it and routes it to the dashboard.
type SettingsHandler struct {
	store  *db.ServerSettingsStore
	tls    *tlsengine.Engine
	domain atomic.Value // current panel domain, "" when unset
}

func NewSettingsHandler(store *db.ServerSettingsStore, tls *tlsengine.Engine) (*SettingsHandler, error) {
	h := &SettingsHandler{store: store, tls: tls}
	d, err := store.Get(panelDomainKey)
	if err != nil {
		return nil, err
	}
	h.domain.Store(d)
	return h, nil
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
	writeAuthJSON(w, http.StatusOK, map[string]string{"panel_domain": h.PanelDomain()})
}

type settingsBody struct {
	PanelDomain string `json:"panel_domain"`
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

	// An empty panel domain turns the setting off. The dashboard stays on the server IP.
	if strings.TrimSpace(body.PanelDomain) == "" {
		if err := h.store.Set(panelDomainKey, ""); err != nil {
			writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		h.domain.Store("")
		writeAuthJSON(w, http.StatusOK, map[string]string{"panel_domain": ""})
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
