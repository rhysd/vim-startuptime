package main

import (
	"fmt"
	"io"
	"math"
	"strings"
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

func alignFloatColumn(data []float64, header string) []string {
	intPart := 0
	for _, f := range data {
		d := int(math.Log10(f)) + 1
		if d > intPart {
			intPart = d
		}
	}
	width := intPart + 1 /*'.'*/ + 6 /*decimal part*/

	aligned := make([]string, 0, len(data)+1)

	if len(header) < width {
		header = fmt.Sprintf("%s%s", strings.Repeat(" ", width-len(header)), header)
	}
	aligned = append(aligned, header)

	for _, f := range data {
		s := fmt.Sprintf("%f", f)
		if len(s) < width {
			s = fmt.Sprintf("%s%s", strings.Repeat(" ", width-len(s)), s)
		}
		aligned = append(aligned, s)
	}

	return aligned
}

func (ave *averageMeasurement) printSummary(w io.Writer) {
	fmt.Fprintf(w, "Total: %f msec\n\n", ave.total.Seconds()*1000)

	averages := make([]float64, 0, len(ave.sortedEntries))
	for _, e := range ave.sortedEntries {
		averages = append(averages, e.duration.Seconds()*1000)
	}

	averageColumn := alignFloatColumn(averages, "AVERAGE")
	fmt.Fprintln(w, averageColumn[0])
	for i, e := range ave.sortedEntries {
		fmt.Fprintf(w, "%s: %s\n", averageColumn[i+1], e.name)
	}
}
