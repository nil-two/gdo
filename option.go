package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

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
