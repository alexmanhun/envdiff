package envset_test

import (
	"testing"

	"envdiff/internal/envset"
)

func TestUnion_DisjointMaps(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	got := envset.Union(a, b)
	if got["A"] != "1" || got["B"] != "2" || len(got) != 2 {
		t.Fatalf("unexpected union: %v", got)
	}
}

func TestUnion_LastValueWins(t *testing.T) {
	a := map[string]string{"K": "first"}
	b := map[string]string{"K": "last"}
	got := envset.Union(a, b)
	if got["K"] != "last" {
		t.Fatalf("expected 'last', got %q", got["K"])
	}
}

func TestUnion_Empty(t *testing.T) {
	got := envset.Union()
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestDifference_Basic(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"B": "2", "C": "3"}
	got := envset.Difference(left, right)
	if len(got) != 1 || got["A"] != "1" {
		t.Fatalf("unexpected difference: %v", got)
	}
}

func TestDifference_NoDiff(t *testing.T) {
	env := map[string]string{"X": "1"}
	got := envset.Difference(env, env)
	if len(got) != 0 {
		t.Fatalf("expected empty difference, got %v", got)
	}
}

func TestSymmetricDifference(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"B": "2", "C": "3"}
	got := envset.SymmetricDifference(left, right)
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %v", got)
	}
	if _, ok := got["A"]; !ok {
		t.Error("expected key A")
	}
	if _, ok := got["C"]; !ok {
		t.Error("expected key C")
	}
}

func TestIntersection_TwoEnvs(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "99", "C": "3"}
	got := envset.Intersection(a, b)
	if len(got) != 1 || got["B"] != "2" {
		t.Fatalf("unexpected intersection: %v", got)
	}
}

func TestIntersection_Empty(t *testing.T) {
	got := envset.Intersection()
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestSortedKeys(t *testing.T) {
	env := map[string]string{"Z": "1", "A": "2", "M": "3"}
	keys := envset.SortedKeys(env)
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Fatalf("unexpected order: %v", keys)
	}
}
