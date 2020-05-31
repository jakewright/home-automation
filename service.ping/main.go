package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/router"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

func main() {
	svc, err := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.ping",
	})
	if err != nil {
		slog.Panicf("Failed to initialise service: %v", err)
	}

	// The router has a default ping handler defined
	// in: libraries/go/router/middleware.go
	r := router.New()
	svc.Run(r)
}
