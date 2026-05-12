package defaults_test

import (
	"reflect"
	"testing"

	"envdiff/internal/defaults"
)

func TestApply_NoFallbacks(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	out := defaults.Apply(dst, defaults.DefaultOptions())
	if !reflect.DeepEqual(out, dst) {
		t.Fatalf("expected %v, got %v", dst, out)
	}
}

func TestApply_FillsMissingKeys(t *testing.T) {
	dst := map[string]string{"A": "1"}
	fb := map[string]string{"A": "99", "B": "2", "C": "3"}
	out := defaults.Apply(dst, defaults.DefaultOptions(), fb)
	if out["A"] != "1" {
		t.Errorf("A should not be overwritten, got %q", out["A"])
	}
	if out["B"] != "2" {
		t.Errorf("B should be filled from fallback, got %q", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("C should be filled from fallback, got %q", out["C"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	dst := map[string]string{"A": "1"}
	fb := map[string]string{"B": "2"}
	defaults.Apply(dst, defaults.DefaultOptions(), fb)
	if _, ok := dst["B"]; ok {
		t.Error("Apply must not mutate dst")
	}
}

func TestApply_OverwriteEmptyValues(t *testing.T) {
	dst := map[string]string{"A": "", "B": "keep"}
	fb := map[string]string{"A": "filled", "B": "ignored"}
	opts := defaults.Options{Overwrite: true}
	out := defaults.Apply(dst, opts, fb)
	if out["A"] != "filled" {
		t.Errorf("expected A=filled, got %q", out["A"])
	}
	if out["B"] != "keep" {
		t.Errorf("expected B=keep, got %q", out["B"])
	}
}

func TestApply_SkipEmptyFallbackValues(t *testing.T) {
	dst := map[string]string{}
	fb := map[string]string{"A": "", "B": "val"}
	opts := defaults.Options{SkipEmpty: true}
	out := defaults.Apply(dst, opts, fb)
	if _, ok := out["A"]; ok {
		t.Error("A should not be set when SkipEmpty is true and fallback is empty")
	}
	if out["B"] != "val" {
		t.Errorf("expected B=val, got %q", out["B"])
	}
}

func TestApply_MultipleFallbacksPriority(t *testing.T) {
	dst := map[string]string{}
	fb1 := map[string]string{"X": "first"}
	fb2 := map[string]string{"X": "second", "Y": "only"}
	out := defaults.Apply(dst, defaults.DefaultOptions(), fb1, fb2)
	if out["X"] != "first" {
		t.Errorf("expected X=first (first fallback wins), got %q", out["X"])
	}
	if out["Y"] != "only" {
		t.Errorf("expected Y=only, got %q", out["Y"])
	}
}

func TestMissing_ReturnsAbsentKeys(t *testing.T) {
	env := map[string]string{"A": "1", "B": ""}
	ref := map[string]string{"A": "x", "B": "y", "C": "z"}
	got := defaults.Missing(env, ref)
	want := []string{"B", "C"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMissing_NoneAbsent(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	ref := map[string]string{"A": "x", "B": "y"}
	got := defaults.Missing(env, ref)
	if len(got) != 0 {
		t.Errorf("expected no missing keys, got %v", got)
	}
}
