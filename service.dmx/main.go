package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/router"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.dmx/handler"
)

func main() {
	svc, err := bootstrap.Init("service.dmx")
	if err != nil {
		slog.Panic("Failed to initialise service: %v", err)
	}

	r := router.New()
	r.Patch("/{device_id}", handler.Update)

	svc.Run(r)
}
