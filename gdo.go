package main

import (
	"os/exec"
	"regexp"
)

type Matcher struct {
	re *regexp.Regexp
}

func NewMatcher(expr string) (m *Matcher, err error) {
	m = &Matcher{}
	m.re, err = regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Matcher) MatchString(s string) bool {
	return m.re.MatchString(s)
}

type Processor struct {
	cmd *exec.Cmd
}
