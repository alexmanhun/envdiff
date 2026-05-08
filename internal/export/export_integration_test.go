package export_test

import (
	"strings"
	"testing"

	"envdiff/internal/export"
	"envdiff/internal/parser"
	"envdiff/internal/redact"
	"os"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestExportAfterParse_Dotenv(t *testing.T) {
	path := writeTempEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=production\n")
	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var sb strings.Builder
	if err := export.Write(&sb, env, export.DefaultOptions()); err != nil {
		t.Fatal(err)
	}
	out := sb.String()
	for _, want := range []string{"DB_HOST=localhost", "DB_PORT=5432", "APP_ENV=production"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output:\n%s", want, out)
		}
	}
}

func TestExportAfterRedact_MasksSecrets(t *testing.T) {
	path := writeTempEnvFile(t, "API_KEY=supersecret\nHOST=example.com\n")
	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	rules := redact.NewDefaultRules()
	redacted := rules.Apply(env)

	var sb strings.Builder
	if err := export.Write(&sb, redacted, export.DefaultOptions()); err != nil {
		t.Fatal(err)
	}
	out := sb.String()
	if strings.Contains(out, "supersecret") {
		t.Errorf("secret value should have been redacted; got:\n%s", out)
	}
	if !strings.Contains(out, "HOST=example.com") {
		t.Errorf("non-secret value should be preserved; got:\n%s", out)
	}
}
