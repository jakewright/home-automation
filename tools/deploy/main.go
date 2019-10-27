package main

import (
	"log"
	"os"

	"github.com/jakewright/home-automation/tools/deploy/build"
)

// BuildDirectory is injected at compile time
var BuildDirectory string

func main() {
	if cwd, err := os.Getwd(); err != nil {
		log.Fatal(err)
	} else if cwd != BuildDirectory {
		//log.Fatalf("Must be run from home-automation root: %s\n", BuildDirectory)
	}

	log.Printf("Checkout out code")
	if err := build.Checkout("master"); err != nil {
		log.Fatal(err)
	}
}
