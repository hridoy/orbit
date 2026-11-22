package jail

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"path/filepath"
)

type Handler struct {
	jail *Jail
}

func NewHandler(jail *Jail) *Handler {
	return &Handler{jail: jail}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/jail/list", h.handleList)
	mux.HandleFunc("GET /api/v1/jail/file", h.handleReadFile)
	mux.HandleFunc("POST /api/v1/jail/file", h.handleWriteFile)
	mux.HandleFunc("DELETE /api/v1/jail/file", h.handleDeleteFile)
	mux.HandleFunc("GET /api/v1/jail/stat", h.handleStat)
}

func parseScopeAndID(r *http.Request) (ScopeType, string, string, error) {
	q := r.URL.Query()
	scopeStr := q.Get("scope")
	scopeID := q.Get("id")
	path := q.Get("path")

	if scopeStr != string(ScopeApp) && scopeStr != string(ScopeStorage) {
		return "", "", "", ErrInvalidScope
	}
	if scopeID == "" {
		return "", "", "", ErrTargetNotFound
	}
	return ScopeType(scopeStr), scopeID, path, nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	scope, id, path, err := parseScopeAndID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	items, err := h.jail.ListDirectory(scope, id, path)
	if err != nil {
		if errors.Is(err, ErrAccessDenied) {
			writeError(w, http.StatusForbidden, err)
			return
		}
		if errors.Is(err, ErrFileNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"scope": scope,
		"id":    id,
		"path":  path,
		"items": items,
	})
}

func (h *Handler) handleReadFile(w http.ResponseWriter, r *http.Request) {
	scope, id, path, err := parseScopeAndID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reader, size, err := h.jail.OpenFileRead(scope, id, path)
	if err != nil {
		if errors.Is(err, ErrAccessDenied) {
			writeError(w, http.StatusForbidden, err)
			return
		}
		if errors.Is(err, ErrFileNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer reader.Close()

	// Sniff or guess MIME type
	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", http.DetectContentType([]byte(""))) // Will be overwritten
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Stream directly to client with minimal memory buffer (32KB buffer default in io.Copy)
	_, _ = io.Copy(w, reader)
	_ = size
}

func (h *Handler) handleWriteFile(w http.ResponseWriter, r *http.Request) {
	scope, id, path, err := parseScopeAndID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if path == "" || path == "/" {
		writeError(w, http.StatusBadRequest, errors.New("cannot overwrite root directory"))
		return
	}

	written, err := h.jail.WriteFile(scope, id, path, r.Body)
	if err != nil {
		if errors.Is(err, ErrAccessDenied) {
			writeError(w, http.StatusForbidden, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "file uploaded successfully",
		"path":    path,
		"bytes":   written,
	})
}

func (h *Handler) handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	scope, id, path, err := parseScopeAndID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if path == "" || path == "/" {
		writeError(w, http.StatusBadRequest, ErrCannotDeleteRoot)
		return
	}

	err = h.jail.Delete(scope, id, path)
	if err != nil {
		if errors.Is(err, ErrAccessDenied) || errors.Is(err, ErrCannotDeleteRoot) {
			writeError(w, http.StatusForbidden, err)
			return
		}
		if errors.Is(err, ErrFileNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "deleted successfully",
		"path":    path,
	})
}

func (h *Handler) handleStat(w http.ResponseWriter, r *http.Request) {
	scope, id, path, err := parseScopeAndID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	item, err := h.jail.Stat(scope, id, path)
	if err != nil {
		if errors.Is(err, ErrAccessDenied) {
			writeError(w, http.StatusForbidden, err)
			return
		}
		if errors.Is(err, ErrFileNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}
