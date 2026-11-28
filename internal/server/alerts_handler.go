package server

import (
	"fmt"
	"net/http"
	"time"

	"orbit/internal/config"
	"orbit/internal/sysinfo"
	"orbit/internal/tlsengine"
)

// Alert is one thing the operator should look at. Level is "warning" or "critical".
type Alert struct {
	Level  string `json:"level"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// AlertsHandler reports disk, memory, certificate, and backup problems on a single VPS, so the
// dashboard can warn before something breaks.
type AlertsHandler struct {
	rootPath string
	tls      *tlsengine.Engine
	backups  *BackupHandler
}

func NewAlertsHandler(cfg *config.Config, tls *tlsengine.Engine, backups *BackupHandler) *AlertsHandler {
	return &AlertsHandler{rootPath: cfg.RootPath, tls: tls, backups: backups}
}

func (h *AlertsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/alerts", h.handleList)
}

func (h *AlertsHandler) handleList(w http.ResponseWriter, r *http.Request) {
	writeAuthJSON(w, http.StatusOK, h.Check(time.Now()))
}

// Check evaluates every rule at the given time. It is separate from the HTTP handler so the
// thresholds can be tested.
func (h *AlertsHandler) Check(now time.Time) []Alert {
	alerts := []Alert{}

	if used := sysinfo.DiskUsedPercent(h.rootPath); used >= 0 {
		switch {
		case used >= 95:
			alerts = append(alerts, Alert{"critical", "Disk almost full", fmt.Sprintf("%d%% of the disk is used. Free space before databases or backups fail.", used)})
		case used >= 85:
			alerts = append(alerts, Alert{"warning", "Disk filling up", fmt.Sprintf("%d%% of the disk is used.", used)})
		}
	}

	if avail := sysinfo.MemoryAvailablePercent(); avail >= 0 && avail < 10 {
		alerts = append(alerts, Alert{"warning", "Low memory", fmt.Sprintf("Only %d%% of memory is available.", avail)})
	}

	if h.tls != nil {
		if domains, err := h.tls.List(); err == nil {
			for _, d := range domains {
				if d.Status != "active" && d.Status != "" {
					alerts = append(alerts, Alert{"critical", "Certificate not active", fmt.Sprintf("%s is %s.", d.Domain, d.Status)})
					continue
				}
				if d.ExpiresAt != nil && d.ExpiresAt.Sub(now) < 14*24*time.Hour {
					alerts = append(alerts, Alert{"warning", "Certificate expiring", fmt.Sprintf("%s expires %s.", d.Domain, d.ExpiresAt.Format("2006-01-02"))})
				}
			}
		}
	}

	if h.backups != nil {
		names := h.backups.backupsOnDisk()
		if len(names) == 0 {
			alerts = append(alerts, Alert{"critical", "No backups", "No control-plane backup exists yet. The nightly timer should create one."})
		} else if t, ok := parseBackupName(names[0]); ok && now.Sub(t) > 36*time.Hour {
			alerts = append(alerts, Alert{"warning", "Backup is stale", fmt.Sprintf("The newest backup is from %s.", t.Format("2006-01-02 15:04 UTC"))})
		}
		h.backups.mu.Lock()
		last := h.backups.last
		h.backups.mu.Unlock()
		if last != nil && !last.OK {
			alerts = append(alerts, Alert{"critical", "Last backup failed", "The most recent backup run returned an error. Open Credentials for the output."})
		}
	}
	return alerts
}

// parseBackupName reads the UTC timestamp from orbit-backup-YYYYMMDD-HHMMSS.
func parseBackupName(name string) (time.Time, bool) {
	const prefix = "orbit-backup-"
	if len(name) < len(prefix)+15 {
		return time.Time{}, false
	}
	t, err := time.Parse("20060102-150405", name[len(prefix):len(prefix)+15])
	return t, err == nil
}
