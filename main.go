package main

import (
	"flag"
	"fmt"
	"os"
)

type options struct {
	count     uint
	vimPath   string
	script    bool
	extraArgs []string
	warmup    uint
	verbose   bool
}

const usageHeader = `Usage: vim-startuptime [flags] [-- VIMARGS...]

  vim-startuptime is a command which provides better --startuptime option of Vim
  or Neovim. It starts Vim with --startuptime multiple times, collects the
  results and outputs summary of the measurements to stdout.

Flags:`

func usage() {
	fmt.Fprintln(os.Stderr, usageHeader)
	flag.PrintDefaults()
}

func main() {
	opts := options{}

	flag.UintVar(&opts.count, "count", 10, "How many times measure startup time")
	flag.StringVar(&opts.vimPath, "vimpath", "vim", "Command to run Vim or Neovim")
	flag.BoolVar(&opts.script, "script", false, "Only collects script loading times")
	flag.UintVar(&opts.warmup, "warmup", 1, "How many times start Vim at warm-up phase")
	flag.BoolVar(&opts.verbose, "verbose", false, "Verbose output to stderr while measurements")

	flag.Usage = usage
	flag.Parse()
	opts.extraArgs = flag.Args()

	collected, err := collectMeasurements(&opts)
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
