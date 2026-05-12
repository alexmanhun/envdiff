// Package envset implements classical set operations for environment variable
// maps (map[string]string).
//
// Available operations:
//
//   - Union           – all keys from one or more envs (last value wins)
//   - Difference      – keys in left that are absent in right
//   - SymmetricDifference – keys that appear in exactly one of two envs
//   - Intersection    – keys present in every provided env
//
// All operations are non-destructive: they never modify their inputs.
//
// Example:
//
//	prod := map[string]string{"DB_HOST": "prod-db", "PORT": "5432"}
//	staging := map[string]string{"DB_HOST": "staging-db", "DEBUG": "true"}
//
//	onlyInProd := envset.Difference(prod, staging)
//	// {"PORT": "5432"}
//
//	common := envset.Intersection(prod, staging)
//	// {"DB_HOST": "prod-db"}
package envset
