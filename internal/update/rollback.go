package update

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// StateFile tracks whether the running binary has proven healthy since the last update.
const StateFile = "update-state.json"

// maxStartsBeforeRollback is how many starts a new binary gets before the agent restores the
// previous one. A start counts as healthy once it has run for healthyAfter.
const (
	maxStartsBeforeRollback = 3
	healthyAfter            = 60 * time.Second
)

type state struct {
	Status   string `json:"status"` // "pending" (new binary not yet proven) or "ok"
	Attempts int    `json:"attempts"`
	Version  string `json:"version"`
}

func readState(dir string) (state, error) {
	var s state
	raw, err := os.ReadFile(filepath.Join(dir, StateFile))
	if errors.Is(err, os.ErrNotExist) {
		return state{Status: "ok"}, nil
	}
	if err != nil {
		return s, err
	}
	return s, json.Unmarshal(raw, &s)
}

func writeState(dir string, s state) error {
	raw, err := json.Marshal(s)
	if err != nil {
		return err
	}
	tmp := filepath.Join(dir, StateFile+".tmp")
	if err := os.WriteFile(tmp, raw, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, StateFile))
}

// MarkPending records that exePath was just installed at version and has not proven itself.
func MarkPending(dir, version string) error {
	return writeState(dir, state{Status: "pending", Attempts: 0, Version: version})
}

// StartupCheck runs when the agent boots. If a new binary has failed to stay up several times,
// it restores exePath+".prev" and returns restored=true. The caller should exit so systemd starts
// the restored binary. Otherwise it counts this start and returns false.
func StartupCheck(dir, exePath string) (restored bool, err error) {
	s, err := readState(dir)
	if err != nil || s.Status != "pending" {
		return false, err
	}
	s.Attempts++
	if s.Attempts <= maxStartsBeforeRollback {
		return false, writeState(dir, s)
	}

	prev := exePath + ".prev"
	if _, err := os.Stat(prev); err != nil {
		return false, errors.New("update failed repeatedly but no previous binary exists")
	}
	bad := exePath + ".bad"
	_ = os.Remove(bad)
	if err := os.Rename(exePath, bad); err != nil {
		return false, err
	}
	if err := os.Rename(prev, exePath); err != nil {
		return false, err
	}
	return true, writeState(dir, state{Status: "ok"})
}

// MarkHealthyAfterDelay clears the pending state once the process has stayed up for healthyAfter.
func MarkHealthyAfterDelay(dir string) {
	go func() {
		time.Sleep(healthyAfter)
		s, err := readState(dir)
		if err != nil || s.Status != "pending" {
			return
		}
		_ = writeState(dir, state{Status: "ok", Version: s.Version})
	}()
}
