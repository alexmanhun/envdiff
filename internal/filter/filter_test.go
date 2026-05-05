package filter_test

import (
	"testing"

	"envdiff/internal/filter"
)

var sampleEnv = map[string]string{
	"APP_NAME":    "myapp",
	"APP_VERSION": "1.0",
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"SECRET_KEY":  "abc123",
}

func TestApply_NoFilter(t *testing.T) {
	result, err := filter.Apply(sampleEnv, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(sampleEnv) {
		t.Errorf("expected %d keys, got %d", len(sampleEnv), len(result))
	}
}

func TestApply_Prefix(t *testing.T) {
	result, err := filter.Apply(sampleEnv, filter.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["APP_NAME"]; !ok {
		t.Error("expected APP_NAME in result")
	}
}

func TestApply_Pattern(t *testing.T) {
	result, err := filter.Apply(sampleEnv, filter.Options{Pattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestApply_ExcludeKeys(t *testing.T) {
	result, err := filter.Apply(sampleEnv, filter.Options{ExcludeKeys: []string{"SECRET_KEY", "DB_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, found := result["SECRET_KEY"]; found {
		t.Error("SECRET_KEY should be excluded")
	}
	if len(result) != 3 {
		t.Errorf("expected 3 keys, got %d", len(result))
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	_, err := filter.Apply(sampleEnv, filter.Options{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestApply_PrefixAndExclude(t *testing.T) {
	result, err := filter.Apply(sampleEnv, filter.Options{
		Prefix:      "APP_",
		ExcludeKeys: []string{"APP_VERSION"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
	if _, ok := result["APP_NAME"]; !ok {
		t.Error("expected APP_NAME in result")
	}
}
