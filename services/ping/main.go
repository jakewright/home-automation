package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
)

func main() {
	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.ping",
	})

	// The router has a default ping handler defined
	// in: libraries/go/router/middleware.go
	svc.Run()
}
