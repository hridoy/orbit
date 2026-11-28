package server

import (
	"testing"
	"time"
)

func TestParseBackupName(t *testing.T) {
	got, ok := parseBackupName("orbit-backup-20261128-111541")
	if !ok || got.Format("2006-01-02 15:04:05") != "2026-11-28 11:15:41" {
		t.Fatalf("parse failed: %v %v", got, ok)
	}
	if _, ok := parseBackupName("random"); ok {
		t.Fatal("expected failure for a non-backup name")
	}
}

func TestStaleBackupAlert(t *testing.T) {
	h := &AlertsHandler{rootPath: t.TempDir(), backups: &BackupHandler{dir: t.TempDir()}}
	// No backups directory means "No backups" is raised, and nothing else from the backup rule.
	alerts := h.Check(time.Now())
	found := false
	for _, a := range alerts {
		if a.Title == "No backups" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected 'No backups' alert, got %+v", alerts)
	}
}
