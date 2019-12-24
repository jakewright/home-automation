package main

import (
	"log"
	"os"

	"github.com/jakewright/home-automation/tools/deploy/cmd"
)

// BuildDirectory is injected at compile time
var BuildDirectory string

func main() {
	if cwd, err := os.Getwd(); err != nil {
		log.Fatal(err)
	} else if cwd != BuildDirectory {
		//log.Fatalf("Must be run from home-automation root: %s\n", BuildDirectory)
	}

	// Load the configuration

	cmd.Execute()
}
