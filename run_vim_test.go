package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunVimOK(t *testing.T) {
	for _, exe := range []string{"vim", "nvim"} {
		t.Run(exe, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "__vim_run_test_")
			if err != nil {
				panic(err)
			}
			defer os.RemoveAll(dir)

			f, err := runVimStartuptime(exe, dir, 3, []string{"-u", "NONE", "-N"})
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			_, fname := filepath.Split(f.Name())
			if fname != "3" {
				t.Error("Invalid result file name", fname, "Wanted '3'")
			}
			bytes, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal("Cannot open startup result file", err)
			}
			content := string(bytes)
			if !strings.Contains(content, "--- "+strings.ToUpper(exe)+" STARTING ---") {
				t.Error("Invalid result file:", content)
			}
		})
	}
}

func TestStartError(t *testing.T) {
	for _, exe := range []string{"vim", "nvim"} {
		t.Run(exe, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "__vim_run_test_error_")
			if err != nil {
				panic(err)
			}
			defer os.RemoveAll(dir)

			_, err = runVimStartuptime(exe, dir, 3, []string{"--foo", "--bar"})
			if err == nil {
				t.Fatal("Invalid extra args should cause an error")
			}
			want := "Failed to run \"" + exe + "\" with args [--foo --bar"
			if !strings.Contains(err.Error(), want) {
				t.Errorf("Wanted %q but got %q", want, err)
			}
		})
	}
}
