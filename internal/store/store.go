package store

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ActivitySession struct {
	ID              string
	App             string
	WindowTitle     string
	StartTime       time.Time
	LastActivity    time.Time
	EndTime         time.Time
	DurationSeconds time.Duration
}

type Store interface {
	Initialize() error
	SaveActivitySession(ActivitySession) error
}

func Initialize() (Store, error) {
	db, err := NewSqliteStore("shallowdepth.db")
	if err != nil {
		return nil, err
	}

	err = db.Initialize()
	if err != nil {
		return nil, err
	}

	return db, nil
}
