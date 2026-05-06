package merge_test

import (
	"testing"

	"envdiff/internal/merge"
)

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"HOST": "localhost", "PORT": "5432"}
	b := map[string]string{"DEBUG": "true", "LOG_LEVEL": "info"}

	res, err := merge.Merge([]map[string]string{a, b}, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if res.Env["HOST"] != "localhost" || res.Env["DEBUG"] != "true" {
		t.Errorf("unexpected env map: %v", res.Env)
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	res, err := merge.Merge([]map[string]string{a, b}, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", res.Env["KEY"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0].Key != "KEY" {
		t.Errorf("expected 1 conflict for KEY, got %v", res.Conflicts)
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}

	res, err := merge.Merge([]map[string]string{a, b}, merge.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", res.Env["KEY"])
	}
}

func TestMerge_StrategyError(t *testing.T) {
	a := map[string]string{"KEY": "val1"}
	b := map[string]string{"KEY": "val2"}

	_, err := merge.Merge([]map[string]string{a, b}, merge.StrategyError)
	if err == nil {
		t.Fatal("expected error for conflicting key, got nil")
	}
}

func TestMerge_StrategyError_SameValue(t *testing.T) {
	a := map[string]string{"KEY": "same"}
	b := map[string]string{"KEY": "same"}

	res, err := merge.Merge([]map[string]string{a, b}, merge.StrategyError)
	if err != nil {
		t.Fatalf("unexpected error for identical values: %v", err)
	}
	if res.Env["KEY"] != "same" {
		t.Errorf("expected 'same', got %q", res.Env["KEY"])
	}
}

func TestMerge_EmptySources(t *testing.T) {
	res, err := merge.Merge([]map[string]string{}, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %v", res.Env)
	}
}
