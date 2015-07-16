package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func guideToHelp() {
	os.Stderr.WriteString(`
Try 'gdo --help' for more information.
`[1:])
}

func usage() {
	os.Stderr.WriteString(`
Usage: gdo [OPTION]... PATTERN COMMAND [ARGS]...
Process matched lines with COMMAND
`[1:])
}

func _main() int {
	f := flag.NewFlagSet("gdo", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	var isHelp bool
	f.BoolVar(&isHelp, "h", false, "")
	f.BoolVar(&isHelp, "help", false, "")
	if err := f.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "gdo:", err)
		guideToHelp()
		return 2
	}
	if isHelp {
		usage()
		return 0
	}
	switch f.NArg() {
	case 0:
		fmt.Fprintln(os.Stderr, "gdo:", "no specify PATTERN and COMMAND")
		return 2
	case 1:
		fmt.Fprintln(os.Stderr, "gdo:", "no specify COMMAND")
		return 2
	}
	pattern, name, args := f.Arg(0), f.Arg(1), f.Args()[2:]

	m, err := NewMatcher(pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gdo:", err)
		return 2
	}
	p, err := NewProcessor(name, args...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gdo:", err)
		return 2
	}

	l := NewLines()
	if err = l.LoadLines(os.Stdin, m); err != nil {
		fmt.Fprintln(os.Stderr, "gdo:", err)
		return 1
	}
	if err = l.Flush(os.Stdout, p); err != nil {
		fmt.Fprintln(os.Stderr, "gdo:", err)
		return 1
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
