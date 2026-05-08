// Package trim removes keys from one or more env maps whose values match
// a given set of conditions (empty, whitespace-only, or a custom predicate).
package trim

import "strings"

// Options controls which entries are removed by Apply.
type Options struct {
	// RemoveEmpty removes keys whose value is the empty string.
	RemoveEmpty bool
	// RemoveWhitespace removes keys whose value is entirely whitespace.
	RemoveWhitespace bool
	// Keys is an explicit list of key names to always remove.
	Keys []string
}

// DefaultOptions returns an Options that removes empty and whitespace-only values.
func DefaultOptions() Options {
	return Options{
		RemoveEmpty:      true,
		RemoveWhitespace: true,
	}
}

// Apply returns a new map with entries removed according to opts.
// The original map is never mutated.
func Apply(env map[string]string, opts Options) map[string]string {
	explicit := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[k] = struct{}{}
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if _, drop := explicit[k]; drop {
			continue
		}
		if opts.RemoveEmpty && v == "" {
			continue
		}
		if opts.RemoveWhitespace && strings.TrimSpace(v) == "" {
			continue
		}
		out[k] = v
	}
	return out
}

// Count returns the number of entries that would be removed from env by Apply.
func Count(env map[string]string, opts Options) int {
	return len(env) - len(Apply(env, opts))
}
