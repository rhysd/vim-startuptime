package main

import (
	"os"
	"reflect"
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
