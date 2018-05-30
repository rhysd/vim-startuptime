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
			collected.entries[e.name] = append(collected.entries[e.name], e.total)
		}
		f.Close()
	}

	return collected, nil
}

type averageEntry struct {
	name     string
	duration time.Duration
}

type averageMeasurement struct {
	total         time.Duration
	sortedEntries []averageEntry
}

func averageDuration(ds []time.Duration) time.Duration {
	total := time.Duration(0)
	for _, d := range ds {
		total += d
	}
	return time.Duration(total.Nanoseconds()/int64(len(ds))) * time.Nanosecond
}

func measureAverageStartuptime(collected *collectedMeasurements) (*averageMeasurement, error) {
	average := &averageMeasurement{}
	average.sortedEntries = make([]averageEntry, 0, len(collected.entries))
	if len(collected.total) == 0 {
		return nil, fmt.Errorf("No total time was collected")
	}
	average.total = averageDuration(collected.total)
	for n, ds := range collected.entries {
		if len(ds) == 0 {
			return nil, fmt.Errorf("No profile was collected for '%s'", n)
		}
		average.sortedEntries = append(average.sortedEntries, averageEntry{n, averageDuration(ds)})
	}

	// Sort in decending order by duration
	sort.Slice(average.sortedEntries, func(i, j int) bool {
		return average.sortedEntries[i].duration > average.sortedEntries[j].duration
	})

	return average, nil
}
