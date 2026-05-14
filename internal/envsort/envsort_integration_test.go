package envSort_test

import (
	"os"
	"testing"

	envSort "github.com/example/envdiff/internal/envsort"
	"github.com/example/envdiff/internal/parser"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestSortFromParsedFile_ByKey(t *testing.T) {
	path := writeTempEnvFile(t, "ZEBRA=z\nALPHA=a\nMIDDLE=m\n")
	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	pairs := envSort.Apply(env, envSort.DefaultOptions())
	want := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	if len(pairs) != len(want) {
		t.Fatalf("got %d pairs, want %d", len(pairs), len(want))
	}
	for i, p := range pairs {
		if p.Key != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, p.Key, want[i])
		}
	}
}

func TestSortFromParsedFile_ByValue_Desc(t *testing.T) {
	path := writeTempEnvFile(t, "A=charlie\nB=alpha\nC=bravo\n")
	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	pairs := envSort.Apply(env, envSort.Options{By: "value", Order: envSort.Desc})
	want := []string{"charlie", "bravo", "alpha"}
	if len(pairs) != len(want) {
		t.Fatalf("got %d pairs, want %d", len(pairs), len(want))
	}
	for i, p := range pairs {
		if p.Value != want[i] {
			t.Errorf("pos %d value: got %q, want %q", i, p.Value, want[i])
		}
	}
}
