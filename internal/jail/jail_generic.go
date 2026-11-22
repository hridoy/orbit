//go:build !linux

package jail

import (
	"io"
	"os"
)

// OpenUnderJail opens a file using userspace canonical path & symlink verification.
// Used as portable fallback for macOS / BSD / Windows development.
func OpenUnderJail(root, relPath string) (io.ReadCloser, int64, error) {
	safePath, err := SafeResolve(root, relPath)
	if err != nil {
		return nil, 0, err
	}

	f, err := os.Open(safePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0, ErrFileNotFound
		}
		return nil, 0, err
	}

	st, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, 0, err
	}

	return f, st.Size(), nil
}
