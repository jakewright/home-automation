package main

import (
	"github.com/jakewright/home-automation/libraries/go/bootstrap"
	"github.com/jakewright/home-automation/service.dmx-proxy/dmx"
	"github.com/jakewright/home-automation/service.dmx-proxy/handler"
)

//go:generate jrpc dmxproxy.def

func main() {
	svc := bootstrap.Init(&bootstrap.Opts{
		ServiceName: "service.dmx-proxy",
	})

	r := handler.NewRouter(&handler.Controller{
		Setter: &dmx.OLA{},
	})

	svc.Run(r)
}
