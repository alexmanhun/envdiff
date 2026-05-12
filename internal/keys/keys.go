// Package keys provides utilities for extracting, sorting, and
// inspecting key sets across one or more parsed env maps.
package keys

import "sort"

// All returns a sorted, deduplicated slice of every key present in
// any of the supplied env maps.
func All(envs ...map[string]string) []string {
	seen := make(map[string]struct{})
	for _, env := range envs {
		for k := range env {
			seen[k] = struct{}{}
		}
	}
	return sortedKeys(seen)
}

// Common returns a sorted slice of keys that appear in every supplied
// env map. Returns nil when fewer than two maps are provided.
func Common(envs ...map[string]string) []string {
	if len(envs) < 2 {
		return nil
	}
	result := make(map[string]struct{})
	for k := range envs[0] {
		result[k] = struct{}{}
	}
	for _, env := range envs[1:] {
		for k := range result {
			if _, ok := env[k]; !ok {
				delete(result, k)
			}
		}
	}
	return sortedKeys(result)
}

// Unique returns a sorted slice of keys present in base but absent
// from every other supplied env map.
func Unique(base map[string]string, others ...map[string]string) []string {
	result := make(map[string]struct{})
outer:
	for k := range base {
		for _, env := range others {
			if _, ok := env[k]; ok {
				continue outer
			}
		}
		result[k] = struct{}{}
	}
	return sortedKeys(result)
}

// Count returns the total number of distinct keys across all envs.
func Count(envs ...map[string]string) int {
	return len(All(envs...))
}

func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
