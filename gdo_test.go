package main

import (
	"os/exec"
	"reflect"
	"regexp"
	"testing"
)

func TestNewMatcher(t *testing.T) {
	expr := "a[bc](def)"
	m, err := NewMatcher(expr)
	if err != nil {
		t.Errorf("NewMatcher(%q) returns %q, want nil",
			expr, err)
	}

	expect, err := regexp.Compile(expr)
	if err != nil {
		t.Errorf("regexp.Compile(%q) returns %q, want nil",
			expr, err)
	}
	actual := m.re
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %q, want %q", actual, expect)
	}
}

func TestMatch(t *testing.T) {
	expr := `\d+`
	m, err := NewMatcher(expr)
	if err != nil {
		t.Errorf("NewMatcher(%q) returns %q, want nil",
			expr, err)
	}

	src1 := "xxx123"
	expect1 := true
	actual1 := m.MatchString(src1)
	if expect1 != actual1 {
		t.Errorf("MatchString(%q) = %v, want %v",
			expr, actual1, expect1)
	}

	src2 := "xxxabc"
	expect2 := false
	actual2 := m.MatchString(src2)
	if expect2 != actual2 {
		t.Errorf("MatchString(%q) = %v, want %v",
			expr, actual2, expect2)
	}
}

func TestNewProcessor(t *testing.T) {
	name := "mkdir"
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

func TestNewLines(t *testing.T) {
	expr := `\d+`
	name, arg := "sed", "s/true/false/g"
	if _, err := exec.LookPath(name); err != nil {
		t.Skipf("%q: doesn't exist", name)
	}

	m, err := NewMatcher(expr)
	if err != nil {
		t.Errorf("NewMatcher(%q) returns %q, want nil",
			name, err)
	}
	p, err := NewProcessor(name, arg)
	if err != nil {
		t.Errorf("NewProcessor(%q, %q) returns %q, want nil",
			name, arg, err)
	}

	expect := &Lines{
		matcher:        m,
		processor:      p,
		lines:          []string{},
		matchedLines:   []string{},
		matchedIndexes: make(map[int]bool),
	}
	actual := NewLines(m, p)
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %v, want %v", actual, expect)
	}
}
