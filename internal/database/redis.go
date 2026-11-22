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

type RedisEngine struct {
	redisBin string
}

func NewRedisEngine() *RedisEngine {
	redisBin, _ := findBinary("redis-server", []string{
		"/usr/bin/redis-server",
		"/usr/local/bin/redis-server",
		"/opt/homebrew/bin/redis-server",
		"/opt/homebrew/opt/redis/bin/redis-server",
	})
	return &RedisEngine{redisBin: redisBin}
}

func (re *RedisEngine) IsAvailable() bool {
	return re.redisBin != ""
}

func (re *RedisEngine) Provision(
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

	socketFile := filepath.Join(runDir, "redis.sock")
	confPath := filepath.Join(confDir, "redis.conf")
	confContent := GenerateRedisConf(port, socketFile, dataDir, adminPassword, profile)

	if err := os.WriteFile(confPath, []byte(confContent), 0600); err != nil {
		return fmt.Errorf("failed to write redis.conf: %w", err)
	}

	return nil
}

func (re *RedisEngine) Start(
	ctx context.Context,
	baseDir string,
	rb *compute.RingBuffer,
) (*exec.Cmd, error) {
	if re.redisBin == "" {
		return nil, fmt.Errorf("native redis-server binary not found on host (install with 'apt install redis-server' or 'brew install redis')")
	}

	confPath := filepath.Join(baseDir, "conf", "redis.conf")
	cmd := exec.CommandContext(ctx, re.redisBin, confPath)

	cmd.Stdout = rb
	cmd.Stderr = rb
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start redis-server: %w", err)
	}

	rb.Append(fmt.Sprintf("[redis] Process started (PID: %d)", cmd.Process.Pid))
	return cmd, nil
}
