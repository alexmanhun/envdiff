// Package export writes env maps to various output formats (shell, dotenv, JSON).
package export

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents a supported export format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatShell  Format = "shell"
	FormatInline Format = "inline"
)

// Options controls export behaviour.
type Options struct {
	Format  Format
	Sorted  bool
	Export  bool   // prefix with "export " (dotenv/shell)
	Comment string // optional header comment
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Format: FormatDotenv,
		Sorted: true,
	}
}

// Write serialises env to w using opts.
func Write(w io.Writer, env map[string]string, opts Options) error {
	keys := keys(env, opts.Sorted)

	if opts.Comment != "" {
		if _, err := fmt.Fprintf(w, "# %s\n", opts.Comment); err != nil {
			return err
		}
	}

	switch opts.Format {
	case FormatShell, FormatDotenv:
		return writeDotenv(w, env, keys, opts.Export)
	case FormatInline:
		return writeInline(w, env, keys)
	default:
		return fmt.Errorf("export: unknown format %q", opts.Format)
	}
}

func writeDotenv(w io.Writer, env map[string]string, keys []string, export bool) error {
	prefix := ""
	if export {
		prefix = "export "
	}
	for _, k := range keys {
		v := quoteIfNeeded(env[k])
		if _, err := fmt.Fprintf(w, "%s%s=%s\n", prefix, k, v); err != nil {
			return err
		}
	}
	return nil
}

func writeInline(w io.Writer, env map[string]string, keys []string) error {
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+quoteIfNeeded(env[k]))
	}
	_, err := fmt.Fprintln(w, strings.Join(parts, " "))
	return err
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n#") {
		return `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	}
	return v
}

func keys(env map[string]string, sorted bool) []string {
	out := make([]string, 0, len(env))
	for k := range env {
		out = append(out, k)
	}
	if sorted {
		sort.Strings(out)
	}
	return out
}
