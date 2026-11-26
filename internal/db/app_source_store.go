package db

import (
	"database/sql"
	"time"
)

type AppSource struct {
	AppID           string     `json:"app_id"`
	RepoURL         string     `json:"repo_url"`
	GitRef          string     `json:"git_ref"`
	InstalledCommit string     `json:"installed_commit"`
	LatestCommit    string     `json:"latest_commit"`
	CheckedAt       *time.Time `json:"checked_at,omitempty"`
}

type AppSourceStore struct {
	db *DB
}

func NewAppSourceStore(db *DB) *AppSourceStore {
	return &AppSourceStore{db: db}
}

func (s *AppSourceStore) Upsert(src *AppSource) error {
	_, err := s.db.Exec(`
	INSERT INTO app_sources (app_id, repo_url, git_ref, installed_commit, latest_commit)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(app_id) DO UPDATE SET
		repo_url = excluded.repo_url,
		git_ref = excluded.git_ref,
		installed_commit = excluded.installed_commit,
		latest_commit = excluded.latest_commit;`,
		src.AppID, src.RepoURL, src.GitRef, src.InstalledCommit, src.LatestCommit)
	return err
}

func (s *AppSourceStore) Get(appID string) (*AppSource, error) {
	var src AppSource
	var checked sql.NullTime
	err := s.db.QueryRow(`
	SELECT app_id, repo_url, git_ref, installed_commit, latest_commit, checked_at
	FROM app_sources WHERE app_id = ?;`, appID).
		Scan(&src.AppID, &src.RepoURL, &src.GitRef, &src.InstalledCommit, &src.LatestCommit, &checked)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if checked.Valid {
		t := checked.Time
		src.CheckedAt = &t
	}
	return &src, nil
}

func (s *AppSourceStore) List() ([]*AppSource, error) {
	rows, err := s.db.Query(`SELECT app_id FROM app_sources ORDER BY app_id;`)
	if err != nil {
		return nil, err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()

	var out []*AppSource
	for _, id := range ids {
		src, err := s.Get(id)
		if err != nil {
			return nil, err
		}
		if src != nil {
			out = append(out, src)
		}
	}
	return out, nil
}

func (s *AppSourceStore) SetLatest(appID, commit string) error {
	_, err := s.db.Exec(`UPDATE app_sources SET latest_commit = ?, checked_at = CURRENT_TIMESTAMP WHERE app_id = ?;`, commit, appID)
	return err
}

func (s *AppSourceStore) SetInstalled(appID, commit string) error {
	_, err := s.db.Exec(`UPDATE app_sources SET installed_commit = ?, latest_commit = ? WHERE app_id = ?;`, commit, commit, appID)
	return err
}

// SetEntrypoint updates the app's binary entrypoint from the repo manifest.
func (s *AppSourceStore) SetEntrypoint(appID, entrypoint string) error {
	_, err := s.db.Exec(`UPDATE apps SET binary_entrypoint = ?, updated_at = CURRENT_TIMESTAMP WHERE app_id = ?;`, entrypoint, appID)
	return err
}
