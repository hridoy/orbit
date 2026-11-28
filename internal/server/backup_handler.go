package server

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"orbit/internal/config"
	"orbit/internal/integrations"
)

// BackupHandler shows backup status and starts a backup on demand. A run is the same
// `orbit-agent backup` command the nightly timer uses, so both paths upload the same way.
type BackupHandler struct {
	dir string
	mgr *integrations.Manager

	mu      sync.Mutex
	running bool
	last    *backupRun
}

type backupRun struct {
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	OK         bool      `json:"ok"`
	Output     string    `json:"output"`
}

func NewBackupHandler(cfg *config.Config, mgr *integrations.Manager) *BackupHandler {
	return &BackupHandler{dir: filepath.Join(cfg.RootPath, "backups"), mgr: mgr}
}

func (h *BackupHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/backup", h.handleStatus)
	mux.HandleFunc("POST /api/v1/backup/run", h.handleRun)
}

// backupsOnDisk returns backup directory names, newest first.
func (h *BackupHandler) backupsOnDisk() []string {
	entries, err := os.ReadDir(h.dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "orbit-backup-") {
			names = append(names, e.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	return names
}

func (h *BackupHandler) handleStatus(w http.ResponseWriter, r *http.Request) {
	backups := h.backupsOnDisk()
	cloud := map[string]any{"configured": false}
	if list, err := h.mgr.List(); err == nil {
		for _, s := range list {
			if s.Provider == integrations.CloudBackup && s.Configured {
				cloud = map[string]any{"configured": true, "remote_path": s.ExtraValue}
			}
		}
	}
	h.mu.Lock()
	running, last := h.running, h.last
	h.mu.Unlock()
	writeAuthJSON(w, http.StatusOK, map[string]any{
		"local_count": len(backups),
		"latest":      first(backups),
		"cloud":       cloud,
		"running":     running,
		"last_run":    last,
	})
}

func (h *BackupHandler) handleRun(w http.ResponseWriter, r *http.Request) {
	exe, err := os.Executable()
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.mu.Lock()
	if h.running {
		h.mu.Unlock()
		writeAuthJSON(w, http.StatusConflict, map[string]string{"error": "a backup is already running"})
		return
	}
	h.running = true
	h.mu.Unlock()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
		defer cancel()
		run := &backupRun{StartedAt: time.Now().UTC()}
		out, err := exec.CommandContext(ctx, exe, "backup", h.dir).CombinedOutput()
		run.FinishedAt = time.Now().UTC()
		run.OK = err == nil
		run.Output = tail(string(out), 4000)
		if err != nil {
			run.Output += "\n" + err.Error()
		}
		h.mu.Lock()
		h.running = false
		h.last = run
		h.mu.Unlock()
	}()
	writeAuthJSON(w, http.StatusAccepted, map[string]string{"message": "backup started"})
}

func first(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

func tail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}
