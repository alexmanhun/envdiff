package stats_test

import (
	"strings"
	"testing"

	"envdiff/internal/stats"
)

func TestWriteText_ContainsCounts(t *testing.T) {
	dev := map[string]string{"A": "1", "B": "2", "C": ""}
	prod := map[string]string{"A": "99", "B": "2", "D": "4"}
	r := stats.Compute(dev, prod)

	var sb strings.Builder
	if err := stats.WriteText(&sb, r); err != nil {
		t.Fatalf("WriteText: %v", err)
	}

	out := sb.String()
	for _, want := range []string{
		"Total keys",
		"Unique keys",
		"Empty values",
		"common",
		"differing",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\noutput:\n%s", want, out)
		}
	}
}

func TestWriteText_NoDiff(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	r := stats.Compute(env)

	var sb strings.Builder
	if err := stats.WriteText(&sb, r); err != nil {
		t.Fatalf("WriteText: %v", err)
	}

	out := sb.String()
	if strings.Contains(out, "differing: ") {
		t.Errorf("expected no differing keys line, got:\n%s", out)
	}
}

func TestWriteText_Empty(t *testing.T) {
	r := stats.Compute()
	var sb strings.Builder
	if err := stats.WriteText(&sb, r); err != nil {
		t.Fatalf("WriteText: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "Total keys") {
		t.Errorf("expected header line, got:\n%s", out)
	}
}
