package s3

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// signV4Header signs a request with the SigV4 header scheme (host not signed, test-only).
func signV4Header(req *http.Request, ak, sk string, payloadHash string, when time.Time) {
	amzDate := when.UTC().Format("20060102T150405Z")
	dateStamp := when.UTC().Format("20060102")
	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)

	signedHeaders := "x-amz-content-sha256;x-amz-date"
	canonicalHeaders := "x-amz-content-sha256:" + payloadHash + "\n" + "x-amz-date:" + amzDate + "\n"
	canonicalRequest := strings.Join([]string{
		req.Method,
		req.URL.Path,
		CanonicalQueryString(req.URL.Query()),
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")

	scope := dateStamp + "/us-east-1/s3/aws4_request"
	stringToSign := "AWS4-HMAC-SHA256\n" + amzDate + "\n" + scope + "\n" + sha256Hex([]byte(canonicalRequest))
	sig := hex.EncodeToString(hmacSHA256(DeriveSigningKey(sk, dateStamp, "us-east-1", "s3"), []byte(stringToSign)))

	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential="+ak+"/"+scope+
		", SignedHeaders="+signedHeaders+", Signature="+sig)
}

func shaHex(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func newSecurityEnv(t *testing.T) (*S3Engine, *http.ServeMux, string, string) {
	t.Helper()
	cfg, engine, tokenMgr, tmpDir := setupS3TestEnv(t)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })
	_ = cfg
	mux := http.NewServeMux()
	engine.RegisterRoutes(mux)
	token, _ := tokenMgr.IssueToken("test_admin", "admin", time.Hour)
	return engine, mux, "Bearer " + token, tmpDir
}

func TestSigV4PayloadMustMatchSignedHash(t *testing.T) {
	engine, mux, auth, _ := newSecurityEnv(t)
	mux.ServeHTTP(httptest.NewRecorder(), withAuth(httptest.NewRequest("PUT", "/s3/photos", nil), auth))

	body := []byte("hello world")
	// Signed hash is for different content than the body actually sent.
	req := httptest.NewRequest("PUT", "/s3/photos/a.txt", bytes.NewReader(body))
	signV4Header(req, engine.sigv4.AccessKey, engine.sigv4.SecretKey, shaHex([]byte("other")), time.Now())
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("tampered body: want 400, got %d", rec.Code)
	}

	// Correct hash is accepted.
	req = httptest.NewRequest("PUT", "/s3/photos/a.txt", bytes.NewReader(body))
	signV4Header(req, engine.sigv4.AccessKey, engine.sigv4.SecretKey, shaHex(body), time.Now())
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("valid signed PUT: want 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestSigV4RejectsSkewAndUnsignedWrites(t *testing.T) {
	engine, mux, auth, _ := newSecurityEnv(t)
	mux.ServeHTTP(httptest.NewRecorder(), withAuth(httptest.NewRequest("PUT", "/s3/photos", nil), auth))

	// Signed an hour ago: replay window must be closed.
	old := httptest.NewRequest("GET", "/s3/photos", nil)
	signV4Header(old, engine.sigv4.AccessKey, engine.sigv4.SecretKey, shaHex(nil), time.Now().Add(-time.Hour))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, old)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("stale signature: want 403, got %d", rec.Code)
	}

	// UNSIGNED-PAYLOAD is not allowed on writes.
	unsigned := httptest.NewRequest("PUT", "/s3/photos/b.txt", strings.NewReader("data"))
	signV4Header(unsigned, engine.sigv4.AccessKey, engine.sigv4.SecretKey, "UNSIGNED-PAYLOAD", time.Now())
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, unsigned)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("unsigned write: want 403, got %d", rec.Code)
	}
}

func TestCredentialsAreNotMasterKey(t *testing.T) {
	engine, _, _, tmpDir := newSecurityEnv(t)
	if engine.sigv4 == nil {
		t.Fatal("s3 credentials not loaded")
	}
	if engine.sigv4.SecretKey == hex.EncodeToString(engine.cfg.MasterKey) {
		t.Fatal("s3 secret must not equal the master key")
	}

	raw, err := os.ReadFile(filepath.Join(tmpDir, "data", S3CredentialsFile))
	if err != nil {
		t.Fatalf("credentials file missing: %v", err)
	}
	if strings.Contains(string(raw), engine.sigv4.SecretKey) {
		t.Fatal("secret stored in plaintext")
	}

	// Reloading yields the same pair.
	again, err := loadOrCreateS3Credentials(filepath.Join(tmpDir, "data"), engine.cfg.MasterKey)
	if err != nil || again.SecretKey != engine.sigv4.SecretKey {
		t.Fatalf("reload mismatch: %v", err)
	}
}

func TestBucketAndKeyValidation(t *testing.T) {
	_, mux, auth, _ := newSecurityEnv(t)

	cases := []struct{ method, path string }{
		{"PUT", "/s3/Bad_Name"},
		{"PUT", "/s3/a"},
		{"PUT", "/s3/bucket..x/k"},
		{"PUT", "/s3/validbucket/.multipart/x"},
		{"PUT", "/s3/validbucket/.s3_upload_evil"},
	}
	for _, c := range cases {
		req := withAuth(httptest.NewRequest(c.method, c.path, strings.NewReader("x")), auth)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("%s %s: want 400, got %d", c.method, c.path, rec.Code)
		}
	}
}

func TestActiveContentIsNeutralized(t *testing.T) {
	_, mux, auth, _ := newSecurityEnv(t)
	mux.ServeHTTP(httptest.NewRecorder(), withAuth(httptest.NewRequest("PUT", "/s3/public-site", nil), auth))

	req := withAuth(httptest.NewRequest("PUT", "/s3/public-site/evil.html", strings.NewReader("<script>x</script>")), auth)
	req.Header.Set("Content-Type", "text/html")
	mux.ServeHTTP(httptest.NewRecorder(), req)

	get := withAuth(httptest.NewRequest("GET", "/s3/public-site/evil.html", nil), auth)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, get)
	if ct := rec.Header().Get("Content-Type"); ct != "application/octet-stream" {
		t.Fatalf("html served as %q, want application/octet-stream", ct)
	}
	if rec.Header().Get("Content-Disposition") != "attachment" {
		t.Fatal("active content must be served as attachment")
	}
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("missing nosniff header")
	}
	if !strings.Contains(rec.Header().Get("Content-Security-Policy"), "sandbox") {
		t.Fatal("missing sandbox CSP")
	}
}

func TestMultipartBoundToBucketAndKey(t *testing.T) {
	_, mux, auth, _ := newSecurityEnv(t)
	mux.ServeHTTP(httptest.NewRecorder(), withAuth(httptest.NewRequest("PUT", "/s3/store", nil), auth))

	init := withAuth(httptest.NewRequest("POST", "/s3/store/a.bin?uploads", nil), auth)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, init)
	body := rec.Body.String()
	id := body[strings.Index(body, "<UploadId>")+len("<UploadId>") : strings.Index(body, "</UploadId>")]

	// Using the upload for a different key must fail.
	wrongKey := withAuth(httptest.NewRequest("PUT", "/s3/store/b.bin?uploadId="+id+"&partNumber=1", strings.NewReader("x")), auth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, wrongKey)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("cross-key part upload: want 404, got %d", rec.Code)
	}

	// Malformed upload IDs (path tricks) are rejected.
	bad := withAuth(httptest.NewRequest("DELETE", "/s3/store/a.bin?uploadId=../../x", nil), auth)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, bad)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("malformed upload id: want 404, got %d", rec.Code)
	}
}

func TestObjectSizeLimit(t *testing.T) {
	_, mux, auth, _ := newSecurityEnv(t)
	mux.ServeHTTP(httptest.NewRecorder(), withAuth(httptest.NewRequest("PUT", "/s3/limits", nil), auth))

	prev := MaxObjectSize
	MaxObjectSize = 10
	defer func() { MaxObjectSize = prev }()

	req := withAuth(httptest.NewRequest("PUT", "/s3/limits/big.bin", bytes.NewReader(make([]byte, 100))), auth)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversize PUT: want 413, got %d", rec.Code)
	}
}

func withAuth(r *http.Request, auth string) *http.Request {
	r.Header.Set("Authorization", auth)
	return r
}
