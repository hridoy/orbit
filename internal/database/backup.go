package database

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// keepBackups is how many backups per database are retained. Older ones are pruned.
const keepBackups = 7

// BackupDatabase writes a logical dump of a running postgres, mysql, or mariadb instance to
// DataDir/backups/<dbID>/ and returns the file path. The instance must be running.
func (m *DatabaseManager) BackupDatabase(parent context.Context, dbID string) (string, error) {
	instance, err := m.dbStore.GetDatabase(dbID)
	if err != nil || instance == nil {
		return "", fmt.Errorf("database %s not found", dbID)
	}

	m.mu.RLock()
	_, running := m.running[dbID]
	m.mu.RUnlock()
	if !running {
		return "", fmt.Errorf("database %s is not running; start it before backing up", dbID)
	}

	parts := strings.Split(instance.EncryptedPassword, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("stored credentials for %s are malformed", dbID)
	}
	pw, err := m.vault.Decrypt(parts[0], parts[1])
	if err != nil {
		return "", fmt.Errorf("decrypt credentials: %w", err)
	}

	dir := filepath.Join(m.cfg.DataDir, "backups", dbID)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create backup dir: %w", err)
	}
	stamp := time.Now().UTC().Format("20060102-150405")
	outPath := filepath.Join(dir, stamp+".sql")

	ctx, cancel := context.WithTimeout(parent, 30*time.Minute)
	defer cancel()

	var cmd *exec.Cmd
	switch instance.Engine {
	case "postgres":
		bin, err := findBinary("pg_dumpall", []string{
			"/usr/lib/postgresql/16/bin/pg_dumpall",
			"/usr/lib/postgresql/15/bin/pg_dumpall",
			"/usr/lib/postgresql/14/bin/pg_dumpall",
			"/usr/local/opt/postgresql@16/bin/pg_dumpall",
			"/opt/homebrew/opt/postgresql@16/bin/pg_dumpall",
		})
		if err != nil {
			return "", err
		}
		// pg_dumpall covers every database and role on the instance.
		cmd = exec.CommandContext(ctx, bin,
			"-h", "127.0.0.1", "-p", fmt.Sprint(instance.Port), "-U", instance.AdminUser)
		cmd.Env = append(os.Environ(), "PGPASSWORD="+string(pw))
	case "mysql", "mariadb":
		bin, err := findBinary("mysqldump", []string{
			"/usr/bin/mysqldump", "/usr/bin/mariadb-dump",
			"/usr/local/opt/mysql/bin/mysqldump", "/opt/homebrew/opt/mysql/bin/mysqldump",
		})
		if err != nil {
			return "", err
		}
		cmd = exec.CommandContext(ctx, bin,
			"-h", "127.0.0.1", "-P", fmt.Sprint(instance.Port), "-u", instance.AdminUser,
			"--all-databases", "--single-transaction", "--routines", "--events")
		cmd.Env = append(os.Environ(), "MYSQL_PWD="+string(pw))
	default:
		return "", fmt.Errorf("backup is not supported for engine %q yet", instance.Engine)
	}

	out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return "", err
	}
	var stderr strings.Builder
	cmd.Stdout = out
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	closeErr := out.Close()
	if runErr != nil {
		_ = os.Remove(outPath)
		return "", fmt.Errorf("dump failed: %w: %s", runErr, strings.TrimSpace(stderr.String()))
	}
	if closeErr != nil {
		_ = os.Remove(outPath)
		return "", closeErr
	}

	pruneBackups(dir, keepBackups)
	return outPath, nil
}

// pruneBackups keeps the newest `keep` dump files in dir.
func pruneBackups(dir string, keep int) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	if len(names) <= keep {
		return
	}
	// Timestamped names sort chronologically.
	sort.Strings(names)
	for _, n := range names[:len(names)-keep] {
		_ = os.Remove(filepath.Join(dir, n))
	}
}
