package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"
	"time"

	"orbit/internal/compute"
	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/security"
)

type RunningDB struct {
	DBID       string
	Engine     string
	Port       int
	Cmd        *exec.Cmd
	Cancel     context.CancelFunc
	RingBuffer *compute.RingBuffer
	Exited     chan struct{}
}

type DatabaseManager struct {
	mu        sync.RWMutex
	cfg       *config.Config
	dbStore   *db.DatabaseStore
	secStore  *db.SecretStore
	vault     *security.Vault
	postgres  *PostgresEngine
	mysql     *MySQLEngine
	mongodb   *MongoEngine
	redis     *RedisEngine
	libsql    *LibSQLEngine
	installed []Installed
	running   map[string]*RunningDB
	buffers   map[string]*compute.RingBuffer
}

func NewDatabaseManager(
	cfg *config.Config,
	dbStore *db.DatabaseStore,
	secStore *db.SecretStore,
	vault *security.Vault,
) *DatabaseManager {
	return &DatabaseManager{
		cfg:       cfg,
		dbStore:   dbStore,
		secStore:  secStore,
		vault:     vault,
		postgres:  NewPostgresEngine(),
		mysql:     NewMySQLEngine(),
		mongodb:   NewMongoEngine(),
		redis:     NewRedisEngine(),
		libsql:    NewLibSQLEngine(),
		installed: discoverInstalled(),
		running:   make(map[string]*RunningDB),
		buffers:   make(map[string]*compute.RingBuffer),
	}
}

func (m *DatabaseManager) DBBaseDir(dbID string) string {
	return filepath.Join(m.cfg.DatabasesDir, dbID)
}

func (m *DatabaseManager) GetRingBuffer(dbID string) *compute.RingBuffer {
	m.mu.Lock()
	defer m.mu.Unlock()

	rb, exists := m.buffers[dbID]
	if !exists {
		rb = compute.NewRingBuffer(1000)
		m.buffers[dbID] = rb
	}
	return rb
}

// FindFreePort searches for an available loopback port not already registered in SQLite.
func (m *DatabaseManager) FindFreePort(startPort int) (int, error) {
	allocatedPorts := make(map[int]bool)
	if existing, err := m.dbStore.ListDatabases(); err == nil {
		for _, d := range existing {
			allocatedPorts[d.Port] = true
		}
	}

	for port := startPort; port < startPort+1000; port++ {
		if allocatedPorts[port] {
			continue
		}
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			_ = ln.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no free ports available above %d", startPort)
}

// ProvisionDatabase creates a new dedicated, isolated database instance.
// Supported engines (matching Dokploy): postgres, mysql, mariadb, mongodb, redis, libsql
func (m *DatabaseManager) ProvisionDatabase(
	ctx context.Context,
	name string,
	engine string,
	version string,
	username string,
	password string,
	ramProfile RAMProfile,
) (*db.DatabaseInstance, string, error) {
	var startPort int
	switch engine {
	case "postgres":
		startPort = 5432
	case "mysql", "mariadb":
		startPort = 3306
	case "mongodb":
		startPort = 27017
	case "redis":
		startPort = 6379
	case "libsql":
		startPort = 8080
	default:
		return nil, "", fmt.Errorf("unsupported database engine: %s (must be postgres, mysql, mariadb, mongodb, redis, or libsql)", engine)
	}

	dbID := fmt.Sprintf("db_%s_%s", engine, generateRandomSlug(6))
	baseDir := m.DBBaseDir(dbID)

	port, err := m.FindFreePort(startPort)
	if err != nil {
		return nil, "", err
	}

	if version == "" {
		version = nativeVersion
	}
	if _, err := m.pick(engine, version); err != nil {
		return nil, "", err
	}

	// Username and password come from the operator or are generated. Redis always uses "default".
	adminUser := "orbit_admin"
	if username != "" {
		if !dbUserRe.MatchString(username) {
			return nil, "", fmt.Errorf("username must be 2 to 32 characters: lowercase letters, digits, underscore, starting with a letter")
		}
		adminUser = username
	}
	if engine == "redis" {
		adminUser = "default"
	}
	rawPassword := password
	if rawPassword == "" {
		rawPassword = generateRandomPassword(24)
	} else if len(rawPassword) < 8 || len(rawPassword) > 128 {
		return nil, "", fmt.Errorf("password must be 8 to 128 characters")
	}

	// Encrypt password via AES-256-GCM
	cipherB64, nonceB64, err := m.vault.Encrypt([]byte(rawPassword))
	if err != nil {
		return nil, "", fmt.Errorf("failed to encrypt password: %w", err)
	}
	encryptedPassword := fmt.Sprintf("%s:%s", cipherB64, nonceB64)

	socketPath := filepath.Join(baseDir, "run")
	dataDir := filepath.Join(baseDir, "data")

	// Run engine-specific provisioning
	switch engine {
	case "postgres":
		pe, err := m.postgresFor(version)
		if err != nil {
			return nil, "", err
		}
		if err := pe.Provision(ctx, baseDir, adminUser, rawPassword, port, ramProfile); err != nil {
			return nil, "", fmt.Errorf("postgres provision error: %w", err)
		}
	case "mysql":
		socketPath = filepath.Join(baseDir, "run", "mysql.sock")
		me, err := m.mysqlFor("mysql", version)
		if err != nil {
			return nil, "", err
		}
		if err := me.Provision(ctx, baseDir, adminUser, rawPassword, port, ramProfile, false); err != nil {
			return nil, "", fmt.Errorf("mysql provision error: %w", err)
		}
	case "mariadb":
		socketPath = filepath.Join(baseDir, "run", "mariadb.sock")
		me, err := m.mysqlFor("mariadb", version)
		if err != nil {
			return nil, "", err
		}
		if err := me.Provision(ctx, baseDir, adminUser, rawPassword, port, ramProfile, true); err != nil {
			return nil, "", fmt.Errorf("mariadb provision error: %w", err)
		}
	case "mongodb":
		mo, err := m.mongoFor(version)
		if err != nil {
			return nil, "", err
		}
		if err := mo.Provision(ctx, baseDir, adminUser, rawPassword, port, ramProfile); err != nil {
			return nil, "", fmt.Errorf("mongodb provision error: %w", err)
		}
	case "redis":
		socketPath = filepath.Join(baseDir, "run", "redis.sock")
		re, err := m.redisFor(version)
		if err != nil {
			return nil, "", err
		}
		if err := re.Provision(ctx, baseDir, rawPassword, port, ramProfile); err != nil {
			return nil, "", fmt.Errorf("redis provision error: %w", err)
		}
	case "libsql":
		socketPath = filepath.Join(baseDir, "data", "libsql.db")
		le, err := m.libsqlFor(version)
		if err != nil {
			return nil, "", err
		}
		if err := le.Provision(ctx, baseDir, rawPassword, port, ramProfile); err != nil {
			return nil, "", fmt.Errorf("libsql provision error: %w", err)
		}
	}

	instance := &db.DatabaseInstance{
		DBID:              dbID,
		Name:              name,
		Engine:            engine,
		Version:           version,
		NodeID:            m.cfg.NodeID,
		Port:              port,
		SocketPath:        socketPath,
		RAMProfile:        string(ramProfile),
		StorageLimitMB:    10240, // 10GB default
		DataDirectory:     dataDir,
		AdminUser:         adminUser,
		EncryptedPassword: encryptedPassword,
		Status:            "stopped",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	if err := m.dbStore.CreateDatabase(instance); err != nil {
		return nil, "", fmt.Errorf("failed to save database record: %w", err)
	}

	return instance, rawPassword, nil
}

// StartDatabase launches the dedicated database engine process.
func (m *DatabaseManager) StartDatabase(ctx context.Context, dbID string) error {
	m.mu.Lock()
	if _, running := m.running[dbID]; running {
		m.mu.Unlock()
		return nil // Already running
	}
	m.mu.Unlock()

	instance, err := m.dbStore.GetDatabase(dbID)
	if err != nil || instance == nil {
		return fmt.Errorf("database %s not found: %w", dbID, err)
	}

	baseDir := m.DBBaseDir(dbID)
	if err := chownDBTree(baseDir); err != nil {
		return fmt.Errorf("prepare data dir: %w", err)
	}
	rb := m.GetRingBuffer(dbID)

	childCtx, cancel := context.WithCancel(context.Background())
	var cmd *exec.Cmd

	switch instance.Engine {
	case "postgres":
		pe, perr := m.postgresFor(instance.Version)
		if perr != nil {
			cancel()
			return perr
		}
		cmd, err = pe.Start(childCtx, baseDir, rb)
	case "mysql":
		me, merr := m.mysqlFor("mysql", instance.Version)
		if merr != nil {
			cancel()
			return merr
		}
		cmd, err = me.Start(childCtx, baseDir, rb, false)
	case "mariadb":
		me, merr := m.mysqlFor("mariadb", instance.Version)
		if merr != nil {
			cancel()
			return merr
		}
		cmd, err = me.Start(childCtx, baseDir, rb, true)
	case "mongodb":
		mo, merr := m.mongoFor(instance.Version)
		if merr != nil {
			cancel()
			return merr
		}
		cmd, err = mo.Start(childCtx, baseDir, rb)
	case "redis":
		re, rerr := m.redisFor(instance.Version)
		if rerr != nil {
			cancel()
			return rerr
		}
		cmd, err = re.Start(childCtx, baseDir, rb)
	case "libsql":
		le, lerr := m.libsqlFor(instance.Version)
		if lerr != nil {
			cancel()
			return lerr
		}
		cmd, err = le.Start(childCtx, baseDir, rb, instance.Port)
	default:
		cancel()
		return fmt.Errorf("unknown engine %s", instance.Engine)
	}

	if err != nil {
		cancel()
		return fmt.Errorf("failed to start %s: %w", instance.Engine, err)
	}

	running := &RunningDB{
		DBID:       dbID,
		Engine:     instance.Engine,
		Port:       instance.Port,
		Cmd:        cmd,
		Cancel:     cancel,
		RingBuffer: rb,
		Exited:     make(chan struct{}),
	}

	m.mu.Lock()
	m.running[dbID] = running
	m.mu.Unlock()

	_ = m.dbStore.UpdateDatabaseStatus(dbID, "running")

	go func() {
		_ = cmd.Wait()
		close(running.Exited)
		m.mu.Lock()
		delete(m.running, dbID)
		m.mu.Unlock()
		_ = m.dbStore.UpdateDatabaseStatus(dbID, "stopped")
		rb.Append(fmt.Sprintf("[%s] Process stopped cleanly.", instance.Engine))
	}()

	return nil
}

// StopDatabase gracefully terminates a running database process.
func (m *DatabaseManager) StopDatabase(dbID string) error {
	m.mu.Lock()
	running, ok := m.running[dbID]
	m.mu.Unlock()

	if !ok || running == nil || running.Cmd == nil || running.Cmd.Process == nil {
		_ = m.dbStore.UpdateDatabaseStatus(dbID, "stopped")
		return nil
	}

	// Send SIGINT / SIGTERM for clean DB shutdown (checkpoint/flush)
	_ = running.Cmd.Process.Signal(syscall.SIGINT)

	select {
	case <-running.Exited:
		return nil
	case <-time.After(10 * time.Second):
		_ = running.Cmd.Process.Kill()
		running.Cancel()
		<-running.Exited
		return nil
	}
}

// RestartDatabase stops and restarts the database.
func (m *DatabaseManager) RestartDatabase(ctx context.Context, dbID string) error {
	_ = m.StopDatabase(dbID)
	time.Sleep(500 * time.Millisecond)
	return m.StartDatabase(ctx, dbID)
}

// DeleteDatabase stops the database and removes all data.
// DeleteDatabase stops the engine and removes its data. It refuses to delete data while the
// process is still alive, and it reports removal failures instead of ignoring them.
func (m *DatabaseManager) DeleteDatabase(ctx context.Context, dbID string) error {
	if err := m.StopDatabase(dbID); err != nil {
		return fmt.Errorf("stop database before delete: %w", err)
	}

	m.mu.RLock()
	_, stillRunning := m.running[dbID]
	m.mu.RUnlock()
	if stillRunning {
		return fmt.Errorf("database %s is still running; refusing to delete its data", dbID)
	}

	if err := os.RemoveAll(m.DBBaseDir(dbID)); err != nil {
		return fmt.Errorf("remove data directory for %s: %w", dbID, err)
	}
	return m.dbStore.DeleteDatabase(dbID)
}

// GetConnectionURIs returns formatted TCP and Unix domain socket URIs for zero-network IPC.
func (m *DatabaseManager) GetConnectionURIs(instance *db.DatabaseInstance, password string) (tcpURI, socketURI string) {
	switch instance.Engine {
	case "postgres":
		tcpURI = fmt.Sprintf("postgres://%s:%s@127.0.0.1:%d/%s?sslmode=disable",
			instance.AdminUser, password, instance.Port, instance.Name)
		socketURI = fmt.Sprintf("postgres://%s:%s@/%s?host=%s",
			instance.AdminUser, password, instance.Name, instance.SocketPath)
	case "mysql", "mariadb":
		tcpURI = fmt.Sprintf("mysql://%s:%s@127.0.0.1:%d/%s",
			instance.AdminUser, password, instance.Port, instance.Name)
		socketURI = fmt.Sprintf("mysql://%s:%s@localhost/%s?socket=%s",
			instance.AdminUser, password, instance.Name, instance.SocketPath)
	case "mongodb":
		tcpURI = fmt.Sprintf("mongodb://%s:%s@127.0.0.1:%d/%s?authSource=admin",
			instance.AdminUser, password, instance.Port, instance.Name)
		socketURI = fmt.Sprintf("mongodb://%s:%s@%%2F%s%%2Fmongodb-%d.sock/%s",
			instance.AdminUser, password, instance.SocketPath, instance.Port, instance.Name)
	case "redis":
		tcpURI = fmt.Sprintf("redis://:%s@127.0.0.1:%d/0", password, instance.Port)
		socketURI = fmt.Sprintf("redis://:%s@%s?db=0", password, instance.SocketPath)
	case "libsql":
		tcpURI = fmt.Sprintf("libsql://127.0.0.1:%d?authToken=%s", instance.Port, password)
		socketURI = fmt.Sprintf("file://%s", instance.SocketPath)
	}
	return tcpURI, socketURI
}

func generateRandomSlug(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)[:n]
}

func generateRandomPassword(length int) string {
	b := make([]byte, length/2+1)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)[:length]
}

var dbUserRe = regexp.MustCompile(`^[a-z][a-z0-9_]{1,31}$`)
