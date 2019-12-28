package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/router"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.dmx/device"
	"github.com/jakewright/home-automation/service.dmx/domain"
	"github.com/jakewright/home-automation/service.dmx/handler"
)

func main() {
	svc, err := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.dmx",
		Firehose:    true,
	})

	if err != nil {
		slog.Panic("Failed to initialise service: %v", err)
	}

	universeNumber := config.Get("universe.number").Int()
	u := &domain.Universe{Number: universeNumber}

	l := device.Loader{
		ServiceName: "service.dmx",
		Universe:    u,
	}

	if err := l.FetchDevices(); err != nil {
		slog.Panic("Failed to load devices: %v", err)
	}

	h := handler.DMXHandler{Universe: u}

	r := router.New()
	r.Get("/{device_id}", h.Read)
	r.Patch("/{device_id}", h.Update)

	svc.Run(r)
}
