package clone_test

import (
	"testing"

	"envdiff/internal/clone"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "db.example.com",
		"DB_PASS":  "secret",
		"LOG_LEVEL": "info",
	}
}

func TestClone_NilOpts(t *testing.T) {
	src := baseEnv()
	dst := clone.Clone(src, nil)

	if len(dst) != len(src) {
		t.Fatalf("expected %d keys, got %d", len(src), len(dst))
	}
	for k, v := range src {
		if dst[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, dst[k])
		}
	}
}

func TestClone_DoesNotMutateSource(t *testing.T) {
	src := baseEnv()
	opts := &clone.Options{
		Overrides: map[string]string{"APP_ENV": "staging"},
	}
	clone.Clone(src, opts)

	if src["APP_ENV"] != "production" {
		t.Errorf("source was mutated: APP_ENV = %q", src["APP_ENV"])
	}
}

func TestClone_StripKeys(t *testing.T) {
	src := baseEnv()
	opts := &clone.Options{
		StripKeys: []string{"DB_PASS", "DB_HOST"},
	}
	dst := clone.Clone(src, opts)

	if _, ok := dst["DB_PASS"]; ok {
		t.Error("expected DB_PASS to be stripped")
	}
	if _, ok := dst["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be stripped")
	}
	if dst["APP_ENV"] != "production" {
		t.Errorf("unexpected APP_ENV value: %q", dst["APP_ENV"])
	}
}

func TestClone_Overrides(t *testing.T) {
	src := baseEnv()
	opts := &clone.Options{
		Overrides: map[string]string{
			"APP_ENV":  "staging",
			"NEW_KEY":  "new_value",
		},
	}
	dst := clone.Clone(src, opts)

	if dst["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", dst["APP_ENV"])
	}
	if dst["NEW_KEY"] != "new_value" {
		t.Errorf("expected NEW_KEY=new_value, got %q", dst["NEW_KEY"])
	}
}

func TestClone_StripAndOverride(t *testing.T) {
	src := baseEnv()
	opts := &clone.Options{
		StripKeys: []string{"DB_PASS"},
		Overrides: map[string]string{"LOG_LEVEL": "debug"},
	}
	dst := clone.Clone(src, opts)

	if _, ok := dst["DB_PASS"]; ok {
		t.Error("expected DB_PASS to be stripped")
	}
	if dst["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %q", dst["LOG_LEVEL"])
	}
	if dst["DB_HOST"] != "db.example.com" {
		t.Errorf("expected DB_HOST unchanged, got %q", dst["DB_HOST"])
	}
}
