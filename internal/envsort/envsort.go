// Package envSort provides utilities for sorting environment variable maps
// by key, value, or custom criteria.
package envSort

import (
	"cmp"
	"maps"
	"slices"
	"strings"
)

// Order defines the sort direction.
type Order string

const (
	Asc  Order = "asc"
	Desc Order = "desc"
)

// Options controls how the sort is performed.
type Options struct {
	// By is the sort field: "key" (default) or "value".
	By string
	// Order is "asc" (default) or "desc".
	Order Order
	// CaseInsensitive performs a case-folded comparison when true.
	CaseInsensitive bool
}

// DefaultOptions returns the default sort options (by key, ascending).
func DefaultOptions() Options {
	return Options{By: "key", Order: Asc}
}

// SortedPair holds a key/value pair after sorting.
type SortedPair struct {
	Key   string
	Value string
}

// Apply sorts the given env map and returns an ordered slice of SortedPair.
// The original map is not mutated.
func Apply(env map[string]string, opts Options) []SortedPair {
	keys := slices.Collect(maps.Keys(env))

	cmpFn := func(a, b string) int {
		var ka, kb string
		if opts.By == "value" {
			ka, kb = env[a], env[b]
		} else {
			ka, kb = a, b
		}
		if opts.CaseInsensitive {
			ka = strings.ToLower(ka)
			kb = strings.ToLower(kb)
		}
		return cmp.Compare(ka, kb)
	}

	slices.SortStableFunc(keys, func(a, b string) int {
		v := cmpFn(a, b)
		if opts.Order == Desc {
			v = -v
		}
		return v
	})

	pairs := make([]SortedPair, len(keys))
	for i, k := range keys {
		pairs[i] = SortedPair{Key: k, Value: env[k]}
	}
	return pairs
}

// Keys returns only the sorted keys from the env map.
func Keys(env map[string]string, opts Options) []string {
	pairs := Apply(env, opts)
	keys := make([]string, len(pairs))
	for i, p := range pairs {
		keys[i] = p.Key
	}
	return keys
}
