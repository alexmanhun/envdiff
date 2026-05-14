package envSort_test

import (
	"testing"

	envSort "github.com/example/envdiff/internal/envsort"
)

func TestApply_DefaultSortByKey(t *testing.T) {
	env := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MANGO": "m"}
	pairs := envSort.Apply(env, envSort.DefaultOptions())
	want := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, p := range pairs {
		if p.Key != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, p.Key, want[i])
		}
	}
}

func TestApply_SortByKeyDesc(t *testing.T) {
	env := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MANGO": "m"}
	pairs := envSort.Apply(env, envSort.Options{By: "key", Order: envSort.Desc})
	want := []string{"ZEBRA", "MANGO", "ALPHA"}
	for i, p := range pairs {
		if p.Key != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, p.Key, want[i])
		}
	}
}

func TestApply_SortByValue(t *testing.T) {
	env := map[string]string{"A": "charlie", "B": "alpha", "C": "bravo"}
	pairs := envSort.Apply(env, envSort.Options{By: "value", Order: envSort.Asc})
	want := []string{"alpha", "bravo", "charlie"}
	for i, p := range pairs {
		if p.Value != want[i] {
			t.Errorf("pos %d: got value %q, want %q", i, p.Value, want[i])
		}
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	env := map[string]string{"beta": "1", "ALPHA": "2", "Gamma": "3"}
	pairs := envSort.Apply(env, envSort.Options{By: "key", Order: envSort.Asc, CaseInsensitive: true})
	want := []string{"ALPHA", "beta", "Gamma"}
	for i, p := range pairs {
		if p.Key != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, p.Key, want[i])
		}
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"B": "2", "A": "1"}
	origKeys := []string{"A", "B"}
	_ = envSort.Apply(env, envSort.DefaultOptions())
	for _, k := range origKeys {
		if _, ok := env[k]; !ok {
			t.Errorf("key %q missing after Apply", k)
		}
	}
}

func TestKeys_ReturnsSortedKeys(t *testing.T) {
	env := map[string]string{"C": "3", "A": "1", "B": "2"}
	keys := envSort.Keys(env, envSort.DefaultOptions())
	want := []string{"A", "B", "C"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestApply_EmptyEnv(t *testing.T) {
	pairs := envSort.Apply(map[string]string{}, envSort.DefaultOptions())
	if len(pairs) != 0 {
		t.Errorf("expected empty result, got %d pairs", len(pairs))
	}
}
