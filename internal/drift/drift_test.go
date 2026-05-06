package drift_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/drift"
)

func writeSnapshot(t *testing.T, data map[string]string) string {
	t.Helper()
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	p := filepath.Join(t.TempDir(), "snap.json")
	if err := os.WriteFile(p, b, 0o644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}
	return p
}

func TestDetect_NoDrift(t *testing.T) {
	snap := map[string]string{"HOST": "localhost", "PORT": "5432"}
	p := writeSnapshot(t, snap)

	report, err := drift.Detect(p, map[string]string{"HOST": "localhost", "PORT": "5432"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.HasDrift() {
		t.Errorf("expected no drift, got added=%v removed=%v changed=%v",
			report.Added, report.Removed, report.Changed)
	}
}

func TestDetect_AddedKey(t *testing.T) {
	p := writeSnapshot(t, map[string]string{"HOST": "localhost"})

	report, err := drift.Detect(p, map[string]string{"HOST": "localhost", "NEW_KEY": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := report.Added["NEW_KEY"]; !ok || v != "value" {
		t.Errorf("expected NEW_KEY in Added, got %v", report.Added)
	}
}

func TestDetect_RemovedKey(t *testing.T) {
	p := writeSnapshot(t, map[string]string{"HOST": "localhost", "OLD_KEY": "gone"})

	report, err := drift.Detect(p, map[string]string{"HOST": "localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := report.Removed["OLD_KEY"]; !ok || v != "gone" {
		t.Errorf("expected OLD_KEY in Removed, got %v", report.Removed)
	}
}

func TestDetect_ChangedValue(t *testing.T) {
	p := writeSnapshot(t, map[string]string{"HOST": "localhost"})

	report, err := drift.Detect(p, map[string]string{"HOST": "production.host"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pair, ok := report.Changed["HOST"]
	if !ok {
		t.Fatalf("expected HOST in Changed")
	}
	if pair[0] != "localhost" || pair[1] != "production.host" {
		t.Errorf("unexpected changed pair: %v", pair)
	}
}

func TestDetect_SnapshotNotExist(t *testing.T) {
	_, err := drift.Detect("/nonexistent/snap.json", map[string]string{})
	if err == nil {
		t.Error("expected error for missing snapshot, got nil")
	}
}
