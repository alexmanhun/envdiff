// Package redact provides utilities for masking sensitive values
// in .env file comparisons before display or logging.
package redact

import "strings"

// DefaultPatterns are key substrings that trigger redaction by default.
var DefaultPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

const mask = "***REDACTED***"

// Rules holds the set of key patterns that should have their values masked.
type Rules struct {
	patterns []string
}

// NewRules returns a Rules instance using the provided patterns.
// Patterns are matched case-insensitively against key names.
func NewRules(patterns []string) Rules {
	upper := make([]string, len(patterns))
	for i, p := range patterns {
		upper[i] = strings.ToUpper(p)
	}
	return Rules{patterns: upper}
}

// NewDefaultRules returns a Rules instance using DefaultPatterns.
func NewDefaultRules() Rules {
	return NewRules(DefaultPatterns)
}

// IsSensitive reports whether the given key matches any redaction pattern.
func (r Rules) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range r.patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// Apply returns a copy of the provided env map with sensitive values masked.
func (r Rules) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.IsSensitive(k) {
			out[k] = mask
		} else {
			out[k] = v
		}
	}
	return out
}

// Value masks a single value if the key is sensitive, otherwise returns it as-is.
func (r Rules) Value(key, value string) string {
	if r.IsSensitive(key) {
		return mask
	}
	return value
}
