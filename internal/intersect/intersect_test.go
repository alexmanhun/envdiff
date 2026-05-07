package intersect_test

import (
	"testing"

	"github.com/user/envdiff/internal/intersect"
)

func TestFind_FewerThanTwoEnvs(t *testing.T) {
	result := intersect.Find(map[string]string{"A": "1"})
	if len(result.Common) != 0 || len(result.Disagreed) != 0 {
		t.Fatal("expected empty result for fewer than two envs")
	}
}

func TestFind_AllAgree(t *testing.T) {
	a := map[string]string{"HOST": "localhost", "PORT": "8080"}
	b := map[string]string{"HOST": "localhost", "PORT": "8080"}
	c := map[string]string{"HOST": "localhost", "PORT": "8080"}

	result := intersect.Find(a, b, c)

	if len(result.Disagreed) != 0 {
		t.Fatalf("expected no disagreed keys, got %v", result.Disagreed)
	}
	if result.Common["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", result.Common["HOST"])
	}
	if result.Common["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", result.Common["PORT"])
	}
}

func TestFind_SomeDisagreed(t *testing.T) {
	a := map[string]string{"HOST": "localhost", "PORT": "8080"}
	b := map[string]string{"HOST": "prod.example.com", "PORT": "8080"}

	result := intersect.Find(a, b)

	if _, ok := result.Common["PORT"]; !ok {
		t.Error("expected PORT in Common")
	}
	if _, ok := result.Disagreed["HOST"]; !ok {
		t.Error("expected HOST in Disagreed")
	}
	vals := result.Disagreed["HOST"]
	if len(vals) != 2 {
		t.Fatalf("expected 2 distinct values for HOST, got %d", len(vals))
	}
}

func TestFind_KeyMissingInOneEnv(t *testing.T) {
	a := map[string]string{"HOST": "localhost", "SECRET": "abc"}
	b := map[string]string{"HOST": "localhost"}

	result := intersect.Find(a, b)

	if _, ok := result.Common["SECRET"]; ok {
		t.Error("SECRET should not be in Common — it is absent from one env")
	}
	if _, ok := result.Common["HOST"]; !ok {
		t.Error("HOST should be in Common")
	}
}

func TestFind_DeduplicatesDisagreedValues(t *testing.T) {
	a := map[string]string{"MODE": "debug"}
	b := map[string]string{"MODE": "release"}
	c := map[string]string{"MODE": "debug"}

	result := intersect.Find(a, b, c)

	vals, ok := result.Disagreed["MODE"]
	if !ok {
		t.Fatal("expected MODE in Disagreed")
	}
	// Only two distinct values: "debug" and "release"
	if len(vals) != 2 {
		t.Errorf("expected 2 distinct values, got %d: %v", len(vals), vals)
	}
}

func TestFind_EmptyEnvs(t *testing.T) {
	result := intersect.Find(map[string]string{}, map[string]string{})
	if len(result.Common) != 0 || len(result.Disagreed) != 0 {
		t.Fatal("expected empty result for empty envs")
	}
}
