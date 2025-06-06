package store

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Activity struct {
	ID        int
	App       string
	StartTime time.Time
	EndTime   time.Time
}

type Store interface {
	Initialize() error
	SaveActivity(Activity) error
}

func Initialize() (Store, error) {
	db, err := NewSqliteStore("./shallowdepth.db")
	if err != nil {
		return nil, err
	}

	err = db.Initialize()
	if err != nil {
		return nil, err
	}

	return db, nil

}
