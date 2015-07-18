package main

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	name := "mkdir"
	if _, err := exec.LookPath(name); err != nil {
		t.Skipf("%q: doesn't exist", name)
	}
	p, err := NewProcessor(name)
	if err != nil {
		t.Errorf("NewMatcher(%q) returns %q, want nil",
			name, err)
	}

	expect := exec.Command(name)
	actual := p.cmd
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("NewProcessor(%q).cmd = %v, want %v",
			name, actual, expect)
	}
}

func TestProcess(t *testing.T) {
	name, arg := "sed", "s/false/true/g"
	if _, err := exec.LookPath(name); err != nil {
		t.Skipf("%q: doesn't exist", name)
	}
	p, err := NewProcessor(name, arg)
	if err != nil {
		t.Errorf("NewProcessor(%q, %q) returns %q, want nil",
			name, arg, err)
	}

	expect := []string{"true", "true", "true", "nil"}
	actual := []string{"true", "false", "false", "nil"}
	if err = p.Process(actual); err != nil {
		t.Errorf("NewProcessor(%q, %q).Process(%q) returns %q, want nil",
			name, arg, actual, err)
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %q, want %q", actual, expect)
	}
}
