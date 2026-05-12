package envchain_test

import (
	"strings"
	"testing"

	"envdiff/internal/coerce"
	"envdiff/internal/envchain"
	"envdiff/internal/trim"
)

// TestChain_TrimThenCoerce exercises envchain with real internal packages:
// first trim empty/whitespace values, then coerce booleans to canonical form.
func TestChain_TrimThenCoerce(t *testing.T) {
	env := map[string]string{
		"DEBUG":   "yes",
		"VERBOSE": "  ",
		"ENABLED": "1",
		"EMPTY":   "",
	}

	trimStep := func(e map[string]string) map[string]string {
		return trim.Apply(e, trim.DefaultOptions())
	}
	coerceStep := func(e map[string]string) map[string]string {
		return coerce.Apply(e, coerce.DefaultOptions())
	}

	c := envchain.New(trimStep, coerceStep)
	out := c.Run(env)

	if _, ok := out["VERBOSE"]; ok {
		t.Error("expected VERBOSE to be trimmed away")
	}
	if _, ok := out["EMPTY"]; ok {
		t.Error("expected EMPTY to be trimmed away")
	}
	if out["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true after coerce, got %q", out["DEBUG"])
	}
	if out["ENABLED"] != "true" {
		t.Errorf("expected ENABLED=true after coerce, got %q", out["ENABLED"])
	}
}

// TestChain_DoesNotMutateOriginal ensures the original env is untouched.
func TestChain_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": ""}
	orig := map[string]string{"FOO": "bar", "BAZ": ""}

	trimStep := func(e map[string]string) map[string]string {
		return trim.Apply(e, trim.DefaultOptions())
	}
	upperStep := func(e map[string]string) map[string]string {
		out := make(map[string]string, len(e))
		for k, v := range e {
			out[k] = strings.ToUpper(v)
		}
		return out
	}

	c := envchain.New(trimStep, upperStep)
	c.Run(env)

	for k, v := range orig {
		if env[k] != v {
			t.Errorf("original env mutated: key %s changed to %q", k, env[k])
		}
	}
}
