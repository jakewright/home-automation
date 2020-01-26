package main

import (
	"os"

	"github.com/jakewright/home-automation/libraries/go/svcdef"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: jrpc file.def")
		os.Exit(1)
	}

	defPath := os.Args[1]

	f, err := svcdef.Parse(defPath)
	if err != nil {
		panic(err)
	}

	if err := generate(defPath, f); err != nil {
		panic(err)
	}
}
