package mask

import (
	"fmt"
	"testing"
)

func BenchmarkApply_SmallEnv(b *testing.B) {
	env := map[string]string{
		"DB_PASS":   "secret",
		"API_TOKEN": "tok_live_abc123",
		"APP_NAME":  "myapp",
	}
	opts := DefaultOptions()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Apply(env, opts)
	}
}

func BenchmarkApply_LargeEnv(b *testing.B) {
	env := make(map[string]string, 500)
	for i := 0; i < 500; i++ {
		env[fmt.Sprintf("KEY_%d", i)] = fmt.Sprintf("value_number_%d_secret", i)
	}
	opts := Options{Mode: ModePartial, VisibleChars: 4}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Apply(env, opts)
	}
}

func BenchmarkValue_ModeFull(b *testing.B) {
	opts := Options{Mode: ModeFull, Placeholder: "****"}
	for i := 0; i < b.N; i++ {
		_ = Value("supersecretvalue", opts)
	}
}
