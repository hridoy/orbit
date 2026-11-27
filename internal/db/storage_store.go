package db

import (
	"database/sql"
	"fmt"
	"time"
)

type StorageBucket struct {
	BucketID          string    `json:"bucket_id"`
	OwnerID           string    `json:"owner_id"`
	IsPublic          bool      `json:"is_public"`
	QuotaBytes        int64     `json:"quota_bytes"`
	UsedBytes         int64     `json:"used_bytes"`
	VersioningEnabled bool      `json:"versioning_enabled"`
	ObjectLockEnabled bool      `json:"object_lock_enabled"`
	RetentionDays     int       `json:"retention_days"`
	LifecycleDays     int       `json:"lifecycle_days"`
	WebhookURL        string    `json:"webhook_url"`
	ErasureEnabled    bool      `json:"erasure_enabled"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type StorageIAMKey struct {
	KeyID        string    `json:"key_id"`
	AccessKey    string    `json:"access_key"`
	SecretKeyEnc string    `json:"-"`
	Label        string    `json:"label"`
	BucketScope  string    `json:"bucket_scope"`
	PrefixScope  string    `json:"prefix_scope"`
	Permissions  string    `json:"permissions"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
}

type StorageObject struct {
	ObjectID    string     `json:"object_id"`
	BucketID    string     `json:"bucket_id"`
	ObjectKey   string     `json:"object_key"`
	SizeBytes   int64      `json:"size_bytes"`
	ContentType string     `json:"content_type"`
	ETag        string     `json:"etag"`
	DiskRelPath string     `json:"disk_rel_path"`
	RetainUntil *time.Time `json:"retain_until,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type StorageObjectVersion struct {
	VersionID      string    `json:"version_id"`
	BucketID       string    `json:"bucket_id"`
	ObjectKey      string    `json:"object_key"`
	SizeBytes      int64     `json:"size_bytes"`
	ContentType    string    `json:"content_type"`
	ETag           string    `json:"etag"`
	DiskRelPath    string    `json:"disk_rel_path"`
	IsDeleteMarker bool      `json:"is_delete_marker"`
	CreatedAt      time.Time `json:"created_at"`
}

type StorageStore struct {
	db *DB
}

func NewStorageStore(db *DB) *StorageStore {
	return &StorageStore{db: db}
}

func (s *StorageStore) CreateBucket(b *StorageBucket) error {
	query := `
	INSERT INTO storage_buckets (bucket_id, owner_id, is_public, quota_bytes, used_bytes)
	VALUES (?, ?, ?, ?, ?);
	`
	isPublicInt := 0
	if b.IsPublic {
		isPublicInt = 1
	}
	_, err := s.db.Exec(query, b.BucketID, b.OwnerID, isPublicInt, b.QuotaBytes, b.UsedBytes)
	return err
}

func (s *StorageStore) DeleteBucket(bucketID string) error {
	_, _ = s.db.Exec(`DELETE FROM storage_objects WHERE bucket_id = ?;`, bucketID)
	_, err := s.db.Exec(`DELETE FROM storage_buckets WHERE bucket_id = ?;`, bucketID)
	return err
}

func (s *StorageStore) GetBucket(bucketID string) (*StorageBucket, error) {
	query := `
	SELECT bucket_id, owner_id, is_public, quota_bytes, used_bytes,
	       versioning_enabled, object_lock_enabled, retention_days, lifecycle_days, webhook_url, erasure_enabled,
	       created_at, updated_at
	FROM storage_buckets WHERE bucket_id = ?;
	`
	row := s.db.QueryRow(query, bucketID)
	var b StorageBucket
	var isPublicInt, versioningInt, objectLockInt, erasureInt int
	err := row.Scan(
		&b.BucketID, &b.OwnerID, &isPublicInt, &b.QuotaBytes, &b.UsedBytes,
		&versioningInt, &objectLockInt, &b.RetentionDays, &b.LifecycleDays, &b.WebhookURL, &erasureInt,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	b.IsPublic = isPublicInt == 1
	b.VersioningEnabled = versioningInt == 1
	b.ObjectLockEnabled = objectLockInt == 1
	b.ErasureEnabled = erasureInt == 1
	return &b, nil
}

func (s *StorageStore) ListBuckets() ([]*StorageBucket, error) {
	query := `
	SELECT bucket_id, owner_id, is_public, quota_bytes, used_bytes,
	       versioning_enabled, object_lock_enabled, retention_days, lifecycle_days, webhook_url, erasure_enabled,
	       created_at, updated_at
	FROM storage_buckets ORDER BY created_at DESC;
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buckets []*StorageBucket
	for rows.Next() {
		var b StorageBucket
		var isPublicInt, versioningInt, objectLockInt, erasureInt int
		if err := rows.Scan(
			&b.BucketID, &b.OwnerID, &isPublicInt, &b.QuotaBytes, &b.UsedBytes,
			&versioningInt, &objectLockInt, &b.RetentionDays, &b.LifecycleDays, &b.WebhookURL, &erasureInt,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		b.IsPublic = isPublicInt == 1
		b.VersioningEnabled = versioningInt == 1
		b.ObjectLockEnabled = objectLockInt == 1
		b.ErasureEnabled = erasureInt == 1
		buckets = append(buckets, &b)
	}
	return buckets, nil
}

func (s *StorageStore) UpdateBucketSettings(
	bucketID string,
	isPublic, versioning, objectLock, erasure bool,
	retentionDays, lifecycleDays int,
	quotaBytes int64,
	webhookURL string,
) error {
	b2i := func(v bool) int {
		if v {
			return 1
		}
		return 0
	}
	query := `
	UPDATE storage_buckets
	SET is_public = ?, versioning_enabled = ?, object_lock_enabled = ?, erasure_enabled = ?,
	    retention_days = ?, lifecycle_days = ?, quota_bytes = ?, webhook_url = ?,
	    updated_at = CURRENT_TIMESTAMP
	WHERE bucket_id = ?;
	`
	_, err := s.db.Exec(query,
		b2i(isPublic), b2i(versioning), b2i(objectLock), b2i(erasure),
		retentionDays, lifecycleDays, quotaBytes, webhookURL, bucketID,
	)
	return err
}

func (s *StorageStore) UpsertObject(o *StorageObject) error {
	query := `
	INSERT INTO storage_objects (
		object_id, bucket_id, object_key, size_bytes, content_type, etag, disk_rel_path, retain_until, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(bucket_id, object_key) DO UPDATE SET
		size_bytes = excluded.size_bytes,
		content_type = excluded.content_type,
		etag = excluded.etag,
		disk_rel_path = excluded.disk_rel_path,
		retain_until = excluded.retain_until,
		updated_at = CURRENT_TIMESTAMP;
	`
	_, err := s.db.Exec(query,
		o.ObjectID, o.BucketID, o.ObjectKey, o.SizeBytes, o.ContentType, o.ETag, o.DiskRelPath, o.RetainUntil,
	)
	if err != nil {
		return fmt.Errorf("upsert object: %w", err)
	}

	// Update used_bytes in bucket
	updateBucketQuery := `
	UPDATE storage_buckets 
	SET used_bytes = (SELECT COALESCE(SUM(size_bytes), 0) FROM storage_objects WHERE bucket_id = ?)
	WHERE bucket_id = ?;
	`
	_, _ = s.db.Exec(updateBucketQuery, o.BucketID, o.BucketID)
	return nil
}

func (s *StorageStore) GetObject(bucketID, key string) (*StorageObject, error) {
	query := `
	SELECT object_id, bucket_id, object_key, size_bytes, content_type, etag, disk_rel_path, retain_until, created_at, updated_at
	FROM storage_objects WHERE bucket_id = ? AND object_key = ?;
	`
	row := s.db.QueryRow(query, bucketID, key)
	var o StorageObject
	var retainUntil sql.NullTime
	err := row.Scan(&o.ObjectID, &o.BucketID, &o.ObjectKey, &o.SizeBytes, &o.ContentType, &o.ETag, &o.DiskRelPath, &retainUntil, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if retainUntil.Valid {
		t := retainUntil.Time
		o.RetainUntil = &t
	}
	return &o, nil
}

func (s *StorageStore) DeleteObject(bucketID, key string) error {
	query := `DELETE FROM storage_objects WHERE bucket_id = ? AND object_key = ?;`
	_, err := s.db.Exec(query, bucketID, key)
	if err != nil {
		return err
	}

	updateBucketQuery := `
	UPDATE storage_buckets 
	SET used_bytes = (SELECT COALESCE(SUM(size_bytes), 0) FROM storage_objects WHERE bucket_id = ?)
	WHERE bucket_id = ?;
	`
	_, _ = s.db.Exec(updateBucketQuery, bucketID, bucketID)
	return nil
}

func (s *StorageStore) ListObjects(bucketID, prefix string, limit int) ([]*StorageObject, error) {
	if limit <= 0 || limit > 1000 {
		limit = 1000
	}
	query := `
	SELECT object_id, bucket_id, object_key, size_bytes, content_type, etag, disk_rel_path, retain_until, created_at, updated_at
	FROM storage_objects 
	WHERE bucket_id = ? AND object_key LIKE ?
	ORDER BY object_key ASC
	LIMIT ?;
	`
	rows, err := s.db.Query(query, bucketID, prefix+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objects []*StorageObject
	for rows.Next() {
		var o StorageObject
		var retainUntil sql.NullTime
		if err := rows.Scan(&o.ObjectID, &o.BucketID, &o.ObjectKey, &o.SizeBytes, &o.ContentType, &o.ETag, &o.DiskRelPath, &retainUntil, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		if retainUntil.Valid {
			t := retainUntil.Time
			o.RetainUntil = &t
		}
		objects = append(objects, &o)
	}
	return objects, nil
}

func (s *StorageStore) CountObjects(bucketID string) (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM storage_objects WHERE bucket_id = ?;`, bucketID).Scan(&count)
	return count, err
}

func (s *StorageStore) SaveObjectVersion(v *StorageObjectVersion) error {
	dm := 0
	if v.IsDeleteMarker {
		dm = 1
	}
	query := `
	INSERT INTO storage_object_versions (
		version_id, bucket_id, object_key, size_bytes, content_type, etag, disk_rel_path, is_delete_marker
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err := s.db.Exec(query, v.VersionID, v.BucketID, v.ObjectKey, v.SizeBytes, v.ContentType, v.ETag, v.DiskRelPath, dm)
	return err
}

func (s *StorageStore) ListObjectVersions(bucketID, objectKey string) ([]*StorageObjectVersion, error) {
	var query string
	var args []any
	if objectKey != "" {
		query = `
		SELECT version_id, bucket_id, object_key, size_bytes, content_type, etag, disk_rel_path, is_delete_marker, created_at
		FROM storage_object_versions
		WHERE bucket_id = ? AND object_key = ?
		ORDER BY created_at DESC LIMIT 200;`
		args = []any{bucketID, objectKey}
	} else {
		query = `
		SELECT version_id, bucket_id, object_key, size_bytes, content_type, etag, disk_rel_path, is_delete_marker, created_at
		FROM storage_object_versions
		WHERE bucket_id = ?
		ORDER BY created_at DESC LIMIT 200;`
		args = []any{bucketID}
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*StorageObjectVersion
	for rows.Next() {
		var v StorageObjectVersion
		var dm int
		if err := rows.Scan(&v.VersionID, &v.BucketID, &v.ObjectKey, &v.SizeBytes, &v.ContentType, &v.ETag, &v.DiskRelPath, &dm, &v.CreatedAt); err != nil {
			return nil, err
		}
		v.IsDeleteMarker = dm == 1
		versions = append(versions, &v)
	}
	return versions, nil
}

func (s *StorageStore) GetObjectVersion(bucketID, versionID string) (*StorageObjectVersion, error) {
	query := `
	SELECT version_id, bucket_id, object_key, size_bytes, content_type, etag, disk_rel_path, is_delete_marker, created_at
	FROM storage_object_versions
	WHERE bucket_id = ? AND version_id = ?;`
	var v StorageObjectVersion
	var dm int
	err := s.db.QueryRow(query, bucketID, versionID).Scan(
		&v.VersionID, &v.BucketID, &v.ObjectKey, &v.SizeBytes, &v.ContentType, &v.ETag, &v.DiskRelPath, &dm, &v.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	v.IsDeleteMarker = dm == 1
	return &v, nil
}

func (s *StorageStore) DeleteObjectVersion(bucketID, versionID string) error {
	_, err := s.db.Exec(`DELETE FROM storage_object_versions WHERE bucket_id = ? AND version_id = ?;`, bucketID, versionID)
	return err
}

func (s *StorageStore) ListExpiredObjects(bucketID string, olderThan time.Time) ([]*StorageObject, error) {
	query := `
	SELECT object_id, bucket_id, object_key, size_bytes, content_type, etag, disk_rel_path, retain_until, created_at, updated_at
	FROM storage_objects
	WHERE bucket_id = ? AND updated_at < ?;`
	rows, err := s.db.Query(query, bucketID, olderThan.UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*StorageObject
	for rows.Next() {
		var o StorageObject
		var retainUntil sql.NullTime
		if err := rows.Scan(&o.ObjectID, &o.BucketID, &o.ObjectKey, &o.SizeBytes, &o.ContentType, &o.ETag, &o.DiskRelPath, &retainUntil, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		if retainUntil.Valid {
			t := retainUntil.Time
			o.RetainUntil = &t
		}
		out = append(out, &o)
	}
	return out, nil
}

func (s *StorageStore) CreateIAMKey(k *StorageIAMKey) error {
	activeInt := 0
	if k.Active {
		activeInt = 1
	}
	query := `
	INSERT INTO storage_iam_keys (
		key_id, access_key, secret_key_enc, label, bucket_scope, prefix_scope, permissions, active
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err := s.db.Exec(query, k.KeyID, k.AccessKey, k.SecretKeyEnc, k.Label, k.BucketScope, k.PrefixScope, k.Permissions, activeInt)
	return err
}

func (s *StorageStore) GetIAMKeyByAccessKey(accessKey string) (*StorageIAMKey, error) {
	query := `
	SELECT key_id, access_key, secret_key_enc, label, bucket_scope, prefix_scope, permissions, active, created_at
	FROM storage_iam_keys WHERE access_key = ? AND active = 1;
	`
	var k StorageIAMKey
	var activeInt int
	err := s.db.QueryRow(query, accessKey).Scan(
		&k.KeyID, &k.AccessKey, &k.SecretKeyEnc, &k.Label, &k.BucketScope, &k.PrefixScope, &k.Permissions, &activeInt, &k.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	k.Active = activeInt == 1
	return &k, nil
}

func (s *StorageStore) ListIAMKeys() ([]*StorageIAMKey, error) {
	query := `
	SELECT key_id, access_key, label, bucket_scope, prefix_scope, permissions, active, created_at
	FROM storage_iam_keys ORDER BY created_at DESC;
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*StorageIAMKey
	for rows.Next() {
		var k StorageIAMKey
		var activeInt int
		if err := rows.Scan(&k.KeyID, &k.AccessKey, &k.Label, &k.BucketScope, &k.PrefixScope, &k.Permissions, &activeInt, &k.CreatedAt); err != nil {
			return nil, err
		}
		k.Active = activeInt == 1
		keys = append(keys, &k)
	}
	return keys, nil
}

func (s *StorageStore) DeleteIAMKey(keyID string) error {
	_, err := s.db.Exec(`DELETE FROM storage_iam_keys WHERE key_id = ?;`, keyID)
	return err
}

type MultipartPartRecord struct {
	UploadID   string `json:"upload_id"`
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
	SizeBytes  int64  `json:"size_bytes"`
	PartPath   string `json:"part_path"`
}

func (s *StorageStore) CreateMultipartUpload(uploadID, bucketID, objectKey string) error {
	query := `INSERT INTO storage_multipart_uploads (upload_id, bucket_id, object_key) VALUES (?, ?, ?);`
	_, err := s.db.Exec(query, uploadID, bucketID, objectKey)
	return err
}

// GetMultipartUpload returns the bucket and key an upload was started for.
func (s *StorageStore) GetMultipartUpload(uploadID string) (bucketID, objectKey string, err error) {
	query := `SELECT bucket_id, object_key FROM storage_multipart_uploads WHERE upload_id = ?;`
	err = s.db.QueryRow(query, uploadID).Scan(&bucketID, &objectKey)
	return bucketID, objectKey, err
}

func (s *StorageStore) SaveMultipartPart(uploadID string, partNumber int, etag string, sizeBytes int64, partPath string) error {
	query := `
	INSERT INTO storage_multipart_parts (upload_id, part_number, etag, size_bytes, part_path)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(upload_id, part_number) DO UPDATE SET
		etag = excluded.etag,
		size_bytes = excluded.size_bytes,
		part_path = excluded.part_path;
	`
	_, err := s.db.Exec(query, uploadID, partNumber, etag, sizeBytes, partPath)
	return err
}

func (s *StorageStore) GetMultipartParts(uploadID string) ([]MultipartPartRecord, error) {
	query := `
	SELECT upload_id, part_number, etag, size_bytes, part_path
	FROM storage_multipart_parts
	WHERE upload_id = ?
	ORDER BY part_number ASC;
	`
	rows, err := s.db.Query(query, uploadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []MultipartPartRecord
	for rows.Next() {
		var p MultipartPartRecord
		if err := rows.Scan(&p.UploadID, &p.PartNumber, &p.ETag, &p.SizeBytes, &p.PartPath); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	return parts, nil
}

func (s *StorageStore) DeleteMultipartUpload(uploadID string) error {
	query := `DELETE FROM storage_multipart_uploads WHERE upload_id = ?;`
	_, err := s.db.Exec(query, uploadID)
	return err
}
