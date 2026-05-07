package template

import (
	"io"
	"testing"
)

func BenchmarkGenerate_SmallEnv(b *testing.B) {
	env := map[string]string{
		"APP_PORT": "8080",
		"DB_HOST":  "localhost",
		"DB_PASS":  "secret",
		"APP_ENV":  "production",
	}
	envs := []map[string]string{env}
	opts := Options{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Generate(io.Discard, envs, opts)
	}
}

func BenchmarkGenerate_LargeEnv(b *testing.B) {
	env := make(map[string]string, 200)
	for i := 0; i < 200; i++ {
		key := fmt.Sprintf("KEY_%03d", i)
		env[key] = fmt.Sprintf("value_%d", i)
	}
	envs := []map[string]string{env}
	opts := Options{IncludeValues: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Generate(io.Discard, envs, opts)
	}
}
