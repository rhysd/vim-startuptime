package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseErrorAt(lineno uint, format string, args ...interface{}) error {
	s := fmt.Sprintf("parse error at line:%d: %s", lineno, format)
	return fmt.Errorf(s, args...)
}

func parseDuration(s string, lineno uint) (time.Duration, error) {
	s = strings.TrimSuffix(s, ":")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return time.Duration(0), parseErrorAt(lineno, "cannot parse field '%s' as millisec duration", s)
	}
	return time.Duration(f*1000) * time.Microsecond, nil
}

func parseStartuptimeEntity(line string, lineno uint) (*measurementEntry, error) {
	e := &measurementEntry{}

	ss := strings.Fields(line)
	if len(ss) <= 2 {
		return nil, parseErrorAt(lineno, "lack of fields: '%s'", line)
	}

	d, err := parseDuration(ss[0], lineno)
	if err != nil {
		return nil, err
	}
	e.elapsed = d

	d, err = parseDuration(ss[1], lineno)
	if err != nil {
		return nil, err
	}
	e.total = d

	if strings.HasSuffix(ss[1], ":") {
		e.script = false
		e.self = time.Duration(0)
		e.name = strings.Join(ss[2:], " ")
		return e, nil
	}

	e.script = true
	if len(ss) < 4 {
		return nil, parseErrorAt(lineno, "failed to parse script measurement line '%s'. too few fields", line)
	}

	d, err = parseDuration(ss[2], lineno)
	if err != nil {
		return nil, err
	}
	e.self = d

	if ss[3] == "sourcing" {
		e.name = strings.Join(ss[4:], " ")
		if e.name == "" {
			return nil, parseErrorAt(lineno, "failed to parse script measurement line '%s'. script name is missing", line)
		}
	} else if strings.HasPrefix(ss[3], "require(") {
		e.name = strings.Join(ss[3:], " ")
	} else {
		return nil, parseErrorAt(lineno, "'sourcing' token or 'require(...)' token is expected but got '%s'", ss[3])
	}

	return e, nil
}

func parseStartuptime(file *os.File) (*measurement, error) {
	m := &measurement{}

	s := bufio.NewScanner(file)
	l := uint(1)
	for s.Scan() {
		if l < 7 {
			// Skip header
			l++
			continue
		}
		t := s.Text()
		if t == "" {
			// Neovim appends an extra empty line at the end of input (#4)
			continue
		}
		e, err := parseStartuptimeEntity(t, l)
		if err != nil {
			return nil, err
		}
		m.entries = append(m.entries, e)
		l++
	}

	if len(m.entries) == 0 {
		return nil, fmt.Errorf("broken --startuptime output while parsing file %s. no entry was parsed", file.Name())
	}
	m.elapsedTotal = m.entries[len(m.entries)-1].elapsed

	return m, nil
}
