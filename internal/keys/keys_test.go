package keys_test

import (
	"testing"

	"github.com/user/envdiff/internal/keys"
)

func TestAll_SingleEnv(t *testing.T) {
	env := map[string]string{"B": "2", "A": "1", "C": "3"}
	got := keys.All(env)
	want := []string{"A", "B", "C"}
	assertEqual(t, want, got)
}

func TestAll_MultipleEnvs_Deduplicates(t *testing.T) {
	a := map[string]string{"X": "1", "Y": "2"}
	b := map[string]string{"Y": "3", "Z": "4"}
	got := keys.All(a, b)
	want := []string{"X", "Y", "Z"}
	assertEqual(t, want, got)
}

func TestAll_Empty(t *testing.T) {
	got := keys.All()
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestCommon_TwoEnvs(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2", "C": "3"}
	b := map[string]string{"B": "x", "C": "y", "D": "z"}
	got := keys.Common(a, b)
	want := []string{"B", "C"}
	assertEqual(t, want, got)
}

func TestCommon_FewerThanTwo(t *testing.T) {
	got := keys.Common(map[string]string{"A": "1"})
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestCommon_NoOverlap(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	got := keys.Common(a, b)
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestUnique_BasicCase(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	other := map[string]string{"B": "x", "D": "y"}
	got := keys.Unique(base, other)
	want := []string{"A", "C"}
	assertEqual(t, want, got)
}

func TestUnique_NoOthers(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	got := keys.Unique(base)
	want := []string{"A", "B"}
	assertEqual(t, want, got)
}

func TestCount_MultipleEnvs(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "3", "C": "4"}
	if got := keys.Count(a, b); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func assertEqual(t *testing.T, want, got []string) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("length mismatch: want %v, got %v", want, got)
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("index %d: want %q, got %q", i, want[i], got[i])
		}
	}
}
