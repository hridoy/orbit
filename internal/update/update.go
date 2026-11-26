// Package update checks GitHub Releases for a newer agent, verifies it, and swaps the binary.
// Trust comes from an Ed25519 signature over manifest.json. The public key is injected at build
// time, so a compromised download server cannot push a binary without the release key.
package update

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// githubBase is overridden in tests only.
var githubBase = "https://github.com"

// Set at build time by the release workflow (-ldflags -X). Empty in local builds, which
// disables updates.
var (
	Repo         = "" // "owner/name"
	PublicKeyHex = "" // hex Ed25519 public key
)

const (
	maxManifestBytes = 1 << 20
	maxBinaryBytes   = 200 << 20
)

var (
	ErrNotConfigured = errors.New("update source is not configured in this build")
	ErrBadSignature  = errors.New("release manifest signature is invalid")
	ErrNoAsset       = errors.New("no release asset for this platform")
	ErrBadChecksum   = errors.New("downloaded binary checksum does not match the manifest")
)

// Manifest is the signed release description published with each release.
type Manifest struct {
	Version string           `json:"version"`
	Assets  map[string]Asset `json:"assets"` // key: "linux/amd64"
}

type Asset struct {
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

// Status is what the dashboard shows.
type Status struct {
	Current   string `json:"current"`
	Latest    string `json:"latest"`
	Available bool   `json:"available"`
	Supported bool   `json:"supported"`
	Reason    string `json:"reason,omitempty"`
}

// VerifyManifest checks the detached signature over the exact manifest bytes.
func VerifyManifest(manifest []byte, sigHex, pubHex string) error {
	pub, err := hex.DecodeString(pubHex)
	if err != nil || len(pub) != ed25519.PublicKeySize {
		return errors.New("invalid public key")
	}
	sig, err := hex.DecodeString(strings.TrimSpace(sigHex))
	if err != nil || len(sig) != ed25519.SignatureSize {
		return ErrBadSignature
	}
	if !ed25519.Verify(pub, manifest, sig) {
		return ErrBadSignature
	}
	return nil
}

// CompareVersions compares dotted numeric versions ("0.2.10" > "0.2.9"). A leading "v" and any
// suffix after "-" are ignored. Returns -1, 0, or 1.
func CompareVersions(a, b string) int {
	pa, pb := parseVersion(a), parseVersion(b)
	for i := 0; i < len(pa) || i < len(pb); i++ {
		var x, y int
		if i < len(pa) {
			x = pa[i]
		}
		if i < len(pb) {
			y = pb[i]
		}
		if x != y {
			if x < y {
				return -1
			}
			return 1
		}
	}
	return 0
}

func parseVersion(v string) []int {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	if i := strings.IndexByte(v, '-'); i >= 0 {
		v = v[:i]
	}
	var out []int
	for _, part := range strings.Split(v, ".") {
		n, err := strconv.Atoi(part)
		if err != nil {
			n = 0
		}
		out = append(out, n)
	}
	return out
}

// ManifestURLs returns where the signed manifest and its signature live.
func ManifestURLs(repo string) (manifestURL, sigURL string) {
	base := githubBase + "/" + repo + "/releases/latest/download/"
	return base + "manifest.json", base + "manifest.sig"
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Minute}
}

func fetch(ctx context.Context, url string, limit int64) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: HTTP %d", url, resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, limit))
}

// LoadManifest downloads and verifies the signed manifest for this build's repo.
func LoadManifest(ctx context.Context) (*Manifest, error) {
	if Repo == "" || PublicKeyHex == "" {
		return nil, ErrNotConfigured
	}
	mURL, sURL := ManifestURLs(Repo)
	raw, err := fetch(ctx, mURL, maxManifestBytes)
	if err != nil {
		return nil, err
	}
	sig, err := fetch(ctx, sURL, 1024)
	if err != nil {
		return nil, err
	}
	if err := VerifyManifest(raw, string(sig), PublicKeyHex); err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("manifest is not valid JSON: %w", err)
	}
	return &m, nil
}

// CheckStatus reports whether a newer signed release exists for this platform.
func CheckStatus(ctx context.Context, current string) Status {
	st := Status{Current: current}
	if Repo == "" || PublicKeyHex == "" {
		st.Reason = ErrNotConfigured.Error()
		return st
	}
	st.Supported = true
	m, err := LoadManifest(ctx)
	if err != nil {
		st.Reason = err.Error()
		return st
	}
	st.Latest = m.Version
	st.Available = CompareVersions(m.Version, current) > 0
	return st
}

// Apply downloads the release for this platform, verifies the checksum, and swaps it in for
// exePath. The previous binary is kept as exePath+".prev". The caller must restart the process.
func Apply(ctx context.Context, exePath string) (string, error) {
	m, err := LoadManifest(ctx)
	if err != nil {
		return "", err
	}
	key := runtime.GOOS + "/" + runtime.GOARCH
	asset, ok := m.Assets[key]
	if !ok {
		return "", ErrNoAsset
	}
	wantPrefix := githubBase + "/" + Repo + "/releases/download/"
	if !strings.HasPrefix(asset.URL, wantPrefix) {
		return "", fmt.Errorf("asset URL %q is outside this repo's releases", asset.URL)
	}

	newPath := exePath + ".new"
	if err := downloadVerified(ctx, asset, newPath); err != nil {
		_ = os.Remove(newPath)
		return "", err
	}

	prevPath := exePath + ".prev"
	_ = os.Remove(prevPath)
	if err := os.Link(exePath, prevPath); err != nil {
		// Hard links can fail on some filesystems; fall back to a copy.
		if err := copyFile(exePath, prevPath); err != nil {
			_ = os.Remove(newPath)
			return "", fmt.Errorf("keep previous binary: %w", err)
		}
	}
	if err := os.Rename(newPath, exePath); err != nil {
		return "", fmt.Errorf("swap binary: %w", err)
	}
	return m.Version, nil
}

func downloadVerified(ctx context.Context, asset Asset, dst string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.URL, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: HTTP %d", asset.URL, resp.StatusCode)
	}

	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(f, h), io.LimitReader(resp.Body, maxBinaryBytes+1))
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return err
	}
	if n > maxBinaryBytes {
		return errors.New("download exceeds size limit")
	}
	if !strings.EqualFold(hex.EncodeToString(h.Sum(nil)), asset.SHA256) {
		return ErrBadChecksum
	}
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, bytes.Clone(data), 0755)
}
