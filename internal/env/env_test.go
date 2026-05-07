package env_test

import (
	"strings"
	"testing"

	"envdiff/internal/env"
)

func TestWrite_FormatPlain(t *testing.T) {
	input := map[string]string{"FOO": "bar", "BAZ": "qux"}
	var sb strings.Builder
	if err := env.Write(&sb, input, env.FormatPlain); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "BAZ=\"") {
		t.Errorf("expected quoted BAZ, got:\n%s", out)
	}
	if !strings.Contains(out, "FOO=\"") {
		t.Errorf("expected quoted FOO, got:\n%s", out)
	}
}

func TestWrite_FormatExport(t *testing.T) {
	input := map[string]string{"API_KEY": "secret"}
	var sb strings.Builder
	if err := env.Write(&sb, input, env.FormatExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.HasPrefix(out, "export API_KEY=") {
		t.Errorf("expected 'export' prefix, got: %s", out)
	}
}

func TestWrite_FormatDockerEnv(t *testing.T) {
	input := map[string]string{"PORT": "8080"}
	var sb strings.Builder
	if err := env.Write(&sb, input, env.FormatDockerEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if out != "PORT=8080\n" {
		t.Errorf("expected unquoted docker-env line, got: %q", out)
	}
}

func TestWrite_SortedOutput(t *testing.T) {
	input := map[string]string{"Z": "last", "A": "first", "M": "mid"}
	var sb strings.Builder
	_ = env.Write(&sb, input, env.FormatDockerEnv)
	lines := strings.Split(strings.TrimRight(sb.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A=") {
		t.Errorf("expected first line to start with A=, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z=") {
		t.Errorf("expected last line to start with Z=, got %s", lines[2])
	}
}

func TestFromSlice_Basic(t *testing.T) {
	pairs := []string{"FOO=bar", "BAZ=qux", "NO_EQUALS"}
	out := env.FromSlice(pairs)
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", out["FOO"])
	}
	if out["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", out["BAZ"])
	}
	if _, ok := out["NO_EQUALS"]; ok {
		t.Error("expected NO_EQUALS to be skipped")
	}
}

func TestFromSlice_ValueContainsEquals(t *testing.T) {
	pairs := []string{"URL=http://example.com?a=1&b=2"}
	out := env.FromSlice(pairs)
	if out["URL"] != "http://example.com?a=1&b=2" {
		t.Errorf("unexpected value: %q", out["URL"])
	}
}

func TestFromSlice_Empty(t *testing.T) {
	out := env.FromSlice(nil)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
