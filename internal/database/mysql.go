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
	mysqldBin, _ := findBinary("mysqld", []string{
		"/usr/sbin/mysqld",
		"/usr/bin/mysqld",
		"/usr/local/bin/mysqld",
		"/usr/local/mysql/bin/mysqld",
		"/opt/homebrew/bin/mysqld",
		"/opt/homebrew/opt/mysql/bin/mysqld",
	})

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

	// Initialize data directory if host binary is present and data dir is empty
	entries, _ := os.ReadDir(dataDir)
	if len(entries) == 0 {
		_ = chownDBTree(baseDir)
		if isMariaDB && e.initDbBin != "" {
			cmd := exec.CommandContext(ctx, e.initDbBin, fmt.Sprintf("--datadir=%s", dataDir), "--auth-root-authentication-method=normal")
			_ = applyDBUser(cmd)
			_ = cmd.Run()
		} else if e.mysqldBin != "" {
			cmd := exec.CommandContext(ctx, e.mysqldBin, "--initialize-insecure", fmt.Sprintf("--datadir=%s", dataDir))
			_ = applyDBUser(cmd)
			_ = cmd.Run()
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
	cmd := exec.CommandContext(ctx, bin, fmt.Sprintf("--defaults-file=%s", confPath))

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
