package main

import (
	"os/exec"
	"reflect"
	"regexp"
	"strings"
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

func TestLoadLines(t *testing.T) {
	src := strings.NewReader(`
abc
123
456
def
ghi
789
jkl
mno
`[1:])
	expr := `\d+`
	expect := &Lines{
		lines:          []string{"abc", "123", "456", "def", "ghi", "789", "jkl", "mno"},
		matchedLines:   []string{"123", "456", "789"},
		matchedIndexes: map[int]bool{1: true, 2: true, 5: true},
	}
	actual, err := LoadLines(src, expr)
	if err != nil {
		t.Errorf("LoadLines(%v) returns %q, want nil",
			src, err)
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("LoadLines(%v) = %v, want %v",
			src, actual, expect)
	}
}
