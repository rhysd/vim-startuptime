package main

import (
	"fmt"
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
		{
			"lua for neovim (#2)",
			[]string{
				"021.963  000.003: parsing arguments",
				"023.515  000.278  000.278: require('vim.shared')",
				"023.699  000.101  000.101: require('vim._meta')",
			},
			&measurement{
				time.Duration(23699) * time.Microsecond,
				[]*measurementEntry{
					{
						false,
						time.Duration(21963) * time.Microsecond,
						time.Duration(3) * time.Microsecond,
						time.Duration(0),
						"parsing arguments",
					},
					{
						true,
						time.Duration(23515) * time.Microsecond,
						time.Duration(278) * time.Microsecond,
						time.Duration(278) * time.Microsecond,
						"require('vim.shared')",
					},
					{
						true,
						time.Duration(23699) * time.Microsecond,
						time.Duration(101) * time.Microsecond,
						time.Duration(101) * time.Microsecond,
						"require('vim._meta')",
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

func TestParseErrors(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "__test_parse_errors")
	if err != nil {
		panic(err)
	}
	tmpname := tmpfile.Name()
	tmpfile.Close()
	defer func() {
		os.Remove(tmpname)
	}()

	header := []string{
		"",
		"",
		"times in msec",
		"clock   self+sourced   self:  sourced script",
		"clock   elapsed:              other lines",
		"",
	}

	for _, tc := range []struct {
		what  string
		lines []string
		msg   string
		line  uint
	}{
		{
			what:  "empty file",
			lines: []string{""},
			msg:   "Broken --startuptime output while parsing file",
		},
		{
			what:  "empty line",
			lines: append(header, ""),
			msg:   "Lack of fields: ''",
			line:  7,
		},
		{
			what: "empty line middle of lines",
			lines: append(header,
				"000.008  000.008: --- VIM STARTING ---",
				"",
				"000.190  000.182: Allocated generic buffers",
			),
			msg:  "Lack of fields: ''",
			line: 8,
		},
		{
			what:  "invalid float at elapsed time",
			lines: append(header, "00-.008  000.008: --- VIM STARTING ---"),
			msg:   "Cannot parse field '00-.008' as millisec duration",
			line:  7,
		},
		{
			what:  "invalid float at self+source",
			lines: append(header, "000.008  000.a08: --- VIM STARTING ---"),
			msg:   "Cannot parse field '000.a08' as millisec duration",
			line:  7,
		},
		{
			what:  "invalid float at self",
			lines: append(header, "198.369  000.161  000.!61: sourcing /foo/bar.vim"),
			msg:   "Cannot parse field '000.!61' as millisec duration",
			line:  7,
		},
		{
			what:  "script name is not existing",
			lines: append(header, "198.369  000.161  000.161: sourcing"),
			msg:   "Script name is missing",
			line:  7,
		},
		{
			what:  "empty description",
			lines: append(header, "198.369  000.161  000.161:"),
			msg:   "Too few fields",
			line:  7,
		},
		{
			what:  "'sourcing' is missing",
			lines: append(header, "198.369  000.161  000.161: /foo/bar.vim foo"),
			msg:   "'sourcing' token or 'require(...)' token is expected but got '/foo/bar.vim'",
			line:  7,
		},
	} {
		t.Run(tc.what, func(t *testing.T) {
			content := []byte(strings.Join(tc.lines, "\n") + "\n")
			if err := ioutil.WriteFile(tmpname, content, 0644); err != nil {
				panic(err)
			}

			f, err := os.Open(tmpname)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			_, err = parseStartuptime(f)
			if err == nil {
				t.Fatal("Error did not happen:", tc.msg)
			}

			msg := err.Error()
			if !strings.Contains(msg, tc.msg) {
				t.Fatalf("Unexpected error. '%s' is not in '%s'", tc.msg, msg)
			}
			if tc.line != 0 && !strings.Contains(msg, fmt.Sprintf("Parse error at line:%d:", tc.line)) {
				t.Fatal("Error should occur at line", tc.line, "(error:", msg, ")")
			}
		})
	}
}
