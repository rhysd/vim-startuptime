package main

import (
	"fmt"
	"os"
)

func main() {
	opts := parseOptions()
	if opts.help {
		usage()
		os.Exit(0)
	}

	collected, err := collectMeasurements(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	summary, err := summarizeStartuptime(collected, opts.verbose)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Extra options: %v\n", opts.extraArgs)
	fmt.Printf("Measured: %d times\n\n", opts.count)
	summary.print(os.Stdout)
}
