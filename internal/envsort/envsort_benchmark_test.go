package envSort_test

import (
	"fmt"
	"testing"

	envSort "github.com/example/envdiff/internal/envsort"
)

func makeEnv(n int) map[string]string {
	m := make(map[string]string, n)
	for i := range n {
		m[fmt.Sprintf("KEY_%04d", i)] = fmt.Sprintf("value_%d", i)
	}
	return m
}

func BenchmarkApply_ByKey_Small(b *testing.B) {
	env := makeEnv(20)
	opts := envSort.DefaultOptions()
	b.ResetTimer()
	for range b.N {
		_ = envSort.Apply(env, opts)
	}
}

func BenchmarkApply_ByKey_Large(b *testing.B) {
	env := makeEnv(500)
	opts := envSort.DefaultOptions()
	b.ResetTimer()
	for range b.N {
		_ = envSort.Apply(env, opts)
	}
}

func BenchmarkApply_ByValue_Large(b *testing.B) {
	env := makeEnv(500)
	opts := envSort.Options{By: "value", Order: envSort.Asc}
	b.ResetTimer()
	for range b.N {
		_ = envSort.Apply(env, opts)
	}
}
