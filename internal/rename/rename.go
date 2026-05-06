// Package rename provides utilities for renaming keys across env file maps.
// It supports bulk renaming via a mapping of old->new key names, with options
// to control behavior when the destination key already exists.
package rename

import "fmt"

// Strategy controls what happens when the destination key already exists.
type Strategy int

const (
	// SkipExisting leaves the destination key untouched if it already exists.
	SkipExisting Strategy = iota
	// OverwriteExisting replaces the destination key value with the source value.
	OverwriteExisting
	// ErrorOnExisting returns an error if the destination key already exists.
	ErrorOnExisting
)

// Result holds the outcome of a Rename operation.
type Result struct {
	// Env is the resulting key/value map after renaming.
	Env map[string]string
	// Skipped contains old keys that were skipped because the destination existed.
	Skipped []string
	// Renamed contains old keys that were successfully renamed.
	Renamed []string
}

// Rename applies the provided mapping (oldKey -> newKey) to env.
// Keys not present in mapping are left unchanged.
// Keys in mapping that do not exist in env are silently ignored.
func Rename(env map[string]string, mapping map[string]string, strategy Strategy) (*Result, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	result := &Result{Env: out}

	for oldKey, newKey := range mapping {
		val, exists := out[oldKey]
		if !exists {
			continue
		}

		_, destExists := out[newKey]
		if destExists && oldKey != newKey {
			switch strategy {
			case ErrorOnExisting:
				return nil, fmt.Errorf("rename: destination key %q already exists", newKey)
			case SkipExisting:
				result.Skipped = append(result.Skipped, oldKey)
				continue
			case OverwriteExisting:
				// fall through to rename
			}
		}

		if oldKey != newKey {
			delete(out, oldKey)
		}
		out[newKey] = val
		result.Renamed = append(result.Renamed, oldKey)
	}

	return result, nil
}
