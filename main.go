package main

import (
	"flag"
	"fmt"
	"io"
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

func parseOptions(out io.Writer, args []string) (*options, int) {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(out)

	o := &options{}
	fs.UintVar(&o.count, "count", 10, "How many times measure startup time")
	fs.StringVar(&o.vimPath, "vimpath", "vim", "Command to run Vim or Neovim")
	fs.BoolVar(&o.script, "script", false, "Only collects script loading times")
	fs.UintVar(&o.warmup, "warmup", 1, "How many times start Vim at warm-up phase")
	fs.BoolVar(&o.verbose, "verbose", false, "Verbose output to stderr while measurements")
	fs.Usage = func() {
		fmt.Fprintln(out, usageHeader)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return nil, 0
		}
		fmt.Fprintf(out, "error while parsing command line arguments: %s\n", err)
		return nil, 1
	}

	o.extraArgs = fs.Args()
	return o, -1
}

func measure(opts *options, out io.Writer) error {
	collected, err := collectMeasurements(opts)
	if err != nil {
		return err
	}

	summary, err := summarizeStartuptime(collected, opts.verbose)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "Extra options: %v\n", opts.extraArgs)
	fmt.Fprintf(out, "Measured: %d times\n\n", opts.count)
	summary.print(out)
	return nil
}

func main() {
	opts, code := parseOptions(os.Stderr, os.Args)
	if code >= 0 {
		os.Exit(code)
	}

	if err := measure(opts, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
