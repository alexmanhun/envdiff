// Package envdiff provides a high-level API that wires together parsing,
// filtering, ignoring, and diffing into a single reusable pipeline.
package envdiff

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/ignore"
	"github.com/user/envdiff/internal/parser"
)

// Options controls how the pipeline behaves.
type Options struct {
	// IgnoreFile is an optional path to a .envdiffignore file.
	IgnoreFile string
	// Prefix restricts comparison to keys that start with this prefix.
	Prefix string
	// Pattern is an optional glob/regex pattern to further filter keys.
	Pattern string
	// ExcludeKeys is a list of exact key names to exclude from comparison.
	ExcludeKeys []string
}

// Result holds the parsed environments and the diff produced by the pipeline.
type Result struct {
	// Left is the parsed and filtered left environment.
	Left map[string]string
	// Right is the parsed and filtered right environment.
	Right map[string]string
	// Diff is the comparison result.
	Diff diff.Result
}

// Run executes the full parse → filter → ignore → diff pipeline for two
// .env files and returns the combined Result.
func Run(leftPath, rightPath string, opts Options) (Result, error) {
	left, err := parser.ParseFile(leftPath)
	if err != nil {
		return Result{}, fmt.Errorf("parsing %s: %w", leftPath, err)
	}

	right, err := parser.ParseFile(rightPath)
	if err != nil {
		return Result{}, fmt.Errorf("parsing %s: %w", rightPath, err)
	}

	// Apply key filters.
	fopts := filter.Options{
		Prefix:      opts.Prefix,
		Pattern:     opts.Pattern,
		ExcludeKeys: opts.ExcludeKeys,
	}
	left, err = filter.Apply(left, fopts)
	if err != nil {
		return Result{}, fmt.Errorf("filtering left env: %w", err)
	}
	right, err = filter.Apply(right, fopts)
	if err != nil {
		return Result{}, fmt.Errorf("filtering right env: %w", err)
	}

	// Apply ignore rules when a file is provided.
	if opts.IgnoreFile != "" {
		rules, err := ignore.LoadFile(opts.IgnoreFile)
		if err != nil {
			return Result{}, fmt.Errorf("loading ignore file: %w", err)
		}
		left = rules.Apply(left)
		right = rules.Apply(right)
	}

	return Result{
		Left:  left,
		Right: right,
		Diff:  diff.Compare(left, right),
	}, nil
}
