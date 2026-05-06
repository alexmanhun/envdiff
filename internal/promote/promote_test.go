package promote_test

import (
	"testing"

	"envdiff/internal/promote"
)

func TestPromote_NewKeysAlwaysAdded(t *testing.T) {
	src := map[string]string{"NEW_KEY": "value"}
	dst := map[string]string{"EXISTING": "old"}

	out, res, err := promote.Promote(src, dst, promote.SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", out["NEW_KEY"])
	}
	if len(res.Promoted) != 1 || res.Promoted[0] != "NEW_KEY" {
		t.Errorf("expected Promoted=[NEW_KEY], got %v", res.Promoted)
	}
}

func TestPromote_SkipExisting(t *testing.T) {
	src := map[string]string{"KEY": "src_val"}
	dst := map[string]string{"KEY": "dst_val"}

	out, res, err := promote.Promote(src, dst, promote.SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "dst_val" {
		t.Errorf("expected dst_val to be preserved, got %q", out["KEY"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "KEY" {
		t.Errorf("expected Skipped=[KEY], got %v", res.Skipped)
	}
}

func TestPromote_Overwrite(t *testing.T) {
	src := map[string]string{"KEY": "src_val"}
	dst := map[string]string{"KEY": "dst_val"}

	out, res, err := promote.Promote(src, dst, promote.Overwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "src_val" {
		t.Errorf("expected src_val after overwrite, got %q", out["KEY"])
	}
	if len(res.Overwritten) != 1 || res.Overwritten[0] != "KEY" {
		t.Errorf("expected Overwritten=[KEY], got %v", res.Overwritten)
	}
}

func TestPromote_ErrorOnConflict_Triggers(t *testing.T) {
	src := map[string]string{"KEY": "a"}
	dst := map[string]string{"KEY": "b"}

	_, _, err := promote.Promote(src, dst, promote.ErrorOnConflict)
	if err == nil {
		t.Fatal("expected error on conflicting values, got nil")
	}
}

func TestPromote_ErrorOnConflict_SameValue(t *testing.T) {
	src := map[string]string{"KEY": "same"}
	dst := map[string]string{"KEY": "same"}

	out, res, err := promote.Promote(src, dst, promote.ErrorOnConflict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "same" {
		t.Errorf("expected same, got %q", out["KEY"])
	}
	if len(res.Promoted) != 1 {
		t.Errorf("expected 1 promoted key, got %v", res.Promoted)
	}
}

func TestPromote_DoesNotMutateSrc(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}
	origLen := len(src)

	promote.Promote(src, dst, promote.Overwrite) //nolint:errcheck
	if len(src) != origLen {
		t.Error("Promote mutated src map")
	}
}
