package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orbit/internal/backup"
	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/integrations"
	"orbit/internal/security"
)

// uploadLatest copies the newest backup in outDir to the cloud destination saved in Credentials.
// It does nothing when no destination is set. The local copy is kept either way.
func uploadLatest(cfg *config.Config, database *db.DB, outDir string) error {
	key, err := resolveMasterKey(cfg, database)
	if err != nil {
		return err
	}
	vault, err := security.NewVault(key)
	if err != nil {
		return err
	}
	conf, remote, ok := integrations.New(db.NewIntegrationStore(database), vault).Get(integrations.CloudBackup)
	if !ok {
		return nil
	}

	entries, err := os.ReadDir(outDir)
	if err != nil {
		return err
	}
	newest := ""
	for _, e := range entries {
		// Names are UTC timestamps, so the lexically largest one is the newest.
		if e.IsDir() && strings.HasPrefix(e.Name(), "orbit-backup-") && e.Name() > newest {
			newest = e.Name()
		}
	}
	if newest == "" {
		return fmt.Errorf("no backup found in %s", outDir)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	fmt.Printf("☁️  uploading %s to %s\n", newest, remote)
	return backup.Upload(ctx, conf, remote, filepath.Join(outDir, newest), cfg.SecretsDir)
}
