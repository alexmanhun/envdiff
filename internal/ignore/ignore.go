// Package ignore provides functionality to load and apply .envdiffignore files,
// allowing users to suppress specific keys from diff output.
package ignore

import (
	"bufio"
	"os"
	"strings"
)

// Rules holds a set of key names that should be ignored during comparison.
type Rules struct {
	keys map[string]struct{}
}

// LoadFile reads an ignore file (one key per line, # for comments)
// and returns a Rules instance. If the file does not exist, an empty
// Rules is returned without error.
func LoadFile(path string) (*Rules, error) {
	r := &Rules{keys: make(map[string]struct{})}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return r, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		r.keys[line] = struct{}{}
	}
	return r, scanner.Err()
}

// NewRules constructs a Rules from a slice of key names.
func NewRules(keys []string) *Rules {
	r := &Rules{keys: make(map[string]struct{}, len(keys))}
	for _, k := range keys {
		r.keys[k] = struct{}{}
	}
	return r
}

// Contains reports whether the given key is ignored.
func (r *Rules) Contains(key string) bool {
	_, ok := r.keys[key]
	return ok
}

// Apply removes ignored keys from the provided env map, returning a new map.
func (r *Rules) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if !r.Contains(k) {
			out[k] = v
		}
	}
	return out
}
