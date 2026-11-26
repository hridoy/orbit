package server

import (
	"encoding/json"
	"net/http"

	"orbit/internal/integrations"
)

// IntegrationsHandler manages saved third-party credentials.
type IntegrationsHandler struct {
	mgr *integrations.Manager
}

func NewIntegrationsHandler(mgr *integrations.Manager) *IntegrationsHandler {
	return &IntegrationsHandler{mgr: mgr}
}

func (h *IntegrationsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/integrations", h.handleList)
	mux.HandleFunc("PUT /api/v1/integrations/{provider}", h.handleSet)
	mux.HandleFunc("DELETE /api/v1/integrations/{provider}", h.handleDelete)
}

func (h *IntegrationsHandler) handleList(w http.ResponseWriter, r *http.Request) {
	list, err := h.mgr.List()
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, list)
}

func (h *IntegrationsHandler) handleSet(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var body struct {
		Secret string `json:"secret"`
		Extra  string `json:"extra"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := h.mgr.Set(r.PathValue("provider"), body.Secret, body.Extra); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]string{"message": "saved"})
}

func (h *IntegrationsHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	if err := h.mgr.Delete(r.PathValue("provider")); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]string{"message": "removed"})
}
