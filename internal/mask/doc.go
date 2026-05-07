// Package mask provides value-masking utilities for env maps.
//
// It supports three modes:
//
//   - ModeNone    — values are returned unchanged.
//   - ModePartial — the first N characters are kept; the rest are replaced
//                   with asterisks. Useful for debugging ("I can see it starts
//                   with 'sk-'") while avoiding full secret exposure.
//   - ModeFull    — the entire value is replaced with a configurable placeholder
//                   string (default "***").
//
// mask is intentionally decoupled from the redact package: redact decides
// *which* keys are sensitive; mask decides *how* to obscure their values.
// The two can be composed freely.
//
// Example:
//
//	env, _ := parser.ParseFile(".env")
//	masked := mask.Apply(env, mask.DefaultOptions())
//	report.Write(os.Stdout, result, masked)
package mask
