package database

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Installed is one database server binary found on this host.
type Installed struct {
	Engine  string `json:"engine"`
	Version string `json:"version"` // "16", "8.4", "6.0", ...
	Bin     string `json:"-"`
}

// nativeVersion marks instances created before versions were chosen. They keep the binary the
// engine picked at startup, so existing data is never handed to a different server version.
const nativeVersion = "native"

// candidate binary locations per engine. Globs let several installed versions appear.
var candidates = map[string][]string{
	"postgres": {"/usr/lib/postgresql/*/bin/postgres"},
	"mysql":    {"/usr/local/mysql*/bin/mysqld"},
	"mariadb":  {"/usr/sbin/mariadbd", "/usr/bin/mariadbd", "/usr/local/mariadb*/bin/mariadbd"},
	"redis":    {"/usr/bin/redis-server", "/usr/local/redis*/bin/redis-server"},
	"mongodb":  {"/usr/bin/mongod", "/usr/local/mongodb*/bin/mongod"},
	"libsql":   {"/usr/bin/sqld", "/usr/local/bin/sqld"},
}

var versionRe = regexp.MustCompile(`(\d+)\.(\d+)`)

// discoverInstalled finds every server binary and reads its version from --version.
func discoverInstalled() []Installed {
	var out []Installed
	for engine, patterns := range candidates {
		seen := map[string]bool{}
		for _, pat := range patterns {
			matches, _ := filepath.Glob(pat)
			for _, bin := range matches {
				info, err := os.Stat(bin)
				if err != nil || info.IsDir() || seen[bin] {
					continue
				}
				seen[bin] = true
				v := readVersion(engine, bin)
				if v == "" {
					continue
				}
				out = append(out, Installed{Engine: engine, Version: v, Bin: bin})
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return versionLess(out[j].Version, out[i].Version) })
	return out
}

func readVersion(engine, bin string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	raw, _ := exec.CommandContext(ctx, bin, "--version").CombinedOutput()
	m := versionRe.FindStringSubmatch(string(raw))
	if m == nil {
		return ""
	}
	if engine == "postgres" {
		return m[1] // PostgreSQL majors are the version users pick: 14, 16
	}
	return m[1] + "." + m[2]
}

func versionLess(a, b string) bool {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(pa) && i < len(pb); i++ {
		var x, y int
		_, _ = fmt.Sscanf(pa[i], "%d", &x)
		_, _ = fmt.Sscanf(pb[i], "%d", &y)
		if x != y {
			return x < y
		}
	}
	return len(pa) < len(pb)
}

// Versions lists the versions installed for each engine, newest first.
func (m *DatabaseManager) Versions() map[string][]string {
	out := map[string][]string{}
	for _, inst := range m.installed {
		out[inst.Engine] = append(out[inst.Engine], inst.Version)
	}
	return out
}

// pick returns the binary for engine and version. An empty version or "native" means the
// engine's startup default.
func (m *DatabaseManager) pick(engine, version string) (Installed, error) {
	if version == "" || version == nativeVersion {
		return Installed{Engine: engine, Version: nativeVersion, Bin: m.defaultBin(engine)}, nil
	}
	for _, inst := range m.installed {
		if inst.Engine == engine && inst.Version == version {
			return inst, nil
		}
	}
	return Installed{}, fmt.Errorf("%s %s is not installed on this server", engine, version)
}

func (m *DatabaseManager) defaultBin(engine string) string {
	switch engine {
	case "postgres":
		return m.postgres.postgresBin
	case "mysql":
		return m.mysql.mysqldBin
	case "mariadb":
		return m.mysql.mariadbdBin
	case "redis":
		return m.redis.redisBin
	case "mongodb":
		return m.mongodb.mongodBin
	case "libsql":
		return m.libsql.sqldBin
	}
	return ""
}

// postgresFor returns a copy of the Postgres engine that runs the chosen server version.
func (m *DatabaseManager) postgresFor(version string) (PostgresEngine, error) {
	inst, err := m.pick("postgres", version)
	if err != nil {
		return PostgresEngine{}, err
	}
	pe := *m.postgres
	pe.postgresBin = inst.Bin
	pe.initdbBin = filepath.Join(filepath.Dir(inst.Bin), "initdb")
	if _, err := os.Stat(pe.initdbBin); err != nil {
		pe.initdbBin = m.postgres.initdbBin
	}
	return pe, nil
}

// mysqlFor returns a copy of the MySQL or MariaDB engine that runs the chosen version.
func (m *DatabaseManager) mysqlFor(engine, version string) (MySQLEngine, error) {
	inst, err := m.pick(engine, version)
	if err != nil {
		return MySQLEngine{}, err
	}
	me := *m.mysql
	if engine == "mariadb" {
		me.mariadbdBin = inst.Bin
	} else {
		me.mysqldBin = inst.Bin
	}
	return me, nil
}

// redisFor, mongoFor and libsqlFor return the engine copy for the chosen version.
func (m *DatabaseManager) redisFor(version string) (RedisEngine, error) {
	inst, err := m.pick("redis", version)
	if err != nil {
		return RedisEngine{}, err
	}
	re := *m.redis
	re.redisBin = inst.Bin
	return re, nil
}

func (m *DatabaseManager) mongoFor(version string) (MongoEngine, error) {
	inst, err := m.pick("mongodb", version)
	if err != nil {
		return MongoEngine{}, err
	}
	me := *m.mongodb
	me.mongodBin = inst.Bin
	return me, nil
}

func (m *DatabaseManager) libsqlFor(version string) (LibSQLEngine, error) {
	inst, err := m.pick("libsql", version)
	if err != nil {
		return LibSQLEngine{}, err
	}
	le := *m.libsql
	le.sqldBin = inst.Bin
	return le, nil
}
