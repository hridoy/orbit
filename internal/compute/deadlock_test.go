package compute

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"orbit/internal/config"
	"orbit/internal/db"
)

// StartProcess must return promptly. A self-deadlock on the supervisor mutex used to hang forever.
func TestStartProcessDoesNotDeadlock(t *testing.T) {
	cfg := &config.Config{NodeID: "node_test"}
	s := NewProcessSupervisor(cfg, NewCgroupController())
	app := &db.App{AppID: "dl", BinaryEntrypoint: "./does-not-exist", RAMLimitMB: 64, CPUQuotaPct: 100, PIDsLimit: 16}

	done := make(chan error, 1)
	go func() {
		_, err := s.StartProcess(context.Background(), app, 18999, t.TempDir(), map[string]string{})
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected an error for a missing binary")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("StartProcess deadlocked")
	}
}

func TestIsRestorable(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "a1", "versions", "v1"), 0700)
	v1, missing := "v1", "gone"
	cases := []struct {
		name string
		app  *db.App
		want bool
	}{
		{"live version on disk", &db.App{AppID: "a1", CurrentVersion: &v1}, true},
		{"version dir removed", &db.App{AppID: "a1", CurrentVersion: &missing}, false},
		{"never deployed", &db.App{AppID: "a1"}, false},
	}
	for _, c := range cases {
		if got := isRestorable(dir, c.app); got != c.want {
			t.Errorf("%s: got %v want %v", c.name, got, c.want)
		}
	}
}
