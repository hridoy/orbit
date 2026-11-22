package database

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"orbit/internal/compute"
)

type LibSQLEngine struct {
	sqldBin string
}

func NewLibSQLEngine() *LibSQLEngine {
	sqldBin, _ := findBinary("sqld", []string{
		"/usr/bin/sqld",
		"/usr/local/bin/sqld",
		"/opt/homebrew/bin/sqld",
	})
	return &LibSQLEngine{sqldBin: sqldBin}
}

func (le *LibSQLEngine) IsAvailable() bool {
	return le.sqldBin != ""
}

func (le *LibSQLEngine) Provision(
	ctx context.Context,
	baseDir string,
	adminPassword string,
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

	confPath := filepath.Join(confDir, "sqld.conf")
	confContent := GenerateLibSQLConf(port, dataDir, profile)
	if err := os.WriteFile(confPath, []byte(confContent), 0600); err != nil {
		return fmt.Errorf("failed to write sqld.conf: %w", err)
	}

	dbPath := filepath.Join(dataDir, "libsql.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if err := os.WriteFile(dbPath, []byte(""), 0600); err != nil {
			return fmt.Errorf("failed to initialize libsql.db: %w", err)
		}
	}

	return nil
}

func (le *LibSQLEngine) Start(
	ctx context.Context,
	baseDir string,
	rb *compute.RingBuffer,
	port int,
) (*exec.Cmd, error) {
	if le.sqldBin == "" {
		return nil, fmt.Errorf("native sqld binary not found on host (install Turso sqld or run 'brew install tursodatabase/tap/sqld')")
	}

	dataDir := filepath.Join(baseDir, "data", "libsql.db")
	cmd := exec.CommandContext(ctx, le.sqldBin,
		"--db-path", dataDir,
		"--http-listen-addr", fmt.Sprintf("127.0.0.1:%d", port),
	)

	cmd.Stdout = rb
	cmd.Stderr = rb
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start sqld: %w", err)
	}

	rb.Append(fmt.Sprintf("[libsql] Process started (PID: %d)", cmd.Process.Pid))
	return cmd, nil
}
