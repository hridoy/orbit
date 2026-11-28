package database

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// enginesRoot is where downloaded server versions live: engines/<engine>/<version>/bin/<server>.
const enginesRoot = "/opt/orbit/engines"

// Install states shown to the dashboard.
const (
	StatusQueued     = "queued"
	StatusInstalling = "installing"
	StatusFailed     = "failed"
)

// installJob tracks one background install.
type installJob struct {
	Status string
	Error  string
	Log    string
}

type installer struct {
	mu    sync.Mutex
	jobs  map[string]*installJob // key: engine/version
	queue chan func()            // installs run one at a time: they are CPU and disk heavy
}

func newInstaller() *installer {
	in := &installer{jobs: map[string]*installJob{}, queue: make(chan func(), 64)}
	go func() {
		for fn := range in.queue {
			fn()
		}
	}()
	return in
}

func jobKey(engine, version string) string { return engine + "/" + version }

// jobStatus returns the status of an install, or "" when none is running or failed.
func (in *installer) jobStatus(engine, version string) (string, string) {
	in.mu.Lock()
	defer in.mu.Unlock()
	j := in.jobs[jobKey(engine, version)]
	if j == nil {
		return "", ""
	}
	return j.Status, j.Error
}

// StartInstall begins installing a catalog version in the background. It returns an error
// straight away for versions that are not in the catalog or already running.
func (m *DatabaseManager) StartInstall(engine, version string) error {
	entry, ok := catalogFor(engine, version)
	if !ok {
		return fmt.Errorf("%s %s is not an available version", engine, version)
	}
	if m.isInstalled(engine, version) {
		return fmt.Errorf("%s %s is already installed", engine, version)
	}
	in := m.installer
	in.mu.Lock()
	if j := in.jobs[jobKey(engine, version)]; j != nil && (j.Status == StatusInstalling || j.Status == StatusQueued) {
		in.mu.Unlock()
		return fmt.Errorf("%s %s is already installing", engine, version)
	}
	job := &installJob{Status: StatusQueued}
	in.jobs[jobKey(engine, version)] = job
	in.mu.Unlock()

	in.queue <- func() {
		in.mu.Lock()
		job.Status = StatusInstalling
		in.mu.Unlock()

		err := runInstall(entry, job)
		in.mu.Lock()
		defer in.mu.Unlock()
		if err != nil {
			job.Status = StatusFailed
			job.Error = err.Error()
			log.Printf("⚠️ install %s %s failed: %v", engine, version, err)
			return
		}
		delete(in.jobs, jobKey(engine, version))
		m.RefreshInstalled()
		log.Printf("✅ installed %s %s", engine, version)
	}
	return nil
}

func (m *DatabaseManager) isInstalled(engine, version string) bool {
	for _, inst := range m.installedSnapshot() {
		if inst.Engine == engine && inst.Version == version {
			return true
		}
	}
	return false
}

// runInstall does the work for one catalog entry. Output is logged for the operator.
func runInstall(e CatalogEntry, job *installJob) error {
	logf := func(format string, args ...any) {
		line := fmt.Sprintf(format, args...)
		job.Log += line + "\n"
		log.Printf("install %s %s: %s", e.Engine, e.Version, line)
	}
	switch e.Method {
	case MethodApt:
		logf("apt-get install %s", e.Pkg)
		return aptInstall(e.Pkg, logf)
	case MethodTarball:
		return installTarball(e, logf)
	case MethodSource:
		return installSource(e, logf)
	}
	return fmt.Errorf("unknown install method %q", e.Method)
}

// aptMu serializes apt: two installs at once fight over the dpkg lock and one fails.
var aptMu sync.Mutex

func aptInstall(pkg string, logf func(string, ...any)) error {
	aptMu.Lock()
	defer aptMu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "apt-get", "install", "-y", "-q", pkg)
	cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("apt-get install %s: %v: %s", pkg, err, tail(out))
	}
	logf("installed package %s", pkg)
	return nil
}

// downloadResumable fetches url into path, resuming a partial file after stalls or dropped
// HTTP/2 streams. Vendor archives are hundreds of MB, and a single stalled connection used to
// fail the whole install. curl gives up on a connection that drops below 10 KB/s for a minute.
func downloadResumable(url, path string) error {
	var err error
	for attempt := 1; attempt <= 6; attempt++ {
		err = run(30*time.Minute, "curl", "-fL", "--http1.1", "-C", "-",
			"--speed-limit", "10000", "--speed-time", "60",
			"--retry", "2", "--retry-delay", "3", "-o", path, url)
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(attempt*5) * time.Second)
	}
	return err
}

// installTarball downloads an official archive, checks its MD5 when the vendor publishes one,
// and unpacks it into engines/<engine>/<version>.
func installTarball(e CatalogEntry, logf func(string, ...any)) error {
	dest := filepath.Join(enginesRoot, e.Engine, e.Version)
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	tmp := filepath.Join(enginesRoot, ".downloads")
	if err := os.MkdirAll(tmp, 0755); err != nil {
		return err
	}
	name := filepath.Base(e.URL)
	archive := filepath.Join(tmp, e.Engine+"-"+e.Version+"-"+name)
	defer os.Remove(archive)

	logf("downloading %s", e.URL)
	if err := downloadResumable(e.URL, archive); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	if e.MD5URL != "" || strings.Contains(e.URL, "cdn.mysql.com") {
		md5URL := e.URL + ".md5"
		sum, err := runOut(time.Minute, "curl", "-fsSL", md5URL)
		if err != nil {
			return fmt.Errorf("could not read checksum: %w", err)
		}
		want := strings.Fields(string(sum))
		if len(want) == 0 {
			return fmt.Errorf("checksum file is empty")
		}
		got, err := runOut(2*time.Minute, "md5sum", archive)
		if err != nil {
			return err
		}
		if !strings.HasPrefix(string(got), want[0]) {
			return fmt.Errorf("checksum mismatch: the download is not the vendor's file")
		}
		logf("checksum verified")
	}

	logf("extracting into %s", dest)
	flags := "-xzf"
	if strings.HasSuffix(name, ".tar.xz") {
		flags = "-xJf"
	}
	if strings.HasSuffix(name, ".tgz") {
		flags = "-xzf"
	}
	if err := run(10*time.Minute, "tar", flags, archive, "-C", dest, "--strip-components=1"); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}
	bin := filepath.Join(dest, "bin", e.Bin)
	if _, err := os.Stat(bin); err != nil {
		return fmt.Errorf("%s not found after extracting", bin)
	}
	logf("installed %s", bin)
	return nil
}

// installSource downloads an official source archive and builds it into engines/<engine>/<version>.
func installSource(e CatalogEntry, logf func(string, ...any)) error {
	dest := filepath.Join(enginesRoot, e.Engine, e.Version)
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	build := filepath.Join(enginesRoot, ".build", e.Engine+"-"+e.Version)
	_ = os.RemoveAll(build)
	if err := os.MkdirAll(build, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(build)

	archive := filepath.Join(build, "src.tar.gz")
	logf("downloading %s", e.URL)
	if err := downloadResumable(e.URL, archive); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	logf("unpacking source")
	if err := run(5*time.Minute, "tar", "-xzf", archive, "-C", build, "--strip-components=1"); err != nil {
		return fmt.Errorf("unpack failed: %w", err)
	}
	logf("building (this takes a few minutes)")
	if err := run(30*time.Minute, "make", "-C", build, "-j2"); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	if err := run(10*time.Minute, "make", "-C", build, "PREFIX="+dest, "install"); err != nil {
		return fmt.Errorf("install failed: %w", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "bin", e.Bin)); err != nil {
		return fmt.Errorf("%s missing after build", e.Bin)
	}
	logf("installed %s", filepath.Join(dest, "bin", e.Bin))
	return nil
}

func run(timeout time.Duration, name string, args ...string) error {
	_, err := runOut(timeout, name, args...)
	return err
}

func runOut(timeout time.Duration, name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.Bytes(), fmt.Errorf("%s: %v: %s", name, err, tail(out.Bytes()))
	}
	return out.Bytes(), nil
}

func tail(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 400 {
		s = s[len(s)-400:]
	}
	return s
}
