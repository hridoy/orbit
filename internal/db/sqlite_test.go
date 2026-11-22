package db

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"orbit/internal/security"
)

func setupTestDB(t *testing.T) (*DB, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "orbit_db_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test.db")
	database, err := Open(dbPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to open database: %v", err)
	}

	return database, tmpDir
}

func TestSQLitePragmasAndMigration(t *testing.T) {
	database, tmpDir := setupTestDB(t)
	defer database.Close()
	defer os.RemoveAll(tmpDir)

	var journalMode string
	err := database.QueryRow("PRAGMA journal_mode;").Scan(&journalMode)
	if err != nil {
		t.Fatalf("failed to query journal mode: %v", err)
	}
	if journalMode != "wal" {
		t.Errorf("expected journal_mode wal, got %s", journalMode)
	}

	var foreignKeys int
	err = database.QueryRow("PRAGMA foreign_keys;").Scan(&foreignKeys)
	if err != nil {
		t.Fatalf("failed to query foreign_keys: %v", err)
	}
	if foreignKeys != 1 {
		t.Errorf("expected foreign_keys 1, got %d", foreignKeys)
	}
}

func TestNodeStore(t *testing.T) {
	database, tmpDir := setupTestDB(t)
	defer database.Close()
	defer os.RemoveAll(tmpDir)

	store := NewNodeStore(database)
	node := &Node{
		NodeID:          "node_test_01",
		Hostname:        "test-host",
		PublicIP:        "192.168.1.10",
		Role:            "standalone",
		CPUCores:        4,
		RAMTotalBytes:   8589934592,
		DiskTotalBytes:  107374182400,
		AgentVersion:    "0.1.0",
		Status:          "online",
		LastHeartbeatAt: time.Now(),
	}

	if err := store.UpsertNode(node); err != nil {
		t.Fatalf("UpsertNode failed: %v", err)
	}

	retrieved, err := store.GetNode("node_test_01")
	if err != nil {
		t.Fatalf("GetNode failed: %v", err)
	}
	if retrieved == nil || retrieved.Hostname != "test-host" {
		t.Errorf("unexpected retrieved node: %+v", retrieved)
	}
}

func TestSecretStoreWithVault(t *testing.T) {
	database, tmpDir := setupTestDB(t)
	defer database.Close()
	defer os.RemoveAll(tmpDir)

	masterKey := make([]byte, 32)
	for i := range masterKey {
		masterKey[i] = byte(i)
	}

	vault, err := security.NewVault(masterKey)
	if err != nil {
		t.Fatalf("NewVault failed: %v", err)
	}

	secretStore := NewSecretStore(database, vault)

	err = secretStore.SetSecret("sec_1", "app", "app_demo", "DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	if err != nil {
		t.Fatalf("SetSecret failed: %v", err)
	}

	val, err := secretStore.GetSecret("app", "app_demo", "DATABASE_URL")
	if err != nil {
		t.Fatalf("GetSecret failed: %v", err)
	}
	if val != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("expected decrypted secret, got: %s", val)
	}

	envMap, err := secretStore.GetScopeEnvMap("app", "app_demo")
	if err != nil {
		t.Fatalf("GetScopeEnvMap failed: %v", err)
	}
	if envMap["DATABASE_URL"] != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("expected env map to contain key, got: %+v", envMap)
	}
}
