package database

import "fmt"

type RAMProfile string

const (
	ProfileMicro     RAMProfile = "micro"     // ~16-32MB RAM (Zero-bloat default)
	ProfileStandard  RAMProfile = "standard"  // ~128MB RAM
	ProfileHeavy     RAMProfile = "heavy"     // ~512MB RAM
	ProfileUnlimited RAMProfile = "unlimited" // 0 limits (Host capacity)
)

// ============================================================================
// PostgreSQL Tuning
// ============================================================================

type PostgresTuning struct {
	SharedBuffers      string
	WorkMem            string
	MaintenanceWorkMem string
	EffectiveCacheSize string
	MaxConnections     int
	WalBuffers         string
}

func GetPostgresTuning(profile RAMProfile) PostgresTuning {
	switch profile {
	case ProfileUnlimited:
		return PostgresTuning{
			SharedBuffers:      "1024MB",
			WorkMem:            "64MB",
			MaintenanceWorkMem: "256MB",
			EffectiveCacheSize: "2048MB",
			MaxConnections:     500,
			WalBuffers:         "64MB",
		}
	case ProfileHeavy:
		return PostgresTuning{
			SharedBuffers:      "256MB",
			WorkMem:            "16MB",
			MaintenanceWorkMem: "64MB",
			EffectiveCacheSize: "512MB",
			MaxConnections:     100,
			WalBuffers:         "16MB",
		}
	case ProfileStandard:
		return PostgresTuning{
			SharedBuffers:      "64MB",
			WorkMem:            "4MB",
			MaintenanceWorkMem: "32MB",
			EffectiveCacheSize: "128MB",
			MaxConnections:     50,
			WalBuffers:         "4MB",
		}
	default: // Micro (Default for zero-bloat sovereign PaaS)
		return PostgresTuning{
			SharedBuffers:      "16MB",
			WorkMem:            "1MB",
			MaintenanceWorkMem: "8MB",
			EffectiveCacheSize: "32MB",
			MaxConnections:     20,
			WalBuffers:         "1MB",
		}
	}
}

// GeneratePostgresConf produces a low-RAM, sovereign-tuned postgresql.conf.
func GeneratePostgresConf(port int, socketDir string, profile RAMProfile) string {
	t := GetPostgresTuning(profile)

	return fmt.Sprintf(`# ==============================================================================
# Orbit Sovereign Auto-Tuned postgresql.conf (Profile: %s)
# ==============================================================================
# Networking & IPC
listen_addresses = '127.0.0.1'
port = %d
unix_socket_directories = '%s'
unix_socket_permissions = 0700

# Memory Tuning
shared_buffers = %s
work_mem = %s
maintenance_work_mem = %s
effective_cache_size = %s
max_connections = %d
wal_buffers = %s

# Disk & Durability
fsync = on
synchronous_commit = off          # Ultra-fast write throughput for sovereign workloads
wal_writer_delay = 200ms
checkpoint_timeout = 15min
min_wal_size = 32MB
max_wal_size = 128MB

# Logging
log_destination = 'stderr'
logging_collector = off
log_min_messages = warning
`, profile, port, socketDir, t.SharedBuffers, t.WorkMem, t.MaintenanceWorkMem, t.EffectiveCacheSize, t.MaxConnections, t.WalBuffers)
}

func GeneratePostgresHbaConf(adminUser string) string {
	return fmt.Sprintf(`# Orbit Sovereign pg_hba.conf
# TYPE  DATABASE        USER            ADDRESS                 METHOD
local   all             all                                     trust
host    all             all             127.0.0.1/32            scram-sha-256
host    all             all             ::1/128                 scram-sha-256
`)
}

// ============================================================================
// MySQL & MariaDB Tuning
// ============================================================================

type MySQLTuning struct {
	InnodbBufferPoolSize string
	InnodbLogFileSize    string
	MaxConnections       int
	PerformanceSchema    string
}

func GetMySQLTuning(profile RAMProfile) MySQLTuning {
	switch profile {
	case ProfileUnlimited:
		return MySQLTuning{
			InnodbBufferPoolSize: "1024M",
			InnodbLogFileSize:    "128M",
			MaxConnections:       500,
			PerformanceSchema:    "on",
		}
	case ProfileHeavy:
		return MySQLTuning{
			InnodbBufferPoolSize: "256M",
			InnodbLogFileSize:    "64M",
			MaxConnections:       100,
			PerformanceSchema:    "off",
		}
	case ProfileStandard:
		return MySQLTuning{
			InnodbBufferPoolSize: "64M",
			InnodbLogFileSize:    "32M",
			MaxConnections:       50,
			PerformanceSchema:    "off",
		}
	default: // Micro (<32MB)
		return MySQLTuning{
			InnodbBufferPoolSize: "16M",
			InnodbLogFileSize:    "16M",
			MaxConnections:       20,
			PerformanceSchema:    "off",
		}
	}
}

func GenerateMySQLConf(port int, socketPath, dataDir string, profile RAMProfile) string {
	t := GetMySQLTuning(profile)

	return fmt.Sprintf(`# ==============================================================================
# Orbit Sovereign Auto-Tuned my.cnf (Profile: %s)
# ==============================================================================
[mysqld]
bind-address = 127.0.0.1
port = %d
socket = %s
datadir = %s

# Low-RAM Optimization (Disable performance schema to save ~150MB RAM)
performance_schema = %s
innodb_buffer_pool_size = %s
innodb_log_file_size = %s
max_connections = %d
key_buffer_size = 8M
table_open_cache = 64
thread_cache_size = 4
innodb_file_per_table = 1
innodb_flush_log_at_trx_commit = 2
skip-name-resolve
`, profile, port, socketPath, dataDir, t.PerformanceSchema, t.InnodbBufferPoolSize, t.InnodbLogFileSize, t.MaxConnections)
}

func GenerateMariaDBConf(port int, socketPath, dataDir string, profile RAMProfile) string {
	t := GetMySQLTuning(profile)
	ariaCache := "8M"
	if profile == ProfileStandard {
		ariaCache = "16M"
	} else if profile == ProfileHeavy {
		ariaCache = "64M"
	} else if profile == ProfileUnlimited {
		ariaCache = "256M"
	}

	return fmt.Sprintf(`# ==============================================================================
# Orbit Sovereign Auto-Tuned MariaDB my.cnf (Profile: %s)
# ==============================================================================
[mariadb]
bind-address = 127.0.0.1
port = %d
socket = %s
datadir = %s

# Low-RAM Optimization
performance_schema = %s
innodb_buffer_pool_size = %s
innodb_log_file_size = %s
max_connections = %d
aria_pagecache_buffer_size = %s
key_buffer_size = 8M
table_open_cache = 64
thread_cache_size = 4
skip-name-resolve
`, profile, port, socketPath, dataDir, t.PerformanceSchema, t.InnodbBufferPoolSize, t.InnodbLogFileSize, t.MaxConnections, ariaCache)
}

// ============================================================================
// MongoDB Tuning
// ============================================================================

func GenerateMongoConf(port int, socketDir, dataDir string, profile RAMProfile) string {
	var cacheSizeGB float64
	switch profile {
	case ProfileUnlimited:
		cacheSizeGB = 4.0
	case ProfileHeavy:
		cacheSizeGB = 1.0
	case ProfileStandard:
		cacheSizeGB = 0.25
	default: // Micro
		cacheSizeGB = 0.05 // ~50MB cache limit
	}
	// MongoDB refuses a WiredTiger cache below 0.25 GB, so small profiles use the minimum.
	if cacheSizeGB < 0.25 {
		cacheSizeGB = 0.25
	}

	return fmt.Sprintf(`# Orbit Sovereign Auto-Tuned mongod.conf (Profile: %s)
systemLog:
  destination: file
  path: %s/mongod.log
  logAppend: true
storage:
  dbPath: %s
  wiredTiger:
    engineConfig:
      cacheSizeGB: %.2f
net:
  port: %d
  bindIp: 127.0.0.1
  unixDomainSocket:
    enabled: true
    pathPrefix: %s
`, profile, dataDir, dataDir, cacheSizeGB, port, socketDir)
}

// ============================================================================
// Redis Tuning
// ============================================================================

func GenerateRedisConf(port int, socketPath, dataDir, password string, profile RAMProfile) string {
	var maxMemoryDirective string
	switch profile {
	case ProfileUnlimited:
		maxMemoryDirective = "# maxmemory unlimited (host capacity)"
	case ProfileHeavy:
		maxMemoryDirective = "maxmemory 512mb\nmaxmemory-policy allkeys-lru"
	case ProfileStandard:
		maxMemoryDirective = "maxmemory 128mb\nmaxmemory-policy allkeys-lru"
	default: // Micro (<32MB)
		maxMemoryDirective = "maxmemory 32mb\nmaxmemory-policy allkeys-lru"
	}

	return fmt.Sprintf(`# ==============================================================================
# Orbit Sovereign Auto-Tuned redis.conf (Profile: %s)
# ==============================================================================
bind 127.0.0.1
port %d
unixsocket %s
unixsocketperm 700

# Working directory & Persistence
dir %s
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec

# Authentication
requirepass %s

# Memory Limits
%s

# Low-RAM Tweaks
tcp-backlog 128
timeout 0
tcp-keepalive 300
activerehashing yes
`, profile, port, socketPath, dataDir, password, maxMemoryDirective)
}

// ============================================================================
// LibSQL Tuning
// ============================================================================

func GenerateLibSQLConf(port int, dataDir string, profile RAMProfile) string {
	return fmt.Sprintf(`# ==============================================================================
# Orbit Sovereign Auto-Tuned libsql / sqld.conf (Profile: %s)
# ==============================================================================
[server]
http_listen_addr = "127.0.0.1:%d"
db_path = "%s/libsql.db"

[wal]
auto_checkpoint = 1000
`, profile, port, dataDir)
}
