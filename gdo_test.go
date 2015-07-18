package main

import (
	"bytes"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

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
	if !reflect.DeepEqual(actual.lines, expect.lines) {
		t.Errorf("lines got %q, want %q", actual, expect)
	}
	if !reflect.DeepEqual(actual.matchedLines, expect.matchedLines) {
		t.Errorf("matchedLines got %q, want %q", actual, expect)
	}
	if !reflect.DeepEqual(actual.matchedIndexes, expect.matchedIndexes) {
		t.Errorf("matchedIndexes got %q, want %q", actual, expect)
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
