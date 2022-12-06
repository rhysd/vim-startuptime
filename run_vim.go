package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	if runtime.GOOS == "windows" {
		if p, err := exec.LookPath(vimpath); err == nil {
			if filepath.Base(p) == "vim.exe" {
				a = append(a, "-e")
			}
		}
	}
	a = append(a, "-c", "qall!")
	a = append(a, args...)

	cmd := exec.Command(vimpath, a...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		for i, b := range out {
			if b == '\n' || b == '\r' {
				out[i] = ' '
			}
		}
		return fmt.Errorf("failed to run %q with args %v: %w. Output: %s", vimpath, a, err, string(out))
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
		return nil, fmt.Errorf("could not open --startuptime result file '%s': %w", outfile, err)
	}
	return f, nil
}
