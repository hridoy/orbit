package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"orbit/internal/buildinfo"
	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/security"
	"orbit/internal/server"
	"orbit/internal/sysinfo"
	"orbit/internal/update"
)

func main() {
	// Subcommand: `orbit-agent rekey <backup-dir> [--dry-run]` rotates the master key. Stop the agent first.
	if len(os.Args) >= 3 && os.Args[1] == "rekey" {
		cfg, err := config.LoadOrInit("")
		if err != nil {
			log.Fatalf("❌ Failed to initialize configuration: %v", err)
		}
		database, err := db.Open(cfg.DBPath)
		if err != nil {
			log.Fatalf("❌ Failed to open database: %v", err)
		}
		defer database.Close()
		dryRun := len(os.Args) == 4 && os.Args[3] == "--dry-run"
		if err := runRekey(cfg, database, os.Args[2], dryRun); err != nil {
			log.Fatalf("❌ Rekey failed: %v", err)
		}
		return
	}

	// Subcommand: `orbit-agent backup <dir>` writes a state backup and exits.
	if len(os.Args) == 3 && os.Args[1] == "backup" {
		cfg, err := config.LoadOrInit("")
		if err != nil {
			log.Fatalf("❌ Failed to initialize configuration: %v", err)
		}
		database, err := db.Open(cfg.DBPath)
		if err != nil {
			log.Fatalf("❌ Failed to open database: %v", err)
		}
		defer database.Close()
		if err := runBackup(cfg, database, os.Args[2]); err != nil {
			log.Fatalf("❌ Backup failed: %v", err)
		}
		if err := uploadLatest(cfg, database, os.Args[2]); err != nil {
			log.Fatalf("❌ Cloud upload failed (local backup kept): %v", err)
		}
		return
	}

	fmt.Printf("🪐 Orbit Agent v%s (%s) [%s/%s]\n", buildinfo.Version, buildinfo.GitCommit, runtime.GOOS, runtime.GOARCH)
	fmt.Println("   Sovereign, Zero-Docker, Single-Binary PaaS")
	fmt.Println("------------------------------------------------------------")

	// 1. Initialize directory layout and configuration
	cfg, err := config.LoadOrInit("")
	if err != nil {
		log.Fatalf("❌ Failed to initialize configuration: %v", err)
	}
	// Roll back a self-update that keeps crashing, before doing any other work.
	if exe, err := os.Executable(); err == nil {
		restored, err := update.StartupCheck(cfg.DataDir, exe)
		if err != nil {
			log.Printf("⚠️ update state check failed: %v", err)
		}
		if restored {
			log.Printf("↩️ new agent version kept crashing; restored the previous binary")
			os.Exit(1)
		}
		update.MarkHealthyAfterDelay(cfg.DataDir)
	}
	fmt.Printf("📁 Orbit Root Directory: %s\n", cfg.RootPath)
	fmt.Printf("🆔 Node ID:              %s\n", cfg.NodeID)

	// 2. Initialize Embedded SQLite (WAL mode, ~8MB cache)
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("❌ Failed to open embedded SQLite database (%s): %v", cfg.DBPath, err)
	}
	defer database.Close()
	fmt.Printf("🗄️  SQLite WAL Store:     %s\n", cfg.DBPath)

	// 3. Single-agent lock, then resolve the active master key (finishing or rolling back any
	// interrupted rotation). `rekey` cannot run while this lock is held.
	agentLock, err := acquireLock(filepath.Join(cfg.SecretsDir, "agent.lock"))
	if err != nil {
		log.Fatalf("❌ %v", err)
	}
	defer agentLock.release()
	masterKey, err := resolveMasterKey(cfg, database)
	if err != nil {
		log.Fatalf("❌ %v", err)
	}
	cfg.MasterKey = masterKey
	if err := ensureCanary(database, cfg.MasterKey); err != nil {
		log.Fatalf("❌ Failed to write vault canary: %v", err)
	}

	// 4. Initialize AES-256-GCM Vault
	vault, err := security.NewVault(cfg.MasterKey)
	if err != nil {
		log.Fatalf("❌ Failed to initialize secrets vault: %v", err)
	}

	// 5. Ed25519 token manager, keyed by a subkey, and the CLI admin token file
	tokenMgr, err := tokenManagerFor(cfg.MasterKey)
	if err != nil {
		log.Fatalf("❌ Failed to initialize token manager: %v", err)
	}
	if err := ensureAdminToken(cfg, tokenMgr); err != nil {
		log.Fatalf("❌ Failed to write admin token: %v", err)
	}
	// Never print the key itself: stdout ends up in journald and logs.
	fingerprint := sha256.Sum256(cfg.MasterKey)
	fmt.Printf("🔑 Master key fingerprint: %s (stored in %s)\n", hex.EncodeToString(fingerprint[:4]), cfg.SecretsDir)

	// 5. Register / Heartbeat local node record
	hostname, _ := os.Hostname()
	nodeStore := db.NewNodeStore(database)
	publicIP := sysinfo.PrimaryIP()
	if publicIP == "" {
		publicIP = "127.0.0.1"
	}
	err = nodeStore.UpsertNode(&db.Node{
		NodeID:          cfg.NodeID,
		Hostname:        hostname,
		PublicIP:        publicIP,
		Role:            "standalone",
		CPUCores:        runtime.NumCPU(),
		RAMTotalBytes:   sysinfo.MemoryTotal(),
		DiskTotalBytes:  sysinfo.DiskTotal(cfg.RootPath),
		AgentVersion:    buildinfo.Version,
		Status:          "online",
		LastHeartbeatAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("⚠️ Failed to record node in SQLite: %v", err)
	}

	// 6. Start HTTP API Server
	srv := server.NewServer(cfg, database, vault, tokenMgr)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("🚀 Orbit Agent listening on http://%s\n", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ HTTP server error: %v", err)
		}
	}()

	<-stopCh
	fmt.Println("\n🛑 Shutting down Orbit Agent gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Close(); err != nil {
		log.Printf("⚠️ Server shutdown error: %v", err)
	}
	_ = ctx

	fmt.Println("👋 Orbit Agent stopped.")
}
