package update

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func setupRelease(t *testing.T, binary []byte, sha string) (pubHex string, cleanup func()) {
	t.Helper()
	pub, priv, _ := ed25519.GenerateKey(nil)
	arch := runtime.GOARCH
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	Repo = "acme/orbit"
	githubBase = srv.URL
	PublicKeyHex = hex.EncodeToString(pub)

	m := Manifest{Version: "0.2.0", Assets: map[string]Asset{
		runtime.GOOS + "/" + arch: {
			URL:    srv.URL + "/acme/orbit/releases/download/v0.2.0/orbit-agent-linux-" + arch,
			SHA256: sha,
		},
	}}
	raw, _ := json.Marshal(m)
	sig := ed25519.Sign(priv, raw)
	mux.HandleFunc("/acme/orbit/releases/latest/download/manifest.json", func(w http.ResponseWriter, r *http.Request) { w.Write(raw) })
	mux.HandleFunc("/acme/orbit/releases/latest/download/manifest.sig", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(hex.EncodeToString(sig))) })
	mux.HandleFunc("/acme/orbit/releases/download/v0.2.0/orbit-agent-linux-"+arch, func(w http.ResponseWriter, r *http.Request) { w.Write(binary) })
	return hex.EncodeToString(pub), srv.Close
}

func TestVerifyManifest(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	msg := []byte(`{"version":"1.0.0"}`)
	sig := hex.EncodeToString(ed25519.Sign(priv, msg))
	pubHex := hex.EncodeToString(pub)

	if err := VerifyManifest(msg, sig, pubHex); err != nil {
		t.Fatalf("valid signature rejected: %v", err)
	}
	if err := VerifyManifest([]byte(`{"version":"9.9.9"}`), sig, pubHex); err == nil {
		t.Fatal("tampered manifest accepted")
	}
	other, _, _ := ed25519.GenerateKey(nil)
	if err := VerifyManifest(msg, sig, hex.EncodeToString(other)); err == nil {
		t.Fatal("signature from wrong key accepted")
	}
}

func TestCompareVersions(t *testing.T) {
	cases := map[[2]string]int{
		{"0.2.10", "0.2.9"}:    1,
		{"v1.0.0", "1.0.0"}:    0,
		{"0.1.0-dev", "0.1.0"}: 0,
		{"0.1.0", "0.2.0"}:     -1,
	}
	for in, want := range cases {
		if got := CompareVersions(in[0], in[1]); got != want {
			t.Errorf("compare(%s,%s)=%d want %d", in[0], in[1], got, want)
		}
	}
}

func TestApplySwapsBinaryAndKeepsPrevious(t *testing.T) {
	newBin := []byte("new-agent-binary")
	sum := sha256.Sum256(newBin)
	pubHex, done := setupRelease(t, newBin, hex.EncodeToString(sum[:]))
	defer done()
	PublicKeyHex = pubHex

	dir := t.TempDir()
	exe := filepath.Join(dir, "orbit-agent")
	os.WriteFile(exe, []byte("old-agent-binary"), 0755)

	version, err := Apply(context.Background(), exe)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if version != "0.2.0" {
		t.Fatalf("version %q", version)
	}
	got, _ := os.ReadFile(exe)
	if string(got) != string(newBin) {
		t.Fatal("binary was not swapped")
	}
	prev, _ := os.ReadFile(exe + ".prev")
	if string(prev) != "old-agent-binary" {
		t.Fatal("previous binary not kept")
	}
}

func TestApplyRejectsChecksumMismatch(t *testing.T) {
	pubHex, done := setupRelease(t, []byte("tampered"), hex.EncodeToString(make([]byte, 32)))
	defer done()
	PublicKeyHex = pubHex

	dir := t.TempDir()
	exe := filepath.Join(dir, "orbit-agent")
	os.WriteFile(exe, []byte("old"), 0755)

	if _, err := Apply(context.Background(), exe); err == nil {
		t.Fatal("checksum mismatch accepted")
	}
	got, _ := os.ReadFile(exe)
	if string(got) != "old" {
		t.Fatal("binary changed despite failed verification")
	}
	if _, err := os.Stat(exe + ".new"); err == nil {
		t.Fatal("partial download left behind")
	}
}

func TestRollbackAfterRepeatedCrashes(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "orbit-agent")
	os.WriteFile(exe, []byte("broken-new"), 0755)
	os.WriteFile(exe+".prev", []byte("good-old"), 0755)
	if err := MarkPending(dir, "0.2.0"); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < maxStartsBeforeRollback; i++ {
		restored, err := StartupCheck(dir, exe)
		if err != nil || restored {
			t.Fatalf("start %d: restored=%v err=%v", i+1, restored, err)
		}
	}
	restored, err := StartupCheck(dir, exe)
	if err != nil || !restored {
		t.Fatalf("expected rollback, restored=%v err=%v", restored, err)
	}
	got, _ := os.ReadFile(exe)
	if string(got) != "good-old" {
		t.Fatal("previous binary not restored")
	}
}

func TestNoRollbackWhenHealthy(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "orbit-agent")
	os.WriteFile(exe, []byte("x"), 0755)
	os.WriteFile(exe+".prev", []byte("y"), 0755)
	// No pending state: startup must never roll back.
	for i := 0; i < 10; i++ {
		if restored, _ := StartupCheck(dir, exe); restored {
			t.Fatal("rolled back without a pending update")
		}
	}
}
