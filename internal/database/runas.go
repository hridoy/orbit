package database

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// dbServiceUser runs database servers when the agent itself runs as root. PostgreSQL and
// MySQL/MariaDB refuse to start as root.
const dbServiceUser = "orbitdb"

// dbCredential returns the unprivileged account to run database processes as. It is nil when
// the agent is not root (for example on a developer's Mac), so processes keep the current user.
func dbCredential() (*syscall.Credential, error) {
	if os.Geteuid() != 0 {
		return nil, nil
	}
	u, err := user.Lookup(dbServiceUser)
	if err != nil {
		if out, cerr := exec.Command("useradd", "--system", "--no-create-home",
			"--shell", "/usr/sbin/nologin", dbServiceUser).CombinedOutput(); cerr != nil {
			return nil, fmt.Errorf("create %s user: %v: %s", dbServiceUser, cerr, strings.TrimSpace(string(out)))
		}
		if u, err = user.Lookup(dbServiceUser); err != nil {
			return nil, err
		}
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return nil, err
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return nil, err
	}
	return &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}, nil
}

// applyDBUser runs cmd in its own process group, as dbServiceUser when the agent is root.
func applyDBUser(cmd *exec.Cmd) error {
	cred, err := dbCredential()
	if err != nil {
		return err
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Credential: cred}
	return nil
}

// chownDBTree hands a database directory to dbServiceUser so its server can write to it.
// It does nothing when the agent is not root.
func chownDBTree(dir string) error {
	cred, err := dbCredential()
	if err != nil || cred == nil {
		return err
	}
	return filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return os.Lchown(p, int(cred.Uid), int(cred.Gid))
	})
}
