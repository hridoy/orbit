//go:build !linux

package compute

import (
	"os/exec"
	"syscall"

	"orbit/internal/db"
)

func configureSysProcAttr(cmd *exec.Cmd, app *db.App) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}
