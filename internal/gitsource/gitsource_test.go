package gitsource

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestNormalizeRepoURL(t *testing.T) {
	good := map[string]string{
		"https://github.com/acme/app":     "https://github.com/acme/app",
		"https://github.com/acme/app.git": "https://github.com/acme/app",
		"https://github.com/acme/app/":    "https://github.com/acme/app",
	}
	for in, want := range good {
		got, err := NormalizeRepoURL(in)
		if err != nil || got != want {
			t.Errorf("%s: got %q, %v", in, got, err)
		}
	}
	bad := []string{
		"http://github.com/acme/app",
		"https://gitlab.com/acme/app",
		"https://github.com.evil.com/acme/app",
		"https://user:pw@github.com/acme/app",
		"https://github.com/acme/app?x=1",
		"https://github.com/acme/../etc",
		"git@github.com:acme/app.git",
		"https://github.com/acme",
	}
	for _, in := range bad {
		if _, err := NormalizeRepoURL(in); err == nil {
			t.Errorf("%s: expected rejection", in)
		}
	}
}

func TestRefAndEntrypointValidation(t *testing.T) {
	for _, r := range []string{"main", "release/v1.2", "v1.0.0"} {
		if err := ValidateRef(r); err != nil {
			t.Errorf("ref %q rejected: %v", r, err)
		}
	}
	for _, r := range []string{"-upload-pack=x", "a..b", "x y", ""} {
		if err := ValidateRef(r); err == nil {
			t.Errorf("ref %q accepted", r)
		}
	}
	for _, e := range []string{"/bin/sh", "../evil", "./server"} {
		err := ValidateEntrypoint(e)
		if e == "./server" && err != nil {
			t.Errorf("entrypoint %q rejected", e)
		}
		if e != "./server" && err == nil {
			t.Errorf("entrypoint %q accepted", e)
		}
	}
}

// TestFetchFromLocalRepo exercises clone, commit lookup, and manifest reading without the network.
func TestFetchFromLocalRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	src := t.TempDir()
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", src}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	run("init", "-b", "main")
	os.WriteFile(filepath.Join(src, "go.mod"), []byte("module demo\n\ngo 1.21\n"), 0644)
	run("add", ".")
	run("commit", "-m", "init")

	ctx := context.Background()
	latest, err := LatestCommit(ctx, src, "main")
	if err != nil {
		t.Fatalf("LatestCommit: %v", err)
	}

	dest := filepath.Join(t.TempDir(), "v")
	commit, err := Fetch(ctx, src, "main", dest)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if commit != latest {
		t.Fatalf("clone commit %s != ls-remote %s", commit, latest)
	}
	if _, err := os.Stat(filepath.Join(dest, ".git")); err == nil {
		t.Fatal(".git was not removed")
	}

	m, err := ReadManifest(dest)
	if err != nil {
		t.Fatalf("ReadManifest: %v", err)
	}
	if m.Build != "go build -o server ." || m.Entrypoint != "./server" {
		t.Fatalf("unexpected inferred manifest: %+v", m)
	}
}
