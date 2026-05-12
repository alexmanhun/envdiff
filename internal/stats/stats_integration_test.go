package stats_test

import (
	"os"
	"path/filepath"
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/stats"
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

func TestStatsFromParsedFiles(t *testing.T) {
	devPath := writeTempEnvFile(t, "APP_ENV=development\nDB_HOST=localhost\nSECRET=\n")
	prodPath := writeTempEnvFile(t, "APP_ENV=production\nDB_HOST=db.prod.example.com\nAPI_KEY=abc123\n")

	dev, err := parser.ParseFile(devPath)
	if err != nil {
		t.Fatalf("parse dev: %v", err)
	}
	prod, err := parser.ParseFile(prodPath)
	if err != nil {
		t.Fatalf("parse prod: %v", err)
	}

	r := stats.Compute(dev, prod)

	if r.UniqueKeys != 4 {
		t.Errorf("UniqueKeys: want 4, got %d", r.UniqueKeys)
	}
	if r.EmptyValues != 1 {
		t.Errorf("EmptyValues: want 1, got %d", r.EmptyValues)
	}
	if len(r.CommonKeys) != 2 {
		t.Errorf("CommonKeys: want 2, got %v", r.CommonKeys)
	}
	if len(r.DuplicateKeys) != 2 {
		t.Errorf("DuplicateKeys: want 2 (APP_ENV, DB_HOST), got %v", r.DuplicateKeys)
	}
}
