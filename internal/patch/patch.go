// Package patch applies a set of key-value changes to an existing env map,
// producing a new map without mutating the original.
package patch

import "fmt"

// Op represents the type of patch operation.
type Op string

const (
	OpSet    Op = "set"    // add or overwrite a key
	OpDelete Op = "delete" // remove a key
)

// Change describes a single patch operation.
type Change struct {
	Op    Op
	Key   string
	Value string // only used for OpSet
}

// Result holds the output of Apply.
type Result struct {
	Env     map[string]string
	Applied []Change // changes that were actually applied
	Skipped []Change // changes skipped (e.g. delete of non-existent key)
}

// Apply applies the given changes to src and returns a Result.
// src is never mutated; a shallow copy is made first.
func Apply(src map[string]string, changes []Change) (*Result, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	result := &Result{Env: out}

	for _, c := range changes {
		if c.Key == "" {
			return nil, fmt.Errorf("patch: change has empty key (op=%s)", c.Op)
		}
		switch c.Op {
		case OpSet:
			out[c.Key] = c.Value
			result.Applied = append(result.Applied, c)
		case OpDelete:
			if _, exists := out[c.Key]; !exists {
				result.Skipped = append(result.Skipped, c)
				continue
			}
			delete(out, c.Key)
			result.Applied = append(result.Applied, c)
		default:
			return nil, fmt.Errorf("patch: unknown op %q for key %q", c.Op, c.Key)
		}
	}

	return result, nil
}
