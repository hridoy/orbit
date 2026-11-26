package database

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"orbit/internal/compute"
)

type PostgresEngine struct {
	postgresBin string
	initdbBin   string
}

func NewPostgresEngine() *PostgresEngine {
	pgBin, _ := findBinary("postgres", []string{
		"/usr/lib/postgresql/16/bin/postgres",
		"/usr/lib/postgresql/15/bin/postgres",
		"/usr/lib/postgresql/14/bin/postgres",
		"/usr/local/opt/postgresql@16/bin/postgres",
		"/opt/homebrew/opt/postgresql@16/bin/postgres",
		"/usr/bin/postgres",
	})

	initdbBin, _ := findBinary("initdb", []string{
		"/usr/lib/postgresql/16/bin/initdb",
		"/usr/lib/postgresql/15/bin/initdb",
		"/usr/lib/postgresql/14/bin/initdb",
		"/usr/local/opt/postgresql@16/bin/initdb",
		"/opt/homebrew/opt/postgresql@16/bin/initdb",
		"/usr/bin/initdb",
	})

	return &PostgresEngine{
		postgresBin: pgBin,
		initdbBin:   initdbBin,
	}
}

func (pe *PostgresEngine) IsAvailable() bool {
	return pe.postgresBin != "" && pe.initdbBin != ""
}

// Provision initializes a fresh cluster using initdb with SCRAM-SHA-256 auth and micro tuning.
func (pe *PostgresEngine) Provision(
	ctx context.Context,
	baseDir string,
	adminUser, adminPassword string,
	port int,
	profile RAMProfile,
) error {
	dataDir := filepath.Join(baseDir, "data")
	runDir := filepath.Join(baseDir, "run")
	confDir := filepath.Join(baseDir, "conf")

	for _, d := range []string{dataDir, runDir, confDir} {
		if err := os.MkdirAll(d, 0700); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}

	// If initdb is available, run it
	if pe.initdbBin != "" {
		pwFile := filepath.Join(baseDir, ".initdb_pw")
		if err := os.WriteFile(pwFile, []byte(adminPassword), 0600); err != nil {
			return fmt.Errorf("failed to write pwfile: %w", err)
		}
		defer os.Remove(pwFile)

		if err := chownDBTree(baseDir); err != nil {
			return fmt.Errorf("prepare data dir: %w", err)
		}
		cmd := exec.CommandContext(ctx, pe.initdbBin,
			"-D", dataDir,
			"-U", adminUser,
			fmt.Sprintf("--pwfile=%s", pwFile),
			"--auth=scram-sha-256",
			"-E", "UTF8",
		)
		if err := applyDBUser(cmd); err != nil {
			return err
		}
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("initdb failed: %s: %w", string(out), err)
		}
	}

	// Write postgresql.conf and pg_hba.conf
	confPath := filepath.Join(confDir, "postgresql.conf")
	confContent := GeneratePostgresConf(port, runDir, profile)
	if err := os.WriteFile(confPath, []byte(confContent), 0600); err != nil {
		return fmt.Errorf("failed to write postgresql.conf: %w", err)
	}

	hbaPath := filepath.Join(dataDir, "pg_hba.conf")
	hbaContent := GeneratePostgresHbaConf(adminUser)
	_ = os.WriteFile(hbaPath, []byte(hbaContent), 0600)

	return nil
}

// Start spawns the postgres server process.
func (pe *PostgresEngine) Start(
	ctx context.Context,
	baseDir string,
	rb *compute.RingBuffer,
) (*exec.Cmd, error) {
	if pe.postgresBin == "" {
		return nil, fmt.Errorf("native postgres binary not found on host")
	}

	dataDir := filepath.Join(baseDir, "data")
	confPath := filepath.Join(baseDir, "conf", "postgresql.conf")

	cmd := exec.CommandContext(ctx, pe.postgresBin,
		"-D", dataDir,
		"-c", fmt.Sprintf("config_file=%s", confPath),
	)

	cmd.Stdout = rb
	cmd.Stderr = rb
	if err := applyDBUser(cmd); err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start postgres: %w", err)
	}

	rb.Append(fmt.Sprintf("[postgres] Process started (PID: %d)", cmd.Process.Pid))
	return cmd, nil
}

func findBinary(name string, fallbacks []string) (string, error) {
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}
	for _, fb := range fallbacks {
		if info, err := os.Stat(fb); err == nil && !info.IsDir() {
			return fb, nil
		}
	}
	return "", fmt.Errorf("binary %s not found", name)
}
