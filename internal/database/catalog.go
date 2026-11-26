package database

// Method says how a version is installed on the server.
type Method string

const (
	MethodApt     Method = "apt"     // package from the PostgreSQL apt repository
	MethodTarball Method = "tarball" // official binary archive, checksum checked when published
	MethodSource  Method = "source"  // official source archive, built with make
)

// CatalogEntry is one version the dashboard offers. Versions are installed only when asked.
type CatalogEntry struct {
	Engine  string
	Version string // major.minor, as shown in the dropdown
	Method  Method
	Pkg     string // apt package
	URL     string // archive URL
	MD5URL  string // published MD5 file, when the vendor publishes one
	Bin     string // server binary name inside the install folder
}

// Catalog lists the versions the dashboard can install. URLs are the vendors' own archives.
var Catalog = []CatalogEntry{
	{Engine: "postgres", Version: "14", Method: MethodApt, Pkg: "postgresql-14", Bin: "postgres"},
	{Engine: "postgres", Version: "15", Method: MethodApt, Pkg: "postgresql-15", Bin: "postgres"},
	{Engine: "postgres", Version: "16", Method: MethodApt, Pkg: "postgresql-16", Bin: "postgres"},
	{Engine: "postgres", Version: "17", Method: MethodApt, Pkg: "postgresql-17", Bin: "postgres"},

	{Engine: "mysql", Version: "8.0", Method: MethodTarball, Bin: "mysqld",
		URL: "https://cdn.mysql.com/archives/mysql-8.0/mysql-8.0.40-linux-glibc2.28-x86_64.tar.xz"},
	{Engine: "mysql", Version: "8.4", Method: MethodTarball, Bin: "mysqld",
		URL: "https://cdn.mysql.com/archives/mysql-8.4/mysql-8.4.3-linux-glibc2.28-x86_64.tar.xz"},
	{Engine: "mysql", Version: "9.1", Method: MethodTarball, Bin: "mysqld",
		URL: "https://cdn.mysql.com/archives/mysql-9.1/mysql-9.1.0-linux-glibc2.28-x86_64.tar.xz"},
	{Engine: "mysql", Version: "9.2", Method: MethodTarball, Bin: "mysqld",
		URL: "https://cdn.mysql.com/archives/mysql-9.2/mysql-9.2.0-linux-glibc2.28-x86_64.tar.xz"},

	{Engine: "mariadb", Version: "10.6", Method: MethodApt, Pkg: "mariadb-server", Bin: "mariadbd"},
	{Engine: "mariadb", Version: "10.11", Method: MethodTarball, Bin: "mariadbd",
		URL: "https://archive.mariadb.org/mariadb-10.11.13/bintar-linux-systemd-x86_64/mariadb-10.11.13-linux-systemd-x86_64.tar.gz"},
	{Engine: "mariadb", Version: "11.4", Method: MethodTarball, Bin: "mariadbd",
		URL: "https://archive.mariadb.org/mariadb-11.4.5/bintar-linux-systemd-x86_64/mariadb-11.4.5-linux-systemd-x86_64.tar.gz"},
	{Engine: "mariadb", Version: "11.8", Method: MethodTarball, Bin: "mariadbd",
		URL: "https://archive.mariadb.org/mariadb-11.8.2/bintar-linux-systemd-x86_64/mariadb-11.8.2-linux-systemd-x86_64.tar.gz"},

	{Engine: "mongodb", Version: "6.0", Method: MethodTarball, Bin: "mongod",
		URL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu2204-6.0.19.tgz"},
	{Engine: "mongodb", Version: "7.0", Method: MethodApt, Pkg: "mongodb-org", Bin: "mongod"},
	{Engine: "mongodb", Version: "8.0", Method: MethodTarball, Bin: "mongod",
		URL: "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu2204-8.0.4.tgz"},

	{Engine: "redis", Version: "6.0", Method: MethodApt, Pkg: "redis-server", Bin: "redis-server"},
	{Engine: "redis", Version: "6.2", Method: MethodSource, Bin: "redis-server",
		URL: "https://download.redis.io/releases/redis-6.2.17.tar.gz"},
	{Engine: "redis", Version: "7.2", Method: MethodSource, Bin: "redis-server",
		URL: "https://download.redis.io/releases/redis-7.2.7.tar.gz"},
	{Engine: "redis", Version: "7.4", Method: MethodSource, Bin: "redis-server",
		URL: "https://download.redis.io/releases/redis-7.4.2.tar.gz"},
	{Engine: "redis", Version: "8.0", Method: MethodSource, Bin: "redis-server",
		URL: "https://download.redis.io/releases/redis-8.0.2.tar.gz"},

	{Engine: "libsql", Version: "0.24", Method: MethodTarball, Bin: "sqld",
		URL: "https://github.com/tursodatabase/libsql/releases/download/libsql-server-v0.24.32/libsql-server-x86_64-unknown-linux-gnu.tar.xz"},
}

// catalogFor returns the entry for an engine and version, if it is in the catalog.
func catalogFor(engine, version string) (CatalogEntry, bool) {
	for _, c := range Catalog {
		if c.Engine == engine && c.Version == version {
			return c, true
		}
	}
	return CatalogEntry{}, false
}
