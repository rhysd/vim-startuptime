package main

import (
	"flag"
	"fmt"
	"os"
)

type options struct {
	help      bool
	count     uint
	vimPath   string
	script    bool
	extraArgs []string
	warmup    uint
	verbose   bool
}

const usageHeader = `Usage: vim-startuptime [flags] [-- vim args]

  vim-startuptime is a command which provides better --startuptime option of Vim
  or Neovim. It starts Vim with --startuptime multiple times, collects the
  results and outputs summary of the measurements to stdout.

Flags:`

func usage() {
	fmt.Fprintln(os.Stderr, usageHeader)
	flag.PrintDefaults()
}

func parseOptions() *options {
	o := &options{}

	flag.BoolVar(&o.help, "help", false, "Show this help")
	flag.UintVar(&o.count, "count", 10, "How many times measure startup time")
	flag.StringVar(&o.vimPath, "vimpath", "vim", "Command to run Vim or Neovim")
	flag.BoolVar(&o.script, "script", false, "Only collects script loading times")
	flag.UintVar(&o.warmup, "warmup", 1, "How many times start Vim at warm-up phase")
	flag.BoolVar(&o.verbose, "verbose", false, "Verbose output to stderr while measurements")

	flag.Usage = usage
	flag.Parse()
	o.extraArgs = flag.Args()
	return o
}
