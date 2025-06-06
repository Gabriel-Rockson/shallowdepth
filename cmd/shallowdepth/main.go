package main

import (
	"log/slog"

	"github.com/Gabriel-Rockson/shallowdepth/internal/store"
	"github.com/Gabriel-Rockson/shallowdepth/internal/tracker"
)

func main() {
	s, err := store.Initialize()
	if err != nil {
		slog.Error("error occurred initializing database: ", "error", err)
	}

	t := tracker.NewTracker(s)
	t.Start()
}
