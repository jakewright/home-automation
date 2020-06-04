package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/libraries/go/slog"
	"github.com/jakewright/home-automation/service.dmx-proxy/dmx"
	"github.com/jakewright/home-automation/service.dmx-proxy/handler"
)

//go:generate jrpc dmxproxy.def

func main() {
	svc, err := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.dmx-proxy",
	})
	if err != nil {
		slog.Panicf("Failed to initialise service: %v", err)
	}

	r := handler.NewRouter(&handler.Controller{
		Setter: &dmx.OLA{},
	})

	svc.Run(r)
}
