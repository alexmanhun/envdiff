package required_test

import (
	"testing"

	"envdiff/internal/required"
)

func TestCheck_NoViolations(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	v := required.Check([]string{"DB_HOST", "DB_PORT"}, []map[string]string{env}, []string{"production"})
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d: %v", len(v), v)
	}
}

func TestCheck_MissingKey(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	v := required.Check([]string{"DB_HOST", "DB_PORT"}, []map[string]string{env}, []string{"staging"})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "DB_PORT" {
		t.Errorf("expected key DB_PORT, got %q", v[0].Key)
	}
	if v[0].EnvName != "staging" {
		t.Errorf("expected env name staging, got %q", v[0].EnvName)
	}
	if v[0].Reason != "is missing" {
		t.Errorf("unexpected reason: %q", v[0].Reason)
	}
}

func TestCheck_EmptyValue(t *testing.T) {
	env := map[string]string{"API_KEY": ""}
	v := required.Check([]string{"API_KEY"}, []map[string]string{env}, nil)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Reason != "is empty" {
		t.Errorf("expected reason 'is empty', got %q", v[0].Reason)
	}
}

func TestCheck_MultipleEnvs(t *testing.T) {
	envA := map[string]string{"SECRET": "abc"}
	envB := map[string]string{"SECRET": ""}
	v := required.Check(
		[]string{"SECRET"},
		[]map[string]string{envA, envB},
		[]string{"dev", "prod"},
	)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].EnvName != "prod" {
		t.Errorf("expected violation in prod, got %q", v[0].EnvName)
	}
}

func TestCheck_FallbackLabel(t *testing.T) {
	env := map[string]string{}
	v := required.Check([]string{"X"}, []map[string]string{env}, nil)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation")
	}
	if v[0].EnvName != "env0" {
		t.Errorf("expected fallback label env0, got %q", v[0].EnvName)
	}
}

func TestViolation_Error(t *testing.T) {
	v := required.Violation{Key: "FOO", EnvName: "test", Reason: "is missing"}
	got := v.Error()
	want := `test: key "FOO" is missing`
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
