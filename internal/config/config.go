package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the core runtime configuration for Orbit.
type Config struct {
	RootPath     string `json:"root_path"`     // Canonical root (/opt/orbit or ~/.orbit)
	DataDir      string `json:"data_dir"`      // RootPath/data
	AppsDir      string `json:"apps_dir"`      // RootPath/apps
	DatabasesDir string `json:"databases_dir"` // RootPath/databases
	StorageDir   string `json:"storage_dir"`   // RootPath/storage
	SecretsDir   string `json:"secrets_dir"`   // RootPath/secrets
	RunDir       string `json:"run_dir"`       // RootPath/run (sockets & pidfiles)
	DBPath       string `json:"db_path"`       // DataDir/orbit.db
	MasterKey    []byte `json:"-"`             // 32-byte AES-GCM master key
	ListenAddr   string `json:"listen_addr"`   // Host:Port (e.g. 127.0.0.1:8787 or :8787)
	NodeID       string `json:"node_id"`       // Unique Node Identifier
}

// DefaultRootPath resolves the sovereign directory root.
// Prefers $ORBIT_ROOT, then /opt/orbit if writable or root, else fallback to ~/.orbit.
func DefaultRootPath() string {
	if envRoot := os.Getenv("ORBIT_ROOT"); envRoot != "" {
		return filepath.Clean(envRoot)
	}

	optOrbit := "/opt/orbit"
	// Check if we are root or can write to /opt
	if os.Geteuid() == 0 {
		return optOrbit
	}
	if info, err := os.Stat(optOrbit); err == nil && info.IsDir() {
		// Test write permission
		testFile := filepath.Join(optOrbit, ".perm_test")
		if err := os.WriteFile(testFile, []byte("ok"), 0600); err == nil {
			_ = os.Remove(testFile)
			return optOrbit
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "./.orbit"
	}
	return filepath.Join(home, ".orbit")
}

// LoadOrInit loads the existing configuration or creates the directory hierarchy.
func LoadOrInit(customRoot string) (*Config, error) {
	root := customRoot
	if root == "" {
		root = DefaultRootPath()
	}

	cfg := &Config{
		RootPath:     root,
		DataDir:      filepath.Join(root, "data"),
		AppsDir:      filepath.Join(root, "apps"),
		DatabasesDir: filepath.Join(root, "databases"),
		StorageDir:   filepath.Join(root, "storage"),
		SecretsDir:   filepath.Join(root, "secrets"),
		RunDir:       filepath.Join(root, "run"),
		ListenAddr:   "127.0.0.1:8787",
	}
	cfg.DBPath = filepath.Join(cfg.DataDir, "orbit.db")

	if addr := os.Getenv("ORBIT_LISTEN_ADDR"); addr != "" {
		cfg.ListenAddr = addr
	}

	// Create required directories with 0700 permissions
	dirs := []string{
		cfg.RootPath,
		cfg.DataDir,
		cfg.AppsDir,
		cfg.DatabasesDir,
		cfg.StorageDir,
		cfg.SecretsDir,
		cfg.RunDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Load or initialize Master Key (32 bytes for AES-256-GCM)
	masterKeyPath := filepath.Join(cfg.SecretsDir, "master.key")
	if keyBytes, err := os.ReadFile(masterKeyPath); err == nil && len(keyBytes) == 32 {
		cfg.MasterKey = keyBytes
	} else {
		newKey := make([]byte, 32)
		if _, err := rand.Read(newKey); err != nil {
			return nil, fmt.Errorf("failed to generate random master key: %w", err)
		}
		if err := os.WriteFile(masterKeyPath, newKey, 0600); err != nil {
			return nil, fmt.Errorf("failed to persist master key to %s: %w", masterKeyPath, err)
		}
		cfg.MasterKey = newKey
	}

	// Load or initialize Node ID
	nodeIDPath := filepath.Join(cfg.DataDir, "node_id")
	if idBytes, err := os.ReadFile(nodeIDPath); err == nil && len(idBytes) > 0 {
		cfg.NodeID = string(idBytes)
	} else {
		raw := make([]byte, 8)
		_, _ = rand.Read(raw)
		cfg.NodeID = "node_" + hex.EncodeToString(raw)
		_ = os.WriteFile(nodeIDPath, []byte(cfg.NodeID), 0600)
	}

	return cfg, nil
}
