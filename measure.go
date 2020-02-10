package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type collectedMeasurements struct {
	total   []time.Duration
	entries map[string][]time.Duration
}

func collectMeasurements(opts *options) (*collectedMeasurements, error) {
	for i := uint(0); i < opts.warmup; i++ {
		if err := runVim(opts.vimPath, opts.extraArgs); err != nil {
			return nil, fmt.Errorf("Failed while warmup: %v", err)
		}
	}

	dir, err := ioutil.TempDir("", "__vim_startuptime_")
	if err != nil {
		return nil, fmt.Errorf("Failed to open temporary directory: %v", err)
	}
	defer os.RemoveAll(dir)

	collected := &collectedMeasurements{entries: map[string][]time.Duration{}}
	for id := 0; id < int(opts.count); id++ {
		f, err := runVimStartuptime(opts.vimPath, dir, id, opts.extraArgs)
		if err != nil {
			return nil, err
		}
		m, err := parseStartuptime(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		collected.total = append(collected.total, m.elapsedTotal)
		for _, e := range m.entries {
			if !opts.script || e.script {
				collected.entries[e.name] = append(collected.entries[e.name], e.total)
			}
		}
		f.Close()
	}

	return collected, nil
}

func summarizeEntry(name string, ds []time.Duration) entrySummary {
	total := time.Duration(0)
	min := ds[0]
	max := ds[0]
	for _, d := range ds {
		total += d
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}
	average := time.Duration(total.Nanoseconds()/int64(len(ds))) * time.Nanosecond
	return entrySummary{name, average, max, min}
}

func summarizeStartuptime(collected *collectedMeasurements) (*measurementSummary, error) {
	summary := &measurementSummary{}
	summary.sortedEntries = make([]entrySummary, 0, len(collected.entries))
	if len(collected.total) == 0 {
		return nil, fmt.Errorf("No total time was collected")
	}
	summary.total = summarizeEntry("Total", collected.total)
	for n, ds := range collected.entries {
		if len(ds) == 0 {
			return nil, fmt.Errorf("No profile was collected for '%s'", n)
		}
		summary.sortedEntries = append(summary.sortedEntries, summarizeEntry(n, ds))
	}

	// Sort in decending order by duration
	sort.Slice(summary.sortedEntries, func(i, j int) bool {
		return summary.sortedEntries[i].average > summary.sortedEntries[j].average
	})

	return summary, nil
}
