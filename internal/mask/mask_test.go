package mask

import (
	"strings"
	"testing"
)

func TestValue_ModeNone(t *testing.T) {
	v := Value("supersecret", Options{Mode: ModeNone})
	if v != "supersecret" {
		t.Errorf("expected unchanged value, got %q", v)
	}
}

func TestValue_ModeFull(t *testing.T) {
	v := Value("supersecret", Options{Mode: ModeFull, Placeholder: "****"})
	if v != "****" {
		t.Errorf("expected placeholder, got %q", v)
	}
}

func TestValue_ModeFull_DefaultPlaceholder(t *testing.T) {
	v := Value("abc", Options{Mode: ModeFull})
	if v != "****" {
		t.Errorf("expected default placeholder, got %q", v)
	}
}

func TestValue_ModePartial_Normal(t *testing.T) {
	v := Value("supersecret", Options{Mode: ModePartial, VisibleChars: 3})
	if !strings.HasPrefix(v, "sup") {
		t.Errorf("expected prefix 'sup', got %q", v)
	}
	if !strings.Contains(v, "*") {
		t.Errorf("expected asterisks in %q", v)
	}
}

func TestValue_ModePartial_ShortValue(t *testing.T) {
	v := Value("ab", Options{Mode: ModePartial, VisibleChars: 5})
	if v != "**" {
		t.Errorf("expected '**', got %q", v)
	}
}

func TestValue_Empty(t *testing.T) {
	for _, mode := range []Mode{ModeNone, ModePartial, ModeFull} {
		v := Value("", Options{Mode: mode, Placeholder: "****"})
		if v != "" {
			t.Errorf("mode %d: expected empty string, got %q", mode, v)
		}
	}
}

func TestApply_MasksAllValues(t *testing.T) {
	env := map[string]string{
		"DB_PASS": "secret123",
		"API_KEY": "abcdef",
	}
	result := Apply(env, Options{Mode: ModeFull, Placeholder: "[redacted]"})
	for k, v := range result {
		if v != "[redacted]" {
			t.Errorf("key %s: expected [redacted], got %q", k, v)
		}
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"SECRET": "myvalue"}
	_ = Apply(env, DefaultOptions())
	if env["SECRET"] != "myvalue" {
		t.Error("Apply mutated the input map")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.Mode != ModePartial {
		t.Errorf("expected ModePartial, got %d", opts.Mode)
	}
	if opts.VisibleChars != 3 {
		t.Errorf("expected 3 visible chars, got %d", opts.VisibleChars)
	}
}
