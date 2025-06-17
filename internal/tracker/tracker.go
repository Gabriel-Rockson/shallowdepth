package tracker

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-vgo/robotgo"

	"github.com/Gabriel-Rockson/shallowdepth/internal/store"
)

type Tracker struct {
	store        store.Store
	lastActivity time.Time
	lastApp      string
	lastPosition robotgo.Point
}

func NewTracker(s store.Store) *Tracker {
	return &Tracker{store: s, lastActivity: time.Now()}
}

// getCurrentFocusedApp gets the currently active window using robotgo
func getCurrentFocusedApp() (string, string, error) {
	pid := robotgo.GetPid()
	if pid == 0 {
		return "", "", fmt.Errorf("failed to get active window PID")
	}

	windowTitle := robotgo.GetTitle()

	appName, err := robotgo.FindName(pid)
	if err != nil {
		return "", "", fmt.Errorf("failed to get process name")

		// Add fallback to get the window title a different way if this one fails
	}

	slog.Info("current active window", "pid", pid, "process", appName, "title", windowTitle)

	return appName, windowTitle, nil
}

func (t *Tracker) detectIdleTime() time.Duration {
	x, y := robotgo.Location()
	currPos := robotgo.Point{
		X: x,
		Y: y,
	}
	now := time.Now()

	if currPos.X != t.lastPosition.X || currPos.Y != t.lastPosition.Y {
		t.lastActivity = now
		t.lastPosition = currPos
	}

	// TODO: What if the user hasn't moved the mouse but has been typing a lot, like how vim users will
	// naturally work?
	// I need to detect keyboard inputs, I don't need to register the input itself, but listen for keydown
	// events and mark the user as active
	return now.Sub(t.lastActivity)
}

func (t *Tracker) isUserActive() bool {
	idleTime := t.detectIdleTime()

	// configurable threshold, you'd want the user to define this
	idleThreshold := 5 * time.Minute

	return idleTime < idleThreshold
}

// Start will commence the tracking processes to know about the currently focused apps and the urls being visited
func (t *Tracker) Start() {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	var currSess *store.ActivitySession
	sessTimeout := 30 * time.Second

	for {
		select {
		case <-ticker.C:
			appName, windowTitle, err := getCurrentFocusedApp()
			if err != nil {
				slog.Error("an error occurred getting the current active app", "error", err)
				continue
			}

			now := time.Now().UTC()
			isActive := t.isUserActive()

			// Handle the session logic
			shouldStartNewSess := currSess == nil ||
				currSess.App != appName ||
				(!isActive && now.Sub(currSess.LastActivity) > sessTimeout)
			if shouldStartNewSess {
				// clean up old session
				if currSess != nil {
					currSess.EndTime = now
					currSess.DurationSeconds = time.Duration(now.Sub(currSess.StartTime).Seconds())
					err := t.store.SaveActivitySession(*currSess)
					if err != nil {
						slog.Error("an error occurred saving the activity", "error", err)
					}
					slog.Info("ended activity session", "app", currSess.App, "duration", now.Sub(currSess.StartTime))
				}

				// if the user is active, then start a new session
				if isActive {
					currSess = &store.ActivitySession{
						App:          appName,
						WindowTitle:  windowTitle,
						StartTime:    now,
						LastActivity: now,
					}
				} else {
					currSess = nil
				}
			} else if isActive && currSess != nil {
				currSess.LastActivity = now
				currSess.WindowTitle = windowTitle // in-case window title has changed
			}
		}
	}
}

// TODO: I want to check the idleness of the user, with this; I am considering using a tool https://pkg.go.dev/github.com/go-vgo/robotgo
// this tool is a Golang Desktop Automation tool that will work for Windows, Linux (X11) and Mac, and maybe this will even help with getting
// to know the current window that the user is interacting with.
