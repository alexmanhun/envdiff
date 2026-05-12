package groupby_test

import (
	"testing"

	"envdiff/internal/groupby"
)

func TestByPrefix_BasicGrouping(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"AWS_KEY":  "abc",
		"APP_ENV":  "prod",
		"UNRELATED": "x",
	}
	groups := groupby.ByPrefix(env, "_", []string{"DB", "AWS"})

	if len(groups) < 2 {
		t.Fatalf("expected at least 2 groups, got %d", len(groups))
	}

	var db groupby.Group
	for _, g := range groups {
		if g.Label == "DB" {
			db = g
		}
	}
	if _, ok := db.Keys["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in DB group")
	}
	if _, ok := db.Keys["DB_PORT"]; !ok {
		t.Error("expected DB_PORT in DB group")
	}
}

func TestByPrefix_OtherGroup(t *testing.T) {
	env := map[string]string{
		"DB_HOST":   "localhost",
		"UNRELATED": "x",
		"ANOTHER":   "y",
	}
	groups := groupby.ByPrefix(env, "_", []string{"DB"})

	var other groupby.Group
	for _, g := range groups {
		if g.Label == "" {
			other = g
		}
	}
	if len(other.Keys) != 2 {
		t.Errorf("expected 2 ungrouped keys, got %d", len(other.Keys))
	}
}

func TestByPrefix_EmptyEnv(t *testing.T) {
	groups := groupby.ByPrefix(map[string]string{}, "_", []string{"DB", "AWS"})
	for _, g := range groups {
		if len(g.Keys) != 0 {
			t.Errorf("expected empty group %q", g.Label)
		}
	}
}

func TestByPrefix_ExactPrefixKey(t *testing.T) {
	env := map[string]string{"DB": "value"}
	groups := groupby.ByPrefix(env, "_", []string{"DB"})
	for _, g := range groups {
		if g.Label == "DB" {
			if _, ok := g.Keys["DB"]; !ok {
				t.Error("expected exact prefix key 'DB' to be grouped under DB")
			}
		}
	}
}

func TestLabels_ReturnsSortedNonEmpty(t *testing.T) {
	groups := []groupby.Group{
		{Label: "ZZ", Keys: map[string]string{"ZZ_A": "1"}},
		{Label: "", Keys: map[string]string{"OTHER": "2"}},
		{Label: "AA", Keys: map[string]string{"AA_B": "3"}},
	}
	labels := groupby.Labels(groups)
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
	if labels[0] != "AA" || labels[1] != "ZZ" {
		t.Errorf("unexpected label order: %v", labels)
	}
}
