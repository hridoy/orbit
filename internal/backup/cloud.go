// Package backup copies control-plane backups off the server through rclone. rclone speaks
// Google Drive, Dropbox, OneDrive, S3-compatible storage, and many more. The operator pastes an
// rclone config in Credentials; it is written to a private file for each run and removed after.
package backup

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RetentionDays is how long copies stay on the remote before they are pruned.
const RetentionDays = 30

// Upload copies the backup directory localDir into remotePath (for example "gdrive:orbit-backups")
// and then prunes remote copies older than RetentionDays. rcloneConf is the full rclone config
// text, and workDir holds the temporary config file (use a private directory).
func Upload(ctx context.Context, rcloneConf, remotePath, localDir, workDir string) error {
	if _, err := exec.LookPath("rclone"); err != nil {
		return errors.New("rclone is not installed; run: apt-get install -y rclone")
	}
	remotePath = strings.TrimRight(strings.TrimSpace(remotePath), "/")
	if !strings.Contains(remotePath, ":") {
		return errors.New("remote path must look like remote:folder")
	}

	conf, err := os.CreateTemp(workDir, "rclone-*.conf")
	if err != nil {
		return fmt.Errorf("create rclone config: %w", err)
	}
	defer os.Remove(conf.Name())
	if err := conf.Chmod(0600); err != nil {
		conf.Close()
		return err
	}
	if _, err := conf.WriteString(rcloneConf); err != nil {
		conf.Close()
		return err
	}
	if err := conf.Close(); err != nil {
		return err
	}

	dest := remotePath + "/" + filepath.Base(localDir)
	if err := rclone(ctx, conf.Name(), "copy", localDir, dest); err != nil {
		return err
	}
	if err := rclone(ctx, conf.Name(), "delete", remotePath, "--min-age", fmt.Sprintf("%dd", RetentionDays)); err != nil {
		return err
	}
	return rclone(ctx, conf.Name(), "rmdirs", remotePath, "--leave-root")
}

func rclone(ctx context.Context, confPath string, args ...string) error {
	full := append([]string{"--config", confPath, "--retries", "3"}, args...)
	out, err := exec.CommandContext(ctx, "rclone", full...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("rclone %s failed: %v: %s", args[0], err, strings.TrimSpace(string(out)))
	}
	return nil
}
