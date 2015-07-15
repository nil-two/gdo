package main

import (
	"regexp"
)

type Matcher struct {
	re *regexp.Regexp
}
