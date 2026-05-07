package scope_test

import (
	"testing"

	"github.com/example/envdiff/internal/diff"
	"github.com/example/envdiff/internal/scope"
)

// TestScopeThenCompare partitions two env maps by scope and then runs a diff
// on each scope independently, verifying that issues are isolated correctly.
func TestScopeThenCompare(t *testing.T) {
	left := map[string]string{
		"APP_HOST":  "localhost",
		"APP_PORT":  "8080",
		"DB_HOST":   "db.local",
		"DB_PASS":   "secret",
	}
	right := map[string]string{
		"APP_HOST":  "prod.example.com",
		"APP_PORT":  "8080",
		"DB_HOST":   "db.local",
		// DB_PASS missing on purpose
	}

	scopes := []string{"APP_", "DB_"}
	leftGroups := scope.Partition(left, scopes)
	rightGroups := scope.Partition(right, scopes)

	// APP_ scope: only a value mismatch on APP_HOST
	appResult := diff.Compare(leftGroups["APP_"].Keys, rightGroups["APP_"].Keys)
	if len(appResult.MissingInRight) != 0 {
		t.Errorf("APP_ scope: unexpected missing-in-right keys: %v", appResult.MissingInRight)
	}
	if len(appResult.Mismatched) != 1 {
		t.Errorf("APP_ scope: expected 1 mismatch, got %d", len(appResult.Mismatched))
	}

	// DB_ scope: DB_PASS missing in right, no mismatches
	dbResult := diff.Compare(leftGroups["DB_"].Keys, rightGroups["DB_"].Keys)
	if len(dbResult.MissingInRight) != 1 {
		t.Errorf("DB_ scope: expected 1 missing-in-right key, got %d", len(dbResult.MissingInRight))
	}
	if dbResult.MissingInRight[0] != "DB_PASS" {
		t.Errorf("DB_ scope: expected DB_PASS missing, got %s", dbResult.MissingInRight[0])
	}
	if len(dbResult.Mismatched) != 0 {
		t.Errorf("DB_ scope: unexpected mismatches: %v", dbResult.Mismatched)
	}
}
