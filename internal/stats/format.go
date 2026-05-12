package stats

import (
	"fmt"
	"io"
	"strings"
)

// WriteText writes a human-readable summary of the Result to w.
func WriteText(w io.Writer, r Result) error {
	lines := []string{
		fmt.Sprintf("Total keys (across all envs): %d", r.TotalKeys),
		fmt.Sprintf("Unique keys:                  %d", r.UniqueKeys),
		fmt.Sprintf("Empty values:                 %d", r.EmptyValues),
		fmt.Sprintf("Keys common to all envs:      %d", len(r.CommonKeys)),
		fmt.Sprintf("Keys with differing values:   %d", len(r.DuplicateKeys)),
	}

	if len(r.CommonKeys) > 0 {
		lines = append(lines, "  common: "+strings.Join(r.CommonKeys, ", "))
	}
	if len(r.DuplicateKeys) > 0 {
		lines = append(lines, "  differing: "+strings.Join(r.DuplicateKeys, ", "))
	}

	for _, l := range lines {
		if _, err := fmt.Fprintln(w, l); err != nil {
			return err
		}
	}
	return nil
}
