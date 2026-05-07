// Package mask provides utilities for partially obscuring env values
// in output, useful when sharing diffs without exposing secrets.
package mask

import "strings"

// Mode controls how values are masked.
type Mode int

const (
	// ModeNone leaves values unchanged.
	ModeNone Mode = iota
	// ModePartial shows the first N characters followed by asterisks.
	ModePartial
	// ModeFull replaces the entire value with asterisks.
	ModeFull
)

// Options configures masking behaviour.
type Options struct {
	Mode       Mode
	// VisibleChars is the number of leading characters to keep in ModePartial.
	VisibleChars int
	// Placeholder is used for ModeFull; defaults to "****".
	Placeholder string
}

// DefaultOptions returns sensible masking defaults.
func DefaultOptions() Options {
	return Options{
		Mode:         ModePartial,
		VisibleChars: 3,
		Placeholder:  "****",
	}
}

// Value masks a single string value according to opts.
func Value(v string, opts Options) string {
	if v == "" {
		return v
	}
	switch opts.Mode {
	case ModeFull:
		p := opts.Placeholder
		if p == "" {
			p = "****"
		}
		return p
	case ModePartial:
		n := opts.VisibleChars
		if n <= 0 {
			n = 3
		}
		if len(v) <= n {
			return strings.Repeat("*", len(v))
		}
		return v[:n] + strings.Repeat("*", len(v)-n)
	default:
		return v
	}
}

// Apply masks all values in the provided env map, returning a new map.
// Keys are never modified.
func Apply(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = Value(v, opts)
	}
	return out
}
