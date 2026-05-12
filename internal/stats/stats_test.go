package stats_test

import (
	"testing"

	"envdiff/internal/stats"
)

func TestCompute_Empty(t *testing.T) {
	r := stats.Compute()
	if r.TotalKeys != 0 || r.UniqueKeys != 0 {
		t.Errorf("expected zero stats, got %+v", r)
	}
}

func TestCompute_SingleEnv(t *testing.T) {
	env := map[string]string{"A": "1", "B": "", "C": "3"}
	r := stats.Compute(env)

	if r.TotalKeys != 3 {
		t.Errorf("TotalKeys: want 3, got %d", r.TotalKeys)
	}
	if r.UniqueKeys != 3 {
		t.Errorf("UniqueKeys: want 3, got %d", r.UniqueKeys)
	}
	if r.EmptyValues != 1 {
		t.Errorf("EmptyValues: want 1, got %d", r.EmptyValues)
	}
	if len(r.DuplicateKeys) != 0 {
		t.Errorf("expected no duplicate keys, got %v", r.DuplicateKeys)
	}
	if len(r.CommonKeys) != 3 {
		t.Errorf("CommonKeys: want 3, got %v", r.CommonKeys)
	}
}

func TestCompute_CommonKeys(t *testing.T) {
	dev := map[string]string{"A": "1", "B": "2"}
	prod := map[string]string{"A": "1", "C": "3"}
	r := stats.Compute(dev, prod)

	if len(r.CommonKeys) != 1 || r.CommonKeys[0] != "A" {
		t.Errorf("CommonKeys: want [A], got %v", r.CommonKeys)
	}
	if r.UniqueKeys != 3 {
		t.Errorf("UniqueKeys: want 3, got %d", r.UniqueKeys)
	}
}

func TestCompute_DuplicateKeys(t *testing.T) {
	dev := map[string]string{"A": "dev", "B": "same"}
	prod := map[string]string{"A": "prod", "B": "same"}
	r := stats.Compute(dev, prod)

	if len(r.DuplicateKeys) != 1 || r.DuplicateKeys[0] != "A" {
		t.Errorf("DuplicateKeys: want [A], got %v", r.DuplicateKeys)
	}
}

func TestCompute_KeyCounts(t *testing.T) {
	dev := map[string]string{"A": "1", "B": "2"}
	prod := map[string]string{"A": "1"}
	r := stats.Compute(dev, prod)

	if r.KeyCounts["A"] != 2 {
		t.Errorf("KeyCounts[A]: want 2, got %d", r.KeyCounts["A"])
	}
	if r.KeyCounts["B"] != 1 {
		t.Errorf("KeyCounts[B]: want 1, got %d", r.KeyCounts["B"])
	}
}

func TestCompute_TotalKeysAcrossEnvs(t *testing.T) {
	a := map[string]string{"X": "1", "Y": "2"}
	b := map[string]string{"X": "3", "Z": "4"}
	r := stats.Compute(a, b)

	if r.TotalKeys != 4 {
		t.Errorf("TotalKeys: want 4, got %d", r.TotalKeys)
	}
}
