package filter

import (
	"regexp"
	"strings"
)

// Options holds filtering configuration.
type Options struct {
	// Prefix filters keys to only those starting with the given prefix.
	Prefix string
	// Pattern filters keys matching the given regex pattern.
	Pattern string
	// ExcludeKeys is a set of exact key names to exclude.
	ExcludeKeys []string
}

// Apply filters a map of env key-value pairs based on the provided Options.
// It returns a new map containing only the keys that pass all active filters.
func Apply(env map[string]string, opts Options) (map[string]string, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}

	excluded := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excluded[k] = struct{}{}
	}

	result := make(map[string]string)
	for k, v := range env {
		if _, skip := excluded[k]; skip {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if re != nil && !re.MatchString(k) {
			continue
		}
		result[k] = v
	}
	return result, nil
}
