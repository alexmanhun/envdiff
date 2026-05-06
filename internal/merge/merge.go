// Package merge provides functionality to merge multiple .env maps into one,
// with configurable conflict resolution strategies.
package merge

import "fmt"

// Strategy defines how key conflicts are resolved during a merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast
	// StrategyError returns an error if the same key appears with different values.
	StrategyError
)

// Conflict records a key that appeared in multiple sources with differing values.
type Conflict struct {
	Key    string
	Values []string // one entry per source, in order
}

// Result holds the merged environment map and any conflicts that were detected.
type Result struct {
	Env       map[string]string
	Conflicts []Conflict
}

// Merge combines multiple env maps into a single map using the given strategy.
// sources is a slice of maps in priority order (index 0 = highest for StrategyFirst,
// lowest for StrategyLast).
func Merge(sources []map[string]string, strategy Strategy) (Result, error) {
	merged := make(map[string]string)
	// track which values each key has seen, preserving order
	seen := make(map[string][]string)

	for _, src := range sources {
		for k, v := range src {
			seen[k] = append(seen[k], v)
		}
	}

	var conflicts []Conflict

	for k, vals := range seen {
		unique := uniqueValues(vals)
		if len(unique) > 1 {
			if strategy == StrategyError {
				return Result{}, fmt.Errorf("merge conflict on key %q: values %v", k, unique)
			}
			conflicts = append(conflicts, Conflict{Key: k, Values: vals})
		}

		switch strategy {
		case StrategyFirst:
			merged[k] = vals[0]
		case StrategyLast:
			merged[k] = vals[len(vals)-1]
		default:
			merged[k] = vals[0]
		}
	}

	return Result{Env: merged, Conflicts: conflicts}, nil
}

func uniqueValues(vals []string) []string {
	seen := make(map[string]struct{})
	out := []string{}
	for _, v := range vals {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}
