// Package envset provides set operations on environment variable maps:
// union, intersection difference, and symmetric difference.
package envset

import "sort"

// Union returns a map containing all keys from all provided envs.
// When a key appears in multiple envs the value from the last env wins.
func Union(envs ...map[string]string) map[string]string {
	out := make(map[string]string)
	for _, env := range envs {
		for k, v := range env {
			out[k] = v
		}
	}
	return out
}

// Difference returns keys (and their values from left) that are present in
// left but absent in right.
func Difference(left, right map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range left {
		if _, ok := right[k]; !ok {
			out[k] = v
		}
	}
	return out
}

// SymmetricDifference returns all keys that appear in exactly one of the two
// envs. Values are taken from whichever env contains the key.
func SymmetricDifference(left, right map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range left {
		if _, ok := right[k]; !ok {
			out[k] = v
		}
	}
	for k, v := range right {
		if _, ok := left[k]; !ok {
			out[k] = v
		}
	}
	return out
}

// Intersection returns keys present in ALL provided envs.
// Values are taken from the first env.
func Intersection(envs ...map[string]string) map[string]string {
	if len(envs) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string)
	for k, v := range envs[0] {
		inAll := true
		for _, env := range envs[1:] {
			if _, ok := env[k]; !ok {
				inAll = false
				break
			}
		}
		if inAll {
			out[k] = v
		}
	}
	return out
}

// SortedKeys returns the keys of env in sorted order.
func SortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
