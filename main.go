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

	fmt.Printf("Extra options: %v\n\n", opts.extraArgs)
	fmt.Printf("Total: %f msec\n\n", measured.total.Seconds()*1000)
	for _, e := range measured.sortedEntries {
		fmt.Printf("%f: %s\n", e.duration.Seconds()*1000, e.name)
	}
}
