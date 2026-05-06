// Package drift detects when a live environment's keys have drifted
// from a saved snapshot, reporting new, removed, or changed keys.
package drift

import (
	"fmt"

	"github.com/user/envdiff/internal/snapshot"
)

// Report holds the result of a drift detection run.
type Report struct {
	Added   map[string]string // keys present in live but not in snapshot
	Removed map[string]string // keys present in snapshot but not in live
	Changed map[string][2]string // key -> [snapshotVal, liveVal]
}

// HasDrift returns true if any drift was detected.
func (r Report) HasDrift() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0
}

// Detect compares a live environment map against a saved snapshot file.
// It returns a Report describing any drift found.
func Detect(snapshotPath string, live map[string]string) (Report, error) {
	snap, err := snapshot.Load(snapshotPath)
	if err != nil {
		return Report{}, fmt.Errorf("drift: loading snapshot: %w", err)
	}

	report := Report{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	// Keys in live but not in snapshot, or changed values.
	for k, liveVal := range live {
		snapVal, ok := snap[k]
		if !ok {
			report.Added[k] = liveVal
		} else if snapVal != liveVal {
			report.Changed[k] = [2]string{snapVal, liveVal}
		}
	}

	// Keys in snapshot but not in live.
	for k, snapVal := range snap {
		if _, ok := live[k]; !ok {
			report.Removed[k] = snapVal
		}
	}

	return report, nil
}
