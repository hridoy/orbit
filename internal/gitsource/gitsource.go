// Package gitsource installs apps from public GitHub repositories: it validates the
// URL, resolves the latest commit, clones a single commit, and reads the build manifest.
package gitsource

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ErrNotGitHub     = errors.New("only https://github.com/<owner>/<repo> URLs are supported")
	ErrBadRef        = errors.New("branch or tag name contains invalid characters")
	ErrBadEntrypoint = errors.New("entrypoint must be a relative path inside the repository")
	ErrNoCommit      = errors.New("ref not found in repository")

	ownerRepoRe = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,100}/[A-Za-z0-9_.-]{1,100}$`)
	refRe       = regexp.MustCompile(`^[A-Za-z0-9._/-]{1,100}$`)
	commitRe    = regexp.MustCompile(`^[0-9a-f]{40}$`)
)

// Manifest is read from orbit.json at the repo root. Without it, a Go module is built with
// `go build -o server .`.
type Manifest struct {
	Build      string `json:"build"`
	Entrypoint string `json:"entrypoint"`
}

// NormalizeRepoURL accepts only https://github.com/owner/repo and returns the canonical form.
func NormalizeRepoURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme != "https" || u.Host != "github.com" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", ErrNotGitHub
	}
	path := strings.TrimSuffix(strings.Trim(u.Path, "/"), ".git")
	if !ownerRepoRe.MatchString(path) {
		return "", ErrNotGitHub
	}
	return "https://github.com/" + path, nil
}

// ValidateRef checks a branch or tag name. The name is also passed as a separate argv entry,
// so it can never be read as a git option.
func ValidateRef(ref string) error {
	if !refRe.MatchString(ref) || strings.Contains(ref, "..") || strings.HasPrefix(ref, "-") {
		return ErrBadRef
	}
	return nil
}

// ValidateEntrypoint requires a relative path that stays inside the version directory.
func ValidateEntrypoint(ep string) error {
	if ep == "" || filepath.IsAbs(ep) || strings.Contains(ep, "..") {
		return ErrBadEntrypoint
	}
	return nil
}

func gitEnv() []string {
	return append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
}

// LatestCommit returns the commit a branch or tag currently points to.
func LatestCommit(ctx context.Context, repoURL, ref string) (string, error) {
	if err := ValidateRef(ref); err != nil {
		return "", err
	}
	for _, name := range []string{"refs/heads/" + ref, "refs/tags/" + ref + "^{}", "refs/tags/" + ref} {
		cmd := exec.CommandContext(ctx, "git", "ls-remote", "--", repoURL, name)
		cmd.Env = gitEnv()
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("git ls-remote failed: %w", err)
		}
		fields := strings.Fields(strings.SplitN(out.String(), "\n", 2)[0])
		if len(fields) >= 1 && commitRe.MatchString(fields[0]) {
			return fields[0], nil
		}
	}
	return "", ErrNoCommit
}

// Fetch clones exactly one commit of ref into dest (shallow, no .git afterward) and
// returns the commit. dest must not exist yet.
func Fetch(ctx context.Context, repoURL, ref, dest string) (string, error) {
	if err := ValidateRef(ref); err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", "--single-branch", "--branch", ref, "--", repoURL, dest)
	cmd.Env = gitEnv()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git clone failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	rev := exec.CommandContext(ctx, "git", "-C", dest, "rev-parse", "HEAD")
	out, err := rev.Output()
	if err != nil {
		return "", fmt.Errorf("read commit: %w", err)
	}
	commit := strings.TrimSpace(string(out))
	if !commitRe.MatchString(commit) {
		return "", fmt.Errorf("unexpected commit id %q", commit)
	}
	if err := os.RemoveAll(filepath.Join(dest, ".git")); err != nil {
		return "", err
	}
	return commit, nil
}

// ReadManifest loads orbit.json, or infers a Go build when go.mod is present.
func ReadManifest(dir string) (Manifest, error) {
	var m Manifest
	raw, err := os.ReadFile(filepath.Join(dir, "orbit.json"))
	if err == nil {
		if err := json.Unmarshal(raw, &m); err != nil {
			return m, fmt.Errorf("orbit.json is not valid JSON: %w", err)
		}
	} else if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
		m = Manifest{Build: "go build -o server .", Entrypoint: "./server"}
	} else {
		return m, errors.New("repo has no orbit.json and is not a Go module")
	}
	if m.Entrypoint == "" {
		m.Entrypoint = "./server"
	}
	if m.Build == "" {
		return m, errors.New("orbit.json must set \"build\"")
	}
	if err := ValidateEntrypoint(m.Entrypoint); err != nil {
		return m, err
	}
	return m, nil
}
