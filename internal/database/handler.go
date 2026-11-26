package database

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"orbit/internal/db"
	"orbit/internal/security"
)

type DatabaseHandler struct {
	dbStore *db.DatabaseStore
	manager *DatabaseManager
	vault   *security.Vault
}

func NewDatabaseHandler(
	dbStore *db.DatabaseStore,
	manager *DatabaseManager,
	vault *security.Vault,
) *DatabaseHandler {
	return &DatabaseHandler{
		dbStore: dbStore,
		manager: manager,
		vault:   vault,
	}
}

func (h *DatabaseHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/databases", h.handleProvision)
	mux.HandleFunc("GET /api/v1/databases/engines", h.handleEngines)
	mux.HandleFunc("GET /api/v1/databases/port-suggestions", h.handlePortSuggestions)
	mux.HandleFunc("POST /api/v1/databases/engines/{engine}/{version}/install", h.handleInstallVersion)
	mux.HandleFunc("GET /api/v1/databases/{id}/access", h.handleAccessGet)
	mux.HandleFunc("PUT /api/v1/databases/{id}/access", h.handleAccessPut)
	mux.HandleFunc("GET /api/v1/databases", h.handleList)
	mux.HandleFunc("GET /api/v1/databases/{id}", h.handleGet)
	mux.HandleFunc("POST /api/v1/databases/{id}/start", h.handleStart)
	mux.HandleFunc("POST /api/v1/databases/{id}/stop", h.handleStop)
	mux.HandleFunc("POST /api/v1/databases/{id}/restart", h.handleRestart)
	mux.HandleFunc("DELETE /api/v1/databases/{id}", h.handleDelete)
	mux.HandleFunc("POST /api/v1/databases/{id}/backup", h.handleBackup)
	mux.HandleFunc("GET /api/v1/databases/{id}/logs", h.handleGetLogs)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *DatabaseHandler) handleProvision(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name       string `json:"name"`
		Engine     string `json:"engine"`      // postgres, mysql, mariadb, mongodb, redis, libsql
		Version    string `json:"version"`     // an installed version, or empty for the default
		Username   string `json:"username"`    // optional; generated when empty
		Password   string `json:"password"`    // optional; generated when empty
		RAMProfile string `json:"ram_profile"` // micro, standard, heavy, unlimited
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Name == "" || body.Engine == "" {
		writeError(w, http.StatusBadRequest, "name and engine are required")
		return
	}

	profile := RAMProfile(body.RAMProfile)
	if profile != ProfileMicro && profile != ProfileStandard && profile != ProfileHeavy && profile != ProfileUnlimited {
		profile = ProfileMicro
	}

	instance, rawPassword, err := h.manager.ProvisionDatabase(r.Context(), body.Name, body.Engine, body.Version, body.Username, body.Password, profile)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tcpURI, socketURI := h.manager.GetConnectionURIs(instance, rawPassword)

	writeJSON(w, http.StatusCreated, map[string]any{
		"database":   instance,
		"password":   rawPassword,
		"tcp_uri":    tcpURI,
		"socket_uri": socketURI,
	})
}

func (h *DatabaseHandler) handleList(w http.ResponseWriter, r *http.Request) {
	list, err := h.dbStore.ListDatabases()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if list == nil {
		list = []*db.DatabaseInstance{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *DatabaseHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	dbID := r.PathValue("id")
	instance, err := h.dbStore.GetDatabase(dbID)
	if err != nil || instance == nil {
		writeError(w, http.StatusNotFound, "database not found")
		return
	}

	// Decrypt password
	var rawPassword string
	parts := strings.Split(instance.EncryptedPassword, ":")
	if len(parts) == 2 {
		if decrypted, err := h.vault.Decrypt(parts[0], parts[1]); err == nil {
			rawPassword = string(decrypted)
		}
	}

	tcpURI, socketURI := h.manager.GetConnectionURIs(instance, rawPassword)

	writeJSON(w, http.StatusOK, map[string]any{
		"database":   instance,
		"password":   rawPassword,
		"tcp_uri":    tcpURI,
		"socket_uri": socketURI,
	})
}

func (h *DatabaseHandler) handleStart(w http.ResponseWriter, r *http.Request) {
	dbID := r.PathValue("id")
	if err := h.manager.StartDatabase(r.Context(), dbID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "database started successfully"})
}

func (h *DatabaseHandler) handleStop(w http.ResponseWriter, r *http.Request) {
	dbID := r.PathValue("id")
	if err := h.manager.StopDatabase(dbID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "database stopped successfully"})
}

func (h *DatabaseHandler) handleRestart(w http.ResponseWriter, r *http.Request) {
	dbID := r.PathValue("id")
	if err := h.manager.RestartDatabase(r.Context(), dbID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "database restarted successfully"})
}

func (h *DatabaseHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	dbID := r.PathValue("id")
	if err := h.manager.DeleteDatabase(r.Context(), dbID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "database deleted successfully"})
}

func (h *DatabaseHandler) handleBackup(w http.ResponseWriter, r *http.Request) {
	dbID := r.PathValue("id")
	path, err := h.manager.BackupDatabase(r.Context(), dbID)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"path": path})
}

func (h *DatabaseHandler) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	dbID := r.PathValue("id")
	tailStr := r.URL.Query().Get("tail")
	tail := 100
	if val, err := strconv.Atoi(tailStr); err == nil && val > 0 {
		tail = val
	}

	rb := h.manager.GetRingBuffer(dbID)
	lines := rb.GetTail(tail)

	writeJSON(w, http.StatusOK, map[string]any{
		"db_id": dbID,
		"lines": lines,
	})
}

// handleEngines lists every engine and each version in the catalog, with its install state:
// installed, installing, failed, or available.
func (h *DatabaseHandler) handleEngines(w http.ResponseWriter, r *http.Request) {
	type version struct {
		Version string `json:"version"`
		Status  string `json:"status"`
		Method  string `json:"method"`
		Error   string `json:"error,omitempty"`
	}
	type engine struct {
		Engine   string    `json:"engine"`
		Name     string    `json:"name"`
		Versions []version `json:"versions"`
	}
	names := map[string]string{
		"postgres": "PostgreSQL", "mysql": "MySQL", "mariadb": "MariaDB",
		"mongodb": "MongoDB", "redis": "Redis", "libsql": "LibSQL (SQLite)",
	}
	installed := h.manager.Versions()
	out := []engine{}
	for _, e := range []string{"postgres", "mysql", "mariadb", "mongodb", "redis", "libsql"} {
		row := engine{Engine: e, Name: names[e], Versions: []version{}}
		for _, c := range Catalog {
			if c.Engine != e {
				continue
			}
			v := version{Version: c.Version, Status: "available", Method: string(c.Method)}
			if contains(installed[e], c.Version) {
				v.Status = "installed"
			} else if st, msg := h.manager.installer.jobStatus(e, c.Version); st != "" {
				v.Status, v.Error = st, msg
			}
			row.Versions = append(row.Versions, v)
		}
		out = append(out, row)
	}
	writeJSON(w, http.StatusOK, out)
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

// handleInstallVersion starts a background install of a catalog version.
func (h *DatabaseHandler) handleInstallVersion(w http.ResponseWriter, r *http.Request) {
	if err := h.manager.StartInstall(r.PathValue("engine"), r.PathValue("version")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "installing"})
}
