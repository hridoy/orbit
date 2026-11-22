package database

import "fmt"

type RAMProfile string

const (
	ProfileMicro    RAMProfile = "micro"    // ~32MB RAM
	ProfileStandard RAMProfile = "standard" // ~128MB RAM
	ProfileHeavy    RAMProfile = "heavy"    // ~512MB RAM
)

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
synchronous_commit = off          # Ultra-fast write throughput for micro-containers
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

// GeneratePostgresHbaConf restricts access strictly to local socket and 127.0.0.1 with SCRAM-SHA-256.
func GeneratePostgresHbaConf(adminUser string) string {
	return fmt.Sprintf(`# Orbit Sovereign pg_hba.conf
# TYPE  DATABASE        USER            ADDRESS                 METHOD
local   all             all                                     trust
host    all             all             127.0.0.1/32            scram-sha-256
host    all             all             ::1/128                 scram-sha-256
`)
}

// GenerateMongoConf produces a tuned mongod.conf with WiredTiger cache memory capping.
func GenerateMongoConf(port int, socketDir, dataDir string, profile RAMProfile) string {
	var cacheSizeGB float64
	switch profile {
	case ProfileHeavy:
		cacheSizeGB = 1.0
	case ProfileStandard:
		cacheSizeGB = 0.25
	default: // Micro
		cacheSizeGB = 0.05 // ~50MB cache limit
	}

	return fmt.Sprintf(`# Orbit Sovereign Auto-Tuned mongod.conf (Profile: %s)
systemLog:
  destination: file
  path: %s/mongod.log
  logAppend: true
storage:
  dbPath: %s
  journal:
    enabled: true
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
