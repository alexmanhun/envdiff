package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Key      string
	Required bool
	Pattern  *regexp.Regexp
}

// Violation describes a single validation failure.
type Violation struct {
	Key     string
	Message string
}

// RuleSet holds a collection of validation rules keyed by variable name.
type RuleSet map[string]Rule

// ParseRules builds a RuleSet from a slice of rule strings.
// Each string has the format: KEY[:required][:pattern=REGEX]
// Example: "PORT:required:pattern=^[0-9]+$"
func ParseRules(specs []string) (RuleSet, error) {
	rs := make(RuleSet, len(specs))
	for _, spec := range specs {
		parts := strings.Split(spec, ":")
		if len(parts) == 0 || parts[0] == "" {
			return nil, fmt.Errorf("invalid rule spec: %q", spec)
		}
		rule := Rule{Key: parts[0]}
		for _, opt := range parts[1:] {
			switch {
			case opt == "required":
				rule.Required = true
			case strings.HasPrefix(opt, "pattern="):
				pat := strings.TrimPrefix(opt, "pattern=")
				re, err := regexp.Compile(pat)
				if err != nil {
					return nil, fmt.Errorf("invalid pattern for key %q: %w", rule.Key, err)
				}
				rule.Pattern = re
			default:
				return nil, fmt.Errorf("unknown rule option %q in spec %q", opt, spec)
			}
		}
		rs[rule.Key] = rule
	}
	return rs, nil
}

// Apply checks the provided env map against the RuleSet and returns any violations.
func Apply(rs RuleSet, env map[string]string) []Violation {
	var violations []Violation
	for key, rule := range rs {
		val, exists := env[key]
		if rule.Required && (!exists || val == "") {
			violations = append(violations, Violation{
				Key:     key,
				Message: "required key is missing or empty",
			})
			continue
		}
		if exists && rule.Pattern != nil && !rule.Pattern.MatchString(val) {
			violations = append(violations, Violation{
				Key:     key,
				Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern.String()),
			})
		}
	}
	return violations
}
