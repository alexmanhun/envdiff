package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of a parsed .env file.
type Snapshot struct {
	File      string            `json:"file"`
	CapturedAt time.Time        `json:"captured_at"`
	Keys      map[string]string `json:"keys"`
}

// Save writes a snapshot of the given key-value map to a JSON file at dest.
func Save(dest, sourceFile string, keys map[string]string) error {
	snap := Snapshot{
		File:       sourceFile,
		CapturedAt: time.Now().UTC(),
		Keys:       keys,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", dest, err)
	}

	return nil
}

// Load reads a snapshot from a JSON file at src.
func Load(src string) (*Snapshot, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", src, err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}

	return &snap, nil
}
