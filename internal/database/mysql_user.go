package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// mysqlInitFileName is the bootstrap script mysqld runs at every start (--init-file). It creates
// the admin user the dashboard hands out. Provisioning used to skip this step, so the user never
// existed and every login returned "Access denied".
const mysqlInitFileName = "init.sql"

// sqlString quotes a value for a single-quoted SQL literal.
func sqlString(v string) string {
	v = strings.ReplaceAll(v, `\`, `\\`)
	return "'" + strings.ReplaceAll(v, "'", "''") + "'"
}

// mysqlInitSQL returns the statements that create or update the admin user for local and TCP
// logins. The forwarder connects over 127.0.0.1, which MySQL matches to the '%' host.
func mysqlInitSQL(user, password string) string {
	u, p := sqlString(user), sqlString(password)
	return fmt.Sprintf(`CREATE USER IF NOT EXISTS %[1]s@'localhost' IDENTIFIED BY %[2]s;
CREATE USER IF NOT EXISTS %[1]s@'%%' IDENTIFIED BY %[2]s;
ALTER USER %[1]s@'localhost' IDENTIFIED BY %[2]s;
ALTER USER %[1]s@'%%' IDENTIFIED BY %[2]s;
GRANT ALL PRIVILEGES ON *.* TO %[1]s@'localhost' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO %[1]s@'%%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
`, u, p)
}

// writeMySQLInitFile stores the bootstrap script in the instance's conf directory, readable only
// by the database service user.
func writeMySQLInitFile(confDir, user, password string) error {
	return os.WriteFile(filepath.Join(confDir, mysqlInitFileName), []byte(mysqlInitSQL(user, password)), 0600)
}

// mysqlInitArgs returns the --init-file argument when the bootstrap script exists.
func mysqlInitArgs(confDir string) []string {
	p := filepath.Join(confDir, mysqlInitFileName)
	if _, err := os.Stat(p); err != nil {
		return nil
	}
	return []string{"--init-file=" + p}
}
