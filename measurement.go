package main

import (
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
