package store

import (
	"database/sql"
	"fmt"
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(path string) (*SqliteStore, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &SqliteStore{db: conn}, nil
}

func (s *SqliteStore) Initialize() error {
	query := `
	CREATE TABLE IF NOT EXISTS activity (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		app TEXT,
		start_time DATETIME,
		end_time DATETIME
	);`

	_, err := s.db.Exec(query)
	return err
}

func (s *SqliteStore) SaveActivity(a Activity) error {
	query := `
	INSERT INTO activity (id, app, start_time, end_time)
	VALUES (?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, a.ID, a.App, a.StartTime, a.EndTime)
	if err != nil {
		return fmt.Errorf("failed to save activity %s: %w", a.App, err)
	}
	return nil
}
