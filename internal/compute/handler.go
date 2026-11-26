package compute

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"orbit/internal/db"
	"orbit/internal/security"
)

type ComputeHandler struct {
	appStore   *db.AppStore
	secStore   *db.SecretStore
	cgroups    CgroupController
	supervisor *ProcessSupervisor
	deployer   *DeployCoordinator
}

func NewComputeHandler(
	appStore *db.AppStore,
	secStore *db.SecretStore,
	cgroups CgroupController,
	supervisor *ProcessSupervisor,
	deployer *DeployCoordinator,
) *ComputeHandler {
	return &ComputeHandler{
		appStore:   appStore,
		secStore:   secStore,
		cgroups:    cgroups,
		supervisor: supervisor,
		deployer:   deployer,
	}
}

// RestoreApps redeploys apps that were live before the last agent restart.
func (h *ComputeHandler) RestoreApps() {
	h.deployer.RestoreApps()
}

func (h *ComputeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/apps", h.handleCreateApp)
	mux.HandleFunc("GET /api/v1/apps", h.handleListApps)
	mux.HandleFunc("GET /api/v1/apps/{id}", h.handleGetApp)
	mux.HandleFunc("POST /api/v1/apps/{id}/deploy/diff", h.handleDeployDiff)
	mux.HandleFunc("POST /api/v1/apps/{id}/deploy/upload", h.handleDeployUpload)
	mux.HandleFunc("POST /api/v1/apps/{id}/deploy/activate", h.handleDeployActivate)
	mux.HandleFunc("GET /api/v1/apps/{id}/logs", h.handleGetLogs)
	mux.HandleFunc("GET /api/v1/apps/{id}/logs/stream", h.handleStreamLogs)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *ComputeHandler) handleCreateApp(w http.ResponseWriter, r *http.Request) {
	var app db.App
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if app.AppID == "" || app.Name == "" {
		writeError(w, http.StatusBadRequest, "app_id and name are required")
		return
	}

	// Default fallback ports if not specified
	if app.ActivePort == 0 {
		app.ActivePort = 8001
	}
	if app.IdlePort == 0 {
		app.IdlePort = 8002
	}
	if app.BinaryEntrypoint == "" {
		app.BinaryEntrypoint = "./server"
	}
	if app.RAMLimitMB == 0 {
		app.RAMLimitMB = 256
	}
	if app.CPUQuotaPct == 0 {
		app.CPUQuotaPct = 100
	}
	if app.PIDsLimit == 0 {
		app.PIDsLimit = 128
	}
	if app.HealthCheckPath == "" {
		app.HealthCheckPath = "/healthz"
	}
	if app.HealthTimeoutSec == 0 {
		app.HealthTimeoutSec = 10
	}
	app.Status = "stopped"
	if app.NodeID == "" {
		app.NodeID = h.deployer.cfg.NodeID // default to this node
	}

	if err := h.appStore.CreateApp(&app); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, app)
}

func (h *ComputeHandler) handleListApps(w http.ResponseWriter, r *http.Request) {
	apps, err := h.appStore.ListApps()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if apps == nil {
		apps = []*db.App{}
	}
	writeJSON(w, http.StatusOK, apps)
}

func (h *ComputeHandler) handleGetApp(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("id")
	app, err := h.appStore.GetApp(appID)
	if err != nil || app == nil {
		writeError(w, http.StatusNotFound, "app not found")
		return
	}

	// Fetch live cgroups stats
	stats, _ := h.cgroups.GetStats(appID)

	writeJSON(w, http.StatusOK, map[string]any{
		"app":   app,
		"stats": stats,
	})
}

func (h *ComputeHandler) handleDeployDiff(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("id")
	var manifest FileManifest
	if err := json.NewDecoder(r.Body).Decode(&manifest); err != nil {
		writeError(w, http.StatusBadRequest, "invalid manifest body")
		return
	}

	diff, err := h.deployer.DiffManifest(appID, manifest)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Prepare version staging folder
	_, _ = h.deployer.StageVersion(appID, manifest.VersionHash, manifest)

	writeJSON(w, http.StatusOK, diff)
}

func (h *ComputeHandler) handleDeployUpload(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("id")
	q := r.URL.Query()
	versionHash := q.Get("version")
	relPath := q.Get("path")

	if versionHash == "" || relPath == "" {
		writeError(w, http.StatusBadRequest, "missing version or path query params")
		return
	}

	versionDir := filepath.Join(h.deployer.AppBaseDir(appID), "versions", versionHash)
	dstPath := filepath.Join(versionDir, relPath)

	_ = os.MkdirAll(filepath.Dir(dstPath), 0700)
	dst, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "uploaded",
		"written": written,
		"path":    relPath,
	})
}

func (h *ComputeHandler) handleDeployActivate(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("id")
	var body struct {
		VersionHash string `json:"version_hash"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.VersionHash == "" {
		writeError(w, http.StatusBadRequest, "version_hash is required")
		return
	}

	claims := security.GetClaims(r.Context())
	deployedBy := "cli"
	if claims != nil {
		deployedBy = claims.TokenID
	}

	deployment, err := h.deployer.ActivateDeployment(r.Context(), appID, body.VersionHash, deployedBy)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, deployment)
}

func (h *ComputeHandler) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("id")
	tailStr := r.URL.Query().Get("tail")
	tail := 100
	if val, err := strconv.Atoi(tailStr); err == nil && val > 0 {
		tail = val
	}

	rb := h.supervisor.GetRingBuffer(appID)
	lines := rb.GetTail(tail)

	writeJSON(w, http.StatusOK, map[string]any{
		"app_id": appID,
		"lines":  lines,
	})
}

func (h *ComputeHandler) handleStreamLogs(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("id")
	rb := h.supervisor.GetRingBuffer(appID)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	// Emit initial tail
	initial := rb.GetTail(50)
	for _, l := range initial {
		data, _ := json.Marshal(l)
		_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
	}
	flusher.Flush()

	ch, unsubscribe := rb.Subscribe()
	defer unsubscribe()

	for {
		select {
		case <-r.Context().Done():
			return
		case line, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(line)
			_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

var _ = strings.TrimSpace
