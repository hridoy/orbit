package s3

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"orbit/internal/db"
)

func (s *S3Engine) registerIAMAndErasureRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/storage/iam-keys", s.handleListIAMKeys)
	mux.HandleFunc("POST /api/v1/storage/iam-keys", s.handleCreateIAMKey)
	mux.HandleFunc("DELETE /api/v1/storage/iam-keys/{id}", s.handleDeleteIAMKey)
	mux.HandleFunc("GET /api/v1/storage/buckets/{name}/erasure/health", s.handleErasureHealth)
	mux.HandleFunc("POST /api/v1/storage/buckets/{name}/erasure/heal", s.handleErasureHeal)
}

func ExtractSigV4AccessKey(r *http.Request) string {
	if cred := r.URL.Query().Get("X-Amz-Credential"); cred != "" {
		if parts := strings.Split(cred, "/"); len(parts) > 0 {
			return parts[0]
		}
	}
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "AWS4-HMAC-SHA256 ") {
		for _, p := range strings.Split(strings.TrimPrefix(authHeader, "AWS4-HMAC-SHA256 "), ", ") {
			if strings.HasPrefix(p, "Credential=") {
				cred := strings.TrimPrefix(p, "Credential=")
				if parts := strings.Split(cred, "/"); len(parts) > 0 {
					return parts[0]
				}
			}
		}
	}
	return ""
}

func (s *S3Engine) sealTenantSecret(secret string) (string, error) {
	aead, err := credentialAEAD(s.cfg.MasterKey)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ct := aead.Seal(nonce, nonce, []byte(secret), nil)
	return base64.StdEncoding.EncodeToString(ct), nil
}

func (s *S3Engine) unsealTenantSecret(enc string) (string, error) {
	aead, err := credentialAEAD(s.cfg.MasterKey)
	if err != nil {
		return "", err
	}
	raw, err := base64.StdEncoding.DecodeString(enc)
	if err != nil || len(raw) < aead.NonceSize() {
		return "", fmt.Errorf("invalid sealed secret")
	}
	nonce, ct := raw[:aead.NonceSize()], raw[aead.NonceSize():]
	plain, err := aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func checkIAMPolicy(k *db.StorageIAMKey, method, bucketName, objectKey string) bool {
	if k == nil || !k.Active {
		return false
	}
	if k.BucketScope != "" && k.BucketScope != "*" && k.BucketScope != bucketName {
		return false
	}
	if k.PrefixScope != "" && objectKey != "" && !strings.HasPrefix(objectKey, k.PrefixScope) {
		return false
	}
	perms := "," + strings.ToLower(strings.ReplaceAll(k.Permissions, " ", "")) + ","
	if strings.Contains(perms, ",admin,") {
		return true
	}
	switch method {
	case http.MethodGet, http.MethodHead:
		return strings.Contains(perms, ",read,") || strings.Contains(perms, ",write,")
	case http.MethodPut, http.MethodPost:
		return strings.Contains(perms, ",write,")
	case http.MethodDelete:
		return strings.Contains(perms, ",delete,")
	default:
		return false
	}
}

func (s *S3Engine) handleListIAMKeys(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	keys, err := s.store.ListIAMKeys()
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	if keys == nil {
		keys = []*db.StorageIAMKey{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"keys": keys})
}

func (s *S3Engine) handleCreateIAMKey(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	var req struct {
		Label       string `json:"label"`
		BucketScope string `json:"bucket_scope"`
		PrefixScope string `json:"prefix_scope"`
		Permissions string `json:"permissions"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<14)).Decode(&req); err != nil {
		writeS3Error(w, http.StatusBadRequest, "InvalidRequest", err.Error())
		return
	}
	req.Label = strings.TrimSpace(req.Label)
	if req.Label == "" {
		req.Label = "tenant-key"
	}
	if req.BucketScope == "" {
		req.BucketScope = "*"
	}
	if req.Permissions == "" {
		req.Permissions = "read"
	}

	idHex, err := randomHex(6)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	akHex, err := randomHex(10)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	skHex, err := randomHex(24)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	accessKey := "ORBIAM" + strings.ToUpper(akHex)
	secretKey := skHex
	sealed, err := s.sealTenantSecret(secretKey)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	rec := &db.StorageIAMKey{
		KeyID:        "iam_" + idHex,
		AccessKey:    accessKey,
		SecretKeyEnc: sealed,
		Label:        req.Label,
		BucketScope:  req.BucketScope,
		PrefixScope:  req.PrefixScope,
		Permissions:  req.Permissions,
		Active:       true,
	}
	if err := s.store.CreateIAMKey(rec); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"key":        rec,
		"secret_key": secretKey, // Returned once on creation
	})
}

func (s *S3Engine) handleDeleteIAMKey(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	id := r.PathValue("id")
	if err := s.store.DeleteIAMKey(id); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *S3Engine) handleErasureHealth(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	name := r.PathValue("name")
	b, err := s.store.GetBucket(name)
	if err != nil || b == nil {
		writeS3Error(w, http.StatusNotFound, "NoSuchBucket", "bucket not found")
		return
	}
	objects, _ := s.store.ListObjects(name, "", 1000)
	drives := s.DriveDirs()

	type driveStat struct {
		DriveIndex int    `json:"drive_index"`
		Path       string `json:"path"`
		Healthy    int    `json:"healthy_shards"`
		Missing    int    `json:"missing_shards"`
		Corrupt    int    `json:"corrupt_shards"`
	}
	stats := make([]driveStat, len(drives))
	for i, d := range drives {
		stats[i] = driveStat{DriveIndex: i, Path: d}
	}

	ecObjects := 0
	degraded := 0
	for _, o := range objects {
		if o.DiskRelPath != "erasure://"+o.ObjectID {
			continue
		}
		ecObjects++
		objMissing := 0
		for i, d := range drives {
			raw, err := os.ReadFile(s.shardPath(d, name, o.ObjectID, i))
			if err != nil {
				stats[i].Missing++
				objMissing++
				continue
			}
			if _, _, ok := unpackShard(raw); !ok {
				stats[i].Corrupt++
				objMissing++
				continue
			}
			stats[i].Healthy++
		}
		if objMissing > 0 {
			degraded++
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"erasure_enabled":  b.ErasureEnabled,
		"data_shards":      DefaultDataShards,
		"parity_shards":    DefaultParityShards,
		"erasure_objects":  ecObjects,
		"degraded_objects": degraded,
		"drives":           stats,
	})
}

func (s *S3Engine) handleErasureHeal(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	name := r.PathValue("name")
	objects, _ := s.store.ListObjects(name, "", 1000)
	healedShards := 0
	failedObjects := 0
	for _, o := range objects {
		if o.DiskRelPath != "erasure://"+o.ObjectID {
			continue
		}
		_, healed, err := s.readErasureObject(name, o.ObjectID)
		if err != nil {
			failedObjects++
		} else {
			healedShards += healed
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"healed_shards":  healedShards,
		"failed_objects": failedObjects,
	})
}
