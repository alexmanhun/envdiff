package report

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"envdiff/internal/diff"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Write writes a human-readable diff report to the given writer.
func Write(w io.Writer, result diff.Result, leftName, rightName string) {
	if !result.HasDiff() {
		fmt.Fprintln(w, "✓ No differences found.")
		return
	}

	if len(result.MissingInRight) > 0 {
		keys := sortedKeys(result.MissingInRight)
		fmt.Fprintf(w, "Keys present in %s but missing in %s:\n", leftName, rightName)
		for _, k := range keys {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(result.MissingInLeft) > 0 {
		keys := sortedKeys(result.MissingInLeft)
		fmt.Fprintf(w, "Keys present in %s but missing in %s:\n", rightName, leftName)
		for _, k := range keys {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(result.Mismatched) > 0 {
		keys := make([]string, 0, len(result.Mismatched))
		for k := range result.Mismatched {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fmt.Fprintln(w, "Mismatched values:")
		for _, k := range keys {
			pair := result.Mismatched[k]
			fmt.Fprintf(w, "  ~ %s: %q (%s) vs %q (%s)\n",
				k, pair[0], leftName, pair[1], rightName)
		}
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Summary returns a one-line summary string of the diff result.
func Summary(result diff.Result) string {
	if !result.HasDiff() {
		return "no differences"
	}
	parts := []string{}
	if n := len(result.MissingInRight); n > 0 {
		parts = append(parts, fmt.Sprintf("%d missing in right", n))
	}
	if n := len(result.MissingInLeft); n > 0 {
		parts = append(parts, fmt.Sprintf("%d missing in left", n))
	}
	if n := len(result.Mismatched); n > 0 {
		parts = append(parts, fmt.Sprintf("%d mismatched", n))
	}
	return strings.Join(parts, ", ")
}
