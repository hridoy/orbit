package s3

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/security"
)

// MaxObjectSize caps a single object or multipart part. Declared as a var so tests can lower it.
var MaxObjectSize int64 = 5 << 30

// minPartSize is the S3 minimum size for every multipart part except the last.
const minPartSize = 5 << 20

type S3Engine struct {
	cfg      *config.Config
	store    *db.StorageStore
	sigv4    *SigV4Auth // nil when the S3 credentials could not be loaded (fail closed)
	credErr  error
	tokenMgr *security.TokenManager
}

func NewS3Engine(
	cfg *config.Config,
	store *db.StorageStore,
	tokenMgr *security.TokenManager,
) *S3Engine {
	engine := &S3Engine{
		cfg:      cfg,
		store:    store,
		tokenMgr: tokenMgr,
	}

	// Random S3 key pair, encrypted at rest. Independent of the vault master key.
	creds, err := loadOrCreateS3Credentials(cfg.DataDir, cfg.MasterKey)
	if err != nil {
		engine.credErr = err
		log.Printf("s3: credentials unavailable, all SigV4 requests will be rejected: %v", err)
		return engine
	}
	engine.sigv4 = NewSigV4Auth(creds.AccessKey, creds.SecretKey, "us-east-1")
	return engine
}

func (s *S3Engine) RegisterRoutes(mux *http.ServeMux) {
	// Dispatcher for /s3/...
	mux.HandleFunc("/s3/", s.handleS3Dispatch)
	// Admin-only: reveal the S3 access key and secret (behind the global bearer auth).
	mux.HandleFunc("GET /api/v1/storage/credentials", s.handleCredentials)
}

func (s *S3Engine) handleCredentials(w http.ResponseWriter, r *http.Request) {
	claims := security.GetClaims(r.Context())
	if claims == nil || claims.Role != "admin" {
		writeS3Error(w, http.StatusForbidden, "AccessDenied", "admin role required")
		return
	}
	if s.sigv4 == nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", "s3 credentials unavailable")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_key": s.sigv4.AccessKey,
		"secret_key": s.sigv4.SecretKey,
		"endpoint":   "/s3",
		"region":     s.sigv4.Region,
	})
}

func (s *S3Engine) BucketDir(bucketID string) string {
	return filepath.Join(s.cfg.StorageDir, bucketID)
}

// authenticate reports whether the request may proceed. The returned string is the
// payload SHA-256 the body must match ("" when no body check applies).
func (s *S3Engine) authenticate(r *http.Request, isPublicBucket bool) (string, bool) {
	// 1. If public bucket and safe read method (GET/HEAD), allow without auth
	if isPublicBucket && (r.Method == http.MethodGet || r.Method == http.MethodHead) {
		return "", true
	}

	// 2. Check Bearer token or X-Orbit-Key (operator/API tokens)
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if _, err := s.tokenMgr.VerifyToken(token); err == nil {
			return "", true
		}
	}
	if xKey := r.Header.Get("X-Orbit-Key"); xKey != "" {
		if _, err := s.tokenMgr.VerifyToken(xKey); err == nil {
			return "", true
		}
	}

	// 3. AWS SigV4 (header or presigned) with the dedicated S3 key pair
	if s.sigv4 == nil {
		return "", false
	}
	payloadHash, err := s.sigv4.VerifyRequest(r)
	if err != nil {
		return "", false
	}
	return payloadHash, true
}

func (s *S3Engine) handleS3Dispatch(w http.ResponseWriter, r *http.Request) {
	// URL Pattern: /s3/{bucket} or /s3/{bucket}/{key...}
	trimmed := strings.TrimPrefix(r.URL.Path, "/s3/")
	if trimmed == "" {
		s.handleListBuckets(w, r)
		return
	}

	parts := strings.SplitN(trimmed, "/", 2)
	bucketName := parts[0]
	var objectKey string
	if len(parts) > 1 {
		objectKey = parts[1]
	}

	if !validBucketName(bucketName) {
		writeS3Error(w, http.StatusBadRequest, "InvalidBucketName", "The specified bucket is not valid.")
		return
	}
	if objectKey != "" && !validObjectKey(objectKey) {
		writeS3Error(w, http.StatusBadRequest, "InvalidArgument", "Object key is not valid.")
		return
	}

	// Fetch bucket to check public access
	bucket, _ := s.store.GetBucket(bucketName)
	isPublic := bucket != nil && bucket.IsPublic

	payloadHash, ok := s.authenticate(r, isPublic)
	if !ok {
		writeS3Error(w, http.StatusForbidden, "AccessDenied", "Access Denied")
		return
	}

	q := r.URL.Query()

	// 1. Bucket-level operations (when objectKey is empty)
	if objectKey == "" {
		switch r.Method {
		case http.MethodPut:
			s.handleCreateBucket(w, r, bucketName)
		case http.MethodGet:
			s.handleListObjectsV2(w, r, bucketName)
		case http.MethodDelete:
			s.handleDeleteBucket(w, r, bucketName)
		default:
			writeS3Error(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "Method not allowed")
		}
		return
	}

	// 2. Multipart operations
	if _, hasUploads := q["uploads"]; hasUploads && r.Method == http.MethodPost {
		s.handleInitiateMultipart(w, r, bucketName, objectKey)
		return
	}
	if uploadID := q.Get("uploadId"); uploadID != "" {
		if !validUploadID(uploadID) {
			writeS3Error(w, http.StatusNotFound, "NoSuchUpload", "The specified upload does not exist.")
			return
		}
		if partNumStr := q.Get("partNumber"); partNumStr != "" && r.Method == http.MethodPut {
			s.handleUploadPart(w, r, bucketName, objectKey, uploadID, partNumStr, payloadHash)
			return
		}
		if r.Method == http.MethodPost {
			s.handleCompleteMultipart(w, r, bucketName, objectKey, uploadID)
			return
		}
		if r.Method == http.MethodDelete {
			s.handleAbortMultipart(w, r, bucketName, objectKey, uploadID)
			return
		}
	}

	// 3. Object-level operations
	switch r.Method {
	case http.MethodPut:
		s.handlePutObject(w, r, bucketName, objectKey, payloadHash)
	case http.MethodGet:
		s.handleGetObject(w, r, bucketName, objectKey)
	case http.MethodHead:
		s.handleHeadObject(w, r, bucketName, objectKey)
	case http.MethodDelete:
		s.handleDeleteObject(w, r, bucketName, objectKey)
	default:
		writeS3Error(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "Method not allowed")
	}
}

var (
	bucketNameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)
	uploadIDRe   = regexp.MustCompile(`^mp_[0-9a-f]{32}$`)
)

// validBucketName enforces the S3 bucket naming rules (also blocks path tricks).
func validBucketName(name string) bool {
	return bucketNameRe.MatchString(name) && !strings.Contains(name, "..")
}

// validObjectKey rejects keys that are empty, too long, or collide with internal paths.
func validObjectKey(key string) bool {
	if key == "" || len(key) > 1024 || strings.ContainsRune(key, 0) {
		return false
	}
	clean := filepath.Clean("/" + key)
	if clean == "/" {
		return false
	}
	first := strings.SplitN(strings.TrimPrefix(clean, "/"), "/", 2)[0]
	return first != ".multipart" && !strings.HasPrefix(first, ".s3_upload_")
}

func validUploadID(id string) bool {
	return uploadIDRe.MatchString(id)
}

// Content types that browsers render as active content. Served as attachments.
var inlineSafeTypes = map[string]bool{
	"image/png":        true,
	"image/jpeg":       true,
	"image/gif":        true,
	"image/webp":       true,
	"text/plain":       true,
	"application/json": true,
}

// normalizeContentType returns a safe stored content type. Active types (html, svg, xml,
// javascript) are downgraded to octet-stream so they cannot execute on the Orbit origin.
func normalizeContentType(ct string) string {
	media, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return "application/octet-stream"
	}
	media = strings.ToLower(media)
	switch media {
	case "text/html", "application/xhtml+xml", "image/svg+xml", "text/xml", "application/xml",
		"text/javascript", "application/javascript", "application/x-javascript":
		return "application/octet-stream"
	}
	if media == "" {
		return "application/octet-stream"
	}
	return media
}

// setObjectHeaders applies response headers that stop stored content from executing.
func setObjectHeaders(w http.ResponseWriter, obj *db.StorageObject) {
	ct := normalizeContentType(obj.ContentType)
	w.Header().Set("Content-Type", ct)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; sandbox")
	if inlineSafeTypes[ct] {
		w.Header().Set("Content-Disposition", "inline")
	} else {
		w.Header().Set("Content-Disposition", "attachment")
	}
}

// copyVerified streams src to dst, enforcing limit and (if expectSHA is set) the declared
// SHA-256. It returns the bytes written and the MD5 and SHA-256 of the body.
func copyVerified(dst io.Writer, src io.Reader, expectSHA string, limit int64) (written int64, md5Hex, sha256Hex string, err error) {
	md5h := md5.New()
	sha := sha256.New()
	lr := &io.LimitedReader{R: src, N: limit + 1}
	written, err = io.Copy(io.MultiWriter(dst, md5h, sha), lr)
	if err != nil {
		return written, "", "", err
	}
	if written > limit {
		return written, "", "", errEntityTooLarge
	}
	sha256Hex = hex.EncodeToString(sha.Sum(nil))
	if expectSHA != "" && sha256Hex != expectSHA {
		return written, "", "", errPayloadMismatch
	}
	return written, hex.EncodeToString(md5h.Sum(nil)), sha256Hex, nil
}

var (
	errEntityTooLarge  = errors.New("entity too large")
	errPayloadMismatch = errors.New("payload hash mismatch")
)

// writeCopyError maps copyVerified failures to S3 error responses.
func writeCopyError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errEntityTooLarge):
		writeS3Error(w, http.StatusRequestEntityTooLarge, "EntityTooLarge", "Your proposed upload exceeds the maximum allowed size.")
	case errors.Is(err, errPayloadMismatch):
		writeS3Error(w, http.StatusBadRequest, "XAmzContentSHA256Mismatch", "The provided x-amz-content-sha256 does not match what was computed.")
	default:
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
	}
}

// ----------------------------------------------------------------------------
// Bucket Operations
// ----------------------------------------------------------------------------

func (s *S3Engine) handleListBuckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := s.store.ListBuckets()
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") || r.URL.Query().Get("format") == "json" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(buckets)
		return
	}

	type BucketXML struct {
		Name         string `xml:"Name"`
		CreationDate string `xml:"CreationDate"`
	}
	type ListAllMyBucketsResult struct {
		XMLName xml.Name    `xml:"ListAllMyBucketsResult"`
		Buckets []BucketXML `xml:"Buckets>Bucket"`
	}

	res := ListAllMyBucketsResult{}
	for _, b := range buckets {
		res.Buckets = append(res.Buckets, BucketXML{
			Name:         b.BucketID,
			CreationDate: b.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	writeXML(w, http.StatusOK, res)
}

func (s *S3Engine) handleCreateBucket(w http.ResponseWriter, r *http.Request, bucketName string) {
	dir := s.BucketDir(bucketName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	b := &db.StorageBucket{
		BucketID:   bucketName,
		OwnerID:    s.cfg.NodeID,
		IsPublic:   false,
		QuotaBytes: 10 * 1024 * 1024 * 1024,
	}
	_ = s.store.CreateBucket(b)

	w.Header().Set("Location", "/s3/"+bucketName)
	w.WriteHeader(http.StatusOK)
}

func (s *S3Engine) handleDeleteBucket(w http.ResponseWriter, r *http.Request, bucketName string) {
	_ = os.RemoveAll(s.BucketDir(bucketName))
	w.WriteHeader(http.StatusNoContent)
}

func (s *S3Engine) handleListObjectsV2(w http.ResponseWriter, r *http.Request, bucketName string) {
	q := r.URL.Query()
	prefix := q.Get("prefix")
	maxKeys := 1000
	if mk := q.Get("max-keys"); mk != "" {
		if val, err := strconv.Atoi(mk); err == nil && val > 0 {
			maxKeys = val
		}
	}

	objects, err := s.store.ListObjects(bucketName, prefix, maxKeys)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	type ContentsXML struct {
		Key          string `xml:"Key"`
		LastModified string `xml:"LastModified"`
		ETag         string `xml:"ETag"`
		Size         int64  `xml:"Size"`
		StorageClass string `xml:"StorageClass"`
	}
	type ListBucketResult struct {
		XMLName     xml.Name      `xml:"ListBucketResult"`
		Name        string        `xml:"Name"`
		Prefix      string        `xml:"Prefix"`
		KeyCount    int           `xml:"KeyCount"`
		MaxKeys     int           `xml:"MaxKeys"`
		IsTruncated bool          `xml:"IsTruncated"`
		Contents    []ContentsXML `xml:"Contents"`
	}

	res := ListBucketResult{
		Name:        bucketName,
		Prefix:      prefix,
		KeyCount:    len(objects),
		MaxKeys:     maxKeys,
		IsTruncated: false,
	}

	for _, o := range objects {
		res.Contents = append(res.Contents, ContentsXML{
			Key:          o.ObjectKey,
			LastModified: o.UpdatedAt.UTC().Format(time.RFC3339),
			ETag:         fmt.Sprintf("\"%s\"", o.ETag),
			Size:         o.SizeBytes,
			StorageClass: "STANDARD",
		})
	}

	writeXML(w, http.StatusOK, res)
}

// ----------------------------------------------------------------------------
// Object Operations (PUT, GET with Range, HEAD, DELETE)
// ----------------------------------------------------------------------------

func (s *S3Engine) handlePutObject(w http.ResponseWriter, r *http.Request, bucketName, objectKey, expectSHA string) {
	bucketDir := s.BucketDir(bucketName)
	if err := os.MkdirAll(bucketDir, 0700); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	// Safe relative destination path
	cleanKey := filepath.Clean("/" + objectKey)
	dstPath := filepath.Join(bucketDir, cleanKey)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0700); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	// Stream directly to temp file
	tmpFile, err := os.CreateTemp(bucketDir, ".s3_upload_*")
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	written, md5Hex, _, err := copyVerified(tmpFile, r.Body, expectSHA, MaxObjectSize)
	if err != nil {
		writeCopyError(w, err)
		return
	}

	if err := tmpFile.Sync(); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	if err := tmpFile.Close(); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	if err := os.Rename(tmpFile.Name(), dstPath); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	contentType := normalizeContentType(r.Header.Get("Content-Type"))
	if r.Header.Get("Content-Type") == "" {
		contentType = normalizeContentType(mime.TypeByExtension(filepath.Ext(objectKey)))
	}

	obj := &db.StorageObject{
		ObjectID:    fmt.Sprintf("obj_%s", md5Hex),
		BucketID:    bucketName,
		ObjectKey:   objectKey,
		SizeBytes:   written,
		ContentType: contentType,
		ETag:        md5Hex,
		DiskRelPath: cleanKey,
	}

	if err := s.store.UpsertObject(obj); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	w.Header().Set("ETag", fmt.Sprintf("\"%s\"", md5Hex))
	w.WriteHeader(http.StatusOK)
}

func (s *S3Engine) handleGetObject(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	obj, err := s.store.GetObject(bucketName, objectKey)
	if err != nil || obj == nil {
		writeS3Error(w, http.StatusNotFound, "NoSuchKey", "The specified key does not exist.")
		return
	}

	filePath := filepath.Join(s.BucketDir(bucketName), obj.DiskRelPath)
	file, err := os.Open(filePath)
	if err != nil {
		writeS3Error(w, http.StatusNotFound, "NoSuchKey", "The specified key does not exist.")
		return
	}
	defer file.Close()

	w.Header().Set("ETag", fmt.Sprintf("\"%s\"", obj.ETag))
	setObjectHeaders(w, obj)
	w.Header().Set("Accept-Ranges", "bytes")

	// Parse HTTP Range header (e.g. bytes=0-1048576)
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" && strings.HasPrefix(rangeHeader, "bytes=") {
		s.handleRangeGet(w, file, obj.SizeBytes, rangeHeader)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(obj.SizeBytes, 10))
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, file)
}

func (s *S3Engine) handleRangeGet(w http.ResponseWriter, file *os.File, totalSize int64, rangeHeader string) {
	byteRange := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(byteRange, "-")

	var start, end int64
	start, _ = strconv.ParseInt(parts[0], 10, 64)
	if len(parts) > 1 && parts[1] != "" {
		end, _ = strconv.ParseInt(parts[1], 10, 64)
	} else {
		end = totalSize - 1
	}

	if start >= totalSize || end >= totalSize || start > end {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", totalSize))
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	contentLength := end - start + 1
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, totalSize))
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.WriteHeader(http.StatusPartialContent)

	_, _ = file.Seek(start, io.SeekStart)
	_, _ = io.CopyN(w, file, contentLength)
}

func (s *S3Engine) handleHeadObject(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	obj, err := s.store.GetObject(bucketName, objectKey)
	if err != nil || obj == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("ETag", fmt.Sprintf("\"%s\"", obj.ETag))
	w.Header().Set("Content-Length", strconv.FormatInt(obj.SizeBytes, 10))
	setObjectHeaders(w, obj)
	w.Header().Set("Last-Modified", obj.UpdatedAt.UTC().Format(http.TimeFormat))
	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusOK)
}

func (s *S3Engine) handleDeleteObject(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	obj, _ := s.store.GetObject(bucketName, objectKey)
	if obj != nil {
		filePath := filepath.Join(s.BucketDir(bucketName), obj.DiskRelPath)
		_ = os.Remove(filePath)
		_ = s.store.DeleteObject(bucketName, objectKey)
	}
	w.WriteHeader(http.StatusNoContent)
}

// ----------------------------------------------------------------------------
// Multipart Upload Handlers
// ----------------------------------------------------------------------------

func (s *S3Engine) handleInitiateMultipart(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	// Unguessable upload ID (16 random bytes), so other callers cannot target it.
	idSuffix, err := randomHex(16)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	uploadID := "mp_" + idSuffix
	if err := s.store.CreateMultipartUpload(uploadID, bucketName, objectKey); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	// Create staging directory for parts
	partsDir := filepath.Join(s.BucketDir(bucketName), ".multipart", uploadID)
	_ = os.MkdirAll(partsDir, 0700)

	type InitiateMultipartUploadResult struct {
		XMLName  xml.Name `xml:"InitiateMultipartUploadResult"`
		Bucket   string   `xml:"Bucket"`
		Key      string   `xml:"Key"`
		UploadId string   `xml:"UploadId"`
	}

	writeXML(w, http.StatusOK, InitiateMultipartUploadResult{
		Bucket:   bucketName,
		Key:      objectKey,
		UploadId: uploadID,
	})
}

// checkUpload ensures uploadID exists and belongs to this bucket and key. It writes the
// S3 NoSuchUpload error and returns false otherwise.
func (s *S3Engine) checkUpload(w http.ResponseWriter, bucketName, objectKey, uploadID string) bool {
	bucketID, key, err := s.store.GetMultipartUpload(uploadID)
	if err != nil || bucketID != bucketName || key != objectKey {
		writeS3Error(w, http.StatusNotFound, "NoSuchUpload", "The specified upload does not exist.")
		return false
	}
	return true
}

func (s *S3Engine) handleUploadPart(w http.ResponseWriter, r *http.Request, bucketName, objectKey, uploadID, partNumStr, expectSHA string) {
	partNum, err := strconv.Atoi(partNumStr)
	if err != nil || partNum < 1 || partNum > 10000 {
		writeS3Error(w, http.StatusBadRequest, "InvalidPart", "Part number must be an integer between 1 and 10000")
		return
	}
	if !s.checkUpload(w, bucketName, objectKey, uploadID) {
		return
	}

	partsDir := filepath.Join(s.BucketDir(bucketName), ".multipart", uploadID)
	if err := os.MkdirAll(partsDir, 0700); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	partFile := filepath.Join(partsDir, fmt.Sprintf("part_%d", partNum))
	f, err := os.Create(partFile)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	defer f.Close()

	written, etag, _, err := copyVerified(f, r.Body, expectSHA, MaxObjectSize)
	if err != nil {
		writeCopyError(w, err)
		return
	}

	if err := s.store.SaveMultipartPart(uploadID, partNum, etag, written, partFile); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	w.Header().Set("ETag", fmt.Sprintf("\"%s\"", etag))
	w.WriteHeader(http.StatusOK)
}

// completeMultipartRequest is the optional body of CompleteMultipartUpload.
type completeMultipartRequest struct {
	XMLName xml.Name `xml:"CompleteMultipartUpload"`
	Parts   []struct {
		PartNumber int    `xml:"PartNumber"`
		ETag       string `xml:"ETag"`
	} `xml:"Part"`
}

func (s *S3Engine) handleCompleteMultipart(w http.ResponseWriter, r *http.Request, bucketName, objectKey, uploadID string) {
	if !s.checkUpload(w, bucketName, objectKey, uploadID) {
		return
	}

	stored, err := s.store.GetMultipartParts(uploadID)
	if err != nil || len(stored) == 0 {
		writeS3Error(w, http.StatusBadRequest, "InvalidPart", "No parts found for upload ID")
		return
	}

	// If the client sends a part list, use exactly those parts, in order, with matching ETags.
	// An empty body means "all uploaded parts" (kept for clients that omit the list).
	selected := stored
	if r.ContentLength != 0 && r.Body != nil {
		var req completeMultipartRequest
		if err := xml.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req); err != nil && err != io.EOF {
			writeS3Error(w, http.StatusBadRequest, "MalformedXML", "The XML you provided was not well-formed.")
			return
		}
		if len(req.Parts) > 0 {
			byNum := make(map[int]db.MultipartPartRecord, len(stored))
			for _, p := range stored {
				byNum[p.PartNumber] = p
			}
			selected = selected[:0:0]
			last := 0
			for _, rp := range req.Parts {
				if rp.PartNumber <= last {
					writeS3Error(w, http.StatusBadRequest, "InvalidPartOrder", "The list of parts was not in ascending order.")
					return
				}
				last = rp.PartNumber
				p, ok := byNum[rp.PartNumber]
				if !ok || p.ETag != strings.Trim(rp.ETag, "\"") {
					writeS3Error(w, http.StatusBadRequest, "InvalidPart", "One or more of the specified parts could not be found or has a different ETag.")
					return
				}
				selected = append(selected, p)
			}
		}
	}

	// Every part except the last must meet the S3 minimum size, and the total must fit the limit.
	var totalSize int64
	for i, p := range selected {
		if i < len(selected)-1 && p.SizeBytes < minPartSize {
			writeS3Error(w, http.StatusBadRequest, "EntityTooSmall", "Your proposed upload is smaller than the minimum allowed size.")
			return
		}
		totalSize += p.SizeBytes
	}
	if totalSize > MaxObjectSize {
		writeS3Error(w, http.StatusRequestEntityTooLarge, "EntityTooLarge", "Your proposed upload exceeds the maximum allowed size.")
		return
	}

	cleanKey := filepath.Clean("/" + objectKey)
	bucketDir := s.BucketDir(bucketName)
	dstPath := filepath.Join(bucketDir, cleanKey)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0700); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	// Assemble into a temp file and rename, so readers never see a half-written object.
	tmpFile, err := os.CreateTemp(bucketDir, ".s3_upload_*")
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	compositeMD5 := md5.New()
	for _, p := range selected {
		pf, err := os.Open(p.PartPath)
		if err != nil {
			writeS3Error(w, http.StatusInternalServerError, "InternalError", "Missing part file")
			return
		}
		_, err = io.Copy(io.MultiWriter(tmpFile, compositeMD5), pf)
		pf.Close()
		if err != nil {
			writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
			return
		}
	}
	if err := tmpFile.Sync(); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	if err := tmpFile.Close(); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	if err := os.Rename(tmpFile.Name(), dstPath); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	compositeETag := fmt.Sprintf("%x-%d", compositeMD5.Sum(nil), len(selected))

	contentType := normalizeContentType(mime.TypeByExtension(filepath.Ext(objectKey)))

	obj := &db.StorageObject{
		ObjectID:    fmt.Sprintf("obj_%x", compositeMD5.Sum(nil)),
		BucketID:    bucketName,
		ObjectKey:   objectKey,
		SizeBytes:   totalSize,
		ContentType: contentType,
		ETag:        compositeETag,
		DiskRelPath: cleanKey,
	}
	if err := s.store.UpsertObject(obj); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	// Clean up staging directory and records only after the object is committed.
	_ = os.RemoveAll(filepath.Join(bucketDir, ".multipart", uploadID))
	_ = s.store.DeleteMultipartUpload(uploadID)

	type CompleteMultipartUploadResult struct {
		XMLName  xml.Name `xml:"CompleteMultipartUploadResult"`
		Location string   `xml:"Location"`
		Bucket   string   `xml:"Bucket"`
		Key      string   `xml:"Key"`
		ETag     string   `xml:"ETag"`
	}

	writeXML(w, http.StatusOK, CompleteMultipartUploadResult{
		Location: fmt.Sprintf("/s3/%s/%s", bucketName, objectKey),
		Bucket:   bucketName,
		Key:      objectKey,
		ETag:     fmt.Sprintf("\"%s\"", compositeETag),
	})
}

func (s *S3Engine) handleAbortMultipart(w http.ResponseWriter, r *http.Request, bucketName, objectKey, uploadID string) {
	if !s.checkUpload(w, bucketName, objectKey, uploadID) {
		return
	}
	partsDir := filepath.Join(s.BucketDir(bucketName), ".multipart", uploadID)
	_ = os.RemoveAll(partsDir)
	_ = s.store.DeleteMultipartUpload(uploadID)
	w.WriteHeader(http.StatusNoContent)
}

func writeXML(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(xml.Header))
	_ = xml.NewEncoder(w).Encode(data)
}

func writeS3Error(w http.ResponseWriter, status int, code, message string) {
	type ErrorResponse struct {
		XMLName xml.Name `xml:"Error"`
		Code    string   `xml:"Code"`
		Message string   `xml:"Message"`
	}
	writeXML(w, status, ErrorResponse{Code: code, Message: message})
}
