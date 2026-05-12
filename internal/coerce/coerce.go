// Package coerce normalises env map values to consistent types or formats.
// For example, it can upper-case boolean-like values, trim surrounding
// whitespace, or expand common shorthands ("yes" → "true").
package coerce

import (
	"strings"
)

// Options controls which coercions are applied.
type Options struct {
	// NormalizeBools converts yes/no/on/off to true/false.
	NormalizeBools bool
	// TrimSpace strips leading and trailing whitespace from every value.
	TrimSpace bool
	// LowercaseBoolValues lowercases the final boolean value (true/false).
	LowercaseBoolValues bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		NormalizeBools:      true,
		TrimSpace:           true,
		LowercaseBoolValues: true,
	}
}

var boolAliases = map[string]string{
	"yes":   "true",
	"no":    "false",
	"on":    "true",
	"off":   "false",
	"1":     "true",
	"0":     "false",
	"true":  "true",
	"false": "false",
}

// Apply returns a new map with coercions applied according to opts.
// The original map is never mutated.
func Apply(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.TrimSpace {
			v = strings.TrimSpace(v)
		}
		if opts.NormalizeBools {
			if canonical, ok := boolAliases[strings.ToLower(v)]; ok {
				if opts.LowercaseBoolValues {
					v = canonical
				} else {
					v = strings.ToUpper(canonical)
				}
			}
		}
		out[k] = v
	}
	return out
}

// Count returns the number of values that were changed by Apply.
func Count(before, after map[string]string) int {
	n := 0
	for k, v := range before {
		if after[k] != v {
			n++
		}
	}
	return n
}
