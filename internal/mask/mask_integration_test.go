package mask_test

import (
	"strings"
	"testing"

	"envdiff/internal/mask"
	"envdiff/internal/redact"
)

// TestMaskAfterRedact verifies that masking can be applied after redaction
// so that sensitive keys get both detection and value obscuring.
func TestMaskAfterRedact(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"APP_NAME":    "myapp",
		"SECRET_KEY":  "topsecret",
	}

	rules := redact.NewDefaultRules()
	redacted := redact.Apply(rules, env)

	// After redact, sensitive values are already masked by redact;
	// apply mask on top to demonstrate composability.
	opts := mask.Options{
		Mode:         mask.ModePartial,
		VisibleChars: 2,
	}
	masked := mask.Apply(redacted, opts)

	// APP_NAME is not sensitive — redact leaves it alone, mask partially obscures.
	appName := masked["APP_NAME"]
	if !strings.HasPrefix(appName, "my") {
		t.Errorf("APP_NAME: expected prefix 'my', got %q", appName)
	}

	// Sensitive keys should have been altered by at least one layer.
	if masked["DB_PASSWORD"] == "hunter2" {
		t.Error("DB_PASSWORD should not appear as plain text after masking pipeline")
	}
	if masked["SECRET_KEY"] == "topsecret" {
		t.Error("SECRET_KEY should not appear as plain text after masking pipeline")
	}
}

func TestMaskPartial_LengthPreserved(t *testing.T) {
	env := map[string]string{"TOKEN": "abcdefgh"}
	masked := mask.Apply(env, mask.Options{Mode: mask.ModePartial, VisibleChars: 3})
	if len(masked["TOKEN"]) != len("abcdefgh") {
		t.Errorf("expected same length after partial mask, got %q", masked["TOKEN"])
	}
}
