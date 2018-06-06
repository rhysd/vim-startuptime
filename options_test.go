package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseExtraArgs(t *testing.T) {
	saved := os.Args
	defer func() {
		os.Args = saved
	}()
	os.Args = []string{os.Args[0], "-count", "1", "--", "--cmd", "Foo"}
	opts := parseOptions()
	want := []string{"--cmd", "Foo"}
	if !reflect.DeepEqual(opts.extraArgs, want) {
		t.Fatal("Wanted:", want, ", but have:", opts.extraArgs)
	}
}

func TestUsage(t *testing.T) {
	saved := os.Stderr
	defer func() {
		os.Stderr = saved
	}()

	dir, err := ioutil.TempDir("", "__vim_usage_test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	tempName := filepath.Join(dir, "stderr")
	temp, err := os.Create(tempName)
	if err != nil {
		panic(err)
	}
	os.Stderr = temp

	usage()

	temp.Close()

	bytes, err := ioutil.ReadFile(tempName)
	if err != nil {
		panic(err)
	}

	out := string(bytes)
	if !strings.Contains(out, "Usage: vim-startuptime") {
		t.Fatal("Unexpected usage output:", out)
	}
}
