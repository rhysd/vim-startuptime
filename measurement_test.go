package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrintSummary(t *testing.T) {
	for _, tc := range []struct {
		what    string
		summary *measurementSummary
		want    []string
	}{
		{
			"different number of digits in entries",
			&measurementSummary{
				entrySummary{
					"Total",
					time.Duration(200) * time.Millisecond,
					time.Duration(300) * time.Millisecond,
					time.Duration(100) * time.Millisecond,
				},
				[]entrySummary{
					{
						"/foo/bar.vim",
						time.Duration(12345) * time.Microsecond,
						time.Duration(13334) * time.Microsecond,
						time.Duration(11112) * time.Microsecond,
					},
					{
						"$VIM/vimrc",
						time.Duration(1234) * time.Microsecond,
						time.Duration(1334) * time.Microsecond,
						time.Duration(1112) * time.Microsecond,
					},
				},
			},
			[]string{
				"Total Average: 200.000000 msec",
				"Total Max:     300.000000 msec",
				"Total Min:     100.000000 msec",
				"",
				"  AVERAGE       MAX       MIN",
				"------------------------------",
				"12.345000 13.334000 11.112000: /foo/bar.vim",
				" 1.234000  1.334000  1.112000: $VIM/vimrc",
			},
		},
		{
			"same number of digits in entries",
			&measurementSummary{
				entrySummary{
					"Total",
					time.Duration(200) * time.Millisecond,
					time.Duration(1000) * time.Millisecond,
					time.Duration(10) * time.Millisecond,
				},
				[]entrySummary{
					{
						"/foo/bar.vim",
						time.Duration(5678) * time.Microsecond,
						time.Duration(7890) * time.Microsecond,
						time.Duration(1234) * time.Microsecond,
					},
					{
						"$VIM/vimrc",
						time.Duration(1234) * time.Microsecond,
						time.Duration(2345) * time.Microsecond,
						time.Duration(1000) * time.Microsecond,
					},
				},
			},
			[]string{
				"Total Average: 200.000000 msec",
				"Total Max:     1000.000000 msec",
				"Total Min:     10.000000 msec",
				"",
				" AVERAGE      MAX      MIN",
				"---------------------------",
				"5.678000 7.890000 1.234000: /foo/bar.vim",
				"1.234000 2.345000 1.000000: $VIM/vimrc",
			},
		},
	} {
		t.Run(tc.what, func(t *testing.T) {
			var buf bytes.Buffer
			tc.summary.print(&buf)
			lines := strings.Split(buf.String(), "\n")

			if lines[len(lines)-1] != "" {
				t.Error("Output does not end with newline")
			}
			lines = lines[:len(lines)-1]

			if len(lines) != len(tc.want) {
				t.Fatalf("Number of lines does not match: %d v.s. %d. ('%s' v.s. '%s')", len(tc.want), len(lines), tc.want, lines)
			}
			for i := range lines {
				want, have := tc.want[i], lines[i]
				if have != want {
					t.Errorf("Line %d does not match: Wanted '%s' but have '%s'", i+1, want, have)
				}
			}
		})
	}
}
