package s3

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"orbit/internal/config"
	"orbit/internal/db"
)

func TestReedSolomonTwoDriveFailureAndAutoHeal(t *testing.T) {
	es, err := NewErasureSet(DefaultDataShards, DefaultParityShards)
	if err != nil {
		t.Fatalf("NewErasureSet: %v", err)
	}

	payload := bytes.Repeat([]byte("OrbitSovereignErasureCoding-2026!"), 300) // ~9.9 KB
	shards, err := es.Encode(payload)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if len(shards) != 6 {
		t.Fatalf("expected 6 shards, got %d", len(shards))
	}

	// Simulate 2 simultaneous drive failures (data shard 1 + parity shard 4 lost)
	shards[1] = nil
	shards[4] = nil

	rebuilt, err := es.Reconstruct(shards, int64(len(payload)))
	if err != nil {
		t.Fatalf("Reconstruct after 2 failures: %v", err)
	}
	if !bytes.Equal(rebuilt, payload) {
		t.Fatalf("rebuilt payload mismatch (len %d vs %d)", len(rebuilt), len(payload))
	}
	if shards[1] == nil || shards[4] == nil {
		t.Fatalf("expected missing shards to be reconstructed in-place")
	}
}

func TestErasureDiskWriteReadAndAutoHeal(t *testing.T) {
	tmpDir := t.TempDir()
	eng := &S3Engine{
		cfg: &config.Config{StorageDir: tmpDir},
	}

	data := []byte("Critical production backup payload for RAID-less multi-drive parity test!")
	if err := eng.writeErasureObject("assets", "obj_test1", data); err != nil {
		t.Fatalf("writeErasureObject: %v", err)
	}

	drives := eng.DriveDirs()
	// Delete drive 0 and corrupt drive 2
	_ = os.Remove(eng.shardPath(drives[0], "assets", "obj_test1", 0))
	_ = os.WriteFile(eng.shardPath(drives[2], "assets", "obj_test1", 2), []byte("corrupted-bytes"), 0600)

	got, healed, err := eng.readErasureObject("assets", "obj_test1")
	if err != nil {
		t.Fatalf("readErasureObject: %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("data mismatch after 2-drive failure")
	}
	if healed != 2 {
		t.Fatalf("expected 2 auto-healed shards, got %d", healed)
	}

	// Second read should have 0 healed because both were repaired on disk
	_, healed2, err := eng.readErasureObject("assets", "obj_test1")
	if err != nil || healed2 != 0 {
		t.Fatalf("expected 0 healed after repair, got healed=%d err=%v", healed2, err)
	}
}

func TestIAMPolicyEnforcement(t *testing.T) {
	k := &db.StorageIAMKey{
		BucketScope: "backups",
		PrefixScope: "prod/",
		Permissions: "read,write",
		Active:      true,
	}

	if !checkIAMPolicy(k, http.MethodGet, "backups", "prod/db.sql") {
		t.Fatal("expected GET on scoped bucket/prefix to succeed")
	}
	if !checkIAMPolicy(k, http.MethodPut, "backups", "prod/db.sql") {
		t.Fatal("expected PUT on scoped bucket/prefix to succeed")
	}
	if checkIAMPolicy(k, http.MethodDelete, "backups", "prod/db.sql") {
		t.Fatal("expected DELETE to be denied without delete/admin permission")
	}
	if checkIAMPolicy(k, http.MethodGet, "other-bucket", "prod/db.sql") {
		t.Fatal("expected wrong bucket to be denied")
	}
	if checkIAMPolicy(k, http.MethodGet, "backups", "staging/db.sql") {
		t.Fatal("expected wrong prefix to be denied")
	}
}
