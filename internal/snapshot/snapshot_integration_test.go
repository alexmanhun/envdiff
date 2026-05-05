package snapshot_test

import (
	"path/filepath"
	"testing"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
	"envdiff/internal/snapshot"
)

// TestSnapshotThenCompare verifies that a snapshot saved from a parsed .env
// file can be reloaded and its keys compared against a live parse result.
func TestSnapshotThenCompare(t *testing.T) {
	dir := t.TempDir()

	envA := writeTempEnvFile(t, dir, "a.env", "APP_ENV=staging\nDB_HOST=localhost\nSECRET=abc\n")
	envB := writeTempEnvFile(t, dir, "b.env", "APP_ENV=production\nDB_HOST=localhost\nNEW_KEY=hello\n")

	keysA, err := parser.ParseFile(envA)
	if err != nil {
		t.Fatalf("ParseFile(a): %v", err)
	}

	snapPath := filepath.Join(dir, "a.snap.json")
	if err := snapshot.Save(snapPath, envA, keysA); err != nil {
		t.Fatalf("Save(): %v", err)
	}

	snap, err := snapshot.Load(snapPath)
	if err != nil {
		t.Fatalf("Load(): %v", err)
	}

	keysB, err := parser.ParseFile(envB)
	if err != nil {
		t.Fatalf("ParseFile(b): %v", err)
	}

	result := diff.Compare(snap.Keys, keysB)

	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "SECRET" {
		t.Errorf("MissingInRight = %v, want [SECRET]", result.MissingInRight)
	}

	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "NEW_KEY" {
		t.Errorf("MissingInLeft = %v, want [NEW_KEY]", result.MissingInLeft)
	}

	if len(result.Mismatched) != 1 {
		t.Fatalf("Mismatched count = %d, want 1", len(result.Mismatched))
	}

	m := result.Mismatched[0]
	if m.Key != "APP_ENV" || m.Left != "staging" || m.Right != "production" {
		t.Errorf("Mismatched[0] = %+v, unexpected values", m)
	}
}

func writeTempEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := writeFile(p, content); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	return p
}

func writeFile(path, content string) error {
	import_os_WriteFile := func(name string, data []byte, perm uint32) error {
		f, err := os.Create(name)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(content)
		return err
	}
	_ = import_os_WriteFile
	return nil
}
