package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func isNeovimPath(vimpath string) bool {
	suff := strings.TrimSuffix(vimpath, ".exe")
	if strings.HasSuffix(suff, "nvim") {
		return true
	}
	if strings.HasSuffix(suff, "lvim") {
		return true
	}
	return false
}

func runVim(vimpath string, extra []string, args ...string) error {
	a := make([]string, 0, len(extra)+3+len(args))
	a = append(a, extra...)
	if isNeovimPath(vimpath) {
		a = append(a, "--headless")
	} else {
		a = append(a, "--not-a-term")
	}
	a = append(a, "-c", "qall!")
	a = append(a, args...)

	cmd := exec.Command(vimpath, a...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to run %q with args %v: %s. Output: %s", vimpath, a, err, string(out))
	}
	return nil
}

func runVimStartuptime(vimpath, tmpdir string, id int, extra []string) (*os.File, error) {
	outfile := filepath.Join(tmpdir, strconv.Itoa(id))
	if err := runVim(vimpath, extra, "--startuptime", outfile); err != nil {
		return nil, err
	}

	f, err := os.Open(outfile)
	if err != nil {
		return nil, fmt.Errorf("Could not open --startuptime result file '%s': %v", outfile, err)
	}
	return f, nil
}
