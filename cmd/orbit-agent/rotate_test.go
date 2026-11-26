package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/s3"
	"orbit/internal/security"
)

func randKey(t *testing.T, fill byte) []byte {
	t.Helper()
	k := bytes.Repeat([]byte{fill}, 32)
	return k
}

func setupRotation(t *testing.T) (*config.Config, *db.DB, []byte) {
	t.Helper()
	root := t.TempDir()
	cfg := &config.Config{
		RootPath:   root,
		SecretsDir: filepath.Join(root, "secrets"),
		DataDir:    filepath.Join(root, "data"),
		DBPath:     filepath.Join(root, "data", "orbit.db"),
	}
	for _, d := range []string{cfg.SecretsDir, cfg.DataDir} {
		if err := os.MkdirAll(d, 0700); err != nil {
			t.Fatal(err)
		}
	}
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })

	oldKey := randKey(t, 7)
	cfg.MasterKey = oldKey
	if err := os.WriteFile(filepath.Join(cfg.SecretsDir, "master.key"), oldKey, 0600); err != nil {
		t.Fatal(err)
	}
	if err := ensureCanary(database, oldKey); err != nil {
		t.Fatal(err)
	}

	vault, _ := security.NewVault(oldKey)
	if err := db.NewSecretStore(database, vault).SetSecret("s1", "app", "a1", "DB_URL", "postgres://secret"); err != nil {
		t.Fatal(err)
	}
	creds := &s3.S3Credentials{AccessKey: "ORBTEST", SecretKey: "sekret"}
	if err := s3.SealS3Credentials(filepath.Join(cfg.DataDir, s3.S3CredentialsFile), oldKey, creds); err != nil {
		t.Fatal(err)
	}
	return cfg, database, oldKey
}

func TestRotateReencryptsEverything(t *testing.T) {
	cfg, database, oldKey := setupRotation(t)
	newKey := randKey(t, 9)

	counts, err := rotateMasterKey(cfg, database, oldKey, newKey, false)
	if err != nil {
		t.Fatalf("rotate: %v", err)
	}
	if counts["secrets"] != 1 || counts["s3_credentials"] != 1 {
		t.Fatalf("unexpected counts: %v", counts)
	}

	// Secret decrypts only with the new key.
	newV, _ := security.NewVault(newKey)
	oldV, _ := security.NewVault(oldKey)
	var c, n string
	if err := database.QueryRow(`SELECT ciphertext_base64, nonce_base64 FROM secrets WHERE secret_id='s1'`).Scan(&c, &n); err != nil {
		t.Fatal(err)
	}
	if plain, err := newV.Decrypt(c, n); err != nil || string(plain) != "postgres://secret" {
		t.Fatalf("secret not readable with new key: %v", err)
	}
	if _, err := oldV.Decrypt(c, n); err == nil {
		t.Fatal("secret still decrypts with old key")
	}

	if !canaryOK(database, newKey) || canaryOK(database, oldKey) {
		t.Fatal("canary does not reflect new key")
	}
	if _, err := s3.OpenS3Credentials(filepath.Join(cfg.DataDir, s3.S3CredentialsFile), newKey); err != nil {
		t.Fatalf("s3 creds not readable with new key: %v", err)
	}

	got, _ := os.ReadFile(filepath.Join(cfg.SecretsDir, "master.key"))
	if !bytes.Equal(got, newKey) {
		t.Fatal("master.key was not replaced")
	}
}

func TestDryRunChangesNothing(t *testing.T) {
	cfg, database, oldKey := setupRotation(t)
	if _, err := rotateMasterKey(cfg, database, oldKey, randKey(t, 9), true); err != nil {
		t.Fatal(err)
	}
	if !canaryOK(database, oldKey) {
		t.Fatal("dry run changed the canary")
	}
	if _, err := os.Stat(filepath.Join(cfg.DataDir, s3.S3CredentialsFile+".new")); err == nil {
		t.Fatal("dry run left a .new file")
	}
}

func TestResolveFinishesInterruptedRotation(t *testing.T) {
	cfg, database, oldKey := setupRotation(t)
	newKey := randKey(t, 9)
	if _, err := rotateMasterKey(cfg, database, oldKey, newKey, false); err != nil {
		t.Fatal(err)
	}
	// Simulate a crash after DB commit: the key file is still the old one.
	masterPath := filepath.Join(cfg.SecretsDir, "master.key")
	os.WriteFile(masterPath+".new", newKey, 0600)
	os.WriteFile(masterPath, oldKey, 0600)

	got, err := resolveMasterKey(cfg, database)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if !bytes.Equal(got, newKey) {
		t.Fatal("resolve did not pick the new key")
	}
	if _, err := os.Stat(masterPath + ".new"); err == nil {
		t.Fatal("pending key file not cleaned up")
	}
}

func TestResolveRollsBackUncommittedRotation(t *testing.T) {
	cfg, database, oldKey := setupRotation(t)
	masterPath := filepath.Join(cfg.SecretsDir, "master.key")
	// Crash before DB commit: database still uses the old key, a pending key file exists.
	os.WriteFile(masterPath+".new", randKey(t, 9), 0600)

	got, err := resolveMasterKey(cfg, database)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if !bytes.Equal(got, oldKey) {
		t.Fatal("resolve did not keep the old key")
	}
	if _, err := os.Stat(masterPath + ".new"); err == nil {
		t.Fatal("uncommitted pending key not removed")
	}
}

func TestResolveRefusesMismatchedKey(t *testing.T) {
	cfg, database, _ := setupRotation(t)
	os.WriteFile(filepath.Join(cfg.SecretsDir, "master.key"), randKey(t, 3), 0600)
	if _, err := resolveMasterKey(cfg, database); err == nil {
		t.Fatal("expected refusal when key does not match the database")
	}
}
