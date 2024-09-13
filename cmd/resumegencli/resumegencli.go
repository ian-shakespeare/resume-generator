package main

import (
	"log"
	"os"
	"resumegenerator/internal/cli"
)

func main() {
	so := log.New(os.Stdout, "", 0)
	se := log.New(os.Stderr, "", 0)

	p, err := cli.NewArgParser([]cli.Flag{
		{Name: "version", Markers: []string{"-v", "--version"}, Description: "Print version", HasValue: false},
		{Name: "help", Markers: []string{"-h", "--help"}, Description: "Print usage and this help message.", HasValue: false},
	})
	if err != nil {
		se.Fatal(err)
	}

	args, err := p.Parse(os.Args)
	if err != nil {
		se.Fatal(err)
	}

	if args.Flags["help"] == "true" {
		so.Println("Help my ass")
	}
}
