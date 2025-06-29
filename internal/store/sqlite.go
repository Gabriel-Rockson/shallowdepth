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
		window_title TEXT,
		start_time TIMESTAMP NOT NULL,
		last_activity TIMESTAMP,
		end_time DATETIME,
		duration_seconds INTEGER
	);`

	_, err := s.db.Exec(query)
	return err
}

func (s *SqliteStore) SaveActivitySession(a ActivitySession) error {
	query := `
	INSERT INTO activity (app, window_title, start_time, last_activity, end_time, duration_seconds)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, a.App, a.WindowTitle, a.StartTime, a.LastActivity, a.EndTime, a.DurationSeconds)
	if err != nil {
		return fmt.Errorf("failed to save activity %s: %w", a.App, err)
	}
	slog.Info("successfully saved activity", "App", a.App, "Window", a.WindowTitle)
	return nil
}
