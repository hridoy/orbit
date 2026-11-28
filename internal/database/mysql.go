package database

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"orbit/internal/compute"
)

type MySQLEngine struct {
	mysqldBin   string
	mariadbdBin string
	initDbBin   string
}

func NewMySQLEngine() *MySQLEngine {
	// Oracle MySQL 8 is looked up only at the paths where it is installed. The mysqld name on
	// Ubuntu is usually MariaDB's, which must not be used for the MySQL engine.
	mysqldBin := firstExisting(
		"/usr/local/mysql/bin/mysqld",
		"/opt/homebrew/opt/mysql/bin/mysqld",
	)

	mariadbdBin, _ := findBinary("mariadbd", []string{
		"/usr/sbin/mariadbd",
		"/usr/bin/mariadbd",
		"/usr/local/bin/mariadbd",
		"/opt/homebrew/bin/mariadbd",
		"/opt/homebrew/opt/mariadb/bin/mariadbd",
	})

	initDbBin, _ := findBinary("mariadb-install-db", []string{
		"/usr/bin/mariadb-install-db",
		"/usr/bin/mysql_install_db",
		"/usr/local/bin/mysql_install_db",
		"/opt/homebrew/bin/mariadb-install-db",
	})

	return &MySQLEngine{
		mysqldBin:   mysqldBin,
		mariadbdBin: mariadbdBin,
		initDbBin:   initDbBin,
	}
}

func (e *MySQLEngine) IsAvailable(isMariaDB bool) bool {
	if isMariaDB {
		return e.mariadbdBin != "" || e.mysqldBin != ""
	}
	return e.mysqldBin != ""
}

func (e *MySQLEngine) Provision(
	ctx context.Context,
	baseDir string,
	adminUser, adminPassword string,
	port int,
	profile RAMProfile,
	isMariaDB bool,
) error {
	dataDir := filepath.Join(baseDir, "data")
	runDir := filepath.Join(baseDir, "run")
	confDir := filepath.Join(baseDir, "conf")

	for _, d := range []string{dataDir, runDir, confDir} {
		if err := os.MkdirAll(d, 0700); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}

	socketFile := filepath.Join(runDir, "mysql.sock")
	if isMariaDB {
		socketFile = filepath.Join(runDir, "mariadb.sock")
	}

	confPath := filepath.Join(confDir, "my.cnf")
	var confContent string
	if isMariaDB {
		confContent = GenerateMariaDBConf(port, socketFile, dataDir, profile)
	} else {
		confContent = GenerateMySQLConf(port, socketFile, dataDir, profile)
	}

	if err := os.WriteFile(confPath, []byte(confContent), 0600); err != nil {
		return fmt.Errorf("failed to write my.cnf: %w", err)
	}
	if adminUser != "" && adminPassword != "" {
		if err := writeMySQLInitFile(confDir, adminUser, adminPassword); err != nil {
			return fmt.Errorf("failed to write admin bootstrap: %w", err)
		}
	}

	// Initialize data directory if host binary is present and data dir is empty
	entries, _ := os.ReadDir(dataDir)
	if len(entries) == 0 {
		_ = chownDBTree(baseDir)
		// --no-defaults keeps host config (such as /etc/mysql) from breaking the first start.
		// Failures are returned: an uninitialized data dir starts with no system tables and
		// every later start fails.
		var cmd *exec.Cmd
		if isMariaDB && e.initDbBin != "" {
			base := filepath.Dir(filepath.Dir(e.mariadbdBin))
			cmd = exec.CommandContext(ctx, e.initDbBin, "--no-defaults", "--auth-root-authentication-method=normal",
				fmt.Sprintf("--basedir=%s", base), fmt.Sprintf("--datadir=%s", dataDir))
		} else if e.mysqldBin != "" {
			cmd = exec.CommandContext(ctx, e.mysqldBin, "--no-defaults", "--initialize-insecure", fmt.Sprintf("--datadir=%s", dataDir))
		}
		if cmd != nil {
			if err := applyDBUser(cmd); err != nil {
				return err
			}
			if out, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("initialize data directory: %w: %s", err, tail(out))
			}
		}
	}

	return nil
}

func (e *MySQLEngine) Start(
	ctx context.Context,
	baseDir string,
	rb *compute.RingBuffer,
	isMariaDB bool,
) (*exec.Cmd, error) {
	bin := e.mysqldBin
	engineName := "mysql"
	if isMariaDB {
		engineName = "mariadb"
		if e.mariadbdBin != "" {
			bin = e.mariadbdBin
		}
	}

	if bin == "" {
		return nil, fmt.Errorf("native %s binary not found on host (install with 'apt install mysql-server' or 'apt install mariadb-server')", engineName)
	}

	confPath := filepath.Join(baseDir, "conf", "my.cnf")
	args := append([]string{fmt.Sprintf("--defaults-file=%s", confPath)}, mysqlInitArgs(filepath.Join(baseDir, "conf"))...)
	cmd := exec.CommandContext(ctx, bin, args...)

	cmd.Stdout = rb
	cmd.Stderr = rb
	if err := applyDBUser(cmd); err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start %s: %w", engineName, err)
	}

	rb.Append(fmt.Sprintf("[%s] Process started (PID: %d)", engineName, cmd.Process.Pid))
	return cmd, nil
}
