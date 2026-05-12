// Package groupby groups environment keys by a shared prefix or custom classifier.
package groupby

import (
	"sort"
	"strings"
)

// Group holds the keys and values that share a common label.
type Group struct {
	Label string
	Keys  map[string]string
}

// ByPrefix splits an env map into groups based on key prefixes separated by
// the given delimiter (e.g. "_"). Keys that do not match any prefix land in a
// group whose label is the empty string.
func ByPrefix(env map[string]string, delimiter string, prefixes []string) []Group {
	index := make(map[string]*Group)
	for _, p := range prefixes {
		index[p] = &Group{Label: p, Keys: make(map[string]string)}
	}
	other := &Group{Label: "", Keys: make(map[string]string)}

	for k, v := range env {
		matched := false
		for _, p := range prefixes {
			if strings.HasPrefix(k, p+delimiter) || k == p {
				index[p].Keys[k] = v
				matched = true
				break
			}
		}
		if !matched {
			other.Keys[k] = v
		}
	}

	result := make([]Group, 0, len(prefixes)+1)
	for _, p := range prefixes {
		result = append(result, *index[p])
	}
	if len(other.Keys) > 0 {
		result = append(result, *other)
	}
	return result
}

// Labels returns the sorted, non-empty labels present in a slice of groups.
func Labels(groups []Group) []string {
	out := make([]string, 0, len(groups))
	for _, g := range groups {
		if g.Label != "" {
			out = append(out, g.Label)
		}
	}
	sort.Strings(out)
	return out
}
