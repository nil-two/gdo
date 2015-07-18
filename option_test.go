package main

import (
	"reflect"
	"testing"
)

func TestParseOption(t *testing.T) {
	args := []string{`\d+`, "perl", "-ple", "$_*=2", "--", "foo", "bar"}
	expect := &Option{
		IsHelp:  false,
		Pattern: `\d+`,
		Command: "perl",
		Arg:     []string{"-ple", "$_*=2"},
		Files:   []string{"foo", "bar"},
	}
	actual, err := ParseOption(args)
	if err != nil {
		t.Errorf("ParseOption(%q) returns %q, want nil",
			args, err)
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %q, want %q", actual, expect)
	}
}
