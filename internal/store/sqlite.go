package store

import (
	"database/sql"
	"fmt"
	"log/slog"
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
		timestamp DATETIME
	);`

	_, err := s.db.Exec(query)
	return err
}

func (s *SqliteStore) SaveActivity(a Activity) error {
	query := `
	INSERT INTO activity (app, timestamp)
	VALUES (?, ?)
	`
	_, err := s.db.Exec(query, a.App, a.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to save activity %s: %w", a.App, err)
	}
	slog.Info("successfully saved activity", "App", a.App, "Timestamp", a.Timestamp)
	return nil
}
