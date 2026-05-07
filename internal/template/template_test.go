package template

import (
	"strings"
	"testing"
)

func TestGenerate_BasicPlaceholder(t *testing.T) {
	env := map[string]string{"APP_PORT": "8080", "DB_HOST": "localhost"}
	var sb strings.Builder
	if err := Generate(&sb, []map[string]string{env}, Options{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "APP_PORT=<value>") {
		t.Errorf("expected APP_PORT=<value>, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=<value>") {
		t.Errorf("expected DB_HOST=<value>, got:\n%s", out)
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	env := map[string]string{"SECRET": "abc"}
	var sb strings.Builder
	_ = Generate(&sb, []map[string]string{env}, Options{Placeholder: "CHANGEME"})
	if !strings.Contains(sb.String(), "SECRET=CHANGEME") {
		t.Errorf("expected SECRET=CHANGEME, got: %s", sb.String())
	}
}

func TestGenerate_IncludeValues(t *testing.T) {
	env := map[string]string{"APP_ENV": "production"}
	var sb strings.Builder
	_ = Generate(&sb, []map[string]string{env}, Options{IncludeValues: true})
	out := sb.String()
	if !strings.Contains(out, "# example: production") {
		t.Errorf("expected example comment, got:\n%s", out)
	}
}

func TestGenerate_MergesMultipleEnvs(t *testing.T) {
	env1 := map[string]string{"A": "1", "B": "2"}
	env2 := map[string]string{"B": "2", "C": "3"}
	var sb strings.Builder
	_ = Generate(&sb, []map[string]string{env1, env2}, Options{})
	out := sb.String()
	for _, key := range []string{"A", "B", "C"} {
		if !strings.Contains(out, key+"=<value>") {
			t.Errorf("expected key %s in output, got:\n%s", key, out)
		}
	}
}

func TestGenerate_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	var sb strings.Builder
	_ = Generate(&sb, []map[string]string{env}, Options{})
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected A_KEY first, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected Z_KEY last, got %s", lines[2])
	}
}

func TestGenerate_DeduplicatesExampleValues(t *testing.T) {
	env1 := map[string]string{"HOST": "localhost"}
	env2 := map[string]string{"HOST": "localhost"}
	var sb strings.Builder
	_ = Generate(&sb, []map[string]string{env1, env2}, Options{IncludeValues: true})
	out := sb.String()
	// "localhost" should appear only once in the comment
	count := strings.Count(out, "localhost")
	if count != 1 {
		t.Errorf("expected localhost once, got %d occurrences in:\n%s", count, out)
	}
}
