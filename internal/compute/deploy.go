package compute

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"orbit/internal/config"
	"orbit/internal/db"
)

type FileManifest struct {
	VersionHash string            `json:"version_hash"`
	Files       map[string]string `json:"files"` // relPath -> sha256
}

type DiffResult struct {
	MissingFiles []string `json:"missing_files"`
	VersionHash  string   `json:"version_hash"`
}

type DeployCoordinator struct {
	cfg        *config.Config
	appStore   *db.AppStore
	secStore   *db.SecretStore
	supervisor *ProcessSupervisor
}

func NewDeployCoordinator(
	cfg *config.Config,
	appStore *db.AppStore,
	secStore *db.SecretStore,
	supervisor *ProcessSupervisor,
) *DeployCoordinator {
	return &DeployCoordinator{
		cfg:        cfg,
		appStore:   appStore,
		secStore:   secStore,
		supervisor: supervisor,
	}
}

// AppBaseDir returns /opt/orbit/apps/<app_id>
func (d *DeployCoordinator) AppBaseDir(appID string) string {
	return filepath.Join(d.cfg.AppsDir, appID)
}

// DiffManifest calculates which files must be uploaded for a new release.
func (d *DeployCoordinator) DiffManifest(appID string, clientManifest FileManifest) (*DiffResult, error) {
	appDir := d.AppBaseDir(appID)
	currentDir := filepath.Join(appDir, "current")

	missing := make([]string, 0)

	for relPath, clientHash := range clientManifest.Files {
		localPath := filepath.Join(currentDir, relPath)
		localHash, err := hashFileSHA256(localPath)
		if err != nil || localHash != clientHash {
			missing = append(missing, relPath)
		}
	}

	return &DiffResult{
		MissingFiles: missing,
		VersionHash:  clientManifest.VersionHash,
	}, nil
}

// StageVersion prepares the new version directory, copying unchanged files via hard links.
func (d *DeployCoordinator) StageVersion(appID, versionHash string, clientManifest FileManifest) (string, error) {
	appDir := d.AppBaseDir(appID)
	versionDir := filepath.Join(appDir, "versions", versionHash)
	currentDir := filepath.Join(appDir, "current")

	if err := os.MkdirAll(versionDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create version directory: %w", err)
	}

	// For unchanged files, create hardlinks from current version to conserve space and I/O
	for relPath, expectedHash := range clientManifest.Files {
		srcFile := filepath.Join(currentDir, relPath)
		dstFile := filepath.Join(versionDir, relPath)

		_ = os.MkdirAll(filepath.Dir(dstFile), 0700)

		if srcHash, err := hashFileSHA256(srcFile); err == nil && srcHash == expectedHash {
			// Hardlink unchanged file
			if _, err := os.Stat(dstFile); os.IsNotExist(err) {
				_ = os.Link(srcFile, dstFile)
			}
		}
	}

	return versionDir, nil
}

// ActivateDeployment launches candidate process on idle_port, checks health, and performs zero-downtime swap.
func (d *DeployCoordinator) ActivateDeployment(
	ctx context.Context,
	appID, versionHash, deployedBy string,
) (*db.AppDeployment, error) {
	startTime := time.Now()

	app, err := d.appStore.GetApp(appID)
	if err != nil || app == nil {
		return nil, fmt.Errorf("app %s not found: %w", appID, err)
	}

	versionDir := filepath.Join(d.AppBaseDir(appID), "versions", versionHash)
	if _, err := os.Stat(versionDir); err != nil {
		return nil, fmt.Errorf("version directory %s does not exist", versionHash)
	}

	// 1. Fetch in-memory decrypted secrets
	envVars, err := d.secStore.GetScopeEnvMap("app", appID)
	if err != nil {
		return nil, fmt.Errorf("failed to load secrets: %w", err)
	}

	candidatePort := app.IdlePort
	oldActivePort := app.ActivePort
	rb := d.supervisor.GetRingBuffer(appID)
	rb.Append(fmt.Sprintf("[deploy] Launching version %s on idle port %d...", versionHash, candidatePort))

	// 2. Launch candidate process on idle_port
	candidateInstance, err := d.supervisor.StartProcess(ctx, app, candidatePort, versionDir, envVars)
	if err != nil {
		return nil, fmt.Errorf("failed to start candidate process: %w", err)
	}

	// 3. Health check verification
	rb.Append(fmt.Sprintf("[deploy] Running health check against port %d (path: %s)...", candidatePort, app.HealthCheckPath))
	healthErr := CheckHealth(ctx, candidatePort, app.HealthCheckPath, app.HealthTimeoutSec)
	if healthErr != nil {
		rb.Append(fmt.Sprintf("[deploy] ❌ Health check failed: %v. Rolling back candidate process.", healthErr))
		_ = d.supervisor.StopProcess(candidateInstance, 2*time.Second)
		return nil, fmt.Errorf("health check failed: %w", healthErr)
	}

	rb.Append(fmt.Sprintf("[deploy] ✅ Health check passed! Performing zero-downtime port swap (%d -> %d)...", oldActivePort, candidatePort))

	// 4. Atomically swap current symlink
	currentSymlink := filepath.Join(d.AppBaseDir(appID), "current")
	tmpSymlink := filepath.Join(d.AppBaseDir(appID), ".current_tmp")
	_ = os.Remove(tmpSymlink)
	if err := os.Symlink(versionDir, tmpSymlink); err == nil {
		_ = os.Rename(tmpSymlink, currentSymlink)
	}

	// 5. Update SQLite state with swapped ports
	if err := d.appStore.SwapPorts(appID, candidatePort, oldActivePort, versionHash); err != nil {
		rb.Append(fmt.Sprintf("[deploy] ⚠️ Failed to update SQLite port state: %v", err))
	}

	// 6. Gracefully drain previous active process
	oldInstance := d.supervisor.GetActive(appID)
	d.supervisor.RegisterActive(appID, candidateInstance)

	if oldInstance != nil {
		rb.Append(fmt.Sprintf("[deploy] Draining old process (PID %d) on port %d...", oldInstance.PID, oldActivePort))
		go func(inst *ProcessInstance) {
			_ = d.supervisor.StopProcess(inst, 10*time.Second)
			rb.Append(fmt.Sprintf("[deploy] Old process on port %d cleanly stopped.", oldActivePort))
		}(oldInstance)
	}

	durationMS := time.Since(startTime).Milliseconds()
	deployment := &db.AppDeployment{
		DeploymentID: fmt.Sprintf("dep_%d", time.Now().UnixNano()),
		AppID:        appID,
		VersionHash:  versionHash,
		PortUsed:     candidatePort,
		Status:       "live",
		DeployedBy:   deployedBy,
		DurationMS:   durationMS,
		Logs:         fmt.Sprintf("Successfully deployed in %dms", durationMS),
		CreatedAt:    time.Now().UTC(),
	}

	_ = d.appStore.RecordDeployment(deployment)
	rb.Append(fmt.Sprintf("[deploy] 🚀 Deployment %s is LIVE on port %d in %dms!", versionHash, candidatePort, durationMS))

	return deployment, nil
}

func hashFileSHA256(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

var _ = errors.New

// isRestorable reports whether an app's last live version still exists on disk.
func isRestorable(appsDir string, app *db.App) bool {
	if app.CurrentVersion == nil || *app.CurrentVersion == "" {
		return false
	}
	versionDir := filepath.Join(appsDir, app.AppID, "versions", *app.CurrentVersion)
	info, err := os.Stat(versionDir)
	return err == nil && info.IsDir()
}

// RestoreApps starts every app that was live before the agent last stopped. Processes are
// children of the agent and die with it, so each one is redeployed to the same version through
// the normal health-checked activation path. It runs apps one at a time.
func (d *DeployCoordinator) RestoreApps() {
	apps, err := d.appStore.ListApps()
	if err != nil {
		log.Printf("restore: list apps: %v", err)
		return
	}
	for _, app := range apps {
		if !isRestorable(d.cfg.AppsDir, app) {
			continue
		}
		version := *app.CurrentVersion
		rb := d.supervisor.GetRingBuffer(app.AppID)
		rb.Append(fmt.Sprintf("[restore] agent restarted; starting version %s", version))
		if _, err := d.ActivateDeployment(context.Background(), app.AppID, version, "restore"); err != nil {
			rb.Append(fmt.Sprintf("[restore] failed: %v", err))
			log.Printf("restore %s: %v", app.AppID, err)
			continue
		}
		log.Printf("restore %s: live on version %s", app.AppID, version)
	}
}
