//go:build linux

package compute

import (
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"

	"orbit/internal/db"
)

func configureSysProcAttr(cmd *exec.Cmd, app *db.App) {
	attr := &syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.SIGKILL, // Child process dies if agent dies
	}

	// If running as root, apply namespaces and user privilege drop
	if os.Geteuid() == 0 {
		// Isolated PID, Mount, UTS (hostname), and IPC namespaces
		attr.Cloneflags = syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC

		// Switch to dedicated non-root user if created
		if app.LinuxUser != "" && app.LinuxUser != "root" {
			if u, err := user.Lookup(app.LinuxUser); err == nil {
				uid, _ := strconv.ParseUint(u.Uid, 10, 32)
				gid, _ := strconv.ParseUint(u.Gid, 10, 32)
				attr.Credential = &syscall.Credential{
					Uid: uint32(uid),
					Gid: uint32(gid),
				}
			}
		}
	}

	cmd.SysProcAttr = attr
}
