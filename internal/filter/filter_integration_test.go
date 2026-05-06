package filter_test

import (
	"testing"

	"envdiff/internal/diff"
	"envdiff/internal/filter"
)

// TestFilterThenCompare verifies that filtering env maps before diffing
// produces results scoped to the filtered keys only.
func TestFilterThenCompare(t *testing.T) {
	left := map[string]string{
		"APP_NAME":    "myapp",
		"APP_VERSION": "1.0",
		"DB_HOST":     "localhost",
		"DB_PASS":     "secret",
	}
	right := map[string]string{
		"APP_NAME":    "myapp",
		"APP_VERSION": "2.0",
		"DB_HOST":     "prod-db",
	}

	opts := filter.Options{Prefix: "APP_"}

	filteredLeft, err := filter.Apply(left, opts)
	if err != nil {
		t.Fatalf("filter left: %v", err)
	}
	filteredRight, err := filter.Apply(right, opts)
	if err != nil {
		t.Fatalf("filter right: %v", err)
	}

	result := diff.Compare(filteredLeft, filteredRight)

	if len(result.MissingInRight) != 0 {
		t.Errorf("expected no missing-in-right after prefix filter, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected no missing-in-left after prefix filter, got %v", result.MissingInLeft)
	}
	if len(result.Mismatched) != 1 {
		t.Errorf("expected 1 mismatch (APP_VERSION), got %d", len(result.Mismatched))
	}
	if result.Mismatched["APP_VERSION"].Left != "1.0" || result.Mismatched["APP_VERSION"].Right != "2.0" {
		t.Errorf("unexpected mismatch values: %+v", result.Mismatched["APP_VERSION"])
	}
}

// TestFilterThenCompare_MissingKey verifies that a key present only in the left
// map is reported as missing-in-right after filtering.
func TestFilterThenCompare_MissingKey(t *testing.T) {
	left := map[string]string{
		"APP_NAME":    "myapp",
		"APP_TIMEOUT": "30s",
	}
	right := map[string]string{
		"APP_NAME": "myapp",
	}

	opts := filter.Options{Prefix: "APP_"}

	filteredLeft, err := filter.Apply(left, opts)
	if err != nil {
		t.Fatalf("filter left: %v", err)
	}
	filteredRight, err := filter.Apply(right, opts)
	if err != nil {
		t.Fatalf("filter right: %v", err)
	}

	result := diff.Compare(filteredLeft, filteredRight)

	if len(result.MissingInRight) != 1 {
		t.Errorf("expected 1 missing-in-right (APP_TIMEOUT), got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected no missing-in-left, got %v", result.MissingInLeft)
	}
	if len(result.Mismatched) != 0 {
		t.Errorf("expected no mismatches, got %v", result.Mismatched)
	}
}
