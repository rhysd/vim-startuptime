package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseOK(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "__test_parse_ok")
	if err != nil {
		panic(err)
	}
	tmpname := tmpfile.Name()
	tmpfile.Close()
	defer func() {
		os.Remove(tmpname)
	}()

	header := `

times in msec
 clock   self+sourced   self:  sourced script
 clock   elapsed:              other lines

`

	for _, tc := range []struct {
		what     string
		lines    []string
		expected *measurement
	}{
		{
			"system only",
			[]string{
				"000.008  000.008: --- VIM STARTING ---",
				"000.190  000.182: Allocated generic buffers",
			},
			&measurement{
				time.Duration(190) * time.Microsecond,
				[]*measurementEntry{
					{
						false,
						time.Duration(8) * time.Microsecond,
						time.Duration(8) * time.Microsecond,
						time.Duration(0),
						"--- VIM STARTING ---",
					},
					{
						false,
						time.Duration(190) * time.Microsecond,
						time.Duration(182) * time.Microsecond,
						time.Duration(0),
						"Allocated generic buffers",
					},
				},
			},
		},
		{
			"script only",
			[]string{
				"012.696  000.429  000.429: sourcing /foo/bar.vim",
				"013.270  001.106  000.677: sourcing $VIM/vimrc",
			},
			&measurement{
				time.Duration(13270) * time.Microsecond,
				[]*measurementEntry{
					{
						true,
						time.Duration(12696) * time.Microsecond,
						time.Duration(429) * time.Microsecond,
						time.Duration(429) * time.Microsecond,
						"/foo/bar.vim",
					},
					{
						true,
						time.Duration(13270) * time.Microsecond,
						time.Duration(1106) * time.Microsecond,
						time.Duration(677) * time.Microsecond,
						"$VIM/vimrc",
					},
				},
			},
		},
		{
			"mixed",
			[]string{
				"198.369  000.161  000.161: sourcing /foo/bar.vim",
				"198.465  001.679: BufEnter autocommands",
				"198.467  000.002: editing files in windows",
				"200.107  001.135  001.135: sourcing $VIM/vimrc",
			},
			&measurement{
				time.Duration(200107) * time.Microsecond,
				[]*measurementEntry{
					{
						true,
						time.Duration(198369) * time.Microsecond,
						time.Duration(161) * time.Microsecond,
						time.Duration(161) * time.Microsecond,
						"/foo/bar.vim",
					},
					{
						false,
						time.Duration(198465) * time.Microsecond,
						time.Duration(1679) * time.Microsecond,
						time.Duration(0),
						"BufEnter autocommands",
					},
					{
						false,
						time.Duration(198467) * time.Microsecond,
						time.Duration(2) * time.Microsecond,
						time.Duration(0),
						"editing files in windows",
					},
					{
						true,
						time.Duration(200107) * time.Microsecond,
						time.Duration(1135) * time.Microsecond,
						time.Duration(1135) * time.Microsecond,
						"$VIM/vimrc",
					},
				},
			},
		},
	} {
		t.Run(tc.what, func(t *testing.T) {
			content := []byte(header + strings.Join(tc.lines, "\n") + "\n")
			if err := ioutil.WriteFile(tmpname, content, 0644); err != nil {
				panic(err)
			}

			f, err := os.Open(tmpname)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			m, err := parseStartuptime(f)
			if err != nil {
				t.Fatal(err)
			}

			w := tc.expected
			if m.elapsedTotal != w.elapsedTotal {
				t.Error("Want total", w.elapsedTotal, "but have", m.elapsedTotal)
			}

			if len(w.entries) != len(m.entries) {
				t.Fatal("Want #entries", w.entries, "but have", m.entries)
			}
			for i := range w.entries {
				have := m.entries[i]
				want := w.entries[i]
				if !reflect.DeepEqual(have, want) {
					t.Errorf("%dth entry not match. Want '%+v', but have '%+v'", i, want, have)
				}
			}
		})
	}
}
