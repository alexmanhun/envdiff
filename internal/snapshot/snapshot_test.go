package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"envdiff/internal/snapshot"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "snap.json")

	keys := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"SECRET":   "abc123",
	}

	if err := snapshot.Save(dest, ".env.production", keys); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := snapshot.Load(dest)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if snap.File != ".env.production" {
		t.Errorf("File = %q, want %q", snap.File, ".env.production")
	}

	if snap.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}

	if snap.CapturedAt.After(time.Now().Add(time.Second)) {
		t.Error("CapturedAt is in the future")
	}

	for k, v := range keys {
		if got := snap.Keys[k]; got != v {
			t.Errorf("Keys[%q] = %q, want %q", k, got, v)
		}
	}

	if len(snap.Keys) != len(keys) {
		t.Errorf("len(Keys) = %d, want %d", len(snap.Keys), len(keys))
	}
}

func TestLoad_NotExist(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("Load() expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	badFile := filepath.Join(dir, "bad.json")

	if err := os.WriteFile(badFile, []byte("not valid json{"), 0o644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	_, err := snapshot.Load(badFile)
	if err == nil {
		t.Fatal("Load() expected error for invalid JSON, got nil")
	}
}

func TestSave_EmptyKeys(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "empty.json")

	if err := snapshot.Save(dest, ".env", map[string]string{}); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := snapshot.Load(dest)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(snap.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(snap.Keys))
	}
}
