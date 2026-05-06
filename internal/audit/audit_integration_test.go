package audit_test

import (
	"path/filepath"
	"testing"
	"time"

	"envdiff/internal/audit"
	"envdiff/internal/diff"
)

// TestAuditAfterCompare verifies that a real diff result can be recorded and
// retrieved from the audit log.
func TestAuditAfterCompare(t *testing.T) {
	left := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
		"USER": "admin",
	}
	right := map[string]string{
		"HOST": "prod.example.com",
		"PORT": "5432",
	}

	result := diff.Compare(left, right)

	entry := audit.Entry{
		Timestamp:  time.Now().UTC(),
		LeftFile:   ".env.dev",
		RightFile:  ".env.prod",
		Missing:    len(result.MissingInRight) + len(result.MissingInLeft),
		Mismatched: len(result.Mismatched),
		HasDiff:    result.HasDiff(),
	}

	path := filepath.Join(t.TempDir(), "audit.json")
	if err := audit.Append(path, entry); err != nil {
		t.Fatalf("Append: %v", err)
	}

	log, err := audit.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	got := log.Entries[0]
	if !got.HasDiff {
		t.Error("expected HasDiff=true")
	}
	if got.Missing != 1 {
		t.Errorf("expected Missing=1, got %d", got.Missing)
	}
	if got.Mismatched != 1 {
		t.Errorf("expected Mismatched=1, got %d", got.Mismatched)
	}
}
