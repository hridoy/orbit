package update

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"

	"orbit/internal/security"
)

// Handler exposes update status and a manual apply action.
type Handler struct {
	stateDir string
	version  string

	mu sync.Mutex // one apply at a time
}

func NewHandler(stateDir, version string) *Handler {
	return &Handler{stateDir: stateDir, version: version}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/platform/update", h.handleStatus)
	mux.HandleFunc("POST /api/v1/platform/update", h.handleApply)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) handleStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	writeJSON(w, http.StatusOK, CheckStatus(ctx, h.version))
}

// handleApply downloads and verifies the latest release, swaps the binary, and exits so systemd
// restarts the agent on the new version. Admin only.
func (h *Handler) handleApply(w http.ResponseWriter, r *http.Request) {
	claims := security.GetClaims(r.Context())
	if claims == nil || claims.Role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "admin role required"})
		return
	}
	if !h.mu.TryLock() {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "an update is already in progress"})
		return
	}
	defer h.mu.Unlock()

	exe, err := os.Executable()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()
	newVersion, err := Apply(ctx, exe)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	if err := MarkPending(h.stateDir, newVersion); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{"status": "restarting", "version": newVersion})
	// Give the response time to flush, then exit. systemd (Restart=always) starts the new binary.
	go func() {
		time.Sleep(1500 * time.Millisecond)
		os.Exit(0)
	}()
}
