package main

import (
	"flag"
	"fmt"
	"os"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
	"envdiff/internal/report"
)

const usage = `envdiff — Compare .env files across environments.

Usage:
  envdiff [flags] <file1> <file2>

Flags:
`

func main() {
	quiet := flag.Bool("quiet", false, "only print summary line")
	exitCode := flag.Bool("exit-code", false, "exit with code 1 if differences found")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}

	leftPath, rightPath := args[0], args[1]

	leftEnv, err := parser.ParseFile(leftPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", leftPath, err)
		os.Exit(2)
	}

	rightEnv, err := parser.ParseFile(rightPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", rightPath, err)
		os.Exit(2)
	}

	result := diff.Compare(leftEnv, rightEnv)

	if *quiet {
		fmt.Println(report.Summary(result))
	} else {
		report.Write(os.Stdout, result, leftPath, rightPath)
	}

	if *exitCode && result.HasDiff() {
		os.Exit(1)
	}
}
