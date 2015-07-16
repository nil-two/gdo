package main

import (
	"os"
)

func guideToHelp() {
	os.Stderr.WriteString(`
Try 'gdo --help' for more information.
`[1:])
}
