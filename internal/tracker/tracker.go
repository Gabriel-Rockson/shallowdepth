package tracker

import (
	"errors"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Gabriel-Rockson/shallowdepth/internal/store"
	"github.com/Gabriel-Rockson/shallowdepth/pkg/utils"
)

type Tracker struct {
	store store.Store
}

func NewTracker(s store.Store) *Tracker {
	return &Tracker{store: s}
}

func runCommandOutput(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// -> First I want to get the current opened application
// -> How do I deal with a scenario that many apps are opened?
// -> One thing I know is that at anytime, multiple apps can be opened
// but only one is going to have focus.
// - is there a way to determine if a window has focus?
// -> What about when a user is using 25 :) monitors? How do I
// tell how many windows are opened across all monitors?

func getCurrentFocusedApp() (string, error) {
	winID, err := runCommandOutput("xdotool", "getactivewindow")
	if err != nil {
		slog.Error("failed to get active window ID", "error", err)
		return "", err
	}

	xpropOut, err := runCommandOutput("xprop", "-id", winID, "_NET_WM_PID")
	if err != nil {
		slog.Error("failed to get _NET_WM_PID", "error", err)
		return "", err
	}

	parts := strings.Split(xpropOut, "=")
	if len(parts) != 2 {
		err := errors.New("malformed xprop output")
		slog.Error("unexpected xprop output", "output", xpropOut, "error", err)
		return "", err
	}

	pidStr := strings.TrimSpace(parts[1])
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		slog.Error("failed to parse PID", "error", err)
		return "", err
	}

	procName, err := runCommandOutput("ps", "-p", pidStr, "-o", "comm=")
	if err != nil {
		slog.Error("failed to get process name", "PID", pid, "error", err)
		return "", err
	}

	procName = utils.Capitalize(procName)
	slog.Info("current active window", "PID", pid, "Process", procName)
	return procName, nil
}

// Start will commence the tracking processes to know about the currently focused apps and the urls being visited
func (t *Tracker) Start() {
	ticker := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-ticker.C:
			app, err := getCurrentFocusedApp()
			if err != nil {
				slog.Error("an error occurred getting the current active app", "error", err)
			}

			activity := store.Activity{
				App:       app,
				Timestamp: time.Now().UTC(),
			}
			err = t.store.SaveActivity(activity)
			if err != nil {
				slog.Error("an error occurred saving the activity", "error", err)
			}
		}
	}
}
