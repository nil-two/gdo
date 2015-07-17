package main

import (
	"bufio"
	"fmt"
	"io"
)

type Lines struct {
	matcher        *Matcher
	processor      *Processor
	lines          []string
	matchedLines   []string
	matchedIndexes map[int]bool
}

func NewLines(opt *Option) (l *Lines, err error) {
	if opt == nil {
		opt = &Option{}
	}
	m, err := NewMatcher(opt.Pattern)
	if err != nil {
		return nil, err
	}
	p, err := NewProcessor(opt.Command, opt.Arg...)
	if err != nil {
		return nil, err
	}
	return &Lines{
		matcher:        m,
		processor:      p,
		matchedIndexes: make(map[int]bool),
	}, nil
}

func (l *Lines) LoadLines(r io.Reader) error {
	b := bufio.NewScanner(r)
	for i := 0; b.Scan(); i++ {
		line := b.Text()
		if l.matcher.MatchString(line) {
			l.matchedLines = append(l.matchedLines, line)
			l.matchedIndexes[i] = true
		}
		l.lines = append(l.lines, line)
	}
	return b.Err()
}

func (l *Lines) Flush(out io.Writer) error {
	if err := l.processor.Process(l.matchedLines); err != nil {
		return err
	}
	mi := 0
	for li := 0; li < len(l.lines); li++ {
		if l.matchedIndexes[li] {
			fmt.Fprintln(out, l.matchedLines[mi])
			mi++
		} else {
			fmt.Fprintln(out, l.lines[li])
		}
	}
	return nil
}
