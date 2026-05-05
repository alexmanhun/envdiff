package format_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/format"
)

func TestParse_Valid(t *testing.T) {
	for _, tc := range []struct{ in string; want format.Type }{
		{"text", format.Text},
		{"json", format.JSON},
	} {
		got, err := format.Parse(tc.in)
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", tc.in, err)
		}
		if got != tc.want {
			t.Errorf("Parse(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParse_Invalid(t *testing.T) {
	_, err := format.Parse("xml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestWriteJSON_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	err := format.WriteJSON(&buf, diff.Result{})
	if err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty object, got %v", out)
	}
}

func TestWriteJSON_MissingAndMismatched(t *testing.T) {
	result := diff.Result{
		MissingInRight: map[string]string{"FOO": "bar"},
		Mismatched:     map[string][2]string{"PORT": {"8080", "9090"}},
	}
	var buf bytes.Buffer
	if err := format.WriteJSON(&buf, result); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out struct {
		MissingInRight []string            `json:"missing_in_right"`
		Mismatched     map[string][2]string `json:"mismatched"`
	}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out.MissingInRight) != 1 || out.MissingInRight[0] != "FOO" {
		t.Errorf("missing_in_right = %v, want [FOO]", out.MissingInRight)
	}
	if v, ok := out.Mismatched["PORT"]; !ok || v[0] != "8080" || v[1] != "9090" {
		t.Errorf("mismatched PORT = %v", out.Mismatched)
	}
}
