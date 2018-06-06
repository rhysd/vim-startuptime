package main

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCollectMeasurementsOK(t *testing.T) {
	for _, path := range []string{"vim", "nvim"} {
		t.Run(path, func(t *testing.T) {
			opts := &options{count: 2, vimPath: path, extraArgs: []string{"-N", "-u", "NONE"}}
			collected, err := collectMeasurements(opts)
			if err != nil {
				t.Fatal(err)
			}
			if len(collected.total) != 2 {
				t.Error("2 total times should be collected but", len(collected.total))
			}
			if len(collected.entries) == 0 {
				t.Error("Collected entries are empty")
			}
			for s, ds := range collected.entries {
				if len(ds) < 2 {
					t.Error("Source time for", s, " should be collected twice but", ds)
				}
			}
		})
	}
}

func TestCollectMeasurementsVimStartError(t *testing.T) {
	for _, path := range []string{"vim", "nvim"} {
		t.Run(path, func(t *testing.T) {
			opts := &options{count: 2, vimPath: path, extraArgs: []string{"--foo"}}
			_, err := collectMeasurements(opts)
			if err == nil {
				t.Fatal("No error occurred")
			}
			if !strings.Contains(err.Error(), "Failed to run "+path+" with args [--foo") {
				t.Fatal("Unexpected error:", err.Error())
			}
		})
	}
}

func TestSummarizeStartuptime(t *testing.T) {
	collected := &collectedMeasurements{
		total: []time.Duration{
			time.Duration(10000),
			time.Duration(20000),
			time.Duration(30000),
		},
		entries: map[string][]time.Duration{
			"/foo/bar": []time.Duration{
				time.Duration(110),
				time.Duration(130),
				time.Duration(120),
			},
			"$VIM/vimrc": []time.Duration{
				time.Duration(1500),
				time.Duration(1400),
				time.Duration(1300),
			},
		},
	}
	want := &measurementSummary{
		total: entrySummary{
			"Total",
			time.Duration(20000),
			time.Duration(30000),
			time.Duration(10000),
		},
		sortedEntries: []entrySummary{
			entrySummary{
				"$VIM/vimrc",
				time.Duration(1400),
				time.Duration(1500),
				time.Duration(1300),
			},
			entrySummary{
				"/foo/bar",
				time.Duration(120),
				time.Duration(130),
				time.Duration(110),
			},
		},
	}

	have, err := summarizeStartuptime(collected)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, have) {
		t.Fatal("Unexpected average value", have, ", wanted", want)
	}
}

func TestSummarizeStartuptimeError(t *testing.T) {
	for _, tc := range []struct {
		what      string
		collected *collectedMeasurements
		msg       string
	}{
		{
			"No total time",
			&collectedMeasurements{
				total: []time.Duration{},
				entries: map[string][]time.Duration{
					"/foo/bar": []time.Duration{time.Duration(110)},
				},
			},
			"No total time was collected",
		},
		{
			"No entry profile result",
			&collectedMeasurements{
				total: []time.Duration{time.Duration(110)},
				entries: map[string][]time.Duration{
					"$VIM/vimrc": []time.Duration{},
				},
			},
			"No profile was collected for '$VIM/vimrc'",
		},
	} {
		t.Run(tc.what, func(t *testing.T) {
			_, err := summarizeStartuptime(tc.collected)
			if err == nil {
				t.Fatal("Error should happen")
			}
			msg := err.Error()
			if !strings.Contains(msg, tc.msg) {
				t.Fatalf("Unexpected error '%s', it should contain '%s'", msg, tc.msg)
			}
		})
	}
}
