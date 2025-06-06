package tracker

import (
	"fmt"
	"math/rand/v2"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Gabriel-Rockson/shallowdepth/internal/store"
)

type Tracker struct {
	store store.Store
}

func NewTracker(s store.Store) *Tracker {
	return &Tracker{store: s}
}

// -> First I want to get the current opened application
// -> How do I deal with a scenario that many apps are opened?
// -> One thing I know is that at anytime, multiple apps can be opened
// but only one is going to have focus.
// - is there a way to determine if a window has focus?
// -> What about when a user is using 25 :) monitors? How do I
// tell how many windows are opened across all monitors?

func getCurrentFocusedApp() string {
	cmd := exec.Command("xdotool", "getactivewindow")
	winIDBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("There was an error executing the commands")
		return ""
	}
	winID := strings.TrimSpace(string(winIDBytes))

	cmd = exec.Command("xprop", "-id", winID, "_NET_WM_PID")
	xpropOutBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting window PID: ", err)
		return ""
	}

	// Exmaple output
	// _NET_WM_PID(CARDINAL) = 3599996
	xpropOut := string(xpropOutBytes)
	parts := strings.Split(xpropOut, "=")
	if len(parts) != 2 {
		fmt.Println("Unexpected error from xprop output: ", xpropOut)
		return ""
	}

	pidStr := strings.TrimSpace(parts[1])
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Println("Error parsing the PID: ", err)
		return ""
	}

	cmd = exec.Command("ps", "-p", pidStr, "-o", "comm=")
	procNameBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("An error occurred getting process name: ", err)
	}

	procName := strings.TrimSpace(string(procNameBytes))

	fmt.Printf("Active window PID: %d, Process name: %s\n", pid, procName)

	return procName

}

// Start will commence the tracking processes to know about the currently focused apps and the urls being visited
func (t *Tracker) Start() {
	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			app := getCurrentFocusedApp()
			activity := store.Activity{
				ID:        int(rand.Int64()),
				App:       app,
				StartTime: time.Now(),
				EndTime:   time.Now(),
			}
			t.store.SaveActivity(activity)
		}
	}
}
