package diff

import "github.com/user/envdiff/internal/parser"

// Result holds the comparison outcome between two env files.
type Result struct {
	// MissingInRight are keys present in left but absent in right.
	MissingInRight []string
	// MissingInLeft are keys present in right but absent in left.
	MissingInLeft []string
	// Mismatched are keys present in both but with different values.
	Mismatched []MismatchedKey
}

// MismatchedKey describes a key whose value differs between the two files.
type MismatchedKey struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Compare returns the diff between two EnvMaps.
func Compare(left, right parser.EnvMap) Result {
	var result Result

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			result.MissingInRight = append(result.MissingInRight, k)
			continue
		}
		if lv != rv {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:        k,
				LeftValue:  lv,
				RightValue: rv,
			})
		}
	}

	for k := range right {
		if _, ok := left[k]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, k)
		}
	}

	return result
}

// HasDiff returns true if the result contains any differences.
func (r Result) HasDiff() bool {
	return len(r.MissingInRight) > 0 ||
		len(r.MissingInLeft) > 0 ||
		len(r.Mismatched) > 0
}
