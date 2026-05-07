package scope_test

import (
	"testing"

	"github.com/example/envdiff/internal/scope"
)

func TestPartition_BasicPrefixes(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
		"UNRELATED": "value",
	}

	groups := scope.Partition(env, []string{"APP_", "DB_"})

	if len(groups["APP_"].Keys) != 2 {
		t.Errorf("expected 2 APP_ keys, got %d", len(groups["APP_"].Keys))
	}
	if len(groups["DB_"].Keys) != 1 {
		t.Errorf("expected 1 DB_ key, got %d", len(groups["DB_"].Keys))
	}
	if len(groups[""].Keys) != 1 {
		t.Errorf("expected 1 unmatched key, got %d", len(groups[""].Keys))
	}
}

func TestPartition_EmptyScopes(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	groups := scope.Partition(env, nil)

	if len(groups[""].Keys) != 2 {
		t.Errorf("expected all keys in unmatched bucket, got %d", len(groups[""].Keys))
	}
}

func TestPartition_EmptyEnv(t *testing.T) {
	groups := scope.Partition(map[string]string{}, []string{"APP_"})
	if len(groups["APP_"].Keys) != 0 {
		t.Error("expected empty group for APP_")
	}
}

func TestPartition_FirstPrefixWins(t *testing.T) {
	// APP_ is listed before APP_EXTRA_ — APP_EXTRA_KEY should go into APP_
	env := map[string]string{"APP_EXTRA_KEY": "v"}
	groups := scope.Partition(env, []string{"APP_", "APP_EXTRA_"})

	if _, ok := groups["APP_"].Keys["APP_EXTRA_KEY"]; !ok {
		t.Error("expected APP_EXTRA_KEY in APP_ group (first match wins)")
	}
	if len(groups["APP_EXTRA_"].Keys) != 0 {
		t.Error("expected APP_EXTRA_ group to be empty")
	}
}

func TestNames_ReturnsSortedNonEmpty(t *testing.T) {
	env := map[string]string{
		"DB_HOST":  "x",
		"APP_PORT": "y",
	}
	groups := scope.Partition(env, []string{"APP_", "DB_", "AWS_"})
	names := scope.Names(groups)

	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d: %v", len(names), names)
	}
	if names[0] != "APP_" || names[1] != "DB_" {
		t.Errorf("expected sorted [APP_ DB_], got %v", names)
	}
}
