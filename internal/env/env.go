// Package env provides utilities for exporting and importing env maps
// to and from common shell-compatible formats.
package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExportFormat controls the output format used by Export.
type ExportFormat int

const (
	// FormatExport writes lines as: export KEY="VALUE"
	FormatExport ExportFormat = iota
	// FormatPlain writes lines as: KEY=VALUE
	FormatPlain
	// FormatDockerEnv writes lines as: KEY=VALUE (no quoting, suitable for --env-file)
	FormatDockerEnv
)

// Write serialises env to w using the requested format.
// Keys are written in sorted order for deterministic output.
func Write(w io.Writer, env map[string]string, format ExportFormat) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		var line string
		switch format {
		case FormatExport:
			line = fmt.Sprintf("export %s=%q\n", k, v)
		case FormatDockerEnv:
			line = fmt.Sprintf("%s=%s\n", k, v)
		default: // FormatPlain
			line = fmt.Sprintf("%s=%q\n", k, v)
		}
		if _, err := io.WriteString(w, line); err != nil {
			return fmt.Errorf("env: write key %q: %w", k, err)
		}
	}
	return nil
}

// FromSlice parses a slice of "KEY=VALUE" strings (as returned by os.Environ)
// into a map. Entries without '=' are silently skipped.
func FromSlice(pairs []string) map[string]string {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			continue
		}
		out[p[:idx]] = p[idx+1:]
	}
	return out
}
