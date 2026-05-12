package envset_test

import (
	"fmt"
	"testing"

	"envdiff/internal/envset"
)

func makeEnv(n int, prefix string) map[string]string {
	env := make(map[string]string, n)
	for i := 0; i < n; i++ {
		env[fmt.Sprintf("%s_KEY_%04d", prefix, i)] = fmt.Sprintf("value_%d", i)
	}
	return env
}

func BenchmarkUnion_Large(b *testing.B) {
	a := makeEnv(500, "A")
	c := makeEnv(500, "C")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = envset.Union(a, c)
	}
}

func BenchmarkDifference_Large(b *testing.B) {
	left := makeEnv(500, "L")
	right := makeEnv(500, "R")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = envset.Difference(left, right)
	}
}

func BenchmarkIntersection_Large(b *testing.B) {
	a := makeEnv(500, "SHARED")
	c := makeEnv(500, "SHARED")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = envset.Intersection(a, c)
	}
}
