package rename_test

import (
	"testing"

	"envdiff/internal/rename"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
	}
}

func TestRename_SimpleRename(t *testing.T) {
	env := baseEnv()
	mapping := map[string]string{"DB_HOST": "DATABASE_HOST"}

	res, err := rename.Rename(env, mapping, rename.SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["DB_HOST"]; ok {
		t.Error("old key DB_HOST should have been removed")
	}
	if res.Env["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", res.Env["DATABASE_HOST"])
	}
	if len(res.Renamed) != 1 || res.Renamed[0] != "DB_HOST" {
		t.Errorf("expected Renamed=[DB_HOST], got %v", res.Renamed)
	}
}

func TestRename_MissingSourceKey(t *testing.T) {
	env := baseEnv()
	mapping := map[string]string{"NONEXISTENT": "NEW_KEY"}

	res, err := rename.Rename(env, mapping, rename.SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["NEW_KEY"]; ok {
		t.Error("NEW_KEY should not exist when source key is absent")
	}
	if len(res.Renamed) != 0 {
		t.Errorf("expected no renames, got %v", res.Renamed)
	}
}

func TestRename_SkipExisting(t *testing.T) {
	env := map[string]string{"OLD": "val1", "NEW": "val2"}
	mapping := map[string]string{"OLD": "NEW"}

	res, err := rename.Rename(env, mapping, rename.SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["NEW"] != "val2" {
		t.Errorf("expected NEW to remain val2, got %q", res.Env["NEW"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "OLD" {
		t.Errorf("expected Skipped=[OLD], got %v", res.Skipped)
	}
}

func TestRename_OverwriteExisting(t *testing.T) {
	env := map[string]string{"OLD": "val1", "NEW": "val2"}
	mapping := map[string]string{"OLD": "NEW"}

	res, err := rename.Rename(env, mapping, rename.OverwriteExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["NEW"] != "val1" {
		t.Errorf("expected NEW=val1 after overwrite, got %q", res.Env["NEW"])
	}
	if _, ok := res.Env["OLD"]; ok {
		t.Error("OLD key should have been removed")
	}
}

func TestRename_ErrorOnExisting(t *testing.T) {
	env := map[string]string{"OLD": "val1", "NEW": "val2"}
	mapping := map[string]string{"OLD": "NEW"}

	_, err := rename.Rename(env, mapping, rename.ErrorOnExisting)
	if err == nil {
		t.Fatal("expected error when destination key exists, got nil")
	}
}

func TestRename_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	mapping := map[string]string{"DB_HOST": "DATABASE_HOST"}

	_, err := rename.Rename(env, mapping, rename.SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["DB_HOST"]; !ok {
		t.Error("original env map should not be mutated")
	}
}
