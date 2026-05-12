package coerce_test

import (
	"testing"

	"envdiff/internal/coerce"
)

func TestApply_DoesNotMutateInput(t *testing.T) {
	orig := map[string]string{"ENABLED": "yes"}
	opts := coerce.DefaultOptions()
	coerce.Apply(orig, opts)
	if orig["ENABLED"] != "yes" {
		t.Fatal("Apply mutated the input map")
	}
}

func TestApply_TrimSpace(t *testing.T) {
	env := map[string]string{"FOO": "  bar  ", "BAZ": "\thello\n"}
	out := coerce.Apply(env, coerce.Options{TrimSpace: true})
	if out["FOO"] != "bar" {
		t.Errorf("expected 'bar', got %q", out["FOO"])
	}
	if out["BAZ"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["BAZ"])
	}
}

func TestApply_NormalizeBools_Yes_No(t *testing.T) {
	env := map[string]string{"A": "yes", "B": "no", "C": "on", "D": "off"}
	out := coerce.Apply(env, coerce.Options{NormalizeBools: true, LowercaseBoolValues: true})
	cases := map[string]string{"A": "true", "B": "false", "C": "true", "D": "false"}
	for k, want := range cases {
		if out[k] != want {
			t.Errorf("key %s: expected %q, got %q", k, want, out[k])
		}
	}
}

func TestApply_NormalizeBools_Numeric(t *testing.T) {
	env := map[string]string{"X": "1", "Y": "0"}
	out := coerce.Apply(env, coerce.Options{NormalizeBools: true, LowercaseBoolValues: true})
	if out["X"] != "true" {
		t.Errorf("expected 'true', got %q", out["X"])
	}
	if out["Y"] != "false" {
		t.Errorf("expected 'false', got %q", out["Y"])
	}
}

func TestApply_NoBoolNormalization_Passthrough(t *testing.T) {
	env := map[string]string{"FLAG": "yes"}
	out := coerce.Apply(env, coerce.Options{NormalizeBools: false})
	if out["FLAG"] != "yes" {
		t.Errorf("expected 'yes' unchanged, got %q", out["FLAG"])
	}
}

func TestApply_NonBoolValuesUnchanged(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	out := coerce.Apply(env, coerce.DefaultOptions())
	if out["HOST"] != "localhost" || out["PORT"] != "5432" {
		t.Errorf("non-bool values should be unchanged, got %v", out)
	}
}

func TestCount_ChangedValues(t *testing.T) {
	before := map[string]string{"A": "yes", "B": "hello", "C": "off"}
	after := coerce.Apply(before, coerce.DefaultOptions())
	n := coerce.Count(before, after)
	if n != 2 {
		t.Errorf("expected 2 changed values, got %d", n)
	}
}

func TestCount_NoChanges(t *testing.T) {
	before := map[string]string{"HOST": "localhost"}
	after := coerce.Apply(before, coerce.Options{})
	if coerce.Count(before, after) != 0 {
		t.Error("expected 0 changes")
	}
}
