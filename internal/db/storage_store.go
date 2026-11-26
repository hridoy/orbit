package db

import (
	"database/sql"
	"fmt"
	"time"
)

type StorageBucket struct {
	BucketID   string    `json:"bucket_id"`
	OwnerID    string    `json:"owner_id"`
	IsPublic   bool      `json:"is_public"`
	QuotaBytes int64     `json:"quota_bytes"`
	UsedBytes  int64     `json:"used_bytes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type StorageObject struct {
	ObjectID    string    `json:"object_id"`
	BucketID    string    `json:"bucket_id"`
	ObjectKey   string    `json:"object_key"`
	SizeBytes   int64     `json:"size_bytes"`
	ContentType string    `json:"content_type"`
	ETag        string    `json:"etag"`
	DiskRelPath string    `json:"disk_rel_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

func (s *StorageStore) GetBucket(bucketID string) (*StorageBucket, error) {
	query := `
	SELECT bucket_id, owner_id, is_public, quota_bytes, used_bytes, created_at, updated_at
	FROM storage_buckets WHERE bucket_id = ?;
	`
	row := s.db.QueryRow(query, bucketID)
	var b StorageBucket
	var isPublicInt int
	err := row.Scan(&b.BucketID, &b.OwnerID, &isPublicInt, &b.QuotaBytes, &b.UsedBytes, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	b.IsPublic = isPublicInt == 1
	return &b, nil
}

func (s *StorageStore) ListBuckets() ([]*StorageBucket, error) {
	query := `
	SELECT bucket_id, owner_id, is_public, quota_bytes, used_bytes, created_at, updated_at
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
		var isPublicInt int
		if err := rows.Scan(&b.BucketID, &b.OwnerID, &isPublicInt, &b.QuotaBytes, &b.UsedBytes, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		b.IsPublic = isPublicInt == 1
		buckets = append(buckets, &b)
	}
	return buckets, nil
}

func (s *StorageStore) UpsertObject(o *StorageObject) error {
	query := `
	INSERT INTO storage_objects (
		object_id, bucket_id, object_key, size_bytes, content_type, etag, disk_rel_path, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(bucket_id, object_key) DO UPDATE SET
		size_bytes = excluded.size_bytes,
		content_type = excluded.content_type,
		etag = excluded.etag,
		disk_rel_path = excluded.disk_rel_path,
		updated_at = CURRENT_TIMESTAMP;
	`
	_, err := s.db.Exec(query,
		o.ObjectID, o.BucketID, o.ObjectKey, o.SizeBytes, o.ContentType, o.ETag, o.DiskRelPath,
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
	SELECT object_id, bucket_id, object_key, size_bytes, content_type, etag, disk_rel_path, created_at, updated_at
	FROM storage_objects WHERE bucket_id = ? AND object_key = ?;
	`
	row := s.db.QueryRow(query, bucketID, key)
	var o StorageObject
	err := row.Scan(&o.ObjectID, &o.BucketID, &o.ObjectKey, &o.SizeBytes, &o.ContentType, &o.ETag, &o.DiskRelPath, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
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
	SELECT object_id, bucket_id, object_key, size_bytes, content_type, etag, disk_rel_path, created_at, updated_at
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
		if err := rows.Scan(&o.ObjectID, &o.BucketID, &o.ObjectKey, &o.SizeBytes, &o.ContentType, &o.ETag, &o.DiskRelPath, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		objects = append(objects, &o)
	}
	return objects, nil
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
