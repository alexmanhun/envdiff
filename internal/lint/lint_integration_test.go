package lint_test

import (
	"os"
	"path/filepath"
	"testing"

	"envdiff/internal/lint"
	"envdiff/internal/parser"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLintParsedFile_Clean(t *testing.T) {
	p := writeTempEnvFile(t, "DATABASE_URL=postgres://localhost/db\nPORT=8080\n")
	env, err := parser.ParseFile(p)
	if err != nil {
		t.Fatal(err)
	}
	r := lint.Check(env)
	if r.HasIssues() {
		t.Fatalf("expected no issues, got: %v", r.Issues)
	}
}

func TestLintParsedFile_WithProblems(t *testing.T) {
	p := writeTempEnvFile(t, "lowercase_key=value\nEMPTY=\nGOOD=ok\n")
	env, err := parser.ParseFile(p)
	if err != nil {
		t.Fatal(err)
	}
	r := lint.Check(env)
	if !r.HasIssues() {
		t.Fatal("expected lint issues from parsed file")
	}
	found := false
	for _, issue := range r.Issues {
		if issue.Key == "lowercase_key" {
			found = true
		}
	}
	if !found {
		t.Error("expected issue for lowercase_key")
	}
}
