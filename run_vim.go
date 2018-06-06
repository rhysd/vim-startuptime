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
	return strings.HasSuffix(vimpath, "nvim") || strings.HasSuffix(vimpath, "nvim.exe")
}

func runVimStartuptime(vimpath, tmpdir string, id int, extra []string) (*os.File, error) {
	outfile := filepath.Join(tmpdir, strconv.Itoa(id))
	args := make([]string, 0, len(extra)+4)
	args = append(args, extra...)
	if isNeovimPath(vimpath) {
		args = append(args, "--headless")
	} else {
		args = append(args, "--not-a-term")
	}
	args = append(args, "-c", "quit", "--startuptime", outfile)

	cmd := exec.Command(vimpath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Failed to run %s with args %v: %s", vimpath, args, string(out))
	}
	f, err := os.Open(outfile)
	if err != nil {
		return nil, fmt.Errorf("Could not open --startuptime result file '%s': %v", outfile, err)
	}
	return f, nil
}
