package format

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Type represents an output format.
type Type string

const (
	Text Type = "text"
	JSON Type = "json"
)

// Parse converts a string to a Type, returning an error for unknown values.
func Parse(s string) (Type, error) {
	switch Type(s) {
	case Text, JSON:
		return Type(s), nil
	default:
		return "", fmt.Errorf("unknown format %q: must be \"text\" or \"json\"", s)
	}
}

// jsonResult is the serialisable form of a diff.Result.
type jsonResult struct {
	MissingInRight []string          `json:"missing_in_right,omitempty"`
	MissingInLeft  []string          `json:"missing_in_left,omitempty"`
	Mismatched     map[string][2]string `json:"mismatched,omitempty"`
}

// WriteJSON encodes the diff result as indented JSON to w.
func WriteJSON(w io.Writer, result diff.Result) error {
	var jr jsonResult

	if len(result.MissingInRight) > 0 {
		keys := make([]string, 0, len(result.MissingInRight))
		for k := range result.MissingInRight {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		jr.MissingInRight = keys
	}

	if len(result.MissingInLeft) > 0 {
		keys := make([]string, 0, len(result.MissingInLeft))
		for k := range result.MissingInLeft {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		jr.MissingInLeft = keys
	}

	if len(result.Mismatched) > 0 {
		jr.Mismatched = make(map[string][2]string, len(result.Mismatched))
		for k, v := range result.Mismatched {
			jr.Mismatched[k] = [2]string{v[0], v[1]}
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jr)
}
