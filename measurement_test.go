package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestPrintSummary(t *testing.T) {
	ave := &averageMeasurement{
		time.Duration(200) * time.Millisecond,
		[]averageEntry{
			averageEntry{
				"/foo/bar.vim",
				time.Duration(12345) * time.Microsecond,
			},
			averageEntry{
				"$VIM/vimrc",
				time.Duration(1234) * time.Microsecond,
			},
		},
	}

	var buf bytes.Buffer
	ave.printSummary(&buf)
	lines := strings.Split(buf.String(), "\n")

	if lines[len(lines)-1] != "" {
		t.Error("Output does not end with newline")
	}
	lines = lines[:len(lines)-1]

	if !strings.HasPrefix(lines[0], "Total: 200") {
		t.Error("Total average is unexpected:", lines[0])
	}
	lines = lines[2:]

	if lines[0] != "  AVERAGE" {
		t.Error("Header is unexpected", lines[0])
	}

	have := lines
	want := []string{
		"  AVERAGE",
		"12.345000: /foo/bar.vim",
		" 1.234000: $VIM/vimrc",
	}

	if !reflect.DeepEqual(have, want) {
		t.Fatalf("Profile result per entry is unexpected. Have '%v' but want '%v'", have, want)
	}
}
