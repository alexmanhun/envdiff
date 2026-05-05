package ignore_test

import (
	"os"
	"path/filepath"
	"testing"

	"envdiff/internal/ignore"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".envdiffignore")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempIgnore: %v", err)
	}
	return p
}

func TestNewRules_Contains(t *testing.T) {
	r := ignore.NewRules([]string{"SECRET", "TOKEN"})
	if !r.Contains("SECRET") {
		t.Error("expected SECRET to be ignored")
	}
	if r.Contains("DB_HOST") {
		t.Error("expected DB_HOST not to be ignored")
	}
}

func TestLoadFile_NotExist(t *testing.T) {
	r, err := ignore.LoadFile("/nonexistent/.envdiffignore")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Contains("ANYTHING") {
		t.Error("empty rules should not contain any key")
	}
}

func TestLoadFile_CommentsAndBlanks(t *testing.T) {
	p := writeTempIgnore(t, "# this is a comment\n\nSECRET_KEY\nAPI_TOKEN\n")
	r, err := ignore.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Contains("SECRET_KEY") {
		t.Error("expected SECRET_KEY to be ignored")
	}
	if !r.Contains("API_TOKEN") {
		t.Error("expected API_TOKEN to be ignored")
	}
	if r.Contains("DB_URL") {
		t.Error("expected DB_URL not to be ignored")
	}
}

func TestApply_RemovesIgnoredKeys(t *testing.T) {
	r := ignore.NewRules([]string{"PASSWORD", "SECRET"})
	env := map[string]string{
		"DB_HOST":  "localhost",
		"PASSWORD": "hunter2",
		"SECRET":   "abc123",
	}
	out := r.Apply(env)
	if _, ok := out["PASSWORD"]; ok {
		t.Error("PASSWORD should have been removed")
	}
	if _, ok := out["SECRET"]; ok {
		t.Error("SECRET should have been removed")
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST should be preserved, got %q", out["DB_HOST"])
	}
}

func TestApply_EmptyRules(t *testing.T) {
	r := ignore.NewRules(nil)
	env := map[string]string{"A": "1", "B": "2"}
	out := r.Apply(env)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}
