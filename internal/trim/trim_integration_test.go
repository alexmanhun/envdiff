package trim_test

import (
	"os"
	"path/filepath"
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/trim"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestTrimAfterParse_RemovesEmpty(t *testing.T) {
	path := writeTempEnvFile(t, "HOST=localhost\nEMPTY=\nSPACE=   \nPORT=8080\n")

	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	cleaned := trim.Apply(env, trim.DefaultOptions())

	for _, unwanted := range []string{"EMPTY", "SPACE"} {
		if _, ok := cleaned[unwanted]; ok {
			t.Errorf("expected %q to be trimmed", unwanted)
		}
	}
	for _, wanted := range []string{"HOST", "PORT"} {
		if _, ok := cleaned[wanted]; !ok {
			t.Errorf("expected %q to be present", wanted)
		}
	}
}

func TestTrimAfterParse_ExplicitKeyRemoval(t *testing.T) {
	path := writeTempEnvFile(t, "DEBUG=true\nSECRET=abc123\nAPP=myapp\n")

	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	cleaned := trim.Apply(env, trim.Options{Keys: []string{"SECRET"}})

	if _, ok := cleaned["SECRET"]; ok {
		t.Error("expected SECRET to be removed")
	}
	if cleaned["DEBUG"] != "true" {
		t.Error("expected DEBUG to remain")
	}
	if cleaned["APP"] != "myapp" {
		t.Error("expected APP to remain")
	}
}
