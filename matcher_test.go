package main

import (
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
