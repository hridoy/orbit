package database

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMySQLInitSQLCreatesAdminForLocalAndTCP(t *testing.T) {
	sql := mysqlInitSQL("orbit_admin", "p'w\\x")
	for _, want := range []string{
		"CREATE USER IF NOT EXISTS 'orbit_admin'@'localhost' IDENTIFIED BY 'p''w\\\\x';",
		"CREATE USER IF NOT EXISTS 'orbit_admin'@'%' IDENTIFIED BY 'p''w\\\\x';",
		"GRANT ALL PRIVILEGES ON *.* TO 'orbit_admin'@'%' WITH GRANT OPTION;",
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("init SQL missing %q\n%s", want, sql)
		}
	}
}

func TestMySQLInitFileIsPrivateAndPassedToStart(t *testing.T) {
	dir := t.TempDir()
	if args := mysqlInitArgs(dir); args != nil {
		t.Fatalf("no init file yet, got %v", args)
	}
	if err := writeMySQLInitFile(dir, "orbit_admin", "secret123"); err != nil {
		t.Fatal(err)
	}
	st, _ := os.Stat(filepath.Join(dir, mysqlInitFileName))
	if st.Mode().Perm() != 0600 {
		t.Fatalf("init file mode %v, want 0600", st.Mode().Perm())
	}
	args := mysqlInitArgs(dir)
	if len(args) != 1 || !strings.HasPrefix(args[0], "--init-file=") {
		t.Fatalf("unexpected args %v", args)
	}
}
