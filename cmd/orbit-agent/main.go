package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"orbit/internal/buildinfo"
	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/security"
	"orbit/internal/server"
)

func main() {
	fmt.Printf("🪐 Orbit Agent v%s (%s) [%s/%s]\n", buildinfo.Version, buildinfo.GitCommit, runtime.GOOS, runtime.GOARCH)
	fmt.Println("   Sovereign, Zero-Docker, Single-Binary PaaS")
	fmt.Println("------------------------------------------------------------")

	// 1. Initialize directory layout and configuration
	cfg, err := config.LoadOrInit("")
	if err != nil {
		log.Fatalf("❌ Failed to initialize configuration: %v", err)
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

	// 3. Initialize AES-256-GCM Vault
	vault, err := security.NewVault(cfg.MasterKey)
	if err != nil {
		log.Fatalf("❌ Failed to initialize secrets vault: %v", err)
	}

	// 4. Initialize Ed25519 Token Manager
	tokenMgr, err := security.NewTokenManager(cfg.MasterKey)
	if err != nil {
		log.Fatalf("❌ Failed to initialize token manager: %v", err)
	}
	fmt.Printf("🔑 Master Admin Key:     %s\n", tokenMgr.MasterKeyHex())

	// 5. Register / Heartbeat local node record
	hostname, _ := os.Hostname()
	nodeStore := db.NewNodeStore(database)
	err = nodeStore.UpsertNode(&db.Node{
		NodeID:          cfg.NodeID,
		Hostname:        hostname,
		PublicIP:        "127.0.0.1",
		Role:            "standalone",
		CPUCores:        runtime.NumCPU(),
		RAMTotalBytes:   0, // Telemetry collector will populate in Phase 5
		DiskTotalBytes:  0,
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
