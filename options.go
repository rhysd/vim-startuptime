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
	extraArgs []string
}

const usageHeader = `Usage: vim-startuptime [flags] [-- vim options]

  vim-startuptime is a command which provides better --startuptime option of Vim.
  It starts Vim with --startuptime multiple times and collects the results.

Flags:`

func usage() {
	fmt.Fprintln(os.Stderr, usageHeader)
	flag.PrintDefaults()
}

func parseOptions() *options {
	o := &options{}

	flag.BoolVar(&o.help, "help", false, "Show this help")
	flag.UintVar(&o.count, "count", 10, "How many times measure startup time")
	flag.StringVar(&o.vimPath, "vimpath", "vim", "Command to run Vim")

	flag.Usage = usage
	flag.Parse()
	o.extraArgs = flag.Args()
	return o
}
