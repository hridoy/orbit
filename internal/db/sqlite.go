package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

// DB wraps the standard *sql.DB with convenience helpers.
type DB struct {
	*sql.DB
}

// Open creates or connects to an embedded SQLite database using WAL mode and low-RAM pragmas.
func Open(dbPath string) (*DB, error) {
	// DSN with query parameters for connection-level settings
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", dbPath)
	rawDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// SQLite WAL mode handles multiple concurrent readers, but serialize write access
	rawDB.SetMaxOpenConns(16)
	rawDB.SetMaxIdleConns(4)
	rawDB.SetConnMaxLifetime(1 * time.Hour)

	// Apply critical performance & low-memory pragmas
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA temp_store = MEMORY;",
		"PRAGMA cache_size = -8000;", // ~8MB memory cache
	}

	for _, p := range pragmas {
		if _, err := rawDB.Exec(p); err != nil {
			rawDB.Close()
			return nil, fmt.Errorf("failed to execute pragma '%s': %w", p, err)
		}
	}

	database := &DB{DB: rawDB}
	if err := database.Migrate(); err != nil {
		rawDB.Close()
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	return database, nil
}

// Migrate applies the embedded schema.sql idempotently.
func (d *DB) Migrate() error {
	_, err := d.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("schema execution error: %w", err)
	}
	// Idempotent column migrations for existing databases
	alterStmts := []string{
		`ALTER TABLE storage_buckets ADD COLUMN versioning_enabled INTEGER NOT NULL DEFAULT 0;`,
		`ALTER TABLE storage_buckets ADD COLUMN object_lock_enabled INTEGER NOT NULL DEFAULT 0;`,
		`ALTER TABLE storage_buckets ADD COLUMN retention_days INTEGER NOT NULL DEFAULT 0;`,
		`ALTER TABLE storage_buckets ADD COLUMN lifecycle_days INTEGER NOT NULL DEFAULT 0;`,
		`ALTER TABLE storage_buckets ADD COLUMN webhook_url TEXT NOT NULL DEFAULT '';`,
		`ALTER TABLE storage_objects ADD COLUMN retain_until DATETIME;`,
	}
	for _, stmt := range alterStmts {
		_, _ = d.Exec(stmt)
	}
	return nil
}
