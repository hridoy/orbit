package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/security"
)

func setupTestServer(t *testing.T) (*Server, *security.TokenManager, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "orbit_server_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	cfg, err := config.LoadOrInit(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to init config: %v", err)
	}

	database, err := db.Open(filepath.Join(cfg.DataDir, "test.db"))
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to open test db: %v", err)
	}

	vault, err := security.NewVault(cfg.MasterKey)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to init vault: %v", err)
	}

	tokenMgr, err := security.NewTokenManager(cfg.MasterKey)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to init token mgr: %v", err)
	}

	srv := NewServer(cfg, database, vault, tokenMgr)
	return srv, tokenMgr, tmpDir
}

func TestHealthzEndpoint(t *testing.T) {
	srv, _, tmpDir := setupTestServer(t)
	defer os.RemoveAll(tmpDir)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp["status"] != "healthy" {
		t.Errorf("expected status healthy, got %v", resp["status"])
	}
}

func TestAuthProtection(t *testing.T) {
	srv, tokenMgr, tmpDir := setupTestServer(t)
	defer os.RemoveAll(tmpDir)

	// 1. Unauthenticated request to protected endpoint
	req := httptest.NewRequest("GET", "/api/v1/node", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}

	// 2. Authenticated request with Bearer token
	token, err := tokenMgr.IssueToken("test_tok", "admin", 1*time.Hour)
	if err != nil {
		t.Fatalf("IssueToken failed: %v", err)
	}

	reqAuth := httptest.NewRequest("GET", "/api/v1/node", nil)
	reqAuth.Header.Set("Authorization", "Bearer "+token)
	recAuth := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recAuth, reqAuth)

	if recAuth.Code != http.StatusOK {
		t.Fatalf("expected 200 OK with valid token, got %d: %s", recAuth.Code, recAuth.Body.String())
	}
}

func TestJailHTTPFileOperations(t *testing.T) {
	srv, tokenMgr, tmpDir := setupTestServer(t)
	defer os.RemoveAll(tmpDir)

	token, _ := tokenMgr.IssueToken("test_tok", "admin", 1*time.Hour)
	authHeader := "Bearer " + token

	// 1. Upload file
	content := []byte("orbit sovereign app binary or source file")
	uploadReq := httptest.NewRequest("POST", "/api/v1/jail/file?scope=app&id=myapp&path=main.py", bytes.NewReader(content))
	uploadReq.Header.Set("Authorization", authHeader)
	uploadRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(uploadRec, uploadReq)

	if uploadRec.Code != http.StatusCreated {
		t.Fatalf("upload failed with code %d: %s", uploadRec.Code, uploadRec.Body.String())
	}

	// 2. Read file back
	readReq := httptest.NewRequest("GET", "/api/v1/jail/file?scope=app&id=myapp&path=main.py", nil)
	readReq.Header.Set("Authorization", authHeader)
	readRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(readRec, readReq)

	if readRec.Code != http.StatusOK {
		t.Fatalf("read failed with code %d: %s", readRec.Code, readRec.Body.String())
	}
	if !bytes.Equal(readRec.Body.Bytes(), content) {
		t.Fatalf("content mismatch: got '%s'", readRec.Body.String())
	}

	// 3. List directory
	listReq := httptest.NewRequest("GET", "/api/v1/jail/list?scope=app&id=myapp&path=", nil)
	listReq.Header.Set("Authorization", authHeader)
	listRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("list failed with code %d: %s", listRec.Code, listRec.Body.String())
	}

	// 4. Traversal attack over HTTP should return 403 Forbidden
	attackReq := httptest.NewRequest("GET", "/api/v1/jail/file?scope=app&id=myapp&path=../../etc/passwd", nil)
	attackReq.Header.Set("Authorization", authHeader)
	attackRec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(attackRec, attackReq)

	if attackRec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for path traversal attack, got %d", attackRec.Code)
	}
}
