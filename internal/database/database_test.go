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

	mysqlConf := GenerateMySQLConf(3306, "/tmp/run/mysql.sock", "/tmp/data", ProfileMicro)
	if !strings.Contains(mysqlConf, "performance_schema = off") || !strings.Contains(mysqlConf, "innodb_buffer_pool_size = 16M") {
		t.Errorf("mysql conf must disable performance schema and set 16M innodb buffer: %s", mysqlConf)
	}

	redisConf := GenerateRedisConf(6379, "/tmp/run/redis.sock", "/tmp/data", "secret", ProfileMicro)
	if !strings.Contains(redisConf, "maxmemory 32mb") || !strings.Contains(redisConf, "requirepass secret") {
		t.Errorf("redis conf must set maxmemory 32mb and requirepass: %s", redisConf)
	}

	libsqlConf := GenerateLibSQLConf(8080, "/tmp/data", ProfileMicro)
	if !strings.Contains(libsqlConf, `http_listen_addr = "127.0.0.1:8080"`) {
		t.Errorf("libsql conf must listen on 8080: %s", libsqlConf)
	}

	mongoConf := GenerateMongoConf(27017, "/tmp/run", "/tmp/data", ProfileMicro)
	if !strings.Contains(mongoConf, "cacheSizeGB: 0.25") {
		t.Errorf("mongod conf must use the 0.25 GB minimum wiredTiger cache for micro")
	}
}

func TestProvisionAllEngines(t *testing.T) {
	_, mgr, tmpDir := setupDBTestEnv(t)
	defer os.RemoveAll(tmpDir)

	ctx := context.Background()

	engines := []struct {
		engine    string
		startPort int
		uriPrefix string
	}{
		{"postgres", 5432, "postgres://"},
		{"mysql", 3306, "mysql://"},
		{"mariadb", 3306, "mysql://"},
		{"mongodb", 27017, "mongodb://"},
		{"redis", 6379, "redis://"},
		{"libsql", 8080, "libsql://"},
	}

	for _, tc := range engines {
		t.Run(tc.engine, func(t *testing.T) {
			instance, password, err := mgr.ProvisionDatabase(ctx, tc.engine+"_db", tc.engine, ProfileMicro)
			if err != nil {
				t.Fatalf("ProvisionDatabase failed for %s: %v", tc.engine, err)
			}

			if instance.Engine != tc.engine {
				t.Errorf("expected engine %s, got %s", tc.engine, instance.Engine)
			}
			if len(password) != 24 {
				t.Errorf("expected 24-char password, got %d", len(password))
			}
			if instance.Port < tc.startPort {
				t.Errorf("expected port >= %d, got %d", tc.startPort, instance.Port)
			}

			// Verify connection URIs
			tcpURI, socketURI := mgr.GetConnectionURIs(instance, password)
			if !strings.HasPrefix(tcpURI, tc.uriPrefix) || !strings.Contains(tcpURI, "127.0.0.1") {
				t.Errorf("invalid tcp URI for %s: %s", tc.engine, tcpURI)
			}
			if socketURI == "" {
				t.Errorf("empty socket URI for %s", tc.engine)
			}
		})
	}
}
