package jail

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"orbit/internal/config"
)

func setupTestJail(t *testing.T) (*Jail, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "orbit_jail_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	cfg := &config.Config{
		RootPath:   tmpDir,
		AppsDir:    filepath.Join(tmpDir, "apps"),
		StorageDir: filepath.Join(tmpDir, "storage"),
	}
	_ = os.MkdirAll(cfg.AppsDir, 0700)
	_ = os.MkdirAll(cfg.StorageDir, 0700)

	return NewJail(cfg), tmpDir
}

func TestJailWriteReadListDelete(t *testing.T) {
	j, tmpDir := setupTestJail(t)
	defer os.RemoveAll(tmpDir)

	appID := "demo_app"
	testContent := "hello sovereign world"

	// 1. Write file
	written, err := j.WriteFile(ScopeApp, appID, "config/settings.json", bytes.NewBufferString(testContent))
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if written != int64(len(testContent)) {
		t.Errorf("expected %d bytes written, got %d", len(testContent), written)
	}

	// 2. Stat file
	stat, err := j.Stat(ScopeApp, appID, "config/settings.json")
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if stat.Size != int64(len(testContent)) || stat.IsDir {
		t.Errorf("unexpected stat: %+v", stat)
	}

	// 3. Read file
	rc, size, err := j.OpenFileRead(ScopeApp, appID, "config/settings.json")
	if err != nil {
		t.Fatalf("OpenFileRead failed: %v", err)
	}
	defer rc.Close()

	if size != int64(len(testContent)) {
		t.Errorf("expected size %d, got %d", len(testContent), size)
	}

	body, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if string(body) != testContent {
		t.Errorf("expected '%s', got '%s'", testContent, string(body))
	}

	// 4. List directory
	items, err := j.ListDirectory(ScopeApp, appID, "config")
	if err != nil {
		t.Fatalf("ListDirectory failed: %v", err)
	}
	if len(items) != 1 || items[0].Name != "settings.json" {
		t.Errorf("unexpected list items: %+v", items)
	}

	// 5. Delete file
	if err := j.Delete(ScopeApp, appID, "config/settings.json"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, _, err = j.OpenFileRead(ScopeApp, appID, "config/settings.json")
	if err != ErrFileNotFound {
		t.Errorf("expected ErrFileNotFound after delete, got %v", err)
	}
}

func TestJailPathTraversalAttacks(t *testing.T) {
	j, tmpDir := setupTestJail(t)
	defer os.RemoveAll(tmpDir)

	appID := "demo_app"

	// Attack 1: Relative dot-dot escape
	_, err := j.WriteFile(ScopeApp, appID, "../../etc/shadow", bytes.NewBufferString("malicious"))
	if err != ErrAccessDenied {
		t.Errorf("expected ErrAccessDenied for dot-dot path, got %v", err)
	}

	// Attack 2: Escaping via symlink pointing outside jail root
	appRoot, err := j.ResolveRoot(ScopeApp, appID)
	if err != nil {
		t.Fatalf("ResolveRoot failed: %v", err)
	}

	// Create an external secret file outside the jail root
	secretOutside := filepath.Join(tmpDir, "host_secret.txt")
	_ = os.WriteFile(secretOutside, []byte("super_secret_host_key"), 0600)

	// Create a symlink inside the app directory pointing to secretOutside
	symlinkPath := filepath.Join(appRoot, "escape_link")
	if err := os.Symlink(secretOutside, symlinkPath); err != nil {
		t.Skipf("cannot create symlink on this filesystem: %v", err)
	}

	// Attempt to read the file through the symlink
	_, _, err = j.OpenFileRead(ScopeApp, appID, "escape_link")
	if err != ErrAccessDenied {
		t.Errorf("expected ErrAccessDenied when reading through escaping symlink, got %v", err)
	}

	// Attack 3: Scope ID path traversal
	_, err = j.ResolveRoot(ScopeApp, "../../../etc")
	if err != ErrAccessDenied {
		t.Errorf("expected ErrAccessDenied for malicious scopeID, got %v", err)
	}
}
