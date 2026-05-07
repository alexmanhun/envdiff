// Package intersect provides utilities for finding common keys across
// multiple environment maps and computing their value agreement.
package intersect

// Result holds the output of an intersection operation.
type Result struct {
	// Common contains keys present in ALL provided envs.
	Common map[string]string
	// Disagreed contains keys present in all envs but with differing values.
	// The map value is a slice of distinct values observed, in input order.
	Disagreed map[string][]string
}

// Find returns keys that appear in every env map provided.
// For keys where all values agree, they are placed in Common.
// For keys where values differ across envs, they are placed in Disagreed.
// If fewer than two envs are provided, Find returns an empty Result.
func Find(envs ...map[string]string) Result {
	result := Result{
		Common:    make(map[string]string),
		Disagreed: make(map[string][]string),
	}

	if len(envs) < 2 {
		return result
	}

	// Seed candidate keys from the first env.
	for key := range envs[0] {
		presentInAll := true
		for _, env := range envs[1:] {
			if _, ok := env[key]; !ok {
				presentInAll = false
				break
			}
		}
		if !presentInAll {
			continue
		}

		// Key is present in all envs — check value agreement.
		baseVal := envs[0][key]
		agreed := true
		for _, env := range envs[1:] {
			if env[key] != baseVal {
				agreed = false
				break
			}
		}

		if agreed {
			result.Common[key] = baseVal
		} else {
			vals := make([]string, 0, len(envs))
			seen := make(map[string]bool)
			for _, env := range envs {
				v := env[key]
				if !seen[v] {
					seen[v] = true
					vals = append(vals, v)
				}
			}
			result.Disagreed[key] = vals
		}
	}

	return result
}
