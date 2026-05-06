package drift_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/drift"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/snapshot"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write env: %v", err)
	}
	return p
}

// TestDriftFromParsedEnv saves a snapshot from a parsed .env file, then
// detects drift against a modified live environment.
func TestDriftFromParsedEnv(t *testing.T) {
	envPath := writeTempEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\nSECRET=abc\n")

	original, err := parser.ParseFile(envPath)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	snapPath := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(snapPath, original); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}

	// Simulate live env: DB_PORT changed, SECRET removed, NEW_KEY added.
	live := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5433",
		"NEW_KEY": "surprise",
	}

	report, err := drift.Detect(snapPath, live)
	if err != nil {
		t.Fatalf("detect: %v", err)
	}

	if !report.HasDrift() {
		t.Fatal("expected drift but none reported")
	}
	if _, ok := report.Added["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY in Added")
	}
	if _, ok := report.Removed["SECRET"]; !ok {
		t.Error("expected SECRET in Removed")
	}
	if pair, ok := report.Changed["DB_PORT"]; !ok || pair[0] != "5432" || pair[1] != "5433" {
		t.Errorf("unexpected DB_PORT change: %v", report.Changed["DB_PORT"])
	}
}
