package ignore_test

import (
	"testing"

	"envdiff/internal/diff"
	"envdiff/internal/ignore"
)

// TestIgnoreThenCompare verifies that ignored keys are excluded before
// the diff is computed, so they never appear in the result.
func TestIgnoreThenCompare(t *testing.T) {
	left := map[string]string{
		"DB_HOST":   "localhost",
		"DB_PORT":   "5432",
		"SECRET_KEY": "abc",
	}
	right := map[string]string{
		"DB_HOST":   "localhost",
		"DB_PORT":   "5433", // intentional mismatch
		"SECRET_KEY": "xyz", // would be a mismatch, but we ignore it
	}

	rules := ignore.NewRules([]string{"SECRET_KEY"})
	filteredLeft := rules.Apply(left)
	filteredRight := rules.Apply(right)

	result := diff.Compare(filteredLeft, filteredRight)

	for _, m := range result.Mismatched {
		if m.Key == "SECRET_KEY" {
			t.Error("SECRET_KEY should have been ignored and not appear in diff")
		}
	}
	for _, k := range result.MissingInLeft {
		if k == "SECRET_KEY" {
			t.Error("SECRET_KEY should not appear in MissingInLeft")
		}
	}
	for _, k := range result.MissingInRight {
		if k == "SECRET_KEY" {
			t.Error("SECRET_KEY should not appear in MissingInRight")
		}
	}

	if len(result.Mismatched) != 1 || result.Mismatched[0].Key != "DB_PORT" {
		t.Errorf("expected exactly one mismatch on DB_PORT, got %+v", result.Mismatched)
	}
}
