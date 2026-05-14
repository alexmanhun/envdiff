package envdiff_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/envdiff"
)

func writeTempEnv(t *testing.T, content string) string {
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

func TestRun_NoDiff(t *testing.T) {
	left := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	right := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	res, err := envdiff.Run(left, right, envdiff.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Diff.HasDiff() {
		t.Errorf("expected no diff, got %+v", res.Diff)
	}
}

func TestRun_MissingKey(t *testing.T) {
	left := writeTempEnv(t, "FOO=bar\nONLY_LEFT=1\n")
	right := writeTempEnv(t, "FOO=bar\n")

	res, err := envdiff.Run(left, right, envdiff.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Diff.MissingInRight["ONLY_LEFT"]; !ok {
		t.Error("expected ONLY_LEFT to be missing in right")
	}
}

func TestRun_WithPrefix(t *testing.T) {
	left := writeTempEnv(t, "APP_FOO=1\nDB_HOST=localhost\n")
	right := writeTempEnv(t, "APP_FOO=2\nDB_HOST=localhost\n")

	res, err := envdiff.Run(left, right, envdiff.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Diff.Mismatched["APP_FOO"]; !ok {
		t.Error("expected APP_FOO to be mismatched")
	}
	if _, ok := res.Left["DB_HOST"]; ok {
		t.Error("DB_HOST should have been filtered out by prefix")
	}
}

func TestRun_WithIgnoreFile(t *testing.T) {
	left := writeTempEnv(t, "FOO=bar\nSECRET=abc\n")
	right := writeTempEnv(t, "FOO=bar\nSECRET=xyz\n")

	ignoreFile := filepath.Join(t.TempDir(), ".envdiffignore")
	if err := os.WriteFile(ignoreFile, []byte("SECRET\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := envdiff.Run(left, right, envdiff.Options{IgnoreFile: ignoreFile})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Diff.HasDiff() {
		t.Errorf("expected no diff after ignoring SECRET, got %+v", res.Diff)
	}
}

func TestRun_InvalidLeftFile(t *testing.T) {
	right := writeTempEnv(t, "FOO=bar\n")
	_, err := envdiff.Run("/nonexistent/path.env", right, envdiff.Options{})
	if err == nil {
		t.Error("expected error for missing left file")
	}
}
