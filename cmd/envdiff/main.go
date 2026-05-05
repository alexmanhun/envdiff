package main

import (
	"flag"
	"fmt"
	"os"

	"envdiff/internal/diff"
	"envdiff/internal/filter"
	"envdiff/internal/parser"
	"envdiff/internal/report"
)

func main() {
	quiet := flag.Bool("quiet", false, "suppress output, only set exit code")
	prefix := flag.String("prefix", "", "only compare keys with this prefix")
	pattern := flag.String("pattern", "", "only compare keys matching this regex")
	exclude := flag.String("exclude", "", "comma-separated list of keys to exclude")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: envdiff [flags] <file1> <file2>")
		os.Exit(2)
	}

	left, err := parser.ParseFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", args[0], err)
		os.Exit(2)
	}

	right, err := parser.ParseFile(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", args[1], err)
		os.Exit(2)
	}

	var excludeKeys []string
	if *exclude != "" {
		for _, k := range splitCSV(*exclude) {
			excludeKeys = append(excludeKeys, k)
		}
	}

	opts := filter.Options{
		Prefix:      *prefix,
		Pattern:     *pattern,
		ExcludeKeys: excludeKeys,
	}

	left, err = filter.Apply(left, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid filter: %v\n", err)
		os.Exit(2)
	}
	right, err = filter.Apply(right, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid filter: %v\n", err)
		os.Exit(2)
	}

	result := diff.Compare(left, right)

	if !*quiet {
		report.Write(os.Stdout, args[0], args[1], result)
	}

	if result.HasDiff() {
		os.Exit(1)
	}
}

func splitCSV(s string) []string {
	var out []string
	for _, part := range splitOn(s, ',') {
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func splitOn(s string, sep rune) []string {
	var parts []string
	start := 0
	for i, c := range s {
		if c == sep {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}
