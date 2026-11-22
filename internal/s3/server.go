package s3

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"orbit/internal/config"
	"orbit/internal/db"
	"orbit/internal/security"
)

type S3Engine struct {
	cfg      *config.Config
	store    *db.StorageStore
	sigv4    *SigV4Auth
	tokenMgr *security.TokenManager
}

func NewS3Engine(
	cfg *config.Config,
	store *db.StorageStore,
	tokenMgr *security.TokenManager,
) *S3Engine {
	// Root master access key and secret key derived from master key
	accessKey := "ORBIT" + cfg.NodeID
	secretKey := tokenMgr.MasterKeyHex()

	return &S3Engine{
		cfg:      cfg,
		store:    store,
		sigv4:    NewSigV4Auth(accessKey, secretKey, "us-east-1"),
		tokenMgr: tokenMgr,
	}
}

func (s *S3Engine) RegisterRoutes(mux *http.ServeMux) {
	// Dispatcher for /s3/...
	mux.HandleFunc("/s3/", s.handleS3Dispatch)
}

func (s *S3Engine) BucketDir(bucketID string) string {
	return filepath.Join(s.cfg.StorageDir, bucketID)
}

func (s *S3Engine) authenticate(r *http.Request, isPublicBucket bool) bool {
	// 1. If public bucket and safe read method (GET/HEAD), allow without auth
	if isPublicBucket && (r.Method == http.MethodGet || r.Method == http.MethodHead) {
		return true
	}

	// 2. Check Bearer token or master key header
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if _, err := s.tokenMgr.VerifyToken(token); err == nil {
			return true
		}
	}
	if xKey := r.Header.Get("X-Orbit-Key"); xKey != "" {
		if _, err := s.tokenMgr.VerifyToken(xKey); err == nil {
			return true
		}
	}

	// 3. Check AWS SigV4 (Header or Presigned)
	if err := s.sigv4.VerifyRequest(r, ""); err == nil {
		return true
	}

	return false
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

	// Fetch bucket to check public access
	bucket, _ := s.store.GetBucket(bucketName)
	isPublic := bucket != nil && bucket.IsPublic

	if !s.authenticate(r, isPublic) {
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
		if partNumStr := q.Get("partNumber"); partNumStr != "" && r.Method == http.MethodPut {
			s.handleUploadPart(w, r, bucketName, objectKey, uploadID, partNumStr)
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
		s.handlePutObject(w, r, bucketName, objectKey)
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

// ----------------------------------------------------------------------------
// Bucket Operations
// ----------------------------------------------------------------------------

func (s *S3Engine) handleListBuckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := s.store.ListBuckets()
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
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

func (s *S3Engine) handlePutObject(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	bucketDir := s.BucketDir(bucketName)
	_ = os.MkdirAll(bucketDir, 0700)

	// Safe relative destination path
	cleanKey := filepath.Clean("/" + objectKey)
	dstPath := filepath.Join(bucketDir, cleanKey)
	_ = os.MkdirAll(filepath.Dir(dstPath), 0700)

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

	hash := md5.New()
	multiWriter := io.MultiWriter(tmpFile, hash)

	written, err := io.Copy(multiWriter, r.Body)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	_ = tmpFile.Sync()
	_ = tmpFile.Close()

	if err := os.Rename(tmpFile.Name(), dstPath); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	etag := hex.EncodeToString(hash.Sum(nil))
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(objectKey))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	obj := &db.StorageObject{
		ObjectID:    fmt.Sprintf("obj_%x", hash.Sum(nil)),
		BucketID:    bucketName,
		ObjectKey:   objectKey,
		SizeBytes:   written,
		ContentType: contentType,
		ETag:        etag,
		DiskRelPath: cleanKey,
	}

	if err := s.store.UpsertObject(obj); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	w.Header().Set("ETag", fmt.Sprintf("\"%s\"", etag))
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
	w.Header().Set("Content-Type", obj.ContentType)
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
	w.Header().Set("Content-Type", obj.ContentType)
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
	uploadID := fmt.Sprintf("mp_%d", time.Now().UnixNano())
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

func (s *S3Engine) handleUploadPart(w http.ResponseWriter, r *http.Request, bucketName, objectKey, uploadID, partNumStr string) {
	partNum, err := strconv.Atoi(partNumStr)
	if err != nil || partNum < 1 || partNum > 10000 {
		writeS3Error(w, http.StatusBadRequest, "InvalidPart", "Part number must be an integer between 1 and 10000")
		return
	}

	partsDir := filepath.Join(s.BucketDir(bucketName), ".multipart", uploadID)
	_ = os.MkdirAll(partsDir, 0700)

	partFile := filepath.Join(partsDir, fmt.Sprintf("part_%d", partNum))
	f, err := os.Create(partFile)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	defer f.Close()

	hash := md5.New()
	written, err := io.Copy(io.MultiWriter(f, hash), r.Body)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	etag := hex.EncodeToString(hash.Sum(nil))
	_ = s.store.SaveMultipartPart(uploadID, partNum, etag, written, partFile)

	w.Header().Set("ETag", fmt.Sprintf("\"%s\"", etag))
	w.WriteHeader(http.StatusOK)
}

func (s *S3Engine) handleCompleteMultipart(w http.ResponseWriter, r *http.Request, bucketName, objectKey, uploadID string) {
	parts, err := s.store.GetMultipartParts(uploadID)
	if err != nil || len(parts) == 0 {
		writeS3Error(w, http.StatusBadRequest, "InvalidPart", "No parts found for upload ID")
		return
	}

	cleanKey := filepath.Clean("/" + objectKey)
	dstPath := filepath.Join(s.BucketDir(bucketName), cleanKey)
	_ = os.MkdirAll(filepath.Dir(dstPath), 0700)

	finalFile, err := os.Create(dstPath)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	defer finalFile.Close()

	compositeMD5 := md5.New()
	var totalSize int64

	for _, p := range parts {
		pf, err := os.Open(p.PartPath)
		if err != nil {
			writeS3Error(w, http.StatusInternalServerError, "InternalError", "Missing part file")
			return
		}
		wLen, _ := io.Copy(io.MultiWriter(finalFile, compositeMD5), pf)
		pf.Close()
		totalSize += wLen
	}

	compositeETag := fmt.Sprintf("%x-%d", compositeMD5.Sum(nil), len(parts))

	// Clean up multipart staging directory and records
	partsDir := filepath.Join(s.BucketDir(bucketName), ".multipart", uploadID)
	_ = os.RemoveAll(partsDir)
	_ = s.store.DeleteMultipartUpload(uploadID)

	contentType := mime.TypeByExtension(filepath.Ext(objectKey))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	obj := &db.StorageObject{
		ObjectID:    fmt.Sprintf("obj_%x", compositeMD5.Sum(nil)),
		BucketID:    bucketName,
		ObjectKey:   objectKey,
		SizeBytes:   totalSize,
		ContentType: contentType,
		ETag:        compositeETag,
		DiskRelPath: cleanKey,
	}
	_ = s.store.UpsertObject(obj)

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
