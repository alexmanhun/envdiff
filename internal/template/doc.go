// Package template provides functionality for generating .env.example files
// from one or more parsed environment maps.
//
// # Overview
//
// When working across multiple environments it is useful to maintain a
// .env.example (or .env.template) file that lists every required key without
// exposing real values. This package automates that process:
//
//	env, _ := parser.ParseFile(".env.production")
//	template.Generate(os.Stdout, []map[string]string{env}, template.Options{})
//
// Output:
//
//	APP_PORT=<value>
//	DB_HOST=<value>
//	SECRET_KEY=<value>
//
// # Options
//
// Placeholder — override the default "<value>" string.
//
// IncludeValues — emit a comment above each key showing the distinct real
// values found across the provided env maps (useful for documentation).
package template
