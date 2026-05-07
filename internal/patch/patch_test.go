package patch_test

import (
	"testing"

	"envdiff/internal/patch"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "db.prod",
		"LOG_LEVEL": "info",
	}
}

func TestApply_SetNewKey(t *testing.T) {
	result, err := patch.Apply(baseEnv(), []patch.Change{
		{Op: patch.OpSet, Key: "NEW_KEY", Value: "hello"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", result.Env["NEW_KEY"])
	}
	if len(result.Applied) != 1 || result.Applied[0].Key != "NEW_KEY" {
		t.Errorf("expected 1 applied change, got %v", result.Applied)
	}
}

func TestApply_OverwriteExisting(t *testing.T) {
	result, err := patch.Apply(baseEnv(), []patch.Change{
		{Op: patch.OpSet, Key: "APP_ENV", Value: "staging"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Env["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", result.Env["APP_ENV"])
	}
}

func TestApply_DeleteExistingKey(t *testing.T) {
	result, err := patch.Apply(baseEnv(), []patch.Change{
		{Op: patch.OpDelete, Key: "LOG_LEVEL"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.Env["LOG_LEVEL"]; ok {
		t.Error("expected LOG_LEVEL to be deleted")
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied change, got %d", len(result.Applied))
	}
}

func TestApply_DeleteMissingKey_Skipped(t *testing.T) {
	result, err := patch.Apply(baseEnv(), []patch.Change{
		{Op: patch.OpDelete, Key: "DOES_NOT_EXIST"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 || result.Skipped[0].Key != "DOES_NOT_EXIST" {
		t.Errorf("expected 1 skipped change, got %v", result.Skipped)
	}
	if len(result.Applied) != 0 {
		t.Errorf("expected 0 applied changes, got %d", len(result.Applied))
	}
}

func TestApply_DoesNotMutateSource(t *testing.T) {
	src := baseEnv()
	_, err := patch.Apply(src, []patch.Change{
		{Op: patch.OpSet, Key: "APP_ENV", Value: "test"},
		{Op: patch.OpDelete, Key: "DB_HOST"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src["APP_ENV"] != "production" {
		t.Errorf("source mutated: APP_ENV=%q", src["APP_ENV"])
	}
	if _, ok := src["DB_HOST"]; !ok {
		t.Error("source mutated: DB_HOST was deleted")
	}
}

func TestApply_UnknownOp_ReturnsError(t *testing.T) {
	_, err := patch.Apply(baseEnv(), []patch.Change{
		{Op: "upsert", Key: "X"},
	})
	if err == nil {
		t.Error("expected error for unknown op, got nil")
	}
}

func TestApply_EmptyKey_ReturnsError(t *testing.T) {
	_, err := patch.Apply(baseEnv(), []patch.Change{
		{Op: patch.OpSet, Key: "", Value: "val"},
	})
	if err == nil {
		t.Error("expected error for empty key, got nil")
	}
}
