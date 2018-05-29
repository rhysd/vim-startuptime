package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunVimOK(t *testing.T) {
	dir, err := ioutil.TempDir("", "__vim_run_test_")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	f, err := runVimStartuptime("vim", dir, 3, []string{})
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
	if !strings.Contains(content, "--- VIM STARTING ---") {
		t.Error("Invalid result file:", content)
	}
}

func TestStartError(t *testing.T) {
	dir, err := ioutil.TempDir("", "__vim_run_test_")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	_, err = runVimStartuptime("vim", dir, 3, []string{"--foo", "--bar"})
	if err == nil {
		t.Fatal("Invalid extra args should cause an error")
	}
	if !strings.Contains(err.Error(), "Failed to run vim with args [--foo --bar") {
		t.Error("Unexpected error:", err)
	}
}
