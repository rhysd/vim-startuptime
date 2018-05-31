package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestPrintSummary(t *testing.T) {
	for _, tc := range []struct {
		what string
		ave  *averageMeasurement
		want []string
	}{
		{
			"different number of digits in entries",
			&averageMeasurement{
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
			},
			[]string{
				"  AVERAGE",
				"12.345000: /foo/bar.vim",
				" 1.234000: $VIM/vimrc",
			},
		},
		{
			"same number of digits in entries",
			&averageMeasurement{
				time.Duration(200) * time.Millisecond,
				[]averageEntry{
					averageEntry{
						"/foo/bar.vim",
						time.Duration(5678) * time.Microsecond,
					},
					averageEntry{
						"$VIM/vimrc",
						time.Duration(1234) * time.Microsecond,
					},
				},
			},
			[]string{
				" AVERAGE",
				"5.678000: /foo/bar.vim",
				"1.234000: $VIM/vimrc",
			},
		},
	} {
		t.Run(tc.what, func(t *testing.T) {
			var buf bytes.Buffer
			tc.ave.printSummary(&buf)
			lines := strings.Split(buf.String(), "\n")

			if lines[len(lines)-1] != "" {
				t.Error("Output does not end with newline")
			}
			lines = lines[:len(lines)-1]

			if !strings.HasPrefix(lines[0], "Total: 200") {
				t.Error("Total average is unexpected:", lines[0])
			}
			lines = lines[2:]

			if !reflect.DeepEqual(lines, tc.want) {
				t.Fatalf("Profile result per entry is unexpected. Have '%v' but want '%v'", lines, tc.want)
			}
		})
	}
}
