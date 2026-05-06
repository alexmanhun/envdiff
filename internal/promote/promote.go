// Package promote copies keys from a source environment map into a
// destination map, with configurable overwrite behaviour.
package promote

import "fmt"

// Strategy controls what happens when a key already exists in dst.
type Strategy int

const (
	// SkipExisting leaves keys that already exist in dst unchanged.
	SkipExisting Strategy = iota
	// Overwrite replaces keys that already exist in dst.
	Overwrite
	// ErrorOnConflict returns an error if a key exists in both maps with
	// different values.
	ErrorOnConflict
)

// Result holds the outcome of a Promote operation.
type Result struct {
	// Promoted contains keys that were written into dst.
	Promoted []string
	// Skipped contains keys that were left unchanged (SkipExisting).
	Skipped []string
	// Overwritten contains keys whose values were replaced (Overwrite).
	Overwritten []string
}

// Promote copies keys from src into dst according to strategy.
// It never mutates src. The returned Result describes what happened.
func Promote(src, dst map[string]string, strategy Strategy) (map[string]string, Result, error) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var res Result
	for k, srcVal := range src {
		dstVal, exists := out[k]
		switch {
		case !exists:
			out[k] = srcVal
			res.Promoted = append(res.Promoted, k)
		case strategy == SkipExisting:
			res.Skipped = append(res.Skipped, k)
		case strategy == Overwrite:
			out[k] = srcVal
			res.Overwritten = append(res.Overwritten, k)
		case strategy == ErrorOnConflict:
			if srcVal != dstVal {
				return nil, Result{}, fmt.Errorf(
					"promote: conflict on key %q: src=%q dst=%q", k, srcVal, dstVal,
				)
			}
			// same value — treat as promoted without conflict
			res.Promoted = append(res.Promoted, k)
		}
	}
	return out, res, nil
}
