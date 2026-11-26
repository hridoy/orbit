package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orbit/internal/config"
	"orbit/internal/db"
)

// runBackup writes a consistent copy of the control-plane state to a new timestamped
// directory under outDir: the SQLite database (via VACUUM INTO, safe while running),
// the master key, the node ID, and the encrypted S3 credentials. Losing master.key makes
// every vault secret and the S3 credentials unrecoverable, so keep this copy off the host.
func runBackup(cfg *config.Config, database *db.DB, outDir string) error {
	stamp := time.Now().UTC().Format("20060102-150405")
	dest := filepath.Join(outDir, "orbit-backup-"+stamp)
	if err := os.MkdirAll(dest, 0700); err != nil {
		return fmt.Errorf("create backup dir: %w", err)
	}

	// Consistent SQLite snapshot. Quote the path for SQL.
	dbCopy := filepath.Join(dest, "orbit.db")
	quoted := strings.ReplaceAll(dbCopy, "'", "''")
	if _, err := database.Exec("VACUUM INTO '" + quoted + "'"); err != nil {
		return fmt.Errorf("snapshot database: %w", err)
	}

	files := []struct {
		src, name string
		optional  bool
	}{
		{filepath.Join(cfg.SecretsDir, "master.key"), "master.key", false},
		{filepath.Join(cfg.DataDir, "node_id"), "node_id", false},
		{filepath.Join(cfg.DataDir, "s3_credentials.enc"), "s3_credentials.enc", true},
	}
	for _, f := range files {
		if err := copyFile(f.src, filepath.Join(dest, f.name)); err != nil {
			if f.optional && os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("copy %s: %w", f.name, err)
		}
	}

	fmt.Printf("✅ Backup written to %s\n", dest)
	fmt.Println("   Move this directory off the server. It contains the master key.")
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
