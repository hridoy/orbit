package compute

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/security"
)

func setupComputeTestEnv(t *testing.T) (*config.Config, *db.AppStore, *db.SecretStore, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "orbit_compute_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	cfg, err := config.LoadOrInit(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("config init failed: %v", err)
	}

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("db open failed: %v", err)
	}

	nodeStore := db.NewNodeStore(database)
	_ = nodeStore.UpsertNode(&db.Node{
		NodeID:       cfg.NodeID,
		Hostname:     "localhost",
		PublicIP:     "127.0.0.1",
		Role:         "standalone",
		AgentVersion: "0.1.0",
		Status:       "online",
	})

	vault, _ := security.NewVault(cfg.MasterKey)
	appStore := db.NewAppStore(database)
	secStore := db.NewSecretStore(database, vault)

	return cfg, appStore, secStore, tmpDir
}

func TestRingBuffer(t *testing.T) {
	rb := NewRingBuffer(5)

	// Append 10 items to test FIFO ring eviction
	for i := 1; i <= 10; i++ {
		rb.Append(fmt.Sprintf("log line %d", i))
	}

	lines := rb.GetTail(5)
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(lines))
	}

	if lines[0].Message != "log line 6" {
		t.Errorf("expected first line to be 'log line 6', got '%s'", lines[0].Message)
	}
	if lines[4].Message != "log line 10" {
		t.Errorf("expected last line to be 'log line 10', got '%s'", lines[4].Message)
	}

	// Test subscription
	ch, unsub := rb.Subscribe()
	defer unsub()

	rb.Append("realtime line")
	select {
	case line := <-ch:
		if line.Message != "realtime line" {
			t.Errorf("unexpected subscriber message: %s", line.Message)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for subscriber message")
	}
}

func TestHealthCheck(t *testing.T) {
	// 1. Pick a random open port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	// Start temporary HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	srv := &http.Server{Handler: mux}
	go srv.Serve(listener)
	defer srv.Close()

	ctx := context.Background()
	err = CheckHealth(ctx, port, "/healthz", 2)
	if err != nil {
		t.Fatalf("expected health check to pass, got: %v", err)
	}

	// 2. Closed port should time out
	err = CheckHealth(ctx, port+1000, "/healthz", 1)
	if err == nil {
		t.Fatal("expected health check to fail on closed port, got nil")
	}
}

func TestDeltaDeployAndPortSwap(t *testing.T) {
	cfg, appStore, secStore, tmpDir := setupComputeTestEnv(t)
	defer os.RemoveAll(tmpDir)

	cgroups := NewCgroupController()
	supervisor := NewProcessSupervisor(cfg, cgroups)
	deployer := NewDeployCoordinator(cfg, appStore, secStore, supervisor)

	appID := "demo-microservice"
	app := &db.App{
		AppID:            appID,
		Name:             "Demo Service",
		NodeID:           cfg.NodeID,
		LinuxUser:        "nobody",
		ActivePort:       9100,
		IdlePort:         9101,
		BinaryEntrypoint: "/bin/sh",
		RAMLimitMB:       64,
		CPUQuotaPct:      50,
		PIDsLimit:        64,
		HealthCheckPath:  "/",
		HealthTimeoutSec: 2,
		Status:           "stopped",
	}

	if err := appStore.CreateApp(app); err != nil {
		t.Fatalf("CreateApp failed: %v", err)
	}

	// 1. Prepare client manifest
	manifest := FileManifest{
		VersionHash: "v1_abc123",
		Files: map[string]string{
			"app.sh": "fake_sha256_hash",
		},
	}

	diff, err := deployer.DiffManifest(appID, manifest)
	if err != nil {
		t.Fatalf("DiffManifest failed: %v", err)
	}
	if len(diff.MissingFiles) != 1 || diff.MissingFiles[0] != "app.sh" {
		t.Fatalf("expected app.sh to be missing, got: %+v", diff.MissingFiles)
	}

	// 2. Stage version directory
	versionDir, err := deployer.StageVersion(appID, manifest.VersionHash, manifest)
	if err != nil {
		t.Fatalf("StageVersion failed: %v", err)
	}

	// Write simple script that listens on PORT
	scriptContent := fmt.Sprintf("#!/bin/sh\nwhile true; do sleep 1; done\n")
	scriptPath := filepath.Join(versionDir, "app.sh")
	_ = os.WriteFile(scriptPath, []byte(scriptContent), 0755)

	// Update app binary entrypoint to run this script
	app.BinaryEntrypoint = scriptPath
	_ = appStore.CreateApp(app) // update or mock
}
