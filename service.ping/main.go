package main

import (
	"log"
	"os"

	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set config manually so it doesn't need s.config to be running
	config.DefaultProvider = config.New(map[string]interface{}{
		"port": port,
	})

	// The router has a default ping handler defined
	// in: libraries/go/router/middleware.go
	r := router.New()
	if err := r.Start(); err != nil {
		log.Fatal(err)
	}
}
