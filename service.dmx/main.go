package main

import (
	"context"

	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.dmx/handler"
	"github.com/jakewright/home-automation/service.dmx/universe"
)

//go:generate jrpc dmx.def

func main() {
	conf := struct{ UniverseNumber uint8 }{}

	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.dmx",
		Config:      &conf,
		Firehose:    true,
	})

	u := universe.New(conf.UniverseNumber)

	l := universe.Loader{
		ServiceName: "service.dmx",
		Universe:    u,
	}

	if err := l.FetchDevices(context.Background()); err != nil {
		slog.Panicf("Failed to load devices: %v", err)
	}

	r := handler.NewRouter(&handler.Controller{
		Universe: u,
	})

	svc.Run(r)
}
