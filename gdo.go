package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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

func NewProcessor(name string, arg ...string) (p *Processor, err error) {
	if _, err = exec.LookPath(name); err != nil {
		return nil, err
	}
	p = &Processor{}
	p.cmd = exec.Command(name, arg...)
	return p, nil
}

func (p *Processor) Process(a []string) error {
	in, err := p.cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer out.Close()

	if err = p.cmd.Start(); err != nil {
		return err
	}
	for _, s := range a {
		fmt.Fprintln(in, s)
	}
	if err = in.Close(); err != nil {
		return err
	}

	b := bufio.NewScanner(out)
	for i := 0; i < len(a) && b.Scan(); i++ {
		a[i] = b.Text()
	}
	return b.Err()
}

type Option struct {
	IsHelp  bool
	Pattern string
	Command string
	Arg     []string
	Files   []string
}

func ParseOption(args []string) (opt *Option, err error) {
	opt = &Option{}
	f := flag.NewFlagSet("gdo", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	f.BoolVar(&opt.IsHelp, "h", false, "")
	f.BoolVar(&opt.IsHelp, "help", false, "")
	if err = f.Parse(args); err != nil {
		return nil, err
	}
	switch f.NArg() {
	case 0:
		return nil, fmt.Errorf("no specify PATTERN and COMMAND")
	case 1:
		return nil, fmt.Errorf("no specify COMMAND")
	}
	opt.Pattern = f.Arg(0)
	opt.Command = f.Arg(1)

	var finishArg bool
	for _, arg := range f.Args()[2:] {
		switch {
		case finishArg:
			opt.Files = append(opt.Files, arg)
		case arg == "--":
			finishArg = true
		default:
			opt.Arg = append(opt.Arg, arg)
		}
	}
	return opt, nil
}

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

func (l *Lines) Flush(out io.Writer, p *Processor) error {
	if err := p.Process(l.matchedLines); err != nil {
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
