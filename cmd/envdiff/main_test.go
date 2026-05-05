package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "envdiff")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestMain_NoDiff(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	b := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	out, err := exec.Command(bin, a, b).CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "No differences") {
		t.Errorf("expected no-diff output, got: %s", out)
	}
}

func TestMain_ExitCode(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "KEY=value\n")
	b := writeTempEnv(t, "KEY=other\n")
	cmd := exec.Command(bin, "-exit-code", a, b)
	if err := cmd.Run(); err == nil {
		t.Error("expected non-zero exit code when diff exists")
	}
}

func TestMain_QuietFlag(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "A=1\nB=2\n")
	b := writeTempEnv(t, "A=1\n")
	out, _ := exec.Command(bin, "-quiet", a, b).CombinedOutput()
	if !strings.Contains(string(out), "missing in right") {
		t.Errorf("expected summary line, got: %s", out)
	}
}
