package lint_test

import (
	"testing"

	"envdiff/internal/lint"
)

func TestCheck_NoIssues(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	r := lint.Check(env)
	if r.HasIssues() {
		t.Fatalf("expected no issues, got: %v", r.Issues)
	}
}

func TestCheck_LowercaseKey(t *testing.T) {
	env := map[string]string{"database_url": "value"}
	r := lint.Check(env)
	if !r.HasIssues() {
		t.Fatal("expected issue for lowercase key")
	}
	assertIssue(t, r, "database_url", "not uppercase")
}

func TestCheck_WhitespaceValue(t *testing.T) {
	env := map[string]string{"HOST": " localhost "}
	r := lint.Check(env)
	assertIssue(t, r, "HOST", "whitespace")
}

func TestCheck_EmptyValue(t *testing.T) {
	env := map[string]string{"SECRET": ""}
	r := lint.Check(env)
	assertIssue(t, r, "SECRET", "empty")
}

func TestCheck_KeyWithSpaces(t *testing.T) {
	env := map[string]string{"MY KEY": "value"}
	r := lint.Check(env)
	assertIssue(t, r, "MY KEY", "spaces")
}

func TestCheck_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"good_key": "",
		"ANOTHER":  " padded ",
	}
	r := lint.Check(env)
	if len(r.Issues) < 2 {
		t.Fatalf("expected at least 2 issues, got %d", len(r.Issues))
	}
}

func assertIssue(t *testing.T, r lint.Result, key, substr string) {
	t.Helper()
	for _, issue := range r.Issues {
		if issue.Key == key && containsSubstr(issue.Message, substr) {
			return
		}
	}
	t.Errorf("expected issue for key %q containing %q, got: %v", key, substr, r.Issues)
}

func containsSubstr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
