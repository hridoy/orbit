package database

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"

	"orbit/internal/db"
)

// WithExposures attaches the store that remembers public ports across restarts.
func (m *DatabaseManager) WithExposures(store *db.DatabaseExposureStore) *DatabaseManager {
	m.dbExposures = store
	return m
}

// accessInfo is what the database page shows for connecting. Internal stays on this server.
// External is the same database through its public port.
type accessInfo struct {
	InternalURI    string `json:"internal_uri"`
	SocketURI      string `json:"socket_uri"`
	ExternalOpen   bool   `json:"external_open"`
	ExternalPort   int    `json:"external_port,omitempty"`
	ExternalURI    string `json:"external_uri,omitempty"`
	PublicHost     string `json:"public_host"`
	SuggestedPorts []int  `json:"suggested_ports"`
}

func (h *DatabaseHandler) handleAccessGet(w http.ResponseWriter, r *http.Request) {
	instance, err := h.dbStore.GetDatabase(r.PathValue("id"))
	if err != nil || instance == nil {
		writeError(w, http.StatusNotFound, "database not found")
		return
	}
	writeJSON(w, http.StatusOK, h.accessFor(r, instance))
}

func (h *DatabaseHandler) accessFor(r *http.Request, instance *db.DatabaseInstance) accessInfo {
	password := ""
	if parts := strings.Split(instance.EncryptedPassword, ":"); len(parts) == 2 {
		if raw, err := h.vault.Decrypt(parts[0], parts[1]); err == nil {
			password = string(raw)
		}
	}
	tcp, sock := h.manager.GetConnectionURIs(instance, password)
	info := accessInfo{
		InternalURI:    tcp,
		SocketURI:      sock,
		PublicHost:     publicHost(r),
		SuggestedPorts: h.manager.SuggestedPorts(instance.Engine, 5),
	}
	if e, _ := h.manager.dbExposures.Get(instance.DBID); e != nil {
		info.ExternalOpen = true
		info.ExternalPort = e.PublicPort
		info.ExternalURI = strings.Replace(tcp, "127.0.0.1:"+strconv.Itoa(instance.Port),
			info.PublicHost+":"+strconv.Itoa(e.PublicPort), 1)
	}
	return info
}

// publicHost is the name the operator reaches this server by: the host of the request, unless
// it was a loopback address (an SSH tunnel), in which case it is left for the operator to fill.
func publicHost(r *http.Request) string {
	host := r.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	if host == "" || host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return "<your-server-ip>"
	}
	return host
}

func (h *DatabaseHandler) handleAccessPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	instance, err := h.dbStore.GetDatabase(id)
	if err != nil || instance == nil {
		writeError(w, http.StatusNotFound, "database not found")
		return
	}
	var body struct {
		Enabled bool `json:"enabled"`
		Port    int  `json:"port"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10)).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Enabled {
		if err := h.manager.Expose(id, body.Port); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		if err := h.manager.Unexpose(id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	writeJSON(w, http.StatusOK, h.accessFor(r, instance))
}

func (h *DatabaseHandler) handlePortSuggestions(w http.ResponseWriter, r *http.Request) {
	engine := r.URL.Query().Get("engine")
	if _, ok := suggestedBase[engine]; !ok {
		writeError(w, http.StatusBadRequest, "unknown engine")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"engine": engine, "ports": h.manager.SuggestedPorts(engine, 8)})
}
