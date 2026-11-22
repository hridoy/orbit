package database

import (
	"context"
	"os"
	"strings"
	"testing"

	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/security"
)

func setupDBTestEnv(t *testing.T) (*config.Config, *DatabaseManager, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "orbit_db_mgr_test_*")
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

	// Register node for foreign key constraint
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
	dbStore := db.NewDatabaseStore(database)
	secStore := db.NewSecretStore(database, vault)

	mgr := NewDatabaseManager(cfg, dbStore, secStore, vault)
	return cfg, mgr, tmpDir
}

func TestTuningProfiles(t *testing.T) {
	microTuning := GetPostgresTuning(ProfileMicro)
	if microTuning.SharedBuffers != "16MB" || microTuning.MaxConnections != 20 {
		t.Errorf("unexpected micro tuning: %+v", microTuning)
	}

	heavyTuning := GetPostgresTuning(ProfileHeavy)
	if heavyTuning.SharedBuffers != "256MB" || heavyTuning.MaxConnections != 100 {
		t.Errorf("unexpected heavy tuning: %+v", heavyTuning)
	}

	// Test config generation
	conf := GeneratePostgresConf(5432, "/tmp/run", ProfileMicro)
	if !strings.Contains(conf, "listen_addresses = '127.0.0.1'") {
		t.Errorf("conf must bind to 127.0.0.1 only")
	}
	if !strings.Contains(conf, "shared_buffers = 16MB") {
		t.Errorf("conf must contain 16MB shared buffers")
	}

	mongoConf := GenerateMongoConf(27017, "/tmp/run", "/tmp/data", ProfileMicro)
	if !strings.Contains(mongoConf, "cacheSizeGB: 0.05") {
		t.Errorf("mongod conf must limit wiredTiger cache to 0.05")
	}
}

func TestProvisionDatabase(t *testing.T) {
	_, mgr, tmpDir := setupDBTestEnv(t)
	defer os.RemoveAll(tmpDir)

	ctx := context.Background()

	// Provision PostgreSQL instance
	instance, password, err := mgr.ProvisionDatabase(ctx, "billing_db", "postgres", ProfileMicro)
	if err != nil {
		t.Fatalf("ProvisionDatabase failed: %v", err)
	}

	if instance.Engine != "postgres" || instance.Name != "billing_db" {
		t.Errorf("unexpected instance: %+v", instance)
	}
	if len(password) != 24 {
		t.Errorf("expected 24-char password, got %d", len(password))
	}
	if instance.Port < 5432 {
		t.Errorf("expected port >= 5432, got %d", instance.Port)
	}

	// Verify connection URIs
	tcpURI, socketURI := mgr.GetConnectionURIs(instance, password)
	if !strings.HasPrefix(tcpURI, "postgres://orbit_admin:") || !strings.Contains(tcpURI, "127.0.0.1") {
		t.Errorf("invalid tcp URI: %s", tcpURI)
	}
	if !strings.Contains(socketURI, "host=") {
		t.Errorf("invalid socket URI: %s", socketURI)
	}
}
