package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"orbit/internal/db"
	"orbit/internal/tlsengine"
)

// TLSHandler manages the dashboard's domain and certificate API.
type TLSHandler struct {
	engine *tlsengine.Engine
}

func NewTLSHandler(engine *tlsengine.Engine) *TLSHandler {
	return &TLSHandler{engine: engine}
}

func (h *TLSHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/tls/domains", h.handleList)
	mux.HandleFunc("POST /api/v1/tls/domains", h.handleAdd)
	mux.HandleFunc("DELETE /api/v1/tls/domains/{domain}", h.handleDelete)
	mux.HandleFunc("GET /api/v1/tls/detect", h.handleDetect)
}

func (h *TLSHandler) handleList(w http.ResponseWriter, r *http.Request) {
	list, err := h.engine.List()
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if list == nil {
		list = []*db.TLSDomain{}
	}
	writeAuthJSON(w, http.StatusOK, list)
}

func (h *TLSHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	var req tlsengine.AddRequest
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	rec, err := h.engine.Add(req)
	if err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	code := http.StatusAccepted
	if rec.Status == tlsengine.StatusActive {
		code = http.StatusCreated
	}
	writeAuthJSON(w, code, rec)
}

func (h *TLSHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	if err := h.engine.Remove(r.PathValue("domain")); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]string{"message": "domain removed"})
}

func (h *TLSHandler) handleDetect(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	det, err := tlsengine.Detect(ctx, r.URL.Query().Get("domain"))
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, context.DeadlineExceeded) {
			code = http.StatusGatewayTimeout
		}
		writeAuthJSON(w, code, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, det)
}
