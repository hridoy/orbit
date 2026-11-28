package backup

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// The round trip uses rclone's local backend as the "cloud", so it proves the copy and the
// prune commands without needing a Google or Dropbox login.
func TestUploadToLocalRemote(t *testing.T) {
	if _, err := exec.LookPath("rclone"); err != nil {
		t.Skip("rclone not installed")
	}
	root := t.TempDir()
	src := filepath.Join(root, "orbit-backup-20260101-000000")
	if err := os.MkdirAll(src, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "orbit.db"), []byte("snapshot"), 0600); err != nil {
		t.Fatal(err)
	}
	remoteDir := filepath.Join(root, "remote")
	conf := "[dst]\ntype = local\n"

	if err := Upload(context.Background(), conf, "dst:"+remoteDir, src, root); err != nil {
		t.Fatalf("upload: %v", err)
	}
	if _, err := os.Stat(filepath.Join(remoteDir, "orbit-backup-20260101-000000", "orbit.db")); err != nil {
		t.Fatalf("backup file not copied: %v", err)
	}
}

func TestUploadRejectsRemoteWithoutColon(t *testing.T) {
	err := Upload(context.Background(), "[x]\ntype = local\n", "orbit-backups", t.TempDir(), t.TempDir())
	if err == nil {
		t.Fatal("expected an error for a remote path without remote: prefix")
	}
}
