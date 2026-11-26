package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/s3"
	"orbit/internal/security"
)

const (
	canaryName      = "canary"
	canaryPlaintext = "orbit-vault-canary-v1"
	adminTokenFile  = "admin.token"
	adminTokenTTL   = 365 * 24 * time.Hour
)

// lockFile is held exclusively by a running agent and by rekey, so the two never overlap.
type lockFile struct{ f *os.File }

func acquireLock(path string) (*lockFile, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return nil, fmt.Errorf("another orbit process holds %s (is the agent running?)", path)
	}
	return &lockFile{f: f}, nil
}

func (l *lockFile) release() {
	if l != nil && l.f != nil {
		_ = syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
		l.f.Close()
	}
}

// writeFileDurable writes data to path via a temp file, fsyncs it, and renames it into place.
func writeFileDurable(path string, data []byte, perm os.FileMode) error {
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// tokenManagerFor builds the token signer from a subkey, never from the master key itself.
func tokenManagerFor(masterKey []byte) (*security.TokenManager, error) {
	return security.NewTokenManager(security.DeriveKey(masterKey, security.TokenSeedLabel))
}

// ensureAdminToken keeps secrets/admin.token holding a valid admin token for the CLI.
func ensureAdminToken(cfg *config.Config, tm *security.TokenManager) error {
	path := filepath.Join(cfg.SecretsDir, adminTokenFile)
	if b, err := os.ReadFile(path); err == nil {
		if _, err := tm.VerifyToken(strings.TrimSpace(string(b))); err == nil {
			return nil
		}
	}
	tok, err := tm.IssueToken("cli-admin", "admin", adminTokenTTL)
	if err != nil {
		return err
	}
	return writeFileDurable(path, []byte(tok), 0600)
}

func readMasterKey(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(b) != 32 {
		return nil, fmt.Errorf("%s must be 32 bytes, got %d", path, len(b))
	}
	return b, nil
}

// canaryOK reports whether key decrypts the vault canary row.
func canaryOK(database *db.DB, key []byte) bool {
	vault, err := security.NewVault(key)
	if err != nil {
		return false
	}
	var c, n string
	err = database.QueryRow(
		`SELECT ciphertext_base64, nonce_base64 FROM vault_meta WHERE name = ?`, canaryName,
	).Scan(&c, &n)
	if err != nil {
		return false
	}
	plain, err := vault.Decrypt(c, n)
	return err == nil && string(plain) == canaryPlaintext
}

// ensureCanary writes the canary with key if it does not exist yet.
func ensureCanary(database *db.DB, key []byte) error {
	var exists int
	if err := database.QueryRow(`SELECT COUNT(*) FROM vault_meta WHERE name = ?`, canaryName).Scan(&exists); err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}
	vault, err := security.NewVault(key)
	if err != nil {
		return err
	}
	c, n, err := vault.Encrypt([]byte(canaryPlaintext))
	if err != nil {
		return err
	}
	_, err = database.Exec(`INSERT INTO vault_meta (name, ciphertext_base64, nonce_base64) VALUES (?, ?, ?)`, canaryName, c, n)
	return err
}

// resolveMasterKey returns the active master key. It finishes or rolls back a rotation that was
// interrupted, using the canary to decide which key the database is encrypted with.
func resolveMasterKey(cfg *config.Config, database *db.DB) ([]byte, error) {
	masterPath := filepath.Join(cfg.SecretsDir, "master.key")
	newKeyPath := masterPath + ".new"
	s3Path := filepath.Join(cfg.DataDir, s3.S3CredentialsFile)
	s3New := s3Path + ".new"

	if _, err := os.Stat(newKeyPath); err == nil {
		nk, readErr := readMasterKey(newKeyPath)
		if readErr == nil && canaryOK(database, nk) {
			// The database already uses the new key: finish the rotation.
			if _, err := os.Stat(s3New); err == nil {
				if err := os.Rename(s3New, s3Path); err != nil {
					return nil, fmt.Errorf("finish rotation: %w", err)
				}
			}
			if err := os.Rename(newKeyPath, masterPath); err != nil {
				return nil, fmt.Errorf("finish rotation: %w", err)
			}
			log.Printf("🔁 finished an interrupted master key rotation")
			return nk, nil
		}
		// The database still uses the old key: discard the pending rotation.
		_ = os.Remove(newKeyPath)
		_ = os.Remove(s3New)
		log.Printf("🔁 discarded an unfinished master key rotation")
	}

	cur, err := readMasterKey(masterPath)
	if err != nil {
		return nil, err
	}
	var hasCanary int
	if err := database.QueryRow(`SELECT COUNT(*) FROM vault_meta WHERE name = ?`, canaryName).Scan(&hasCanary); err != nil {
		return nil, err
	}
	if hasCanary > 0 && !canaryOK(database, cur) {
		return nil, fmt.Errorf("master key in %s does not match the database; refusing to start", masterPath)
	}
	return cur, nil
}

// rotateMasterKey re-encrypts every stored ciphertext under newKey. The order is:
// write master.key.new and the new S3 file, commit the database, then rename both into place.
// resolveMasterKey completes or rolls back the rename steps if the process dies midway.
// With dryRun it decrypts everything and reports counts, changing nothing.
func rotateMasterKey(cfg *config.Config, database *db.DB, oldKey, newKey []byte, dryRun bool) (map[string]int, error) {
	oldV, err := security.NewVault(oldKey)
	if err != nil {
		return nil, err
	}
	newV, err := security.NewVault(newKey)
	if err != nil {
		return nil, err
	}
	counts := map[string]int{}

	type cipherUpdate struct{ id, c, n string }
	var secretUpdates []cipherUpdate
	rows, err := database.Query(`SELECT secret_id, ciphertext_base64, nonce_base64 FROM secrets`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id, c, n string
		if err := rows.Scan(&id, &c, &n); err != nil {
			rows.Close()
			return nil, err
		}
		plain, err := oldV.Decrypt(c, n)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("secret %s does not decrypt with the current key: %w", id, err)
		}
		nc, nn, err := newV.Encrypt(plain)
		if err != nil {
			rows.Close()
			return nil, err
		}
		secretUpdates = append(secretUpdates, cipherUpdate{id, nc, nn})
	}
	rows.Close()
	counts["secrets"] = len(secretUpdates)

	var dbUpdates []cipherUpdate
	dbRows, err := database.Query(`SELECT db_id, encrypted_password FROM databases`)
	if err != nil {
		return nil, err
	}
	for dbRows.Next() {
		var id, enc string
		if err := dbRows.Scan(&id, &enc); err != nil {
			dbRows.Close()
			return nil, err
		}
		parts := strings.Split(enc, ":")
		if len(parts) != 2 {
			dbRows.Close()
			return nil, fmt.Errorf("database %s has malformed stored password", id)
		}
		plain, err := oldV.Decrypt(parts[0], parts[1])
		if err != nil {
			dbRows.Close()
			return nil, fmt.Errorf("database %s password does not decrypt with the current key: %w", id, err)
		}
		nc, nn, err := newV.Encrypt(plain)
		if err != nil {
			dbRows.Close()
			return nil, err
		}
		dbUpdates = append(dbUpdates, cipherUpdate{id, nc, nn})
	}
	dbRows.Close()
	counts["databases"] = len(dbUpdates)

	canaryC, canaryN, err := newV.Encrypt([]byte(canaryPlaintext))
	if err != nil {
		return nil, err
	}

	// S3 credentials: re-seal to a .new file now, rename after the database commits.
	s3Path := filepath.Join(cfg.DataDir, s3.S3CredentialsFile)
	s3New := s3Path + ".new"
	if _, err := os.Stat(s3Path); err == nil {
		creds, err := s3.OpenS3Credentials(s3Path, oldKey)
		if err != nil {
			return nil, fmt.Errorf("s3 credentials do not decrypt with the current key: %w", err)
		}
		if err := s3.SealS3Credentials(s3New, newKey, creds); err != nil {
			return nil, err
		}
		counts["s3_credentials"] = 1
	}

	if dryRun {
		_ = os.Remove(s3New)
		return counts, nil
	}

	masterPath := filepath.Join(cfg.SecretsDir, "master.key")
	if err := writeFileDurable(masterPath+".new", newKey, 0600); err != nil {
		_ = os.Remove(s3New)
		return nil, err
	}

	tx, err := database.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	for _, u := range secretUpdates {
		if _, err := tx.Exec(`UPDATE secrets SET ciphertext_base64 = ?, nonce_base64 = ? WHERE secret_id = ?`, u.c, u.n, u.id); err != nil {
			return nil, err
		}
	}
	for _, u := range dbUpdates {
		if _, err := tx.Exec(`UPDATE databases SET encrypted_password = ? WHERE db_id = ?`, u.c+":"+u.n, u.id); err != nil {
			return nil, err
		}
	}
	if _, err := tx.Exec(`DELETE FROM vault_meta WHERE name = ?`, canaryName); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(`INSERT INTO vault_meta (name, ciphertext_base64, nonce_base64) VALUES (?, ?, ?)`, canaryName, canaryC, canaryN); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("database commit failed, old key still active: %w", err)
	}

	if _, err := os.Stat(s3New); err == nil {
		if err := os.Rename(s3New, s3Path); err != nil {
			return nil, fmt.Errorf("rename s3 credentials (restart agent to finish): %w", err)
		}
	}
	if err := os.Rename(masterPath+".new", masterPath); err != nil {
		return nil, fmt.Errorf("rename master key (restart agent to finish): %w", err)
	}
	return counts, nil
}

// rekeyFromBackupDir is the `orbit-agent rekey` flow: lock, resolve, back up, rotate, reissue.
func runRekey(cfg *config.Config, database *db.DB, backupDir string, dryRun bool) error {
	lock, err := acquireLock(filepath.Join(cfg.SecretsDir, "agent.lock"))
	if err != nil {
		return fmt.Errorf("%w; stop the agent before rekeying", err)
	}
	defer lock.release()

	oldKey, err := resolveMasterKey(cfg, database)
	if err != nil {
		return err
	}
	if err := ensureCanary(database, oldKey); err != nil {
		return err
	}

	if !dryRun {
		if err := runBackup(cfg, database, backupDir); err != nil {
			return fmt.Errorf("backup before rekey failed, nothing changed: %w", err)
		}
	}

	newKey := make([]byte, 32)
	if _, err := rand.Read(newKey); err != nil {
		return err
	}
	counts, err := rotateMasterKey(cfg, database, oldKey, newKey, dryRun)
	if err != nil {
		return fmt.Errorf("rotation failed, nothing changed: %w", err)
	}
	fmt.Printf("re-encrypted: secrets=%d databases=%d s3_credentials=%d\n",
		counts["secrets"], counts["databases"], counts["s3_credentials"])
	if dryRun {
		fmt.Println("dry run: nothing changed")
		return nil
	}

	tm, err := tokenManagerFor(newKey)
	if err != nil {
		return err
	}
	if err := ensureAdminToken(cfg, tm); err != nil {
		return err
	}
	fmt.Println("✅ master key rotated. Restart the agent; the CLI now uses the new secrets/admin.token.")
	fmt.Println("   Delete old backups and purge old logs: they contain the previous master key.")
	return nil
}
