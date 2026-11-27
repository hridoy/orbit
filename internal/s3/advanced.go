package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orbit/internal/db"
	"orbit/internal/security"
)

func (s *S3Engine) registerAdvancedRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/storage/buckets/{name}", s.handleGetBucketDetails)
	mux.HandleFunc("PATCH /api/v1/storage/buckets/{name}/settings", s.handleUpdateBucketSettings)
	mux.HandleFunc("GET /api/v1/storage/buckets/{name}/versions", s.handleListVersions)
	mux.HandleFunc("POST /api/v1/storage/buckets/{name}/versions/{versionId}/restore", s.handleRestoreVersion)
	mux.HandleFunc("DELETE /api/v1/storage/buckets/{name}/versions/{versionId}", s.handleDeleteVersion)
	mux.HandleFunc("POST /api/v1/storage/buckets/{name}/presign", s.handlePresignObject)
	mux.HandleFunc("POST /api/v1/storage/buckets/{name}/lifecycle/run", s.handleRunLifecycle)
}

func requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	claims := security.GetClaims(r.Context())
	if claims == nil || claims.Role != "admin" {
		writeS3Error(w, http.StatusForbidden, "AccessDenied", "admin role required")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *S3Engine) handleGetBucketDetails(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	name := r.PathValue("name")
	b, err := s.store.GetBucket(name)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	if b == nil {
		writeS3Error(w, http.StatusNotFound, "NoSuchBucket", "bucket not found")
		return
	}
	objCount, _ := s.store.CountObjects(name)
	versions, _ := s.store.ListObjectVersions(name, "")
	writeJSON(w, http.StatusOK, map[string]any{
		"bucket":        b,
		"object_count":  objCount,
		"version_count": len(versions),
		"region":        "us-east-1",
	})
}

type bucketSettingsReq struct {
	IsPublic          bool   `json:"is_public"`
	VersioningEnabled bool   `json:"versioning_enabled"`
	ObjectLockEnabled bool   `json:"object_lock_enabled"`
	RetentionDays     int    `json:"retention_days"`
	LifecycleDays     int    `json:"lifecycle_days"`
	QuotaBytes        int64  `json:"quota_bytes"`
	WebhookURL        string `json:"webhook_url"`
}

func (s *S3Engine) handleUpdateBucketSettings(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	name := r.PathValue("name")
	b, err := s.store.GetBucket(name)
	if err != nil || b == nil {
		writeS3Error(w, http.StatusNotFound, "NoSuchBucket", "bucket not found")
		return
	}

	var req bucketSettingsReq
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<16)).Decode(&req); err != nil {
		writeS3Error(w, http.StatusBadRequest, "InvalidRequest", err.Error())
		return
	}
	if req.RetentionDays < 0 {
		req.RetentionDays = 0
	}
	if req.LifecycleDays < 0 {
		req.LifecycleDays = 0
	}
	if req.QuotaBytes < 0 {
		req.QuotaBytes = 0
	}
	req.WebhookURL = strings.TrimSpace(req.WebhookURL)

	if err := s.store.UpdateBucketSettings(
		name,
		req.IsPublic,
		req.VersioningEnabled,
		req.ObjectLockEnabled,
		req.RetentionDays,
		req.LifecycleDays,
		req.QuotaBytes,
		req.WebhookURL,
	); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	updated, _ := s.store.GetBucket(name)
	writeJSON(w, http.StatusOK, map[string]any{"bucket": updated})
}

func (s *S3Engine) handleListVersions(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	name := r.PathValue("name")
	key := r.URL.Query().Get("key")
	versions, err := s.store.ListObjectVersions(name, key)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	if versions == nil {
		versions = []*db.StorageObjectVersion{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"versions": versions})
}

func (s *S3Engine) handleRestoreVersion(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	name := r.PathValue("name")
	versionID := r.PathValue("versionId")

	ver, err := s.store.GetObjectVersion(name, versionID)
	if err != nil || ver == nil {
		writeS3Error(w, http.StatusNotFound, "NoSuchVersion", "version not found")
		return
	}
	if ver.IsDeleteMarker {
		writeS3Error(w, http.StatusBadRequest, "InvalidRequest", "cannot restore a delete marker")
		return
	}

	srcPath := filepath.Join(s.BucketDir(name), ver.DiskRelPath)
	cleanKey := filepath.Clean("/" + ver.ObjectKey)
	dstPath := filepath.Join(s.BucketDir(name), cleanKey)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0700); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	// If a current version exists, archive it first
	if existing, _ := s.store.GetObject(name, ver.ObjectKey); existing != nil {
		_ = s.archiveObjectVersion(name, existing)
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		writeS3Error(w, http.StatusNotFound, "NoSuchVersion", "version payload missing on disk")
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	_, copyErr := io.Copy(dstFile, srcFile)
	_ = dstFile.Close()
	if copyErr != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", copyErr.Error())
		return
	}

	obj := &db.StorageObject{
		ObjectID:    fmt.Sprintf("obj_%s", ver.ETag),
		BucketID:    name,
		ObjectKey:   ver.ObjectKey,
		SizeBytes:   ver.SizeBytes,
		ContentType: ver.ContentType,
		ETag:        ver.ETag,
		DiskRelPath: cleanKey,
	}
	if err := s.store.UpsertObject(obj); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"restored": obj})
}

func (s *S3Engine) handleDeleteVersion(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	name := r.PathValue("name")
	versionID := r.PathValue("versionId")

	ver, _ := s.store.GetObjectVersion(name, versionID)
	if ver != nil && !ver.IsDeleteMarker && ver.DiskRelPath != "" {
		_ = os.Remove(filepath.Join(s.BucketDir(name), ver.DiskRelPath))
	}
	if err := s.store.DeleteObjectVersion(name, versionID); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *S3Engine) handlePresignObject(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	if s.sigv4 == nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", "s3 credentials unavailable")
		return
	}
	name := r.PathValue("name")
	var req struct {
		Key            string `json:"key"`
		Method         string `json:"method"`
		ExpiresSeconds int    `json:"expires_seconds"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<14)).Decode(&req); err != nil || req.Key == "" {
		writeS3Error(w, http.StatusBadRequest, "InvalidRequest", "key is required")
		return
	}
	if req.Method == "" {
		req.Method = http.MethodGet
	}
	if req.ExpiresSeconds <= 0 {
		req.ExpiresSeconds = 3600
	}
	if req.ExpiresSeconds > 7*24*3600 {
		req.ExpiresSeconds = 7 * 24 * 3600
	}
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	signedPath := s.sigv4.GeneratePresignedURL(req.Method, r.Host, name, req.Key, time.Duration(req.ExpiresSeconds)*time.Second)
	if !strings.HasPrefix(signedPath, "http") {
		signedPath = fmt.Sprintf("%s://%s%s", scheme, r.Host, signedPath)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"url":        signedPath,
		"expires_at": time.Now().Add(time.Duration(req.ExpiresSeconds) * time.Second).UTC().Format(time.RFC3339),
	})
}

func (s *S3Engine) handleRunLifecycle(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	name := r.PathValue("name")
	b, err := s.store.GetBucket(name)
	if err != nil || b == nil {
		writeS3Error(w, http.StatusNotFound, "NoSuchBucket", "bucket not found")
		return
	}
	if b.LifecycleDays <= 0 {
		writeJSON(w, http.StatusOK, map[string]any{"expired_count": 0, "message": "Lifecycle expiration is disabled"})
		return
	}

	cutoff := time.Now().Add(-time.Duration(b.LifecycleDays) * 24 * time.Hour)
	expired, err := s.store.ListExpiredObjects(name, cutoff)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	now := time.Now()
	deleted := 0
	for _, obj := range expired {
		if obj.RetainUntil != nil && obj.RetainUntil.After(now) {
			continue // Object Lock protects it
		}
		_ = os.Remove(filepath.Join(s.BucketDir(name), obj.DiskRelPath))
		_ = s.store.DeleteObject(name, obj.ObjectKey)
		deleted++
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"expired_count": deleted,
		"cutoff":        cutoff.UTC().Format(time.RFC3339),
	})
}

func (s *S3Engine) archiveObjectVersion(bucketName string, existing *db.StorageObject) error {
	idHex, err := randomHex(8)
	if err != nil {
		return err
	}
	versionID := "v_" + idHex
	relVerPath := filepath.Join(".versions", versionID)
	absVerDir := filepath.Join(s.BucketDir(bucketName), ".versions")
	if err := os.MkdirAll(absVerDir, 0700); err != nil {
		return err
	}

	srcPath := filepath.Join(s.BucketDir(bucketName), existing.DiskRelPath)
	dstPath := filepath.Join(s.BucketDir(bucketName), relVerPath)

	srcFile, err := os.Open(srcPath)
	if err == nil {
		defer srcFile.Close()
		dstFile, err := os.Create(dstPath)
		if err == nil {
			_, _ = io.Copy(dstFile, srcFile)
			_ = dstFile.Close()
		}
	}

	return s.store.SaveObjectVersion(&db.StorageObjectVersion{
		VersionID:      versionID,
		BucketID:       bucketName,
		ObjectKey:      existing.ObjectKey,
		SizeBytes:      existing.SizeBytes,
		ContentType:    existing.ContentType,
		ETag:           existing.ETag,
		DiskRelPath:    relVerPath,
		IsDeleteMarker: false,
	})
}

func (s *S3Engine) recordDeleteMarker(bucketName, objectKey string) {
	idHex, err := randomHex(8)
	if err != nil {
		return
	}
	_ = s.store.SaveObjectVersion(&db.StorageObjectVersion{
		VersionID:      "dm_" + idHex,
		BucketID:       bucketName,
		ObjectKey:      objectKey,
		SizeBytes:      0,
		ContentType:    "application/x-delete-marker",
		ETag:           "",
		DiskRelPath:    "",
		IsDeleteMarker: true,
	})
}

func (s *S3Engine) emitWebhook(webhookURL, eventName, bucketName, objectKey string, sizeBytes int64, etag string) {
	if webhookURL == "" {
		return
	}
	payload := map[string]any{
		"event":      eventName,
		"bucket":     bucketName,
		"key":        objectKey,
		"size_bytes": sizeBytes,
		"etag":       etag,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Orbit-S3-Webhook/1.0")
	resp, err := client.Do(req)
	if err == nil {
		_ = resp.Body.Close()
	}
}
