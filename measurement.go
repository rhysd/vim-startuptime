package main

import (
	"fmt"
	"io"
	"time"
)

type measurementEntry struct {
	script  bool
	elapsed time.Duration
	total   time.Duration
	self    time.Duration
	name    string
}

type measurement struct {
	elapsedTotal time.Duration
	entries      []*measurementEntry
}

type averageEntry struct {
	name     string
	duration time.Duration
}

type averageMeasurement struct {
	total         time.Duration
	sortedEntries []averageEntry
}

func (ave *averageMeasurement) printSummary(w io.Writer) {
	fmt.Fprintf(w, "Total: %f msec\n\n", ave.total.Seconds()*1000)
	for _, e := range ave.sortedEntries {
		fmt.Fprintf(w, "%f: %s\n", e.duration.Seconds()*1000, e.name)
	}
}
