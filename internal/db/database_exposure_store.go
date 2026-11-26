package db

import (
	"database/sql"
	"errors"
)

type DatabaseExposure struct {
	DBID       string `json:"db_id"`
	PublicPort int    `json:"public_port"`
}

type DatabaseExposureStore struct {
	db *DB
}

func NewDatabaseExposureStore(db *DB) *DatabaseExposureStore {
	return &DatabaseExposureStore{db: db}
}

func (s *DatabaseExposureStore) Get(dbID string) (*DatabaseExposure, error) {
	var e DatabaseExposure
	err := s.db.QueryRow(`SELECT db_id, public_port FROM database_exposures WHERE db_id = ?;`, dbID).Scan(&e.DBID, &e.PublicPort)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *DatabaseExposureStore) List() ([]*DatabaseExposure, error) {
	rows, err := s.db.Query(`SELECT db_id, public_port FROM database_exposures ORDER BY public_port;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*DatabaseExposure
	for rows.Next() {
		var e DatabaseExposure
		if err := rows.Scan(&e.DBID, &e.PublicPort); err != nil {
			return nil, err
		}
		out = append(out, &e)
	}
	return out, rows.Err()
}

func (s *DatabaseExposureStore) Set(dbID string, port int) error {
	_, err := s.db.Exec(`
	INSERT INTO database_exposures (db_id, public_port) VALUES (?, ?)
	ON CONFLICT(db_id) DO UPDATE SET public_port = excluded.public_port;`, dbID, port)
	return err
}

func (s *DatabaseExposureStore) Delete(dbID string) error {
	_, err := s.db.Exec(`DELETE FROM database_exposures WHERE db_id = ?;`, dbID)
	return err
}
