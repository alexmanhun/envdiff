// Package audit records a history of envdiff comparison runs.
package audit

import (
	"encoding/json"
	"os"
	"time"
)

// Entry represents a single recorded comparison run.
type Entry struct {
	Timestamp   time.Time `json:"timestamp"`
	LeftFile    string    `json:"left_file"`
	RightFile   string    `json:"right_file"`
	Missing     int       `json:"missing_keys"`
	Mismatched  int       `json:"mismatched_keys"`
	HasDiff     bool      `json:"has_diff"`
}

// Log holds a list of audit entries.
type Log struct {
	Entries []Entry `json:"entries"`
}

// Append loads the log at path (creating it if absent), appends e, and saves it.
func Append(path string, e Entry) error {
	l, err := Load(path)
	if err != nil {
		return err
	}
	l.Entries = append(l.Entries, e)
	return save(path, l)
}

// Load reads the audit log from path. Returns an empty Log if the file does
// not exist.
func Load(path string) (Log, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Log{}, nil
	}
	if err != nil {
		return Log{}, err
	}
	var l Log
	if err := json.Unmarshal(data, &l); err != nil {
		return Log{}, err
	}
	return l, nil
}

func save(path string, l Log) error {
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
