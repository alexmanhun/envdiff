// Package template generates a .env.example file from one or more parsed
// env maps, replacing values with placeholder descriptions.
package template

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Options controls how the template is rendered.
type Options struct {
	// Placeholder is written as the value for every key.
	// Defaults to "<value>" if empty.
	Placeholder string

	// IncludeValues keeps the original value as a comment above each key.
	IncludeValues bool
}

// Generate writes an .env.example template to w based on the supplied env maps.
// Keys present in any of the maps are included; values are replaced with the
// configured placeholder.
func Generate(w io.Writer, envs []map[string]string, opts Options) error {
	if opts.Placeholder == "" {
		opts.Placeholder = "<value>"
	}

	merged := mergeKeys(envs)

	for _, key := range merged {
		if opts.IncludeValues {
			originals := collectValues(key, envs)
			if len(originals) > 0 {
				fmt.Fprintf(w, "# example: %s\n", strings.Join(originals, " | "))
			}
		}
		fmt.Fprintf(w, "%s=%s\n", key, opts.Placeholder)
	}
	return nil
}

// mergeKeys returns a sorted, deduplicated list of all keys across envs.
func mergeKeys(envs []map[string]string) []string {
	seen := make(map[string]struct{})
	for _, env := range envs {
		for k := range env {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// collectValues gathers non-empty values for key across all envs.
func collectValues(key string, envs []map[string]string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, env := range envs {
		if v, ok := env[key]; ok && v != "" {
			if _, dup := seen[v]; !dup {
				out = append(out, v)
				seen[v] = struct{}{}
			}
		}
	}
	return out
}
