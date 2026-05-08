package trim_test

import (
	"testing"

	"envdiff/internal/trim"
)

func TestApply_DoesNotMutateInput(t *testing.T) {
	original := map[string]string{"A": "", "B": "hello"}
	_ = trim.Apply(original, trim.DefaultOptions())
	if _, ok := original["A"]; !ok {
		t.Fatal("Apply mutated the input map")
	}
}

func TestApply_RemoveEmpty(t *testing.T) {
	env := map[string]string{"A": "", "B": "value", "C": ""}
	out := trim.Apply(env, trim.Options{RemoveEmpty: true})
	if _, ok := out["A"]; ok {
		t.Error("expected A to be removed")
	}
	if _, ok := out["C"]; ok {
		t.Error("expected C to be removed")
	}
	if out["B"] != "value" {
		t.Errorf("expected B=value, got %q", out["B"])
	}
}

func TestApply_RemoveWhitespace(t *testing.T) {
	env := map[string]string{"A": "   ", "B": "\t", "C": "ok"}
	out := trim.Apply(env, trim.Options{RemoveWhitespace: true})
	if _, ok := out["A"]; ok {
		t.Error("expected A to be removed")
	}
	if _, ok := out["B"]; ok {
		t.Error("expected B to be removed")
	}
	if out["C"] != "ok" {
		t.Errorf("expected C=ok, got %q", out["C"])
	}
}

func TestApply_ExplicitKeys(t *testing.T) {
	env := map[string]string{"KEEP": "yes", "DROP": "nonempty", "ALSO": "fine"}
	out := trim.Apply(env, trim.Options{Keys: []string{"DROP"}})
	if _, ok := out["DROP"]; ok {
		t.Error("expected DROP to be removed")
	}
	if out["KEEP"] != "yes" {
		t.Error("expected KEEP to remain")
	}
}

func TestApply_NoOptions(t *testing.T) {
	env := map[string]string{"A": "", "B": "  ", "C": "val"}
	out := trim.Apply(env, trim.Options{})
	if len(out) != len(env) {
		t.Errorf("expected all %d keys, got %d", len(env), len(out))
	}
}

func TestCount(t *testing.T) {
	env := map[string]string{"A": "", "B": "  ", "C": "val"}
	n := trim.Count(env, trim.DefaultOptions())
	if n != 2 {
		t.Errorf("expected 2 removed, got %d", n)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := trim.DefaultOptions()
	if !opts.RemoveEmpty || !opts.RemoveWhitespace {
		t.Error("DefaultOptions should enable RemoveEmpty and RemoveWhitespace")
	}
}
