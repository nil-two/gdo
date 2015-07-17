package main

import (
	"bytes"
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

func TestNewLines(t *testing.T) {
	expr := `false`
	name, arg := "sed", []string{"s/.*/true/g"}
	if _, err := exec.LookPath(name); err != nil {
		t.Skipf("%q: doesn't exist", name)
	}

	m, err := NewMatcher(expr)
	if err != nil {
		t.Errorf("NewMatcher(%q) returns %q, want nil",
			expr, err)
	}
	p, err := NewProcessor(name, arg...)
	if err != nil {
		t.Errorf("NewProcessor(%q, %q) returns %q, want nil",
			name, arg, err)
	}

	opt := &Option{
		Pattern: expr,
		Command: name,
		Arg:     arg,
	}

	expect := &Lines{
		matcher:        m,
		processor:      p,
		lines:          []string{},
		matchedLines:   []string{},
		matchedIndexes: make(map[int]bool),
	}
	actual, err := NewLines(opt)
	if err != nil {
		t.Errorf("NewLines(%v) returns %q, want nil",
			opt, err)
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %v, want %v", actual, expect)
	}
}

func TestLoadLines(t *testing.T) {
	expr := `\d+`
	m, err := NewMatcher(expr)
	if err != nil {
		t.Errorf("NewMatcher(%q) returns %q, want nil",
			expr, err)
	}
	src := strings.NewReader(`
abc
123
def
456
789
ghi
jkl
mno
`[1:])
	expect := &Lines{
		lines:          []string{"abc", "123", "def", "456", "789", "ghi", "jkl", "mno"},
		matchedLines:   []string{"123", "456", "789"},
		matchedIndexes: map[int]bool{1: true, 3: true, 4: true},
	}
	actual := &Lines{}
	actual.matchedIndexes = make(map[int]bool)
	actual.matcher = m
	if err = actual.LoadLines(src); err != nil {
		t.Errorf("NewLines(%v).LoadLines(%v) returns %q, want nil",
			m, src, err)
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %v, want %v", actual, expect)
	}
}

func TestFlush(t *testing.T) {
	name, arg := "sed", "s/.*/true/"
	if _, err := exec.LookPath(name); err != nil {
		t.Skipf("%q: doesn't exist", name)
	}
	p, err := NewProcessor(name, arg)
	if err != nil {
		t.Errorf("NewProcessor(%q, %q) returns %q, want nil",
			name, arg, err)
	}
	l := &Lines{
		processor:      p,
		lines:          []string{"true", "false", "false", "true", "nil", "true", "false"},
		matchedLines:   []string{"false", "false", "false"},
		matchedIndexes: map[int]bool{1: true, 2: true, 6: true},
	}

	out := bytes.NewBuffer(make([]byte, 0))
	if l.Flush(out); err != nil {
	}
	expect := `
true
true
true
true
nil
true
true
`[1:]
	actual := out.String()
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %q, want %q", actual, expect)
	}
}
