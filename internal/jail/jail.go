package jail

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orbit/internal/config"
)

var (
	ErrAccessDenied     = errors.New("access denied: path escapes jail boundary")
	ErrInvalidScope     = errors.New("invalid jail scope: must be 'app' or 'storage'")
	ErrTargetNotFound   = errors.New("target scope identifier not found")
	ErrFileNotFound     = errors.New("file not found")
	ErrCannotDeleteRoot = errors.New("cannot delete jail root directory")
)

type ScopeType string

const (
	ScopeApp     ScopeType = "app"
	ScopeStorage ScopeType = "storage"
)

type FileItem struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"` // Relative path inside jail
	Size    int64     `json:"size"`
	IsDir   bool      `json:"is_dir"`
	ModTime time.Time `json:"mod_time"`
	Mode    string    `json:"mode"`
}

// Jail enforces path-traversal proof confinement scoped to an app or bucket directory.
type Jail struct {
	cfg *config.Config
}

func NewJail(cfg *config.Config) *Jail {
	return &Jail{cfg: cfg}
}

// ResolveRoot determines the canonical jail root on disk.
func (j *Jail) ResolveRoot(scope ScopeType, scopeID string) (string, error) {
	if scopeID == "" || strings.Contains(scopeID, "..") || strings.ContainsAny(scopeID, "/\\") {
		return "", ErrAccessDenied
	}

	var baseDir string
	switch scope {
	case ScopeApp:
		baseDir = j.cfg.AppsDir
	case ScopeStorage:
		baseDir = j.cfg.StorageDir
	default:
		return "", ErrInvalidScope
	}

	jailRoot := filepath.Join(baseDir, scopeID)
	// Ensure the directory exists
	if err := os.MkdirAll(jailRoot, 0700); err != nil {
		return "", fmt.Errorf("failed to initialize jail root: %w", err)
	}

	canonicalRoot, err := filepath.EvalSymlinks(jailRoot)
	if err != nil {
		return "", fmt.Errorf("failed to resolve canonical jail root: %w", err)
	}

	return canonicalRoot, nil
}

// SafeResolve verifies that relPath resolves strictly beneath canonicalRoot.
func SafeResolve(canonicalRoot, relPath string) (string, error) {
	if strings.Contains(relPath, "\x00") {
		return "", ErrAccessDenied
	}

	cleanRel := filepath.Clean(relPath)
	if strings.HasPrefix(cleanRel, "/") {
		cleanRel = strings.TrimPrefix(cleanRel, "/")
	}

	targetPath := filepath.Clean(filepath.Join(canonicalRoot, cleanRel))

	// Ensure targetPath is strictly within canonicalRoot
	rel, err := filepath.Rel(canonicalRoot, targetPath)
	if err != nil || strings.HasPrefix(rel, "..") || rel == ".." && cleanRel != "" && cleanRel != "." {
		return "", ErrAccessDenied
	}

	// Symlink evaluation check (if path already exists)
	if evalPath, err := filepath.EvalSymlinks(targetPath); err == nil {
		relSym, err := filepath.Rel(canonicalRoot, evalPath)
		if err != nil || strings.HasPrefix(relSym, "..") {
			return "", ErrAccessDenied
		}
		return evalPath, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	// If the file does not exist yet (for writing), evaluate the parent directory
	parent := filepath.Dir(targetPath)
	if evalParent, err := filepath.EvalSymlinks(parent); err == nil {
		relParent, err := filepath.Rel(canonicalRoot, evalParent)
		if err != nil || strings.HasPrefix(relParent, "..") {
			return "", ErrAccessDenied
		}
	}

	return targetPath, nil
}

// ListDirectory returns the contents of relPath within the jailed scope.
func (j *Jail) ListDirectory(scope ScopeType, scopeID, relPath string) ([]FileItem, error) {
	root, err := j.ResolveRoot(scope, scopeID)
	if err != nil {
		return nil, err
	}

	safePath, err := SafeResolve(root, relPath)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(safePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	items := make([]FileItem, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		childRel := filepath.Clean(filepath.Join(relPath, entry.Name()))
		items = append(items, FileItem{
			Name:    entry.Name(),
			Path:    childRel,
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			ModTime: info.ModTime(),
			Mode:    info.Mode().String(),
		})
	}

	return items, nil
}

// Stat returns metadata for a file or directory.
func (j *Jail) Stat(scope ScopeType, scopeID, relPath string) (*FileItem, error) {
	root, err := j.ResolveRoot(scope, scopeID)
	if err != nil {
		return nil, err
	}

	safePath, err := SafeResolve(root, relPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(safePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return &FileItem{
		Name:    info.Name(),
		Path:    filepath.Clean(relPath),
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		ModTime: info.ModTime(),
		Mode:    info.Mode().String(),
	}, nil
}

// OpenFileRead opens a file safely for streaming read.
func (j *Jail) OpenFileRead(scope ScopeType, scopeID, relPath string) (io.ReadCloser, int64, error) {
	root, err := j.ResolveRoot(scope, scopeID)
	if err != nil {
		return nil, 0, err
	}

	return OpenUnderJail(root, relPath)
}

// WriteFile streams content from r directly into the jailed destination with atomic replace.
func (j *Jail) WriteFile(scope ScopeType, scopeID, relPath string, r io.Reader) (int64, error) {
	root, err := j.ResolveRoot(scope, scopeID)
	if err != nil {
		return 0, err
	}

	safePath, err := SafeResolve(root, relPath)
	if err != nil {
		return 0, err
	}

	// Ensure parent dir exists
	parentDir := filepath.Dir(safePath)
	if err := os.MkdirAll(parentDir, 0700); err != nil {
		return 0, fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to atomic temp file first
	tmpFile, err := os.CreateTemp(parentDir, ".orbit_upload_*")
	if err != nil {
		return 0, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	written, err := io.Copy(tmpFile, r)
	if err != nil {
		return 0, fmt.Errorf("failed to write content: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		return 0, fmt.Errorf("failed to sync temp file: %w", err)
	}
	_ = tmpFile.Close()

	if err := os.Rename(tmpFile.Name(), safePath); err != nil {
		return 0, fmt.Errorf("failed to commit file: %w", err)
	}

	return written, nil
}

// Delete removes a file or directory within the jail.
func (j *Jail) Delete(scope ScopeType, scopeID, relPath string) error {
	root, err := j.ResolveRoot(scope, scopeID)
	if err != nil {
		return err
	}

	safePath, err := SafeResolve(root, relPath)
	if err != nil {
		return err
	}

	if safePath == root {
		return ErrCannotDeleteRoot
	}

	err = os.RemoveAll(safePath)
	if err != nil && os.IsNotExist(err) {
		return ErrFileNotFound
	}
	return err
}
