package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/acme/autocert"
	"orbit/internal/config"
	"orbit/internal/db"
)

func setupProxyTestEnv(t *testing.T) (*config.Config, *db.DB, *db.AppStore, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "orbit_proxy_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	cfg, err := config.LoadOrInit(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("config init failed: %v", err)
	}

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("db open failed: %v", err)
	}

	// Upsert node for foreign key constraint
	nodeStore := db.NewNodeStore(database)
	_ = nodeStore.UpsertNode(&db.Node{
		NodeID:       cfg.NodeID,
		Hostname:     "localhost",
		PublicIP:     "127.0.0.1",
		Role:         "standalone",
		AgentVersion: "0.1.0",
		Status:       "online",
	})

	appStore := db.NewAppStore(database)
	return cfg, database, appStore, tmpDir
}

func TestBandwidthMeter(t *testing.T) {
	meter := NewBandwidthMeter()

	meter.RecordRequestStart("app_1")
	meter.RecordBytesIn("app_1", 100)
	meter.RecordBytesOut("app_1", 500)
	meter.RecordRequestEnd("app_1")

	stats := meter.GetAppStats("app_1")
	if stats.BytesIn != 100 {
		t.Errorf("expected 100 bytes in, got %d", stats.BytesIn)
	}
	if stats.BytesOut != 500 {
		t.Errorf("expected 500 bytes out, got %d", stats.BytesOut)
	}
	if stats.TotalRequests != 1 {
		t.Errorf("expected 1 total request, got %d", stats.TotalRequests)
	}
}

func TestReverseProxyRouting(t *testing.T) {
	_, _, appStore, tmpDir := setupProxyTestEnv(t)
	defer os.RemoveAll(tmpDir)

	// 1. Mock upstream backend application
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Upstream-Handled", "true")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello from backend app"))
	}))
	defer backend.Close()

	// Parse backend port
	var backendPort int
	_, _ = fmt.Sscanf(strings.TrimPrefix(backend.URL, "http://127.0.0.1:"), "%d", &backendPort)

	meter := NewBandwidthMeter()
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("control plane fallback"))
	})

	rp := NewReverseProxy(appStore, meter, fallback)
	rp.SetRoute("api.example.com", "app_backend", backendPort)

	// 2. Request for routed domain
	req := httptest.NewRequest("GET", "http://api.example.com/items", nil)
	rec := httptest.NewRecorder()
	rp.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from reverse proxy, got %d", rec.Code)
	}
	if rec.Header().Get("X-Upstream-Handled") != "true" {
		t.Error("expected upstream header X-Upstream-Handled: true")
	}
	if rec.Body.String() != "hello from backend app" {
		t.Errorf("unexpected body: %s", rec.Body.String())
	}

	// 3. Request for unrouted domain falls back to control plane
	reqFallback := httptest.NewRequest("GET", "http://unrouted.example.com/healthz", nil)
	recFallback := httptest.NewRecorder()
	rp.ServeHTTP(recFallback, reqFallback)

	if recFallback.Body.String() != "control plane fallback" {
		t.Errorf("expected fallback response, got: %s", recFallback.Body.String())
	}
}

func TestSQLiteCertCache(t *testing.T) {
	_, database, _, tmpDir := setupProxyTestEnv(t)
	defer os.RemoveAll(tmpDir)

	cache := NewSQLiteCertCache(database)
	ctx := context.Background()

	// 1. Get non-existent cert
	_, err := cache.Get(ctx, "example.com")
	if err != autocert.ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss, got %v", err)
	}

	// 2. Put cert
	certData := []byte("CERTIFICATE_PEM_DATA_BLOCK")
	if err := cache.Put(ctx, "example.com", certData); err != nil {
		t.Fatalf("Put cert failed: %v", err)
	}

	// 3. Get cached cert
	retrieved, err := cache.Get(ctx, "example.com")
	if err != nil {
		t.Fatalf("Get cert failed: %v", err)
	}
	if string(retrieved) != string(certData) {
		t.Fatalf("cert data mismatch: got '%s'", string(retrieved))
	}

	// 4. Delete cert
	if err := cache.Delete(ctx, "example.com"); err != nil {
		t.Fatalf("Delete cert failed: %v", err)
	}
	_, err = cache.Get(ctx, "example.com")
	if err != autocert.ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss after delete, got %v", err)
	}
}
