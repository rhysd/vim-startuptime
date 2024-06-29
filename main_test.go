package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseOptionsOK(t *testing.T) {
	var b bytes.Buffer

	o, s := parseOptions(&b, []string{"vim-startuptime", "-count", "3", "-vimpath", "nvim", "-script", "-warmup", "5", "-verbose", "--", "--foo"})
	if s >= 0 {
		t.Fatal("unexpected exit with status:", s)
	}

	want := &options{
		count:     3,
		vimPath:   "nvim",
		script:    true,
		extraArgs: []string{"--foo"},
		warmup:    5,
		verbose:   true,
	}

	if !cmp.Equal(o, want, cmp.AllowUnexported(options{})) {
		t.Fatal(cmp.Diff(o, want, cmp.AllowUnexported(options{})))
	}

	stderr := b.String()
	if stderr != "" {
		t.Fatalf("Unexpected stderr output %q", stderr)
	}
}

func TestParseOptionsUnknownFlag(t *testing.T) {
	var b bytes.Buffer

	_, s := parseOptions(&b, []string{"vim-startuptime", "-foo"})
	if s <= 0 {
		t.Fatal("unexpected status:", s)
	}

	stderr := b.String()
	if !strings.Contains(stderr, "flag provided but not defined: -foo") {
		t.Fatal("unexpected error output to stderr:", stderr)
	}
}

func TestParseOptionsHelpOutput(t *testing.T) {
	var b bytes.Buffer

	_, s := parseOptions(&b, []string{"vim-startuptime", "-help"})
	if s != 0 {
		t.Fatal("unexpected status:", s)
	}

	stderr := b.String()
	if !strings.HasPrefix(stderr, "Usage: vim-startuptime [flags] [-- VIMARGS...]") {
		t.Fatal("unexpected help output to stderr:", stderr)
	}
}

func TestRunMeasurementMainOK(t *testing.T) {
	o := &options{
		count:     1,
		vimPath:   "vim",
		script:    false,
		extraArgs: []string{},
		warmup:    0,
		verbose:   false,
	}

	var b bytes.Buffer
	if err := measure(o, &b); err != nil {
		t.Fatal(err)
	}

	out := b.String()
	for _, want := range []string{
		"Extra options: []",
		"Measured: 1 times",
		"Total Average:",
		"Total Max:",
		"Total Min:",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output should contain %q but got %q", want, out)
		}
	}
}

func TestErrorWhileRunningMain(t *testing.T) {
	o := &options{
		count:     1,
		vimPath:   "this-executable-does-not-exist",
		script:    false,
		extraArgs: []string{},
		warmup:    0,
		verbose:   false,
	}

	err := measure(o, io.Discard)
	if !strings.Contains(err.Error(), `failed to run "this-executable-does-not-exist"`) {
		t.Fatal("error was unexpected:", err)
	}
}
