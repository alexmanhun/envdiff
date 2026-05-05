package validate_test

import (
	"testing"

	"envdiff/internal/validate"
)

func TestParseRules_Valid(t *testing.T) {
	specs := []string{
		"PORT:required:pattern=^[0-9]+$",
		"HOST:required",
		"DEBUG",
	}
	rs, err := validate.ParseRules(specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rs) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rs))
	}
	if !rs["PORT"].Required {
		t.Error("PORT should be required")
	}
	if rs["PORT"].Pattern == nil {
		t.Error("PORT should have a pattern")
	}
	if !rs["HOST"].Required {
		t.Error("HOST should be required")
	}
	if rs["DEBUG"].Required {
		t.Error("DEBUG should not be required")
	}
}

func TestParseRules_InvalidPattern(t *testing.T) {
	_, err := validate.ParseRules([]string{"KEY:pattern=[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestParseRules_UnknownOption(t *testing.T) {
	_, err := validate.ParseRules([]string{"KEY:unknown"})
	if err == nil {
		t.Fatal("expected error for unknown option")
	}
}

func TestApply_NoViolations(t *testing.T) {
	rs, _ := validate.ParseRules([]string{"PORT:required:pattern=^[0-9]+$"})
	env := map[string]string{"PORT": "8080"}
	violations := validate.Apply(rs, env)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestApply_MissingRequired(t *testing.T) {
	rs, _ := validate.ParseRules([]string{"PORT:required"})
	env := map[string]string{}
	violations := validate.Apply(rs, env)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %q", violations[0].Key)
	}
}

func TestApply_PatternMismatch(t *testing.T) {
	rs, _ := validate.ParseRules([]string{"PORT:pattern=^[0-9]+$"})
	env := map[string]string{"PORT": "not-a-number"}
	violations := validate.Apply(rs, env)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestApply_EmptyRequiredValue(t *testing.T) {
	rs, _ := validate.ParseRules([]string{"SECRET:required"})
	env := map[string]string{"SECRET": ""}
	violations := validate.Apply(rs, env)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for empty required value, got %d", len(violations))
	}
}
