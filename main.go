package main

import (
	"fmt"
	"os"

	"github.com/yuya-takeyama/argf"
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

func printErr(err error) {
	fmt.Fprintln(os.Stderr, "gdo:", err)
}

func _main() int {
	opt, err := ParseOption(os.Args[1:])
	if err != nil {
		printErr(err)
		guideToHelp()
		return 2
	}
	if opt.IsHelp {
		usage()
		return 0
	}

	l, err := NewLines(opt)
	if err != nil {
		printErr(err)
		guideToHelp()
		return 2
	}
	r, err := argf.From(opt.Files)
	if err != nil {
		printErr(err)
		guideToHelp()
		return 2
	}

	if err = l.LoadLines(r); err != nil {
		printErr(err)
		return 1
	}
	if err = l.Flush(os.Stdout); err != nil {
		printErr(err)
		return 1
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
