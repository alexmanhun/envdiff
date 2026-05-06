// Package lint provides checks for common .env file style and quality issues.
package lint

import (
	"fmt"
	"strings"
)

// Issue represents a single lint finding.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

// Result holds all issues found during linting.
type Result struct {
	Issues []Issue
}

func (r *Result) HasIssues() bool {
	return len(r.Issues) > 0
}

// Check runs all lint rules against the provided env map and returns a Result.
func Check(env map[string]string) Result {
	var issues []Issue

	for k, v := range env {
		if k != strings.ToUpper(k) {
			issues = append(issues, Issue{Key: k, Message: "key is not uppercase"})
		}
		if strings.TrimSpace(v) != v {
			issues = append(issues, Issue{Key: k, Message: "value has leading or trailing whitespace"})
		}
		if v == "" {
			issues = append(issues, Issue{Key: k, Message: "value is empty"})
		}
		if strings.Contains(k, " ") {
			issues = append(issues, Issue{Key: k, Message: "key contains spaces"})
		}
	}

	return Result{Issues: issues}
}
