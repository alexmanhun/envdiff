package export

import (
	"fmt"
	"io"
	"testing"
)

func makeEnv(n int) map[string]string {
	m := make(map[string]string, n)
	for i := 0; i < n; i++ {
		m[fmt.Sprintf("KEY_%04d", i)] = fmt.Sprintf("value_%d", i)
	}
	return m
}

func BenchmarkWrite_Dotenv_Small(b *testing.B) {
	env := makeEnv(10)
	opts := DefaultOptions()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Write(io.Discard, env, opts)
	}
}

func BenchmarkWrite_Dotenv_Large(b *testing.B) {
	env := makeEnv(500)
	opts := DefaultOptions()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Write(io.Discard, env, opts)
	}
}

func BenchmarkWrite_Inline_Large(b *testing.B) {
	env := makeEnv(500)
	opts := Options{Format: FormatInline, Sorted: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Write(io.Discard, env, opts)
	}
}
