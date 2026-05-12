// Package required checks that a set of keys are present and non-empty
// in one or more env maps. It is useful for enforcing mandatory configuration
// before an application starts or before promoting an env file.
package required

import "fmt"

// Violation describes a single missing or empty required key.
type Violation struct {
	Key     string
	EnvName string
	Reason  string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: key %q %s", v.EnvName, v.Key, v.Reason)
}

// Check verifies that every key in keys exists and is non-empty in each of
// the supplied envs. envNames provides a human-readable label for each env
// and must have the same length as envs. If envNames is nil, numeric labels
// are used.
//
// All violations are collected and returned; the caller decides whether to
// treat them as fatal.
func Check(keys []string, envs []map[string]string, envNames []string) []Violation {
	var violations []Violation

	for i, env := range envs {
		name := label(envNames, i)
		for _, k := range keys {
			v, ok := env[k]
			if !ok {
				violations = append(violations, Violation{
					Key:     k,
					EnvName: name,
					Reason:  "is missing",
				})
			} else if v == "" {
				violations = append(violations, Violation{
					Key:     k,
					EnvName: name,
					Reason:  "is empty",
				})
			}
		}
	}

	return violations
}

// label returns envNames[i] when available, otherwise "env<i>".
func label(names []string, i int) string {
	if i < len(names) && names[i] != "" {
		return names[i]
	}
	return fmt.Sprintf("env%d", i)
}
