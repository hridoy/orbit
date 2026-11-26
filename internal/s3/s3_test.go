package s3

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/security"
)

func setupS3TestEnv(t *testing.T) (*config.Config, *S3Engine, *security.TokenManager, string) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "orbit_s3_test_*")
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

	storageStore := db.NewStorageStore(database)
	tokenMgr, _ := security.NewTokenManager(cfg.MasterKey)

	engine := NewS3Engine(cfg, storageStore, tokenMgr)
	return cfg, engine, tokenMgr, tmpDir
}

func TestS3ObjectLifecycleAndRange(t *testing.T) {
	_, engine, tokenMgr, tmpDir := setupS3TestEnv(t)
	defer os.RemoveAll(tmpDir)

	mux := http.NewServeMux()
	engine.RegisterRoutes(mux)

	token, _ := tokenMgr.IssueToken("test_admin", "admin", 1*time.Hour)
	authHeader := "Bearer " + token

	// 1. Create Bucket
	reqCreate := httptest.NewRequest("PUT", "/s3/assets", nil)
	reqCreate.Header.Set("Authorization", authHeader)
	recCreate := httptest.NewRecorder()
	mux.ServeHTTP(recCreate, reqCreate)

	if recCreate.Code != http.StatusOK {
		t.Fatalf("create bucket failed: %d, body: %s", recCreate.Code, recCreate.Body.String())
	}

	// 2. PutObject
	payload := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ" // 36 bytes
	reqPut := httptest.NewRequest("PUT", "/s3/assets/data/sample.txt", strings.NewReader(payload))
	reqPut.Header.Set("Authorization", authHeader)
	reqPut.Header.Set("Content-Type", "text/plain")
	recPut := httptest.NewRecorder()
	mux.ServeHTTP(recPut, reqPut)

	if recPut.Code != http.StatusOK {
		t.Fatalf("put object failed: %d, body: %s", recPut.Code, recPut.Body.String())
	}
	etag := recPut.Header().Get("ETag")
	if etag == "" {
		t.Fatal("missing ETag header in PutObject response")
	}

	// 3. GetObject full
	reqGet := httptest.NewRequest("GET", "/s3/assets/data/sample.txt", nil)
	reqGet.Header.Set("Authorization", authHeader)
	recGet := httptest.NewRecorder()
	mux.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Fatalf("get object failed: %d", recGet.Code)
	}
	if recGet.Body.String() != payload {
		t.Fatalf("content mismatch: got '%s'", recGet.Body.String())
	}

	// 4. GetObject with Range header (bytes=0-9) -> should return "0123456789" (206 Partial Content)
	reqRange := httptest.NewRequest("GET", "/s3/assets/data/sample.txt", nil)
	reqRange.Header.Set("Authorization", authHeader)
	reqRange.Header.Set("Range", "bytes=0-9")
	recRange := httptest.NewRecorder()
	mux.ServeHTTP(recRange, reqRange)

	if recRange.Code != http.StatusPartialContent {
		t.Fatalf("expected 206 Partial Content, got %d", recRange.Code)
	}
	if recRange.Body.String() != "0123456789" {
		t.Fatalf("expected '0123456789', got '%s'", recRange.Body.String())
	}
	if recRange.Header().Get("Content-Range") != fmt.Sprintf("bytes 0-9/%d", len(payload)) {
		t.Errorf("unexpected Content-Range: %s", recRange.Header().Get("Content-Range"))
	}

	// 5. HeadObject
	reqHead := httptest.NewRequest("HEAD", "/s3/assets/data/sample.txt", nil)
	reqHead.Header.Set("Authorization", authHeader)
	recHead := httptest.NewRecorder()
	mux.ServeHTTP(recHead, reqHead)

	if recHead.Code != http.StatusOK {
		t.Fatalf("head object failed: %d", recHead.Code)
	}

	// 6. DeleteObject
	reqDel := httptest.NewRequest("DELETE", "/s3/assets/data/sample.txt", nil)
	reqDel.Header.Set("Authorization", authHeader)
	recDel := httptest.NewRecorder()
	mux.ServeHTTP(recDel, reqDel)

	if recDel.Code != http.StatusNoContent {
		t.Fatalf("delete object failed: %d", recDel.Code)
	}
}

func TestS3MultipartUpload(t *testing.T) {
	_, engine, tokenMgr, tmpDir := setupS3TestEnv(t)
	defer os.RemoveAll(tmpDir)

	mux := http.NewServeMux()
	engine.RegisterRoutes(mux)

	token, _ := tokenMgr.IssueToken("test_admin", "admin", 1*time.Hour)
	authHeader := "Bearer " + token

	// Create Bucket
	reqCreate := httptest.NewRequest("PUT", "/s3/large-bucket", nil)
	reqCreate.Header.Set("Authorization", authHeader)
	mux.ServeHTTP(httptest.NewRecorder(), reqCreate)

	// 1. Initiate Multipart Upload
	reqInit := httptest.NewRequest("POST", "/s3/large-bucket/large.bin?uploads", nil)
	reqInit.Header.Set("Authorization", authHeader)
	recInit := httptest.NewRecorder()
	mux.ServeHTTP(recInit, reqInit)

	if recInit.Code != http.StatusOK {
		t.Fatalf("initiate multipart failed: %d, body: %s", recInit.Code, recInit.Body.String())
	}

	body := recInit.Body.String()
	uploadIDStart := strings.Index(body, "<UploadId>") + len("<UploadId>")
	uploadIDEnd := strings.Index(body, "</UploadId>")
	uploadID := body[uploadIDStart:uploadIDEnd]

	// 2. Upload Part 1
	chunk1 := bytes.Repeat([]byte("A"), 5<<20) // S3 minimum part size is 5 MiB
	reqPart1 := httptest.NewRequest("PUT", fmt.Sprintf("/s3/large-bucket/large.bin?uploadId=%s&partNumber=1", uploadID), bytes.NewReader(chunk1))
	reqPart1.Header.Set("Authorization", authHeader)
	recPart1 := httptest.NewRecorder()
	mux.ServeHTTP(recPart1, reqPart1)
	if recPart1.Code != http.StatusOK {
		t.Fatalf("upload part 1 failed: %d", recPart1.Code)
	}

	// 3. Upload Part 2
	chunk2 := bytes.Repeat([]byte("B"), 1024)
	reqPart2 := httptest.NewRequest("PUT", fmt.Sprintf("/s3/large-bucket/large.bin?uploadId=%s&partNumber=2", uploadID), bytes.NewReader(chunk2))
	reqPart2.Header.Set("Authorization", authHeader)
	recPart2 := httptest.NewRecorder()
	mux.ServeHTTP(recPart2, reqPart2)
	if recPart2.Code != http.StatusOK {
		t.Fatalf("upload part 2 failed: %d", recPart2.Code)
	}

	// 4. Complete Multipart Upload
	reqComp := httptest.NewRequest("POST", fmt.Sprintf("/s3/large-bucket/large.bin?uploadId=%s", uploadID), nil)
	reqComp.Header.Set("Authorization", authHeader)
	recComp := httptest.NewRecorder()
	mux.ServeHTTP(recComp, reqComp)
	if recComp.Code != http.StatusOK {
		t.Fatalf("complete multipart failed: %d, body: %s", recComp.Code, recComp.Body.String())
	}

	// 5. Verify the combined object
	reqGet := httptest.NewRequest("GET", "/s3/large-bucket/large.bin", nil)
	reqGet.Header.Set("Authorization", authHeader)
	recGet := httptest.NewRecorder()
	mux.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Fatalf("get completed object failed: %d", recGet.Code)
	}

	combined, _ := io.ReadAll(recGet.Body)
	if len(combined) != len(chunk1)+len(chunk2) {
		t.Fatalf("expected %d bytes, got %d", len(chunk1)+len(chunk2), len(combined))
	}
	if !bytes.Equal(combined[:len(chunk1)], chunk1) || !bytes.Equal(combined[len(chunk1):], chunk2) {
		t.Fatal("assembled content parts do not match original chunks")
	}
}

func TestPresignedURL(t *testing.T) {
	_, engine, _, tmpDir := setupS3TestEnv(t)
	defer os.RemoveAll(tmpDir)

	presignedURL := engine.sigv4.GeneratePresignedURL("GET", "localhost:8787", "assets", "logo.png", 5*time.Minute)
	if !strings.Contains(presignedURL, "X-Amz-Signature=") || !strings.Contains(presignedURL, "X-Amz-Algorithm=AWS4-HMAC-SHA256") {
		t.Fatalf("invalid presigned URL: %s", presignedURL)
	}
}
