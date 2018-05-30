package main

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCollectMeasurementsOK(t *testing.T) {
	opts := &options{count: 2, vimPath: "vim", extraArgs: []string{"-N", "-u", "NONE"}}
	collected, err := collectMeasurements(opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(collected.total) != 2 {
		t.Error("2 total times should be collected but", len(collected.total))
	}
	for _, d := range collected.total {
		if d == time.Duration(0) {
			t.Error("Zero duration in collected total times:", collected.total)
		}
	}
	if len(collected.entries) == 0 {
		t.Error("Collected entries are empty")
	}
	for s, ds := range collected.entries {
		if len(ds) < 2 {
			t.Error("Source time for", s, " should be collected twice but", ds)
		}
		for _, d := range ds {
			if d == time.Duration(0) {
				t.Error("Zero duration in collected times for", s, ":", ds)
			}
		}
	}
}

func TestCollectMeasurementsVimStartError(t *testing.T) {
	opts := &options{count: 2, vimPath: "vim", extraArgs: []string{"--foo"}}
	_, err := collectMeasurements(opts)
	if err == nil {
		t.Fatal("No error occurred")
	}
	if !strings.Contains(err.Error(), "Failed to run vim with args [--foo") {
		t.Fatal("Unexpected error:", err.Error())
	}
}

func TestMeasureAverageStartuptime(t *testing.T) {
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
	want := &averageMeasurement{
		total: time.Duration(20000),
		sortedEntries: []averageEntry{
			averageEntry{"$VIM/vimrc", time.Duration(1400)},
			averageEntry{"/foo/bar", time.Duration(120)},
		},
	}

	have, err := measureAverageStartuptime(collected)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, have) {
		t.Fatal("Unexpected average value", have, ", wanted", want)
	}
}

func TestMeasureAverageStartuptimeError(t *testing.T) {
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
			_, err := measureAverageStartuptime(tc.collected)
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
