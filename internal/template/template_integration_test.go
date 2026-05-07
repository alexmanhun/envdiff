package template_test

import (
	"strings"
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/template"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envdiff-*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestTemplateFromParsedFiles(t *testing.T) {
	path1 := writeTempEnvFile(t, "APP_PORT=8080\nDB_HOST=localhost\n")
	path2 := writeTempEnvFile(t, "APP_PORT=9090\nDB_PASS=secret\n")

	env1, err := parser.ParseFile(path1)
	if err != nil {
		t.Fatalf("parse file1: %v", err)
	}
	env2, err := parser.ParseFile(path2)
	if err != nil {
		t.Fatalf("parse file2: %v", err)
	}

	var sb strings.Builder
	if err := template.Generate(&sb, []map[string]string{env1, env2}, template.Options{}); err != nil {
		t.Fatalf("generate: %v", err)
	}

	out := sb.String()
	for _, key := range []string{"APP_PORT", "DB_HOST", "DB_PASS"} {
		if !strings.Contains(out, key+"=<value>") {
			t.Errorf("expected %s=<value> in output:\n%s", key, out)
		}
	}
}

func TestTemplateIncludeValues_MultipleFiles(t *testing.T) {
	path1 := writeTempEnvFile(t, "APP_ENV=staging\n")
	path2 := writeTempEnvFile(t, "APP_ENV=production\n")

	env1, _ := parser.ParseFile(path1)
	env2, _ := parser.ParseFile(path2)

	var sb strings.Builder
	_ = template.Generate(&sb, []map[string]string{env1, env2}, template.Options{IncludeValues: true})
	out := sb.String()

	if !strings.Contains(out, "staging") || !strings.Contains(out, "production") {
		t.Errorf("expected both values in example comment:\n%s", out)
	}
}
