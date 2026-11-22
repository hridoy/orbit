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

type MongoEngine struct {
	mongodBin string
}

func NewMongoEngine() *MongoEngine {
	mongodBin, _ := findBinary("mongod", []string{
		"/usr/bin/mongod",
		"/usr/local/bin/mongod",
		"/opt/homebrew/bin/mongod",
	})
	return &MongoEngine{mongodBin: mongodBin}
}

func (me *MongoEngine) IsAvailable() bool {
	return me.mongodBin != ""
}

func (me *MongoEngine) Provision(
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

	confPath := filepath.Join(confDir, "mongod.conf")
	confContent := GenerateMongoConf(port, runDir, dataDir, profile)
	if err := os.WriteFile(confPath, []byte(confContent), 0600); err != nil {
		return fmt.Errorf("failed to write mongod.conf: %w", err)
	}

	return nil
}

func (me *MongoEngine) Start(
	ctx context.Context,
	baseDir string,
	rb *compute.RingBuffer,
) (*exec.Cmd, error) {
	if me.mongodBin == "" {
		return nil, fmt.Errorf("native mongod binary not found on host")
	}

	confPath := filepath.Join(baseDir, "conf", "mongod.conf")
	cmd := exec.CommandContext(ctx, me.mongodBin, "--config", confPath)

	cmd.Stdout = rb
	cmd.Stderr = rb
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start mongod: %w", err)
	}

	rb.Append(fmt.Sprintf("[mongod] Process started (PID: %d)", cmd.Process.Pid))
	return cmd, nil
}
