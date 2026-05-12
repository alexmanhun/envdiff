// Package stats provides summary statistics over one or more parsed env maps.
package stats

import (
	"sort"
	"strings"
)

// Result holds aggregate statistics for a collection of env maps.
type Result struct {
	TotalKeys    int
	UniqueKeys   int
	EmptyValues  int
	DuplicateKeys []string // keys that appear in more than one env with differing values
	CommonKeys   []string // keys present in every env
	KeyCounts    map[string]int // how many envs each key appears in
}

// Compute derives statistics from one or more env maps.
// Each map represents one environment (e.g. dev, staging, prod).
func Compute(envs ...map[string]string) Result {
	if len(envs) == 0 {
		return Result{KeyCounts: map[string]int{}}
	}

	keyCounts := map[string]int{}
	keyValues := map[string][]string{} // key -> distinct values across envs
	emptyCount := 0
	totalKeys := 0

	for _, env := range envs {
		for k, v := range env {
			totalKeys++
			keyCounts[k]++
			if strings.TrimSpace(v) == "" {
				emptyCount++
			}
			keyValues[k] = appendUnique(keyValues[k], v)
		}
	}

	n := len(envs)
	var commonKeys []string
	var duplicateKeys []string

	for k, count := range keyCounts {
		if count == n {
			commonKeys = append(commonKeys, k)
		}
		if len(keyValues[k]) > 1 {
			duplicateKeys = append(duplicateKeys, k)
		}
	}

	sort.Strings(commonKeys)
	sort.Strings(duplicateKeys)

	return Result{
		TotalKeys:    totalKeys,
		UniqueKeys:   len(keyCounts),
		EmptyValues:  emptyCount,
		DuplicateKeys: duplicateKeys,
		CommonKeys:   commonKeys,
		KeyCounts:    keyCounts,
	}
}

func appendUnique(slice []string, val string) []string {
	for _, v := range slice {
		if v == val {
			return slice
		}
	}
	return append(slice, val)
}
