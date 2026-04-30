package diff

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
)

func TestCompare_NoDiff(t *testing.T) {
	left := parser.EnvMap{"A": "1", "B": "2"}
	right := parser.EnvMap{"A": "1", "B": "2"}
	result := Compare(left, right)
	if result.HasDiff() {
		t.Error("expected no diff for identical maps")
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := parser.EnvMap{"A": "1", "B": "2"}
	right := parser.EnvMap{"A": "1"}
	result := Compare(left, right)
	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "B" {
		t.Errorf("expected MissingInRight=[B], got %v", result.MissingInRight)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := parser.EnvMap{"A": "1"}
	right := parser.EnvMap{"A": "1", "C": "3"}
	result := Compare(left, right)
	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "C" {
		t.Errorf("expected MissingInLeft=[C], got %v", result.MissingInLeft)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	left := parser.EnvMap{"A": "1", "B": "old"}
	right := parser.EnvMap{"A": "1", "B": "new"}
	result := Compare(left, right)
	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "B" || m.LeftValue != "old" || m.RightValue != "new" {
		t.Errorf("unexpected mismatch: %+v", m)
	}
}

func TestCompare_HasDiff(t *testing.T) {
	left := parser.EnvMap{"X": "1"}
	right := parser.EnvMap{"Y": "2"}
	result := Compare(left, right)
	if !result.HasDiff() {
		t.Error("expected HasDiff() to be true")
	}
}
