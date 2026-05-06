package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tmpLog(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.json")
}

func TestLoad_NotExist(t *testing.T) {
	l, err := Load("/nonexistent/audit.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.Entries) != 0 {
		t.Errorf("expected empty log, got %d entries", len(l.Entries))
	}
}

func TestAppend_CreatesFile(t *testing.T) {
	path := tmpLog(t)
	e := Entry{
		Timestamp:  time.Now().UTC(),
		LeftFile:   ".env.dev",
		RightFile:  ".env.prod",
		Missing:    2,
		Mismatched: 1,
		HasDiff:    true,
	}
	if err := Append(path, e); err != nil {
		t.Fatalf("Append error: %v", err)
	}
	l, err := Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if len(l.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.Entries))
	}
	if l.Entries[0].Missing != 2 {
		t.Errorf("expected Missing=2, got %d", l.Entries[0].Missing)
	}
}

func TestAppend_AccumulatesEntries(t *testing.T) {
	path := tmpLog(t)
	for i := 0; i < 3; i++ {
		if err := Append(path, Entry{Timestamp: time.Now().UTC()}); err != nil {
			t.Fatalf("Append error: %v", err)
		}
	}
	l, _ := Load(path)
	if len(l.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(l.Entries))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := tmpLog(t)
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
