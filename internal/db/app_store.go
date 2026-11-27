package db

import (
	"database/sql"
	"fmt"
	"time"
)

type App struct {
	AppID            string    `json:"app_id"`
	Name             string    `json:"name"`
	Domain           *string   `json:"domain,omitempty"`
	NodeID           string    `json:"node_id"`
	LinuxUser        string    `json:"linux_user"`
	ActivePort       int       `json:"active_port"`
	IdlePort         int       `json:"idle_port"`
	BinaryEntrypoint string    `json:"binary_entrypoint"`
	ExecArgs         string    `json:"exec_args"`
	RAMLimitMB       int       `json:"ram_limit_mb"`
	CPUQuotaPct      int       `json:"cpu_quota_pct"`
	PIDsLimit        int       `json:"pids_limit"`
	HealthCheckPath  string    `json:"health_check_path"`
	HealthTimeoutSec int       `json:"health_timeout_sec"`
	AutoRestart      bool      `json:"auto_restart"`
	Status           string    `json:"status"`
	CurrentVersion   *string   `json:"current_version,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AppDeployment struct {
	DeploymentID string    `json:"deployment_id"`
	AppID        string    `json:"app_id"`
	VersionHash  string    `json:"version_hash"`
	PortUsed     int       `json:"port_used"`
	Status       string    `json:"status"`
	DeployedBy   string    `json:"deployed_by"`
	DurationMS   int64     `json:"duration_ms"`
	Logs         string    `json:"logs"`
	CreatedAt    time.Time `json:"created_at"`
}

type AppStore struct {
	db *DB
}

func NewAppStore(db *DB) *AppStore {
	return &AppStore{db: db}
}

func (s *AppStore) CreateApp(app *App) error {
	query := `
	INSERT INTO apps (
		app_id, name, domain, node_id, linux_user,
		active_port, idle_port, binary_entrypoint, exec_args,
		ram_limit_mb, cpu_quota_pct, pids_limit,
		health_check_path, health_timeout_sec, auto_restart, status
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	autoRestartInt := 0
	if app.AutoRestart {
		autoRestartInt = 1
	}

	_, err := s.db.Exec(query,
		app.AppID, app.Name, app.Domain, app.NodeID, app.LinuxUser,
		app.ActivePort, app.IdlePort, app.BinaryEntrypoint, app.ExecArgs,
		app.RAMLimitMB, app.CPUQuotaPct, app.PIDsLimit,
		app.HealthCheckPath, app.HealthTimeoutSec, autoRestartInt, app.Status,
	)
	if err != nil {
		return fmt.Errorf("create app %s: %w", app.AppID, err)
	}
	return nil
}

func (s *AppStore) GetApp(appID string) (*App, error) {
	query := `
	SELECT app_id, name, domain, node_id, linux_user,
	       active_port, idle_port, binary_entrypoint, exec_args,
	       ram_limit_mb, cpu_quota_pct, pids_limit,
	       health_check_path, health_timeout_sec, auto_restart,
	       status, current_version, created_at, updated_at
	FROM apps WHERE app_id = ?;
	`
	row := s.db.QueryRow(query, appID)
	var a App
	var autoRestartInt int
	err := row.Scan(
		&a.AppID, &a.Name, &a.Domain, &a.NodeID, &a.LinuxUser,
		&a.ActivePort, &a.IdlePort, &a.BinaryEntrypoint, &a.ExecArgs,
		&a.RAMLimitMB, &a.CPUQuotaPct, &a.PIDsLimit,
		&a.HealthCheckPath, &a.HealthTimeoutSec, &autoRestartInt,
		&a.Status, &a.CurrentVersion, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get app %s: %w", appID, err)
	}
	a.AutoRestart = autoRestartInt == 1
	return &a, nil
}

func (s *AppStore) ListApps() ([]*App, error) {
	query := `
	SELECT app_id, name, domain, node_id, linux_user,
	       active_port, idle_port, binary_entrypoint, exec_args,
	       ram_limit_mb, cpu_quota_pct, pids_limit,
	       health_check_path, health_timeout_sec, auto_restart,
	       status, current_version, created_at, updated_at
	FROM apps ORDER BY created_at DESC;
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("list apps: %w", err)
	}
	defer rows.Close()

	var apps []*App
	for rows.Next() {
		var a App
		var autoRestartInt int
		err := rows.Scan(
			&a.AppID, &a.Name, &a.Domain, &a.NodeID, &a.LinuxUser,
			&a.ActivePort, &a.IdlePort, &a.BinaryEntrypoint, &a.ExecArgs,
			&a.RAMLimitMB, &a.CPUQuotaPct, &a.PIDsLimit,
			&a.HealthCheckPath, &a.HealthTimeoutSec, &autoRestartInt,
			&a.Status, &a.CurrentVersion, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan app: %w", err)
		}
		a.AutoRestart = autoRestartInt == 1
		apps = append(apps, &a)
	}
	return apps, nil
}

func (s *AppStore) UpdateAppStatus(appID, status string) error {
	query := `UPDATE apps SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE app_id = ?;`
	_, err := s.db.Exec(query, status, appID)
	return err
}

func (s *AppStore) SwapPorts(appID string, newActivePort, newIdlePort int, newVersion string) error {
	query := `
	UPDATE apps SET 
		active_port = ?, 
		idle_port = ?, 
		current_version = ?, 
		status = 'healthy', 
		updated_at = CURRENT_TIMESTAMP 
	WHERE app_id = ?;
	`
	_, err := s.db.Exec(query, newActivePort, newIdlePort, newVersion, appID)
	return err
}

func (s *AppStore) RecordDeployment(d *AppDeployment) error {
	query := `
	INSERT INTO app_deployments (
		deployment_id, app_id, version_hash, port_used, status, deployed_by, duration_ms, logs
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err := s.db.Exec(query,
		d.DeploymentID, d.AppID, d.VersionHash, d.PortUsed, d.Status, d.DeployedBy, d.DurationMS, d.Logs,
	)
	return err
}

func (s *AppStore) ListDeployments(appID string, limit int) ([]*AppDeployment, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	query := `
	SELECT deployment_id, app_id, version_hash, port_used, status, deployed_by, duration_ms, logs, created_at
	FROM app_deployments
	WHERE app_id = ?
	ORDER BY created_at DESC
	LIMIT ?;
	`
	rows, err := s.db.Query(query, appID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*AppDeployment
	for rows.Next() {
		var d AppDeployment
		if err := rows.Scan(&d.DeploymentID, &d.AppID, &d.VersionHash, &d.PortUsed, &d.Status, &d.DeployedBy, &d.DurationMS, &d.Logs, &d.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &d)
	}
	return list, nil
}
