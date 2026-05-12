// Package defaults fills missing keys in a target env map using values
// from one or more fallback env maps, in priority order.
package defaults

// Options controls how defaults are applied.
type Options struct {
	// Overwrite replaces existing non-empty values in dst with fallback values.
	// By default only missing (absent) keys are filled.
	Overwrite bool

	// SkipEmpty skips fallback values that are the empty string.
	SkipEmpty bool
}

// DefaultOptions returns the recommended zero-value Options.
func DefaultOptions() Options {
	return Options{}
}

// Apply fills keys that are absent (or empty, when Overwrite is true) in dst
// using the first fallback that provides a non-empty value.
// dst is never mutated; a new map is returned.
func Apply(dst map[string]string, opts Options, fallbacks ...map[string]string) map[string]string {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	for _, fb := range fallbacks {
		for k, fv := range fb {
			if opts.SkipEmpty && fv == "" {
				continue
			}
			existing, present := out[k]
			if !present || (opts.Overwrite && existing == "") {
				out[k] = fv
			}
		}
	}

	return out
}

// Missing returns the set of keys that are present in reference but absent
// from env (or present with an empty value).
func Missing(env map[string]string, reference map[string]string) []string {
	var keys []string
	for k := range reference {
		v, ok := env[k]
		if !ok || v == "" {
			keys = append(keys, k)
		}
	}
	sort(keys)
	return keys
}

func sort(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
