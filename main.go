package main

import (
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
