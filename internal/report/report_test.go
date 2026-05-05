package report_test

import (
	"bytes"
	"strings"
	"testing"

	"envdiff/internal/diff"
	"envdiff/internal/report"
)

func makeResult(missingRight, missingLeft map[string]string, mismatched map[string][2]string) diff.Result {
	return diff.Result{
		MissingInRight: missingRight,
		MissingInLeft:  missingLeft,
		Mismatched:     mismatched,
	}
}

func TestWrite_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	result := makeResult(nil, nil, nil)
	report.Write(&buf, result, ".env", ".env.prod")
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestWrite_MissingInRight(t *testing.T) {
	var buf bytes.Buffer
	result := makeResult(map[string]string{"DB_HOST": "localhost"}, nil, nil)
	report.Write(&buf, result, ".env", ".env.prod")
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, ".env.prod") {
		t.Errorf("expected right filename in output, got: %s", out)
	}
}

func TestWrite_Mismatched(t *testing.T) {
	var buf bytes.Buffer
	result := makeResult(nil, nil, map[string][2]string{
		"PORT": {"3000", "8080"},
	})
	report.Write(&buf, result, ".env", ".env.prod")
	out := buf.String()
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %s", out)
	}
	if !strings.Contains(out, "3000") || !strings.Contains(out, "8080") {
		t.Errorf("expected both values in output, got: %s", out)
	}
}

func TestSummary_NoDiff(t *testing.T) {
	result := makeResult(nil, nil, nil)
	if s := report.Summary(result); s != "no differences" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithDiff(t *testing.T) {
	result := makeResult(
		map[string]string{"A": "1", "B": "2"},
		map[string]string{"C": "3"},
		map[string][2]string{"D": {"x", "y"}},
	)
	s := report.Summary(result)
	if !strings.Contains(s, "2 missing in right") {
		t.Errorf("unexpected summary: %s", s)
	}
	if !strings.Contains(s, "1 missing in left") {
		t.Errorf("unexpected summary: %s", s)
	}
	if !strings.Contains(s, "1 mismatched") {
		t.Errorf("unexpected summary: %s", s)
	}
}
