package groupby_test

import (
	"os"
	"path/filepath"
	"testing"

	"envdiff/internal/groupby"
	"envdiff/internal/parser"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnvFile: %v", err)
	}
	return p
}

func TestGroupByAfterParse_BasicPrefixes(t *testing.T) {
	path := writeTempEnvFile(t, `
DB_HOST=localhost
DB_PORT=5432
AWS_KEY=abc123
AWS_REGION=us-east-1
APP_ENV=production
`)
	env, err := parser.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	groups := groupby.ByPrefix(env, "_", []string{"DB", "AWS"})

	groupMap := make(map[string]groupby.Group)
	for _, g := range groups {
		groupMap[g.Label] = g
	}

	db, ok := groupMap["DB"]
	if !ok {
		t.Fatal("expected DB group")
	}
	if len(db.Keys) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(db.Keys))
	}

	aws, ok := groupMap["AWS"]
	if !ok {
		t.Fatal("expected AWS group")
	}
	if len(aws.Keys) != 2 {
		t.Errorf("expected 2 AWS keys, got %d", len(aws.Keys))
	}

	// APP_ENV should be in the ungrouped bucket
	other, ok := groupMap[""]
	if !ok {
		t.Fatal("expected ungrouped bucket")
	}
	if _, found := other.Keys["APP_ENV"]; !found {
		t.Error("expected APP_ENV in ungrouped bucket")
	}
}
