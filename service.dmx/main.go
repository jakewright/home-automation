package main

import (
	"context"

	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/config"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.dmx/handler"
	"github.com/jakewright/home-automation/service.dmx/universe"
)

//go:generate jrpc dmx.def

func main() {
	svc, err := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.dmx",
		Firehose:    true,
	})

	if err != nil {
		slog.Panicf("Failed to initialise service: %v", err)
	}

	universeNumber := config.Get("universe.number").Int()
	u := universe.New(universeNumber)

	l := universe.Loader{
		ServiceName: "service.dmx",
		Universe:    u,
	}

	if err := l.FetchDevices(context.Background()); err != nil {
		slog.Panicf("Failed to load devices: %v", err)
	}

	h := handler.DMXHandler{Universe: u}

	r := handler.NewRouter()
	r.GetDevice = h.Read
	r.UpdateDevice = h.Update

	svc.Run(r)
}
