package envset_test

import (
	"os"
	"path/filepath"
	"testing"

	"envdiff/internal/envset"
	"envdiff/internal/parser"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestUnionFromParsedFiles(t *testing.T) {
	p1 := writeTempEnvFile(t, "A=1\nB=2\n")
	p2 := writeTempEnvFile(t, "B=99\nC=3\n")

	env1, err := parser.ParseFile(p1)
	if err != nil {
		t.Fatalf("parse file 1: %v", err)
	}
	env2, err := parser.ParseFile(p2)
	if err != nil {
		t.Fatalf("parse file 2: %v", err)
	}

	union := envset.Union(env1, env2)
	if len(union) != 3 {
		t.Fatalf("expected 3 keys, got %d: %v", len(union), union)
	}
	if union["B"] != "99" {
		t.Errorf("expected B=99 (last wins), got %q", union["B"])
	}
}

func TestDifferenceFromParsedFiles(t *testing.T) {
	p1 := writeTempEnvFile(t, "A=1\nB=2\nD=4\n")
	p2 := writeTempEnvFile(t, "A=1\nB=2\nC=3\n")

	env1, _ := parser.ParseFile(p1)
	env2, _ := parser.ParseFile(p2)

	diff := envset.Difference(env1, env2)
	if len(diff) != 1 {
		t.Fatalf("expected 1 key in difference, got %v", diff)
	}
	if diff["D"] != "4" {
		t.Errorf("expected D=4, got %v", diff)
	}
}
