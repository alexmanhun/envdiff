// Package clone provides functionality to copy an env map,
// optionally stripping or overriding specific keys.
package clone

// Options controls how the clone operation behaves.
type Options struct {
	// StripKeys removes these keys from the cloned output.
	StripKeys []string
	// Overrides replaces values for the given keys in the cloned output.
	Overrides map[string]string
}

// Clone returns a deep copy of src, applying any Options provided.
// If opts is nil, a plain copy is returned.
func Clone(src map[string]string, opts *Options) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}

	if opts == nil {
		return dst
	}

	for _, key := range opts.StripKeys {
		delete(dst, key)
	}

	for k, v := range opts.Overrides {
		dst[k] = v
	}

	return dst
}
