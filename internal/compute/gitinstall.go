package compute

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/gitsource"
	"orbit/internal/security"
)

const (
	installTimeout = 20 * time.Minute
	checkTimeout   = 20 * time.Second
	buildTimeout   = 15 * time.Minute
)

// GitHandler installs apps from GitHub, checks for new commits, and updates them.
type GitHandler struct {
	cfg        *config.Config
	appStore   *db.AppStore
	sources    *db.AppSourceStore
	supervisor *ProcessSupervisor
	deployer   *DeployCoordinator

	mu   sync.Mutex
	busy map[string]bool // app IDs with an install or update in progress
}

func NewGitHandler(
	cfg *config.Config,
	appStore *db.AppStore,
	sources *db.AppSourceStore,
	supervisor *ProcessSupervisor,
	deployer *DeployCoordinator,
) *GitHandler {
	return &GitHandler{
		cfg:        cfg,
		appStore:   appStore,
		sources:    sources,
		supervisor: supervisor,
		deployer:   deployer,
		busy:       make(map[string]bool),
	}
}

func (h *GitHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/apps/install-git", h.handleInstall)
	mux.HandleFunc("GET /api/v1/apps/{id}/source", h.handleSource)
	mux.HandleFunc("POST /api/v1/apps/{id}/source/update", h.handleUpdate)
	mux.HandleFunc("GET /api/v1/updates", h.handleUpdates)
}

// requireDeployer allows admin and deployer tokens to change apps.
func requireDeployer(w http.ResponseWriter, r *http.Request) bool {
	claims := security.GetClaims(r.Context())
	if claims == nil || (claims.Role != "admin" && claims.Role != "deployer") {
		writeError(w, http.StatusForbidden, "admin or deployer role required")
		return false
	}
	return true
}

func (h *GitHandler) tryLock(appID string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.busy[appID] {
		return false
	}
	h.busy[appID] = true
	return true
}

func (h *GitHandler) unlock(appID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.busy, appID)
}

// pickPorts finds a free pair of loopback ports that no app is using.
func pickPorts(apps []*db.App) (int, int, error) {
	used := map[int]bool{}
	for _, a := range apps {
		used[a.ActivePort] = true
		used[a.IdlePort] = true
	}
	free := func(p int) bool {
		if used[p] {
			return false
		}
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p))
		if err != nil {
			return false
		}
		_ = ln.Close()
		return true
	}
	for p := 8001; p < 9000; p += 2 {
		if free(p) && free(p+1) {
			return p, p + 1, nil
		}
	}
	return 0, 0, fmt.Errorf("no free port pair available")
}

func randomSuffix() string {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeError(w, code, msg)
}

// handleInstall creates an app from a GitHub repo and starts the build in the background.
func (h *GitHandler) handleInstall(w http.ResponseWriter, r *http.Request) {
	if !requireDeployer(w, r) {
		return
	}
	var body struct {
		RepoURL string `json:"repo_url"`
		Ref     string `json:"ref"`
		Name    string `json:"name"`
		AppID   string `json:"app_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid request body")
		return
	}
	repo, err := gitsource.NormalizeRepoURL(body.RepoURL)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	ref := body.Ref
	if ref == "" {
		ref = "main"
	}
	if err := gitsource.ValidateRef(ref); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	appID := body.AppID
	if appID == "" {
		appID = "app_" + randomSuffix()
	}
	if existing, _ := h.appStore.GetApp(appID); existing != nil {
		writeErr(w, http.StatusConflict, "app_id already exists; use the update endpoint instead")
		return
	}
	name := body.Name
	if name == "" {
		parts := strings.Split(repo, "/")
		name = parts[len(parts)-1]
	}

	apps, err := h.appStore.ListApps()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	activePort, idlePort, err := pickPorts(apps)
	if err != nil {
		writeErr(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	app := &db.App{
		AppID:            appID,
		Name:             name,
		NodeID:           h.cfg.NodeID,
		ActivePort:       activePort,
		IdlePort:         idlePort,
		BinaryEntrypoint: "./server",
		RAMLimitMB:       256,
		CPUQuotaPct:      100,
		PIDsLimit:        128,
		HealthCheckPath:  "/healthz",
		HealthTimeoutSec: 10,
		AutoRestart:      true,
		Status:           "stopped",
	}
	if err := h.appStore.CreateApp(app); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.sources.Upsert(&db.AppSource{AppID: appID, RepoURL: repo, GitRef: ref}); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !h.tryLock(appID) {
		writeErr(w, http.StatusConflict, "another operation is running for this app")
		return
	}
	go h.runJob(appID, repo, ref, "install")

	writeJSON(w, http.StatusAccepted, map[string]any{"app": app, "status": "installing"})
}

// handleSource returns where an app came from and its update state.
func (h *GitHandler) handleSource(w http.ResponseWriter, r *http.Request) {
	src, err := h.sources.Get(r.PathValue("id"))
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if src == nil {
		writeErr(w, http.StatusNotFound, "app was not installed from GitHub")
		return
	}
	writeJSON(w, http.StatusOK, src)
}

// handleUpdate deploys the latest commit of an app. The old version stays live until the new one passes its health check.
func (h *GitHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if !requireDeployer(w, r) {
		return
	}
	appID := r.PathValue("id")
	src, err := h.sources.Get(appID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if src == nil {
		writeErr(w, http.StatusNotFound, "app was not installed from GitHub")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), checkTimeout)
	defer cancel()
	latest, err := gitsource.LatestCommit(ctx, src.RepoURL, src.GitRef)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	_ = h.sources.SetLatest(appID, latest)
	if latest == src.InstalledCommit {
		writeJSON(w, http.StatusOK, map[string]string{"status": "up_to_date", "commit": latest})
		return
	}

	if !h.tryLock(appID) {
		writeErr(w, http.StatusConflict, "another operation is running for this app")
		return
	}
	go h.runJob(appID, src.RepoURL, src.GitRef, "update")
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "updating", "to": latest})
}

// handleUpdates lists GitHub apps that have a newer commit than the one installed.
func (h *GitHandler) handleUpdates(w http.ResponseWriter, r *http.Request) {
	sources, err := h.sources.List()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	type item struct {
		AppID           string `json:"app_id"`
		Name            string `json:"name"`
		RepoURL         string `json:"repo_url"`
		GitRef          string `json:"git_ref"`
		InstalledCommit string `json:"installed_commit"`
		LatestCommit    string `json:"latest_commit"`
	}
	updates := []item{}
	for _, src := range sources {
		ctx, cancel := context.WithTimeout(r.Context(), checkTimeout)
		latest, err := gitsource.LatestCommit(ctx, src.RepoURL, src.GitRef)
		cancel()
		if err != nil {
			continue // keep the last known state when GitHub is unreachable
		}
		_ = h.sources.SetLatest(src.AppID, latest)
		if latest == src.InstalledCommit {
			continue
		}
		name := src.AppID
		if app, _ := h.appStore.GetApp(src.AppID); app != nil {
			name = app.Name
		}
		updates = append(updates, item{
			AppID: src.AppID, Name: name, RepoURL: src.RepoURL, GitRef: src.GitRef,
			InstalledCommit: src.InstalledCommit, LatestCommit: latest,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"count": len(updates), "updates": updates})
}

// runJob clones, builds, and activates one commit. It logs into the app's log buffer.
func (h *GitHandler) runJob(appID, repo, ref, kind string) {
	defer h.unlock(appID)
	rb := h.supervisor.GetRingBuffer(appID)
	logf := func(format string, args ...any) {
		rb.Append(fmt.Sprintf("[%s] "+format, append([]any{kind}, args...)...))
	}
	logf("fetching %s@%s", repo, ref)

	versionsDir := filepath.Join(h.deployer.AppBaseDir(appID), "versions")
	if err := os.MkdirAll(versionsDir, 0700); err != nil {
		logf("failed: %v", err)
		return
	}

	// Clone into a temp dir first; rename to the commit name once the commit is known.
	tmp := filepath.Join(versionsDir, ".fetch-"+randomSuffix())
	fetchCtx, cancelFetch := context.WithTimeout(context.Background(), installTimeout)
	commit, err := gitsource.Fetch(fetchCtx, repo, ref, tmp)
	cancelFetch()
	if err != nil {
		_ = os.RemoveAll(tmp)
		logf("fetch failed: %v", err)
		return
	}
	dest := filepath.Join(versionsDir, commit)
	if _, err := os.Stat(dest); err == nil {
		_ = os.RemoveAll(tmp) // already built
	} else if err := os.Rename(tmp, dest); err != nil {
		_ = os.RemoveAll(tmp)
		logf("failed to stage version: %v", err)
		return
	}
	logf("commit %s", commit)

	manifest, err := gitsource.ReadManifest(dest)
	if err != nil {
		logf("manifest error: %v", err)
		return
	}

	logf("building: %s", manifest.Build)
	buildCtx, cancelBuild := context.WithTimeout(context.Background(), buildTimeout)
	cmd := exec.CommandContext(buildCtx, "sh", "-c", manifest.Build)
	cmd.Dir = dest
	cmd.Stdout = rb
	cmd.Stderr = rb
	buildErr := cmd.Run()
	cancelBuild()
	if buildErr != nil {
		logf("build failed: %v (the running version is unchanged)", buildErr)
		return
	}
	if _, err := os.Stat(filepath.Join(dest, manifest.Entrypoint)); err != nil {
		logf("entrypoint %s not found after build", manifest.Entrypoint)
		return
	}

	if err := h.sources.SetEntrypoint(appID, manifest.Entrypoint); err != nil {
		logf("failed to save entrypoint: %v", err)
		return
	}

	// Activation starts the process, so it must not inherit a request or build context.
	if _, err := h.deployer.ActivateDeployment(context.Background(), appID, commit, "github"); err != nil {
		logf("deploy failed: %v (the running version is unchanged)", err)
		return
	}
	if err := h.sources.SetInstalled(appID, commit); err != nil {
		logf("failed to record commit: %v", err)
		return
	}
	logf("live at commit %s", commit[:12])
}
