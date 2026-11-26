package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strconv"
	"testing"

	"orbit/internal/db"
)

func TestSingleLiveAppFallback(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("from-app"))
	}))
	defer backend.Close()
	port, _ := strconv.Atoi(mustPort(t, backend.URL))

	database, err := db.Open(filepath.Join(t.TempDir(), "orbit.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	apps := db.NewAppStore(database)
	h := singleLiveAppFallback(apps)

	// No live app: 404.
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("no apps: want 404, got %d", rec.Code)
	}

	// One live app: forwarded.
	v := "v1"
	nodeID := "node_test"
	if _, err := database.Exec(`INSERT INTO nodes (node_id, hostname, public_ip, agent_version) VALUES (?, 'h', '127.0.0.1', 'test')`, nodeID); err != nil {
		t.Fatal(err)
	}
	if err := apps.CreateApp(&db.App{AppID: "a1", Name: "a1", NodeID: nodeID, LinuxUser: "", ActivePort: port, IdlePort: port + 1,
		BinaryEntrypoint: "./server", RAMLimitMB: 64, CPUQuotaPct: 100, PIDsLimit: 16, HealthCheckPath: "/healthz",
		HealthTimeoutSec: 5, Status: "live"}); err != nil {
		t.Fatal(err)
	}
	if err := apps.SwapPorts("a1", port, port+1, v); err != nil {
		t.Fatal(err)
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Body.String() != "from-app" {
		t.Fatalf("single app: want from-app, got %q (%d)", rec.Body.String(), rec.Code)
	}
}

func mustPort(t *testing.T, raw string) string {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return u.Port()
}
