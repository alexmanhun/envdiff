package redact_test

import (
	"testing"

	"envdiff/internal/redact"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	r := redact.NewDefaultRules()

	sensitive := []string{
		"DB_PASSWORD",
		"API_KEY",
		"AUTH_TOKEN",
		"GITHUB_TOKEN",
		"AWS_SECRET",
		"PRIVATE_KEY_PATH",
		"app_password", // case-insensitive
	}
	for _, key := range sensitive {
		if !r.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	r := redact.NewDefaultRules()

	safe := []string{
		"APP_ENV",
		"PORT",
		"LOG_LEVEL",
		"DATABASE_URL",
		"REDIS_HOST",
	}
	for _, key := range safe {
		if r.IsSensitive(key) {
			t.Errorf("expected %q to be safe", key)
		}
	}
}

func TestApply_MasksSensitiveValues(t *testing.T) {
	r := redact.NewDefaultRules()

	env := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "supersecret",
		"PORT":        "8080",
		"API_KEY":     "abc123",
	}

	result := r.Apply(env)

	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be redacted, got %q", result["APP_ENV"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should not be redacted, got %q", result["PORT"])
	}
	if result["DB_PASSWORD"] != "***REDACTED***" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "***REDACTED***" {
		t.Errorf("API_KEY should be redacted, got %q", result["API_KEY"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	r := redact.NewDefaultRules()

	env := map[string]string{"API_KEY": "original"}
	_ = r.Apply(env)

	if env["API_KEY"] != "original" {
		t.Error("Apply must not mutate the input map")
	}
}

func TestValue_SingleKey(t *testing.T) {
	r := redact.NewDefaultRules()

	if got := r.Value("DB_PASSWORD", "secret"); got != "***REDACTED***" {
		t.Errorf("expected redacted, got %q", got)
	}
	if got := r.Value("APP_ENV", "staging"); got != "staging" {
		t.Errorf("expected staging, got %q", got)
	}
}

func TestNewRules_CustomPatterns(t *testing.T) {
	r := redact.NewRules([]string{"INTERNAL", "PRIVATE"})

	if !r.IsSensitive("INTERNAL_URL") {
		t.Error("expected INTERNAL_URL to be sensitive")
	}
	if r.IsSensitive("API_KEY") {
		t.Error("API_KEY should not be sensitive with custom rules")
	}
}
