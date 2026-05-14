// Package envSort provides deterministic ordering of environment variable maps.
//
// Sorting is useful when writing .env files, producing human-readable diffs,
// or generating templates where a stable key order is required.
//
// # Usage
//
//	env := map[string]string{"ZEBRA": "z", "ALPHA": "a"}
//
//	// Sort by key ascending (default)
//	pairs := envSort.Apply(env, envSort.DefaultOptions())
//
//	// Sort by value, descending, case-insensitive
//	pairs = envSort.Apply(env, envSort.Options{
//		By:              "value",
//		Order:           envSort.Desc,
//		CaseInsensitive: true,
//	})
//
//	// Retrieve only the ordered keys
//	keys := envSort.Keys(env, envSort.DefaultOptions())
//
// The input map is never mutated.
package envSort
