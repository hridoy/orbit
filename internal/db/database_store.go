package db

import (
	"database/sql"
	"fmt"
	"time"
)

type DatabaseInstance struct {
	DBID              string    `json:"db_id"`
	Name              string    `json:"name"`
	Engine            string    `json:"engine"` // 'postgres', 'mongodb'
	Version           string    `json:"version"`
	NodeID            string    `json:"node_id"`
	Port              int       `json:"port"`
	SocketPath        string    `json:"socket_path"`
	RAMProfile        string    `json:"ram_profile"` // 'micro', 'standard', 'heavy'
	StorageLimitMB    int       `json:"storage_limit_mb"`
	DataDirectory     string    `json:"data_directory"`
	AdminUser         string    `json:"admin_user"`
	EncryptedPassword string    `json:"-"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type DatabaseStore struct {
	db *DB
}

func NewDatabaseStore(db *DB) *DatabaseStore {
	return &DatabaseStore{db: db}
}

func (s *DatabaseStore) CreateDatabase(d *DatabaseInstance) error {
	query := `
	INSERT INTO databases (
		db_id, name, engine, version, node_id,
		port, socket_path, ram_profile, storage_limit_mb,
		data_directory, admin_user, encrypted_password, status
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err := s.db.Exec(query,
		d.DBID, d.Name, d.Engine, d.Version, d.NodeID,
		d.Port, d.SocketPath, d.RAMProfile, d.StorageLimitMB,
		d.DataDirectory, d.AdminUser, d.EncryptedPassword, d.Status,
	)
	if err != nil {
		return fmt.Errorf("create database %s: %w", d.DBID, err)
	}
	return nil
}

func (s *DatabaseStore) GetDatabase(dbID string) (*DatabaseInstance, error) {
	query := `
	SELECT db_id, name, engine, version, node_id,
	       port, socket_path, ram_profile, storage_limit_mb,
	       data_directory, admin_user, encrypted_password, status,
	       created_at, updated_at
	FROM databases WHERE db_id = ?;
	`
	row := s.db.QueryRow(query, dbID)
	var d DatabaseInstance
	err := row.Scan(
		&d.DBID, &d.Name, &d.Engine, &d.Version, &d.NodeID,
		&d.Port, &d.SocketPath, &d.RAMProfile, &d.StorageLimitMB,
		&d.DataDirectory, &d.AdminUser, &d.EncryptedPassword, &d.Status,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get database %s: %w", dbID, err)
	}
	return &d, nil
}

func (s *DatabaseStore) ListDatabases() ([]*DatabaseInstance, error) {
	query := `
	SELECT db_id, name, engine, version, node_id,
	       port, socket_path, ram_profile, storage_limit_mb,
	       data_directory, admin_user, encrypted_password, status,
	       created_at, updated_at
	FROM databases ORDER BY created_at DESC;
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("list databases: %w", err)
	}
	defer rows.Close()

	var list []*DatabaseInstance
	for rows.Next() {
		var d DatabaseInstance
		err := rows.Scan(
			&d.DBID, &d.Name, &d.Engine, &d.Version, &d.NodeID,
			&d.Port, &d.SocketPath, &d.RAMProfile, &d.StorageLimitMB,
			&d.DataDirectory, &d.AdminUser, &d.EncryptedPassword, &d.Status,
			&d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan database: %w", err)
		}
		list = append(list, &d)
	}
	return list, nil
}

func (s *DatabaseStore) UpdateDatabaseStatus(dbID, status string) error {
	query := `UPDATE databases SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE db_id = ?;`
	_, err := s.db.Exec(query, status, dbID)
	return err
}

func (s *DatabaseStore) DeleteDatabase(dbID string) error {
	query := `DELETE FROM databases WHERE db_id = ?;`
	_, err := s.db.Exec(query, dbID)
	return err
}
