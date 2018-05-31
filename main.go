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

	measured, err := measureAverageStartuptime(collected)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Extra options: %v\n", opts.extraArgs)
	fmt.Printf("Measured: %d times\n\n", opts.count)
	measured.printSummary(os.Stdout)
}
