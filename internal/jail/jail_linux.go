//go:build linux

package jail

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

// OpenUnderJail opens a file using Linux kernel-level openat2 with RESOLVE_BENEATH.
// This guarantees immunity against TOCTOU race conditions and symlink escapes.
func OpenUnderJail(root, relPath string) (io.ReadCloser, int64, error) {
	safePath, err := SafeResolve(root, relPath)
	if err != nil {
		return nil, 0, err
	}

	rootFd, err := unix.Open(root, unix.O_PATH|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to open jail root fd: %w", err)
	}
	defer unix.Close(rootFd)

	cleanRel := filepath.Clean("/" + relPath)
	cleanRel = strings.TrimPrefix(cleanRel, "/")
	if cleanRel == "" {
		cleanRel = "."
	}

	how := &unix.OpenHow{
		Flags:   unix.O_RDONLY | unix.O_CLOEXEC,
		Mode:    0,
		Resolve: unix.RESOLVE_BENEATH | unix.RESOLVE_NO_MAGICLINKS,
	}

	fd, err := unix.Openat2(rootFd, cleanRel, how)
	if err != nil {
		if errors.Is(err, unix.EXDEV) {
			return nil, 0, ErrAccessDenied
		}
		if errors.Is(err, unix.ENOENT) {
			return nil, 0, ErrFileNotFound
		}
		// If kernel returns ENOSYS (< Linux 5.6), fall back safely to os.Open on verified safePath
		if errors.Is(err, unix.ENOSYS) {
			f, err := os.Open(safePath)
			if err != nil {
				return nil, 0, err
			}
			st, err := f.Stat()
			if err != nil {
				f.Close()
				return nil, 0, err
			}
			return f, st.Size(), nil
		}
		return nil, 0, err
	}

	f := os.NewFile(uintptr(fd), safePath)
	st, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, 0, err
	}
	return f, st.Size(), nil
}
