// Package scope provides utilities for grouping and partitioning env keys
// by a named scope (e.g. "APP_", "DB_", "AWS_") so callers can compare or
// report on isolated subsets of a larger env map.
package scope

import "strings"

// Group holds the keys that belong to a single scope.
type Group struct {
	Name string
	Keys map[string]string
}

// Partition splits env into one Group per prefix in scopes.
// Keys that match no prefix are collected under the empty-string group "".
// Matching is case-sensitive and uses strings.HasPrefix.
func Partition(env map[string]string, scopes []string) map[string]*Group {
	result := make(map[string]*Group, len(scopes)+1)

	for _, s := range scopes {
		result[s] = &Group{Name: s, Keys: make(map[string]string)}
	}
	// bucket for unmatched keys
	result[""] = &Group{Name: "", Keys: make(map[string]string)}

	for k, v := range env {
		matched := false
		for _, s := range scopes {
			if strings.HasPrefix(k, s) {
				result[s].Keys[k] = v
				matched = true
				break
			}
		}
		if !matched {
			result[""].Keys[k] = v
		}
	}

	return result
}

// Names returns the sorted list of non-empty scope names present in groups
// that contain at least one key.
func Names(groups map[string]*Group) []string {
	out := make([]string, 0, len(groups))
	for name, g := range groups {
		if name != "" && len(g.Keys) > 0 {
			out = append(out, name)
		}
	}
	sortStrings(out)
	return out
}

// sortStrings sorts a string slice in-place (avoids importing sort at top-level).
func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
